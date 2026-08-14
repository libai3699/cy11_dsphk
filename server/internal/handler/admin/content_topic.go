package admin

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/hotspottopic"
	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContentTopicHandler struct{}

func NewContentTopicHandler() *ContentTopicHandler {
	return &ContentTopicHandler{}
}

type topicGeneratePayload struct {
	MerchantID       uint64   `json:"merchantId" binding:"required"`
	BenchmarkID      uint64   `json:"benchmarkId"`
	BenchmarkIDs     []uint64 `json:"benchmarkIds"`
	BenchmarkName    string   `json:"benchmarkName"`
	CityHotspots     []string `json:"cityHotspots"`
	IndustryHotspots []string `json:"industryHotspots"`
	NationalHotspots []string `json:"nationalHotspots"`
	ExtraRequirement string   `json:"extraRequirement"`
}

type topicListResponse struct {
	List  []contentTopicResponse `json:"list"`
	Total int64                  `json:"total"`
	Page  int                    `json:"page"`
	Size  int                    `json:"size"`
}

type contentTopicResponse struct {
	model.ContentTopic
	Tags []string `json:"tags"`
}

type hotspotTopicTaskResponse struct {
	model.HotspotTopicTask
	Input  any                    `json:"input,omitempty"`
	Result any                    `json:"result,omitempty"`
	Topics []contentTopicResponse `json:"topics,omitempty"`
}

func (h *ContentTopicHandler) List(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	merchantID := parseUint64Query(c.Query("merchantId"))
	status := strings.TrimSpace(c.Query("status"))

	query := database.DB.Model(&model.ContentTopic{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"merchant_name LIKE ? OR benchmark_name LIKE ? OR title LIKE ? OR hook LIKE ? OR angle LIKE ? OR target LIKE ? OR recommend_reason LIKE ?",
			like,
			like,
			like,
			like,
			like,
			like,
			like,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		basehandler.ServerError(c, "读取选题数量失败")
		return
	}

	var list []model.ContentTopic
	if err := query.Order("id DESC").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取选题失败")
		return
	}

	basehandler.OK(c, topicListResponse{
		List:  buildContentTopicResponses(list),
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func (h *ContentTopicHandler) Generate(c *gin.Context) {
	var payload topicGeneratePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家")
		return
	}

	task, input, err := buildHotspotTopicTask(payload, c.GetUint64("admin_user_id"))
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	if err := database.DB.Create(&task).Error; err != nil {
		basehandler.ServerError(c, "创建找爆款任务失败")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()
	output, runErr := hotspottopic.Agent{}.Run(ctx, input)
	if runErr != nil {
		task.Status = model.HotspotTopicTaskStatusFailed
		task.ErrorMessage = runErr.Error()
		_ = database.DB.Model(&task).Updates(map[string]any{
			"status":        task.Status,
			"error_message": task.ErrorMessage,
			"updated_by":    task.UpdatedBy,
		}).Error
		basehandler.ServerError(c, "找爆款 Agent 执行失败："+runErr.Error())
		return
	}

	resultJSON, err := json.Marshal(output)
	if err != nil {
		basehandler.ServerError(c, "找爆款结果序列化失败")
		return
	}
	task.Status = model.HotspotTopicTaskStatusCompleted
	task.ResultJSON = string(resultJSON)
	if err := database.DB.Model(&task).Updates(map[string]any{
		"status":      task.Status,
		"result_json": task.ResultJSON,
		"updated_by":  task.UpdatedBy,
	}).Error; err != nil {
		basehandler.ServerError(c, "保存找爆款结果失败")
		return
	}

	topics := make([]model.ContentTopic, 0, len(output.Topics))
	for _, topic := range output.Topics {
		tagsJSON, _ := json.Marshal(topic.Tags)
		topics = append(topics, model.ContentTopic{
			MerchantID:      task.MerchantID,
			MerchantName:    task.MerchantName,
			BenchmarkID:     task.BenchmarkID,
			BenchmarkName:   task.BenchmarkName,
			Title:           strings.TrimSpace(topic.Title),
			Hook:            strings.TrimSpace(topic.Hook),
			Angle:           strings.TrimSpace(topic.Angle),
			Target:          strings.TrimSpace(topic.Target),
			RiskLevel:       strings.TrimSpace(topic.RiskLevel),
			RecommendReason: strings.TrimSpace(topic.RecommendReason),
			TagsJSON:        string(tagsJSON),
			Status:          model.ContentTopicStatusPending,
			SourceTaskID:    task.ID,
			CreatedBy:       task.CreatedBy,
			UpdatedBy:       task.UpdatedBy,
		})
	}
	if len(topics) > 0 {
		if err := database.DB.Create(&topics).Error; err != nil {
			basehandler.ServerError(c, "保存选题失败")
			return
		}
	}

	if err := database.DB.First(&task, task.ID).Error; err != nil {
		basehandler.ServerError(c, "读取找爆款任务失败")
		return
	}
	basehandler.OK(c, hotspotTopicTaskResponse{
		HotspotTopicTask: task,
		Input:            parseJSONField(task.InputSnapshot),
		Result:           parseJSONField(task.ResultJSON),
		Topics:           buildContentTopicResponses(topics),
	})
}

func (h *ContentTopicHandler) UpdateStatus(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	var payload struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择状态")
		return
	}
	if !isValidContentTopicStatus(payload.Status) {
		basehandler.BadRequest(c, "选题状态不正确")
		return
	}

	var item model.ContentTopic
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "选题不存在")
			return
		}
		basehandler.ServerError(c, "读取选题失败")
		return
	}
	if err := database.DB.Model(&item).Updates(map[string]any{
		"status":     payload.Status,
		"updated_by": c.GetUint64("admin_user_id"),
	}).Error; err != nil {
		basehandler.ServerError(c, "更新选题状态失败")
		return
	}
	if err := database.DB.First(&item, id).Error; err != nil {
		basehandler.ServerError(c, "读取选题失败")
		return
	}
	basehandler.OK(c, buildContentTopicResponse(item))
}

func buildHotspotTopicTask(payload topicGeneratePayload, operatorID uint64) (model.HotspotTopicTask, hotspottopic.Input, error) {
	if payload.MerchantID == 0 {
		return model.HotspotTopicTask{}, hotspottopic.Input{}, errors.New("请选择商家")
	}

	var merchant model.Merchant
	if err := database.DB.First(&merchant, payload.MerchantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.HotspotTopicTask{}, hotspottopic.Input{}, errors.New("商家不存在")
		}
		return model.HotspotTopicTask{}, hotspottopic.Input{}, errors.New("读取商家失败")
	}

	var packages []model.MerchantPackage
	_ = database.DB.Where("merchant_id = ? AND status = ?", merchant.ID, model.MerchantPackageStatusEnabled).Find(&packages).Error

	benchmarks, err := loadTopicBenchmarks(merchant.ID, payload)
	if err != nil {
		return model.HotspotTopicTask{}, hotspottopic.Input{}, err
	}
	benchmarkID, benchmarkName := topicBenchmarkIdentity(benchmarks, payload)

	input := hotspottopic.Input{
		MerchantID:       merchant.ID,
		MerchantName:     merchant.Name,
		Industry:         merchant.Industry,
		City:             merchant.City,
		BenchmarkAccount: benchmarkName,
		ExtraRequirement: strings.TrimSpace(payload.ExtraRequirement),
		Options:          agent.RunOptions{OperatorID: operatorID, DryRun: true},
	}
	if len(benchmarks) == 1 {
		input.BenchmarkSummary = benchmarks[0].Takeaway
	}
	for _, item := range benchmarks {
		input.BenchmarkAccounts = append(input.BenchmarkAccounts, hotspottopic.BenchmarkAccount{
			Name:           item.AccountName,
			Platform:       item.Platform,
			City:           item.City,
			Industry:       item.Industry,
			FollowerCount:  item.FollowerCount,
			BestPlayCount:  item.BestPlayCount,
			LatestHitTitle: item.LatestHitTitle,
			Takeaway:       item.Takeaway,
			Risk:           item.Risk,
		})
	}
	for _, item := range packages {
		input.Products = append(input.Products, hotspottopic.Product{
			Name:        item.Name,
			Price:       item.SellingPrice,
			GrossMargin: item.SellingPrice - item.CostPrice,
		})
	}
	appendHotspots := func(items []string, source string, scope string) {
		for _, title := range items {
			title = strings.TrimSpace(title)
			if title == "" {
				continue
			}
			input.Hotspots = append(input.Hotspots, hotspottopic.Hotspot{
				Title:  title,
				Source: source,
				Scope:  scope,
			})
		}
	}
	appendHotspots(payload.CityHotspots, "manual", "同城")
	appendHotspots(payload.IndustryHotspots, "manual", "行业")
	appendHotspots(payload.NationalHotspots, "manual", "全国")

	snapshot := map[string]any{
		"merchant":       merchant,
		"packages":       packages,
		"benchmarks":     benchmarks,
		"manualInput":    payload,
		"agentInput":     input,
		"snapshotTime":   time.Now().Format("2006-01-02 15:04:05"),
		"snapshotSource": "admin_manual_input",
	}
	inputSnapshot, err := json.Marshal(snapshot)
	if err != nil {
		return model.HotspotTopicTask{}, hotspottopic.Input{}, errors.New("输入快照生成失败")
	}

	return model.HotspotTopicTask{
		MerchantID:    merchant.ID,
		MerchantName:  merchant.Name,
		BenchmarkID:   benchmarkID,
		BenchmarkName: benchmarkName,
		Status:        model.HotspotTopicTaskStatusRunning,
		InputSnapshot: string(inputSnapshot),
		CreatedBy:     operatorID,
		UpdatedBy:     operatorID,
	}, input, nil
}

func loadTopicBenchmarks(merchantID uint64, payload topicGeneratePayload) ([]model.BenchmarkAccount, error) {
	ids := make([]uint64, 0, len(payload.BenchmarkIDs)+1)
	seen := map[uint64]struct{}{}
	for _, id := range payload.BenchmarkIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if payload.BenchmarkID > 0 {
		if _, ok := seen[payload.BenchmarkID]; !ok {
			ids = append(ids, payload.BenchmarkID)
		}
	}
	if len(ids) == 0 {
		return []model.BenchmarkAccount{}, nil
	}

	var benchmarks []model.BenchmarkAccount
	if err := database.DB.
		Where("merchant_id = ? AND id IN ?", merchantID, ids).
		Order("best_play_count DESC, follower_count DESC, id DESC").
		Find(&benchmarks).Error; err != nil {
		return nil, errors.New("读取对标账号失败")
	}
	if len(benchmarks) == 0 {
		return nil, errors.New("所选对标账号不属于当前商家")
	}
	return benchmarks, nil
}

func topicBenchmarkIdentity(benchmarks []model.BenchmarkAccount, payload topicGeneratePayload) (uint64, string) {
	if len(benchmarks) == 0 {
		return 0, strings.TrimSpace(payload.BenchmarkName)
	}
	if len(benchmarks) == 1 {
		return benchmarks[0].ID, benchmarks[0].AccountName
	}
	names := make([]string, 0, len(benchmarks))
	for _, item := range benchmarks {
		names = append(names, item.AccountName)
	}
	return 0, strings.Join(names, "、")
}

func buildContentTopicResponses(list []model.ContentTopic) []contentTopicResponse {
	result := make([]contentTopicResponse, 0, len(list))
	for _, item := range list {
		result = append(result, buildContentTopicResponse(item))
	}
	return result
}

func buildContentTopicResponse(item model.ContentTopic) contentTopicResponse {
	var tags []string
	_ = json.Unmarshal([]byte(item.TagsJSON), &tags)
	return contentTopicResponse{ContentTopic: item, Tags: tags}
}

func isValidContentTopicStatus(status string) bool {
	switch status {
	case model.ContentTopicStatusPending, model.ContentTopicStatusAccepted, model.ContentTopicStatusDisabled:
		return true
	default:
		return false
	}
}
