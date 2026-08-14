package platformresearch

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
)

const defaultSearchLimit = 5

var platformSearchRules = map[string]struct {
	Name    string
	Domains []string
}{
	"douyin": {
		Name:    "抖音",
		Domains: []string{"douyin.com", "iesdouyin.com"},
	},
	"xiaohongshu": {
		Name:    "小红书",
		Domains: []string{"xiaohongshu.com"},
	},
	"meituan": {
		Name:    "美团",
		Domains: []string{"meituan.com", "dianping.com"},
	},
	"eleme": {
		Name:    "饿了么",
		Domains: []string{"ele.me", "eleme.cn"},
	},
}

type SearchQuery struct {
	Platform  string   `json:"platform"`
	Keyword   string   `json:"keyword"`
	Query     string   `json:"query"`
	MustTerms []string `json:"mustTerms,omitempty"`
	Reason    string   `json:"reason,omitempty"`
}

func BuildKeywords(merchantName, city, industry, address string, products []Product, extra []string) []string {
	keywords := []string{}
	merchantName = strings.TrimSpace(merchantName)
	city = strings.TrimSpace(city)
	industry = strings.TrimSpace(industry)
	address = strings.TrimSpace(address)

	category := inferCategory(merchantName, industry, address, products)
	district := extractDistrict(address)
	road := extractRoad(address)
	nearby := extractNearbyShop(address)

	core := strings.TrimSpace(city + " " + category)
	if core != "" {
		keywords = append(keywords,
			core+" 店铺数量",
			core+" 排行榜",
			core+" 头部店",
			core+" 哪家好",
			core+" 探店",
			core+" 团购",
			core+" 评价",
			core+" 差评",
			core+" 美团",
			core+" 大众点评",
			core+" 小红书",
			core+" 抖音",
		)
	}
	if city != "" && district != "" && category != "" {
		keywords = append(keywords,
			city+" "+district+" "+category+" 排行榜",
			city+" "+district+" "+category+" 团购",
			city+" "+district+" "+category+" 评价",
		)
	}
	if city != "" && road != "" && category != "" {
		keywords = append(keywords,
			city+" "+road+" "+category,
			road+" "+category+" 团购",
		)
	}
	if nearby != "" {
		keywords = append(keywords, nearby, city+" "+nearby+" 评价", city+" "+nearby+" 团购")
	}

	for _, product := range products {
		name := normalizeProductKeyword(product.Name, category)
		if name == "" {
			continue
		}
		if city != "" && category != "" {
			keywords = append(keywords,
				city+" "+category+" "+name,
				city+" "+category+" "+name+" 团购",
				city+" "+category+" "+name+" 评价",
			)
			continue
		}
		keywords = append(keywords, name+" 团购", name+" 评价")
	}

	for _, item := range extra {
		item = normalizeManualKeyword(item, city, category)
		if item != "" {
			keywords = append(keywords, item)
		}
	}
	return dedupeNonEmpty(keywords)
}

func DefaultSources(sources []string) []string {
	if len(sources) == 0 {
		return []string{"douyin", "xiaohongshu", "meituan", "eleme"}
	}
	result := make([]string, 0, len(sources))
	for _, source := range sources {
		source = strings.ToLower(strings.TrimSpace(source))
		if _, ok := platformSearchRules[source]; ok {
			result = append(result, source)
		}
	}
	if len(result) == 0 {
		return []string{"douyin", "xiaohongshu", "meituan", "eleme"}
	}
	return dedupeNonEmpty(result)
}

func BuildSearchQueries(sources []string, keywords []string) []SearchQuery {
	sources = DefaultSources(sources)
	result := make([]SearchQuery, 0, len(sources)*len(keywords))
	for _, source := range sources {
		rule := platformSearchRules[source]
		for _, keyword := range keywords {
			keyword = strings.TrimSpace(keyword)
			if keyword == "" {
				continue
			}
			terms := extractSearchTerms(keyword)
			queryParts := quoteTerms(terms)
			if len(queryParts) == 0 {
				queryParts = []string{keyword}
			}
			for _, domain := range rule.Domains {
				queryParts = append(queryParts, "site:"+domain)
			}
			result = append(result, SearchQuery{
				Platform:  source,
				Keyword:   keyword,
				Query:     strings.Join(queryParts, " "),
				MustTerms: importantTerms(terms),
				Reason:    "system_exact_platform_query",
			})
		}
	}
	return dedupeQueries(result)
}

func SearchPublicResults(ctx context.Context, sources []string, keywords []string, limit int) ([]SearchResult, error) {
	return SearchPublicResultsWithQueries(ctx, BuildSearchQueries(sources, keywords), limit)
}

func SearchPublicResultsWithQueries(ctx context.Context, queries []SearchQuery, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = defaultSearchLimit
	}
	if limit > 10 {
		limit = 10
	}

	client := &http.Client{Timeout: 15 * time.Second}
	results := make([]SearchResult, 0, len(queries)*limit)
	var lastErr error
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 4)

	for _, query := range dedupeQueries(queries) {
		query := normalizeSearchQuery(query)
		if query.Platform == "" || query.Query == "" {
			continue
		}
		if _, ok := platformSearchRules[query.Platform]; !ok {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			rule := platformSearchRules[query.Platform]
			items, err := searchBing(ctx, client, query, rule.Domains, limit)
			if err == nil && len(items) == 0 {
				items, err = searchSogouWeixin(ctx, client, query, limit)
			}
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				lastErr = err
				return
			}
			results = append(results, items...)
		}()
	}
	wg.Wait()

	results = dedupeResults(results)
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	if len(results) == 0 && lastErr != nil {
		return nil, lastErr
	}
	return results, nil
}

func searchBing(ctx context.Context, client *http.Client, query SearchQuery, domains []string, limit int) ([]SearchResult, error) {
	u, _ := url.Parse("https://www.bing.com/search")
	q := u.Query()
	q.Set("q", query.Query)
	q.Set("count", fmt.Sprintf("%d", limit*2))
	q.Set("format", "rss")
	q.Set("mkt", "zh-CN")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/126 Safari/537.36")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.6")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search %s failed: %w", query.Platform, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("search %s status %d", query.Platform, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	raw := string(body)
	results := parseBingRSSResults(query, raw, domains, limit)
	if len(results) > 0 {
		return results, nil
	}
	return parseBingHTMLResults(query, raw, domains, limit), nil
}

func searchSogouWeixin(ctx context.Context, client *http.Client, query SearchQuery, limit int) ([]SearchResult, error) {
	searchText := query.Keyword
	if strings.TrimSpace(searchText) == "" {
		searchText = strings.Join(query.MustTerms, " ")
	}
	u, _ := url.Parse("https://weixin.sogou.com/weixin")
	q := u.Query()
	q.Set("type", "2")
	q.Set("query", searchText)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/126 Safari/537.36")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.6")
	req.Header.Set("Referer", "https://weixin.sogou.com/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sogou weixin search %s failed: %w", query.Platform, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("sogou weixin search %s status %d", query.Platform, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseSogouWeixinResults(query, string(body), limit), nil
}

func parseSogouWeixinResults(query SearchQuery, raw string, limit int) []SearchResult {
	blockRe := regexp.MustCompile(`(?is)<div class="txt-box">.*?</div>\s*</li>`)
	linkRe := regexp.MustCompile(`(?is)<h3>\s*<a[^>]+href="([^"]+)"[^>]*>(.*?)</a>`)
	snippetRe := regexp.MustCompile(`(?is)<p class="txt-info"[^>]*>(.*?)</p>`)

	blocks := blockRe.FindAllString(raw, -1)
	results := make([]SearchResult, 0, minInt(limit, len(blocks)))
	for _, block := range blocks {
		linkMatch := linkRe.FindStringSubmatch(block)
		if len(linkMatch) < 3 {
			continue
		}
		snippet := ""
		if snippetMatch := snippetRe.FindStringSubmatch(block); len(snippetMatch) >= 2 {
			snippet = cleanHTML(snippetMatch[1])
		}
		itemURL := html.UnescapeString(strings.TrimSpace(linkMatch[1]))
		if strings.HasPrefix(itemURL, "/") {
			itemURL = "https://weixin.sogou.com" + itemURL
		}
		item := SearchResult{
			Platform: query.Platform,
			Keyword:  query.Keyword,
			Query:    "sogou_weixin:" + query.Keyword,
			Title:    cleanHTML(linkMatch[2]),
			URL:      itemURL,
			Snippet:  snippet,
			Source:   "sogou_weixin_bridge",
		}
		if !validSearchResult(&item, nil, query.MustTerms) {
			continue
		}
		results = append(results, item)
		if len(results) >= limit {
			break
		}
	}
	return results
}

func parseBingRSSResults(query SearchQuery, raw string, domains []string, limit int) []SearchResult {
	type item struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
	}
	var parsed struct {
		Channel struct {
			Items []item `xml:"item"`
		} `xml:"channel"`
	}
	if err := xml.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil
	}
	results := make([]SearchResult, 0, minInt(limit, len(parsed.Channel.Items)))
	for _, value := range parsed.Channel.Items {
		item := SearchResult{
			Platform: query.Platform,
			Keyword:  query.Keyword,
			Query:    query.Query,
			Title:    cleanHTML(value.Title),
			URL:      html.UnescapeString(strings.TrimSpace(value.Link)),
			Snippet:  cleanHTML(value.Description),
			Source:   "bing_rss",
		}
		if !validSearchResult(&item, domains, query.MustTerms) {
			continue
		}
		results = append(results, item)
		if len(results) >= limit {
			break
		}
	}
	return results
}

func parseBingHTMLResults(query SearchQuery, raw string, domains []string, limit int) []SearchResult {
	blockRe := regexp.MustCompile(`(?is)<li class="b_algo".*?</li>`)
	linkRe := regexp.MustCompile(`(?is)<h2.*?<a[^>]+href="([^"]+)"[^>]*>(.*?)</a>`)
	snippetRe := regexp.MustCompile(`(?is)<p[^>]*>(.*?)</p>`)

	blocks := blockRe.FindAllString(raw, -1)
	results := make([]SearchResult, 0, minInt(limit, len(blocks)))
	for _, block := range blocks {
		linkMatch := linkRe.FindStringSubmatch(block)
		if len(linkMatch) < 3 {
			continue
		}
		snippet := ""
		if snippetMatch := snippetRe.FindStringSubmatch(block); len(snippetMatch) >= 2 {
			snippet = cleanHTML(snippetMatch[1])
		}
		item := SearchResult{
			Platform: query.Platform,
			Keyword:  query.Keyword,
			Query:    query.Query,
			Title:    cleanHTML(linkMatch[2]),
			URL:      html.UnescapeString(linkMatch[1]),
			Snippet:  snippet,
			Source:   "bing",
		}
		if !validSearchResult(&item, domains, query.MustTerms) {
			continue
		}
		results = append(results, item)
		if len(results) >= limit {
			break
		}
	}
	return results
}

func validSearchResult(item *SearchResult, domains []string, mustTerms []string) bool {
	if item.Title == "" || item.URL == "" {
		return false
	}
	if len(domains) > 0 && !urlMatchesDomains(item.URL, domains) {
		return false
	}

	item.Score = relevanceScore(*item, mustTerms)
	minScore := 2
	if len(domains) > 0 {
		minScore = 1
	}
	return item.Score >= minScore
}

func relevanceScore(item SearchResult, mustTerms []string) int {
	text := normalizeText(item.Title + " " + item.Snippet + " " + item.URL)
	score := 0
	for _, term := range mustTerms {
		term = normalizeText(term)
		if term == "" {
			continue
		}
		if strings.Contains(text, term) {
			score += 2
		}
	}
	for _, noise := range []string{
		"花卉", "鲜花", "花语", "花百科", "植物", "百度百科",
		"小游戏", "游戏大全", "steam", "calculator", "autocad", "drawing", "vietnamese",
		"旅游攻略", "景点排行",
	} {
		if strings.Contains(text, noise) {
			score -= 3
		}
	}
	if strings.Contains(text, "火锅") || strings.Contains(text, "餐饮") || strings.Contains(text, "团购") || strings.Contains(text, "探店") {
		score += 2
	}
	return score
}

func urlMatchesDomains(rawURL string, domains []string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	host := strings.ToLower(parsed.Host)
	if host == "" {
		return false
	}
	for _, domain := range domains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if host == domain || strings.HasSuffix(host, "."+domain) {
			return true
		}
	}
	return false
}

func extractSearchTerms(keyword string) []string {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil
	}
	fields := strings.FieldsFunc(keyword, func(r rune) bool {
		return unicode.IsSpace(r) || r == ',' || r == '，' || r == '/' || r == '|' || r == ';' || r == '；' || r == ':'
	})
	result := make([]string, 0, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(strings.Trim(field, `"'“”‘’`))
		if field == "" {
			continue
		}
		if len([]rune(field)) <= 1 && !hasDigit(field) {
			continue
		}
		result = append(result, field)
	}
	return dedupeNonEmpty(result)
}

func inferCategory(merchantName, industry, address string, products []Product) string {
	text := merchantName + " " + industry + " " + address
	for _, product := range products {
		text += " " + product.Name + " " + product.TrafficLabel
	}
	categories := []string{
		"火锅", "酸汤", "烤肉", "烧烤", "串串", "冒菜", "麻辣烫", "烤鱼", "牛肉粉", "羊肉粉",
		"粉面", "小吃", "咖啡", "奶茶", "甜品", "自助餐", "西餐", "日料", "湘菜", "川菜", "黔菜",
	}
	for _, category := range categories {
		if strings.Contains(text, category) {
			if category == "酸汤" {
				return "酸汤火锅"
			}
			return category
		}
	}
	industry = strings.TrimSpace(industry)
	if industry != "" && industry != "餐饮" && industry != "美食" {
		return industry
	}
	return "火锅"
}

func extractDistrict(address string) string {
	re := regexp.MustCompile(`([\p{Han}]{2,8}(区|县|市|镇|街道))`)
	matches := re.FindAllString(address, -1)
	for _, match := range matches {
		if strings.Contains(match, "贵阳市") || strings.Contains(match, "贵州省") || strings.Contains(match, "中国") {
			continue
		}
		return match
	}
	return ""
}

func extractRoad(address string) string {
	re := regexp.MustCompile(`([\p{Han}A-Za-z0-9]{2,16}(路|街|大道|巷|商圈|广场|城|园))`)
	matches := re.FindAllString(address, -1)
	for _, match := range matches {
		match = cleanLocationTerm(match)
		if match == "" || strings.Contains(match, "邮政编码") || strings.Contains(match, "贵州省") || strings.Contains(match, "贵阳市") {
			continue
		}
		return match
	}
	return ""
}

func extractNearbyShop(address string) string {
	re := regexp.MustCompile(`([\p{Han}A-Za-z0-9]{2,20}(火锅店|火锅|餐厅|饭店|店))旁`)
	match := re.FindStringSubmatch(address)
	if len(match) >= 2 {
		return cleanShopName(match[1])
	}
	return ""
}

func cleanLocationTerm(value string) string {
	value = strings.TrimSpace(value)
	for _, marker := range []string{"省", "市", "区", "县", "镇", "街道"} {
		if index := strings.LastIndex(value, marker); index >= 0 && index+len(marker) < len(value) {
			value = value[index+len(marker):]
		}
	}
	return strings.TrimSpace(value)
}

func cleanShopName(value string) string {
	value = cleanLocationTerm(value)
	for _, marker := range []string{"大道", "高新路", "路", "街", "巷", "广场", "商圈"} {
		if index := strings.LastIndex(value, marker); index >= 0 && index+len(marker) < len(value) {
			value = value[index+len(marker):]
		}
	}
	return strings.TrimSpace(value)
}

func normalizeProductKeyword(name, category string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	if name == "单人" {
		return "单人餐"
	}
	if strings.Contains(name, "单人") && !strings.Contains(name, "餐") && category != "" {
		return name + "餐"
	}
	return name
}

func normalizeManualKeyword(keyword, city, category string) string {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return ""
	}
	if keyword == "单人" {
		keyword = "单人餐"
	}
	if category != "" && !strings.Contains(keyword, category) {
		keyword = category + " " + keyword
	}
	if city != "" && !strings.Contains(keyword, city) {
		keyword = city + " " + keyword
	}
	return keyword
}

func importantTerms(terms []string) []string {
	result := make([]string, 0, len(terms))
	for _, term := range terms {
		term = strings.TrimSpace(term)
		if term == "" {
			continue
		}
		if len([]rune(term)) <= 1 && !hasDigit(term) {
			continue
		}
		result = append(result, term)
	}
	if len(result) > 4 {
		result = result[:4]
	}
	return result
}

func quoteTerms(terms []string) []string {
	result := make([]string, 0, len(terms))
	for _, term := range terms {
		term = strings.TrimSpace(term)
		if term == "" {
			continue
		}
		result = append(result, `"`+strings.ReplaceAll(term, `"`, "")+`"`)
	}
	return result
}

func normalizeSearchQuery(query SearchQuery) SearchQuery {
	query.Platform = strings.ToLower(strings.TrimSpace(query.Platform))
	query.Keyword = strings.TrimSpace(query.Keyword)
	query.Query = strings.TrimSpace(query.Query)
	query.MustTerms = importantTerms(query.MustTerms)
	if len(query.MustTerms) == 0 {
		query.MustTerms = importantTerms(extractSearchTerms(query.Keyword))
	}
	if query.Query == "" {
		query.Query = strings.Join(quoteTerms(query.MustTerms), " ")
	}
	return query
}

func cleanHTML(value string) string {
	tagRe := regexp.MustCompile(`(?is)<[^>]+>`)
	value = tagRe.ReplaceAllString(value, "")
	value = html.UnescapeString(value)
	value = strings.Join(strings.Fields(value), " ")
	return strings.TrimSpace(value)
}

func normalizeText(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(value), ""))
}

func hasDigit(value string) bool {
	for _, r := range value {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func dedupeNonEmpty(values []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func dedupeQueries(values []SearchQuery) []SearchQuery {
	seen := map[string]struct{}{}
	result := make([]SearchQuery, 0, len(values))
	for _, value := range values {
		value = normalizeSearchQuery(value)
		key := value.Platform + "|" + value.Query
		if value.Platform == "" || value.Query == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, value)
	}
	return result
}

func dedupeResults(values []SearchResult) []SearchResult {
	seen := map[string]struct{}{}
	result := make([]SearchResult, 0, len(values))
	for _, value := range values {
		key := strings.TrimSpace(value.URL)
		if key == "" {
			key = value.Platform + "|" + value.Title
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, value)
	}
	return result
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
