package admin

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/platformresearch"
	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlatformResearchHandler struct{}

func NewPlatformResearchHandler() *PlatformResearchHandler {
	return &PlatformResearchHandler{}
}

type platformResearchPayload struct {
	MerchantID       uint64   `json:"merchantId" binding:"required"`
	Sources          []string `json:"sources"`
	Keywords         []string `json:"keywords"`
	Limit            int      `json:"limit"`
	ExtraRequirement string   `json:"extraRequirement"`
}

type platformResearchListResponse struct {
	List  []platformResearchResponse `json:"list"`
	Total int64                      `json:"total"`
	Page  int                        `json:"page"`
	Size  int                        `json:"size"`
}

type platformResearchResponse struct {
	model.PlatformResearchTask
	Sources       any `json:"sources,omitempty"`
	Keywords      any `json:"keywords,omitempty"`
	SearchResults any `json:"searchResults,omitempty"`
	GoodCases     any `json:"goodCases,omitempty"`
	BadCases      any `json:"badCases,omitempty"`
	Insights      any `json:"insights,omitempty"`
	Suggestions   any `json:"suggestions,omitempty"`
}

func (h *PlatformResearchHandler) List(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	merchantID := parseUint64Query(c.Query("merchantId"))
	keyword := strings.TrimSpace(c.Query("keyword"))

	query := database.DB.Model(&model.PlatformResearchTask{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("merchant_name LIKE ? OR industry LIKE ? OR city LIKE ? OR summary LIKE ?", like, like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		basehandler.ServerError(c, "读取平台调研数量失败")
		return
	}

	var list []model.PlatformResearchTask
	if err := query.Order("id DESC").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取平台调研列表失败")
		return
	}

	basehandler.OK(c, platformResearchListResponse{
		List:  buildPlatformResearchResponses(list),
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func (h *PlatformResearchHandler) Get(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var item model.PlatformResearchTask
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "平台调研任务不存在")
			return
		}
		basehandler.ServerError(c, "读取平台调研任务失败")
		return
	}
	basehandler.OK(c, buildPlatformResearchResponse(item))
}

func (h *PlatformResearchHandler) Generate(c *gin.Context) {
	var payload platformResearchPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家")
		return
	}

	task, input, err := buildPlatformResearchTask(c.Request.Context(), payload, c.GetUint64("admin_user_id"))
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}

	if err := database.DB.Create(&task).Error; err != nil {
		basehandler.ServerError(c, "创建平台调研任务失败")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	output, runErr := platformresearch.Agent{}.Run(ctx, input)
	if runErr != nil {
		database.DB.Model(&task).Updates(map[string]any{
			"status":        model.PlatformResearchTaskStatusFailed,
			"error_message": runErr.Error(),
			"updated_by":    c.GetUint64("admin_user_id"),
		})
		basehandler.ServerError(c, runErr.Error())
		return
	}

	goodCasesJSON, _ := json.Marshal(output.GoodCases)
	badCasesJSON, _ := json.Marshal(output.BadCases)
	insightsJSON, _ := json.Marshal(output.Insights)
	suggestionsJSON, _ := json.Marshal(output.Suggestions)

	updates := map[string]any{
		"status":           model.PlatformResearchTaskStatusCompleted,
		"summary":          output.Summary,
		"good_cases_json":  string(goodCasesJSON),
		"bad_cases_json":   string(badCasesJSON),
		"insights_json":    string(insightsJSON),
		"suggestions_json": string(suggestionsJSON),
		"updated_by":       c.GetUint64("admin_user_id"),
	}
	if err := database.DB.Model(&task).Updates(updates).Error; err != nil {
		basehandler.ServerError(c, "保存平台调研结果失败")
		return
	}

	task.Status = model.PlatformResearchTaskStatusCompleted
	task.Summary = output.Summary
	task.GoodCasesJSON = string(goodCasesJSON)
	task.BadCasesJSON = string(badCasesJSON)
	task.InsightsJSON = string(insightsJSON)
	task.SuggestionsJSON = string(suggestionsJSON)
	task.SearchResultsJSON = mustJSON(input.SearchResults)
	task.ErrorMessage = ""

	basehandler.OK(c, buildPlatformResearchResponse(task))
}

func buildPlatformResearchTask(ctx context.Context, payload platformResearchPayload, operatorID uint64) (model.PlatformResearchTask, platformresearch.Input, error) {
	var merchant model.Merchant
	if err := database.DB.First(&merchant, payload.MerchantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.PlatformResearchTask{}, platformresearch.Input{}, errors.New("商家不存在")
		}
		return model.PlatformResearchTask{}, platformresearch.Input{}, errors.New("读取商家失败")
	}

	var packages []model.MerchantPackage
	if err := database.DB.Where("merchant_id = ? AND status = ?", merchant.ID, model.MerchantPackageStatusEnabled).
		Order("id DESC").
		Limit(20).
		Find(&packages).Error; err != nil {
		return model.PlatformResearchTask{}, platformresearch.Input{}, errors.New("读取套餐失败")
	}

	products := make([]platformresearch.Product, 0, len(packages))
	for _, item := range packages {
		products = append(products, platformresearch.Product{
			Name:          item.Name,
			SellingPrice:  item.SellingPrice,
			OriginalPrice: item.OriginalPrice,
			TrafficLabel:  item.TrafficLabel,
			ProfitGuard:   item.ProfitGuard,
		})
	}

	sources := platformresearch.DefaultSources(payload.Sources)
	keywords := platformresearch.BuildKeywords(merchant.Name, merchant.City, merchant.Industry, merchant.Address, products, payload.Keywords)
	if len(keywords) == 0 {
		return model.PlatformResearchTask{}, platformresearch.Input{}, errors.New("缺少搜索关键词，请先完善商家城市、行业、地址或手动输入关键词")
	}
	if len(keywords) > 16 {
		keywords = keywords[:16]
	}

	input := platformresearch.Input{
		MerchantID:       merchant.ID,
		MerchantName:     merchant.Name,
		Industry:         merchant.Industry,
		City:             merchant.City,
		Address:          merchant.Address,
		Remark:           merchant.Remark,
		Products:         products,
		Sources:          sources,
		Keywords:         keywords,
		ExtraRequirement: strings.TrimSpace(payload.ExtraRequirement),
		Options:          agent.RunOptions{OperatorID: operatorID, DryRun: true},
	}

	searchQueries := platformresearch.BuildSearchQueries(sources, keywords)
	planCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	if plannedQueries, planErr := platformresearch.PlanSearchQueries(planCtx, input); planErr == nil && len(plannedQueries) > 0 {
		searchQueries = plannedQueries
	}
	cancel()
	if len(searchQueries) > 16 {
		searchQueries = searchQueries[:16]
	}

	limit := payload.Limit
	if limit <= 0 {
		limit = 5
	}
	searchResults, err := platformresearch.SearchPublicResultsWithQueries(ctx, searchQueries, limit)
	if err != nil {
		return model.PlatformResearchTask{}, platformresearch.Input{}, err
	}
	if len(searchResults) == 0 {
		return model.PlatformResearchTask{}, platformresearch.Input{}, errors.New("没有搜索到相关公开结果，请换更具体的关键词")
	}

	input.SearchQueries = searchQueries
	input.SearchResults = searchResults

	return model.PlatformResearchTask{
		MerchantID:        merchant.ID,
		MerchantName:      merchant.Name,
		Industry:          merchant.Industry,
		City:              merchant.City,
		SourcesJSON:       mustJSON(sources),
		KeywordsJSON:      mustJSON(keywords),
		SearchResultsJSON: mustJSON(searchResults),
		Status:            model.PlatformResearchTaskStatusPending,
		CreatedBy:         operatorID,
		UpdatedBy:         operatorID,
	}, input, nil
}

func buildPlatformResearchResponses(list []model.PlatformResearchTask) []platformResearchResponse {
	result := make([]platformResearchResponse, 0, len(list))
	for _, item := range list {
		result = append(result, buildPlatformResearchResponse(item))
	}
	return result
}

func buildPlatformResearchResponse(item model.PlatformResearchTask) platformResearchResponse {
	return platformResearchResponse{
		PlatformResearchTask: item,
		Sources:              parseJSONField(item.SourcesJSON),
		Keywords:             parseJSONField(item.KeywordsJSON),
		SearchResults:        parseJSONField(item.SearchResultsJSON),
		GoodCases:            parseJSONField(item.GoodCasesJSON),
		BadCases:             parseJSONField(item.BadCasesJSON),
		Insights:             parseJSONField(item.InsightsJSON),
		Suggestions:          parseJSONField(item.SuggestionsJSON),
	}
}

func mustJSON(value any) string {
	raw, err := json.Marshal(value)
	if err != nil {
		return "[]"
	}
	return string(raw)
}
