package platformresearch

import (
	"context"
	"encoding/json"
	"fmt"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

const Name = "platform_research"

type Agent struct{}

type Input struct {
	MerchantID       uint64           `json:"merchantId"`
	MerchantName     string           `json:"merchantName"`
	Industry         string           `json:"industry"`
	City             string           `json:"city"`
	Address          string           `json:"address,omitempty"`
	Remark           string           `json:"remark,omitempty"`
	Products         []Product        `json:"products,omitempty"`
	Sources          []string         `json:"sources,omitempty"`
	Keywords         []string         `json:"keywords,omitempty"`
	SearchQueries    []SearchQuery    `json:"searchQueries,omitempty"`
	SearchResults    []SearchResult   `json:"searchResults,omitempty"`
	ExtraRequirement string           `json:"extraRequirement,omitempty"`
	Options          agent.RunOptions `json:"options,omitempty"`
}

type Product struct {
	Name          string  `json:"name"`
	SellingPrice  float64 `json:"sellingPrice"`
	OriginalPrice float64 `json:"originalPrice"`
	TrafficLabel  string  `json:"trafficLabel,omitempty"`
	ProfitGuard   string  `json:"profitGuard,omitempty"`
}

type SearchResult struct {
	Platform string `json:"platform"`
	Keyword  string `json:"keyword"`
	Query    string `json:"query,omitempty"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Snippet  string `json:"snippet"`
	Source   string `json:"source"`
	Score    int    `json:"score,omitempty"`
}

type CaseStudy struct {
	Platform string   `json:"platform"`
	Title    string   `json:"title"`
	URL      string   `json:"url,omitempty"`
	Reason   string   `json:"reason"`
	Takeaway string   `json:"takeaway"`
	Risks    []string `json:"risks,omitempty"`
}

type Output struct {
	agent.Result
	GoodCases []CaseStudy `json:"goodCases,omitempty"`
	BadCases  []CaseStudy `json:"badCases,omitempty"`
	Insights  []string    `json:"insights,omitempty"`
}

func PlanSearchQueries(ctx context.Context, input Input) ([]SearchQuery, error) {
	client, ok := provider.NewStepFunForAgent(Name)
	if !ok {
		return nil, agent.ErrProviderNotConfigured
	}

	payload, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return nil, err
	}

	resp, err := client.Chat(ctx, []provider.Message{
		{
			Role: "system",
			Content: `你是本地生活获客的「搜索词规划 Agent」。
你只负责把商家资料、套餐、人工关键词，改写成更适合搜索引擎的四平台搜索计划。
必须输出严格 JSON，不要 Markdown，不要解释。
JSON 结构：
{
  "queries": [
    {
      "platform": "douyin/xiaohongshu/meituan/eleme",
      "keyword": "给运营人员看的关键词",
      "query": "真正发给搜索引擎的查询词",
      "mustTerms": ["必须相关的词"],
      "reason": "为什么这样搜"
    }
  ]
}
要求：
1. 先识别真实品类。行业是“餐饮/美食”时不能直接用餐饮，要从商家名、地址、套餐里推断火锅/烤肉/粉面等具体品类。
2. 不要把完整地址、邮政编码、中国贵州省这类长地址放进 query；地址只抽城市、区县、商圈、道路、附近竞品名。
3. 搜索目标不是搜这个商家自己，而是先查市场：本地大概有哪些店、头部是谁、团购怎么做、评价和差评集中在哪。
4. query 要尽量精确，例如 "贵阳" "火锅" "排行榜"、"乌当区" "火锅" "团购"、"贵阳" "火锅" "小红书"。
5. 每个平台 2-4 条 query 即可，总数不超过 16 条。
6. 抖音/小红书重点搜探店、团购、爆款标题、同城内容。
7. 美团/饿了么重点搜套餐、评价、差评、价格带、附近商家。
8. mustTerms 必须包含城市/区域/品类等核心词，用来过滤垃圾结果。
9. 不要生成“单人”“推荐”“评价”这种脱离城市和品类的泛词。`,
		},
		{
			Role:    "user",
			Content: "请生成平台调研搜索计划 JSON：\n" + string(payload),
		},
	})
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Queries []SearchQuery `json:"queries"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return nil, err
	}
	if len(parsed.Queries) == 0 {
		return nil, fmt.Errorf("stepfun search plan response is empty")
	}
	return parsed.Queries, nil
}

func (Agent) Run(ctx context.Context, input Input) (Output, error) {
	client, ok := provider.NewStepFunForAgent(Name)
	if !ok {
		return Output{}, agent.ErrProviderNotConfigured
	}

	payload, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return Output{}, err
	}

	resp, err := client.Chat(ctx, []provider.Message{
		{
			Role: "system",
			Content: `你是本地生活短视频获客的「四平台调研 Agent」。
你负责根据抖音、小红书、美团、饿了么的公开搜索结果，先判断本地市场和头部竞品，再判断哪些案例值得学习、哪些案例不适合照抄，并给后续选题/文案/分镜提供依据。
你必须输出严格 JSON 对象，不要 Markdown，不要解释，不要反问。
JSON 结构：
{
  "summary": "一句话总结本轮调研结论",
  "suggestions": ["后续动作"],
  "goodCases": [
    {
      "platform": "douyin/xiaohongshu/meituan/eleme",
      "title": "案例标题",
      "url": "来源链接",
      "reason": "为什么算好案例",
      "takeaway": "这个商家可以学习什么",
      "risks": ["照抄风险"]
    }
  ],
  "badCases": [
    {
      "platform": "douyin/xiaohongshu/meituan/eleme",
      "title": "案例标题",
      "url": "来源链接",
      "reason": "为什么不适合",
      "takeaway": "应该避开什么",
      "risks": ["风险点"]
    }
  ],
  "insights": ["汇总洞察"]
}
要求：
1. 先看“城市 + 品类”的市场，不要只盯商家自己；例如贵阳火锅要优先识别本地热门店、团购打法、评价痛点。
2. goodCases 输出 3-8 条，优先选择同城、同行业、同价位、同场景的结果。
3. badCases 输出 2-6 条，重点识别低价伤利润、夸大宣传、平台违规、转化链路不清晰。
4. insights 输出 5-10 条，必须能直接服务后续选题、文案、分镜和拍摄任务。
5. 如果公开搜索结果质量不足，不要只返回“无关”，要说明下一轮应该怎么搜：城市、区县、品类、竞品、套餐词怎么组合。
6. 如果结果里出现无关内容，要明确归类为无效结果，不允许当成好案例。`,
		},
		{
			Role:    "user",
			Content: "请根据以下输入完成四平台调研分析 JSON：\n" + string(payload),
		},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Summary     string      `json:"summary"`
		Suggestions []string    `json:"suggestions"`
		GoodCases   []CaseStudy `json:"goodCases"`
		BadCases    []CaseStudy `json:"badCases"`
		Insights    []string    `json:"insights"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if parsed.Summary == "" {
		return Output{}, fmt.Errorf("stepfun platform research response is empty")
	}

	return Output{
		Result: agent.Result{
			Agent:       Name,
			Version:     "v0.1.0-stepfun-search",
			Status:      "completed",
			Summary:     parsed.Summary,
			Suggestions: parsed.Suggestions,
			Artifacts: map[string]any{
				"provider":         "stepfun",
				"model":            resp.Model,
				"input":            input,
				"rawModelResponse": resp.Content,
			},
		},
		GoodCases: parsed.GoodCases,
		BadCases:  parsed.BadCases,
		Insights:  parsed.Insights,
	}, nil
}
