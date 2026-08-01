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

	output, runErr := hotspottopic.Agent{}.Run(context.Background(), input)
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

	var benchmark model.BenchmarkAccount
	if payload.BenchmarkID > 0 {
		_ = database.DB.First(&benchmark, payload.BenchmarkID).Error
	}
	benchmarkName := strings.TrimSpace(payload.BenchmarkName)
	if benchmarkName == "" {
		benchmarkName = benchmark.AccountName
	}

	input := hotspottopic.Input{
		MerchantID:       merchant.ID,
		MerchantName:     merchant.Name,
		Industry:         merchant.Industry,
		City:             merchant.City,
		BenchmarkAccount: benchmarkName,
		BenchmarkSummary: benchmark.Takeaway,
		ExtraRequirement: strings.TrimSpace(payload.ExtraRequirement),
		Options:          agent.RunOptions{OperatorID: operatorID, DryRun: true},
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
		"benchmark":      benchmark,
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
		BenchmarkID:   benchmark.ID,
		BenchmarkName: benchmarkName,
		Status:        model.HotspotTopicTaskStatusRunning,
		InputSnapshot: string(inputSnapshot),
		CreatedBy:     operatorID,
		UpdatedBy:     operatorID,
	}, input, nil
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
