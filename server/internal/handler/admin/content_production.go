package admin

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/accountdiagnosis"
	"cy11dsphk/server/internal/agent/benchmark"
	"cy11dsphk/server/internal/agent/copywriting"
	"cy11dsphk/server/internal/agent/hotspottopic"
	"cy11dsphk/server/internal/agent/provider"
	"cy11dsphk/server/internal/agent/review"
	"cy11dsphk/server/internal/agent/storyboard"
	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContentProductionHandler struct{}

func NewContentProductionHandler() *ContentProductionHandler {
	return &ContentProductionHandler{}
}

type listResponse[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

type scriptGeneratePayload struct {
	TopicID          uint64 `json:"topicId" binding:"required"`
	ExtraRequirement string `json:"extraRequirement"`
}

type scriptStatusPayload struct {
	Status string `json:"status"`
}

type contentScriptResponse struct {
	model.ContentScript
	ShootingNotes []string `json:"shootingNotes"`
	Input         any      `json:"input,omitempty"`
	Result        any      `json:"result,omitempty"`
}

type storyboardGeneratePayload struct {
	ScriptID  uint64   `json:"scriptId" binding:"required"`
	Locations []string `json:"locations"`
}

type storyboardStatusPayload struct {
	Status string `json:"status"`
}

type contentStoryboardResponse struct {
	model.ContentStoryboard
	Shots  []storyboard.Shot `json:"shots"`
	Input  any               `json:"input,omitempty"`
	Result any               `json:"result,omitempty"`
}

type schedulePayload struct {
	MerchantID     uint64 `json:"merchantId"`
	StoryboardID   uint64 `json:"storyboardId"`
	VideoTitle     string `json:"videoTitle" binding:"required"`
	PublishTime    string `json:"publishTime"`
	Owner          string `json:"owner"`
	DouyinAccount  string `json:"douyinAccount"`
	MaterialStatus string `json:"materialStatus"`
	Status         string `json:"status"`
	Remark         string `json:"remark"`
}

type shootingTaskPayload struct {
	StoryboardID uint64 `json:"storyboardId" binding:"required"`
	TaskTitle    string `json:"taskTitle"`
	Assignee     string `json:"assignee"`
	ShootTime    string `json:"shootTime"`
	Deadline     string `json:"deadline"`
	Status       string `json:"status"`
	MaterialURL  string `json:"materialUrl"`
	Remark       string `json:"remark"`
}

type shootingTaskStatusPayload struct {
	Status      string `json:"status"`
	MaterialURL string `json:"materialUrl"`
	Remark      string `json:"remark"`
}

type shootingTaskResponse struct {
	model.ShootingTask
	Shots []storyboard.Shot `json:"shots"`
}

type scheduleStatusPayload struct {
	Status         string `json:"status"`
	MaterialStatus string `json:"materialStatus"`
	Remark         string `json:"remark"`
}

type reviewGeneratePayload struct {
	ScheduleID     uint64  `json:"scheduleId" binding:"required"`
	PeriodStart    string  `json:"periodStart"`
	PeriodEnd      string  `json:"periodEnd"`
	PlayCount      int64   `json:"playCount"`
	LikeCount      int64   `json:"likeCount"`
	CommentCount   int64   `json:"commentCount"`
	ShareCount     int64   `json:"shareCount"`
	DealCount      int64   `json:"dealCount"`
	WriteOffAmount float64 `json:"writeOffAmount"`
}

type reviewTaskResponse struct {
	model.ContentReviewTask
	Input  any `json:"input,omitempty"`
	Result any `json:"result,omitempty"`
}

func (h *ContentProductionHandler) ListScripts(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	merchantID := parseUint64Query(c.Query("merchantId"))
	topicID := parseUint64Query(c.Query("topicId"))
	status := strings.TrimSpace(c.Query("status"))

	query := database.DB.Model(&model.ContentScript{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if topicID > 0 {
		query = query.Where("topic_id = ?", topicID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("merchant_name LIKE ? OR topic_title LIKE ? OR title LIKE ? OR opening LIKE ? OR full_script LIKE ?", like, like, like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		basehandler.ServerError(c, "读取文案数量失败")
		return
	}
	var list []model.ContentScript
	if err := query.Order("id DESC").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取文案列表失败")
		return
	}
	basehandler.OK(c, listResponse[contentScriptResponse]{
		List:  buildScriptResponses(list),
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func (h *ContentProductionHandler) GenerateScript(c *gin.Context) {
	var payload scriptGeneratePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择选题")
		return
	}

	input, err := buildCopywritingInput(payload, c.GetUint64("admin_user_id"))
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}

	inputJSON, _ := json.Marshal(input)
	item := model.ContentScript{
		MerchantID:    input.MerchantID,
		MerchantName:  input.Merchant.Name,
		TopicID:       input.TopicID,
		TopicTitle:    input.TopicTitle,
		Status:        model.ContentScriptStatusDraft,
		InputSnapshot: string(inputJSON),
		CreatedBy:     input.Options.OperatorID,
		UpdatedBy:     input.Options.OperatorID,
	}
	if err := database.DB.Create(&item).Error; err != nil {
		basehandler.ServerError(c, "创建文案任务失败")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()
	output, runErr := copywriting.Agent{}.Run(ctx, input)
	if runErr != nil {
		item.ErrorMessage = runErr.Error()
		_ = database.DB.Model(&item).Updates(map[string]any{
			"error_message": item.ErrorMessage,
			"updated_by":    item.UpdatedBy,
		}).Error
		basehandler.ServerError(c, "文案 Agent 执行失败："+runErr.Error())
		return
	}

	resultJSON, _ := json.Marshal(output)
	notesJSON, _ := json.Marshal(output.ShootingNotes)
	if err := database.DB.Model(&item).Updates(map[string]any{
		"title":               strings.TrimSpace(output.Title),
		"opening":             strings.TrimSpace(output.Opening),
		"body":                strings.TrimSpace(output.Body),
		"ending":              strings.TrimSpace(output.Ending),
		"cta":                 strings.TrimSpace(output.CTA),
		"full_script":         strings.TrimSpace(output.FullScript),
		"shooting_notes_json": string(notesJSON),
		"result_json":         string(resultJSON),
		"updated_by":          item.UpdatedBy,
	}).Error; err != nil {
		basehandler.ServerError(c, "保存文案失败")
		return
	}
	if err := database.DB.First(&item, item.ID).Error; err != nil {
		basehandler.ServerError(c, "读取文案失败")
		return
	}
	basehandler.OK(c, buildScriptResponse(item))
}

func (h *ContentProductionHandler) UpdateScriptStatus(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	var payload scriptStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil || !isValidScriptStatus(payload.Status) {
		basehandler.BadRequest(c, "文案状态不正确")
		return
	}
	var item model.ContentScript
	if err := database.DB.First(&item, id).Error; err != nil {
		basehandler.BadRequest(c, "文案不存在")
		return
	}
	if err := database.DB.Model(&item).Updates(map[string]any{"status": payload.Status, "updated_by": c.GetUint64("admin_user_id")}).Error; err != nil {
		basehandler.ServerError(c, "更新文案状态失败")
		return
	}
	_ = database.DB.First(&item, id).Error
	basehandler.OK(c, buildScriptResponse(item))
}

func (h *ContentProductionHandler) ListStoryboards(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	merchantID := parseUint64Query(c.Query("merchantId"))
	scriptID := parseUint64Query(c.Query("scriptId"))
	status := strings.TrimSpace(c.Query("status"))

	query := database.DB.Model(&model.ContentStoryboard{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if scriptID > 0 {
		query = query.Where("script_id = ?", scriptID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("merchant_name LIKE ? OR topic_title LIKE ? OR script_title LIKE ? OR shots_json LIKE ?", like, like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		basehandler.ServerError(c, "读取分镜数量失败")
		return
	}
	var list []model.ContentStoryboard
	if err := query.Order("id DESC").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取分镜列表失败")
		return
	}
	basehandler.OK(c, listResponse[contentStoryboardResponse]{
		List:  buildStoryboardResponses(list),
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func (h *ContentProductionHandler) GenerateStoryboard(c *gin.Context) {
	var payload storyboardGeneratePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择文案")
		return
	}

	input, source, err := buildStoryboardInput(payload, c.GetUint64("admin_user_id"))
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}

	inputJSON, _ := json.Marshal(input)
	item := model.ContentStoryboard{
		MerchantID:    source.MerchantID,
		MerchantName:  source.MerchantName,
		TopicID:       source.TopicID,
		TopicTitle:    source.TopicTitle,
		ScriptID:      source.ID,
		ScriptTitle:   source.Title,
		Status:        model.ContentStoryboardStatusDraft,
		InputSnapshot: string(inputJSON),
		CreatedBy:     input.Options.OperatorID,
		UpdatedBy:     input.Options.OperatorID,
	}
	if err := database.DB.Create(&item).Error; err != nil {
		basehandler.ServerError(c, "创建分镜任务失败")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()
	output, runErr := storyboard.Agent{}.Run(ctx, input)
	if runErr != nil {
		item.ErrorMessage = runErr.Error()
		_ = database.DB.Model(&item).Updates(map[string]any{"error_message": item.ErrorMessage, "updated_by": item.UpdatedBy}).Error
		basehandler.ServerError(c, "分镜 Agent 执行失败："+runErr.Error())
		return
	}

	resultJSON, _ := json.Marshal(output)
	shotsJSON, _ := json.Marshal(output.Shots)
	if err := database.DB.Model(&item).Updates(map[string]any{
		"shots_json":  string(shotsJSON),
		"result_json": string(resultJSON),
		"updated_by":  item.UpdatedBy,
	}).Error; err != nil {
		basehandler.ServerError(c, "保存分镜失败")
		return
	}
	if err := database.DB.First(&item, item.ID).Error; err != nil {
		basehandler.ServerError(c, "读取分镜失败")
		return
	}
	basehandler.OK(c, buildStoryboardResponse(item))
}

func (h *ContentProductionHandler) UpdateStoryboardStatus(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	var payload storyboardStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil || !isValidStoryboardStatus(payload.Status) {
		basehandler.BadRequest(c, "分镜状态不正确")
		return
	}
	var item model.ContentStoryboard
	if err := database.DB.First(&item, id).Error; err != nil {
		basehandler.BadRequest(c, "分镜不存在")
		return
	}
	if err := database.DB.Model(&item).Updates(map[string]any{"status": payload.Status, "updated_by": c.GetUint64("admin_user_id")}).Error; err != nil {
		basehandler.ServerError(c, "更新分镜状态失败")
		return
	}
	_ = database.DB.First(&item, id).Error
	basehandler.OK(c, buildStoryboardResponse(item))
}

func (h *ContentProductionHandler) ListShootingTasks(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	merchantID := parseUint64Query(c.Query("merchantId"))
	storyboardID := parseUint64Query(c.Query("storyboardId"))
	status := strings.TrimSpace(c.Query("status"))

	query := database.DB.Model(&model.ShootingTask{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if storyboardID > 0 {
		query = query.Where("storyboard_id = ?", storyboardID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("merchant_name LIKE ? OR topic_title LIKE ? OR script_title LIKE ? OR task_title LIKE ? OR assignee LIKE ? OR shots_json LIKE ?", like, like, like, like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		basehandler.ServerError(c, "读取拍摄任务数量失败")
		return
	}
	var list []model.ShootingTask
	if err := query.Order("id DESC").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取拍摄任务列表失败")
		return
	}
	basehandler.OK(c, listResponse[shootingTaskResponse]{
		List:  buildShootingTaskResponses(list),
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func (h *ContentProductionHandler) CreateShootingTask(c *gin.Context) {
	var payload shootingTaskPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择分镜")
		return
	}
	item, err := buildShootingTaskFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	item.CreatedBy = c.GetUint64("admin_user_id")
	item.UpdatedBy = item.CreatedBy
	if err := database.DB.Create(&item).Error; err != nil {
		basehandler.ServerError(c, "创建拍摄任务失败")
		return
	}
	basehandler.OK(c, buildShootingTaskResponse(item))
}

func (h *ContentProductionHandler) UpdateShootingTaskStatus(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	var payload shootingTaskStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "拍摄任务状态不正确")
		return
	}
	var item model.ShootingTask
	if err := database.DB.First(&item, id).Error; err != nil {
		basehandler.BadRequest(c, "拍摄任务不存在")
		return
	}
	updates := map[string]any{"updated_by": c.GetUint64("admin_user_id")}
	if payload.Status != "" {
		if !isValidShootingTaskStatus(payload.Status) {
			basehandler.BadRequest(c, "拍摄任务状态不正确")
			return
		}
		updates["status"] = payload.Status
	}
	if payload.MaterialURL != "" {
		updates["material_url"] = strings.TrimSpace(payload.MaterialURL)
	}
	if payload.Remark != "" {
		updates["remark"] = strings.TrimSpace(payload.Remark)
	}
	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		basehandler.ServerError(c, "更新拍摄任务失败")
		return
	}
	_ = database.DB.First(&item, id).Error
	basehandler.OK(c, buildShootingTaskResponse(item))
}

func (h *ContentProductionHandler) ListSchedules(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	merchantID := parseUint64Query(c.Query("merchantId"))
	storyboardID := parseUint64Query(c.Query("storyboardId"))
	status := strings.TrimSpace(c.Query("status"))

	query := database.DB.Model(&model.PublishSchedule{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if storyboardID > 0 {
		query = query.Where("storyboard_id = ?", storyboardID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("merchant_name LIKE ? OR topic_title LIKE ? OR script_title LIKE ? OR video_title LIKE ? OR owner LIKE ? OR douyin_account LIKE ?", like, like, like, like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		basehandler.ServerError(c, "读取排期数量失败")
		return
	}
	var list []model.PublishSchedule
	if err := query.Order("id DESC").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取排期列表失败")
		return
	}
	basehandler.OK(c, listResponse[model.PublishSchedule]{List: list, Total: total, Page: page, Size: size})
}

func (h *ContentProductionHandler) CreateSchedule(c *gin.Context) {
	var payload schedulePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请输入视频标题")
		return
	}
	item, err := buildScheduleFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	item.CreatedBy = c.GetUint64("admin_user_id")
	item.UpdatedBy = item.CreatedBy
	if err := database.DB.Create(&item).Error; err != nil {
		basehandler.ServerError(c, "创建发布排期失败")
		return
	}
	basehandler.OK(c, item)
}

func (h *ContentProductionHandler) UpdateSchedule(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	var payload schedulePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请输入视频标题")
		return
	}
	updates, err := buildScheduleFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	var item model.PublishSchedule
	if err := database.DB.First(&item, id).Error; err != nil {
		basehandler.BadRequest(c, "排期不存在")
		return
	}
	if err := database.DB.Model(&item).Updates(scheduleUpdateMap(updates, c.GetUint64("admin_user_id"))).Error; err != nil {
		basehandler.ServerError(c, "更新发布排期失败")
		return
	}
	_ = database.DB.First(&item, id).Error
	basehandler.OK(c, item)
}

func (h *ContentProductionHandler) UpdateScheduleStatus(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	var payload scheduleStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "排期状态不正确")
		return
	}
	var item model.PublishSchedule
	if err := database.DB.First(&item, id).Error; err != nil {
		basehandler.BadRequest(c, "排期不存在")
		return
	}
	updates := map[string]any{"updated_by": c.GetUint64("admin_user_id")}
	if payload.Status != "" {
		if !isValidScheduleStatus(payload.Status) {
			basehandler.BadRequest(c, "排期状态不正确")
			return
		}
		updates["status"] = payload.Status
	}
	if payload.MaterialStatus != "" {
		updates["material_status"] = strings.TrimSpace(payload.MaterialStatus)
	}
	if payload.Remark != "" {
		updates["remark"] = strings.TrimSpace(payload.Remark)
	}
	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		basehandler.ServerError(c, "更新排期状态失败")
		return
	}
	_ = database.DB.First(&item, id).Error
	basehandler.OK(c, item)
}

func (h *ContentProductionHandler) ListReviews(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	merchantID := parseUint64Query(c.Query("merchantId"))
	query := database.DB.Model(&model.ContentReviewTask{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		basehandler.ServerError(c, "读取复盘数量失败")
		return
	}
	var list []model.ContentReviewTask
	if err := query.Order("id DESC").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取复盘列表失败")
		return
	}
	result := make([]reviewTaskResponse, 0, len(list))
	for _, item := range list {
		result = append(result, buildReviewResponse(item))
	}
	basehandler.OK(c, listResponse[reviewTaskResponse]{List: result, Total: total, Page: page, Size: size})
}

func (h *ContentProductionHandler) GenerateReview(c *gin.Context) {
	var payload reviewGeneratePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择发布视频")
		return
	}
	var schedule model.PublishSchedule
	if err := database.DB.First(&schedule, payload.ScheduleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "发布排期不存在")
			return
		}
		basehandler.ServerError(c, "读取发布排期失败")
		return
	}

	periodStart, periodEnd, err := parseReviewPeriod(payload.PeriodStart, payload.PeriodEnd)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}

	input := review.Input{
		MerchantID:  schedule.MerchantID,
		VideoID:     schedule.ID,
		Title:       schedule.VideoTitle,
		PeriodStart: formatDatePtr(periodStart),
		PeriodEnd:   formatDatePtr(periodEnd),
		Metrics: review.Metrics{
			PlayCount:      payload.PlayCount,
			LikeCount:      payload.LikeCount,
			CommentCount:   payload.CommentCount,
			ShareCount:     payload.ShareCount,
			DealCount:      payload.DealCount,
			WriteOffAmount: payload.WriteOffAmount,
		},
		Options: agent.RunOptions{OperatorID: c.GetUint64("admin_user_id"), DryRun: true},
	}
	inputJSON, _ := json.Marshal(input)
	item := model.ContentReviewTask{
		MerchantID:     schedule.MerchantID,
		MerchantName:   schedule.MerchantName,
		ScheduleID:     schedule.ID,
		VideoTitle:     schedule.VideoTitle,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		Status:         model.ContentReviewStatusCompleted,
		InputSnapshot:  string(inputJSON),
		PlayCount:      payload.PlayCount,
		LikeCount:      payload.LikeCount,
		CommentCount:   payload.CommentCount,
		ShareCount:     payload.ShareCount,
		DealCount:      payload.DealCount,
		WriteOffAmount: payload.WriteOffAmount,
		CreatedBy:      input.Options.OperatorID,
		UpdatedBy:      input.Options.OperatorID,
	}
	if err := database.DB.Create(&item).Error; err != nil {
		basehandler.ServerError(c, "创建复盘任务失败")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()
	output, runErr := review.Agent{}.Run(ctx, input)
	if runErr != nil {
		item.Status = model.ContentReviewStatusFailed
		item.ErrorMessage = runErr.Error()
		_ = database.DB.Model(&item).Updates(map[string]any{"status": item.Status, "error_message": item.ErrorMessage, "updated_by": item.UpdatedBy}).Error
		basehandler.ServerError(c, "复盘 Agent 执行失败："+runErr.Error())
		return
	}
	resultJSON, _ := json.Marshal(output)
	if err := database.DB.Model(&item).Updates(map[string]any{"result_json": string(resultJSON), "updated_by": item.UpdatedBy}).Error; err != nil {
		basehandler.ServerError(c, "保存复盘失败")
		return
	}
	_ = database.DB.Model(&schedule).Updates(map[string]any{"status": model.PublishScheduleStatusReviewed, "updated_by": item.UpdatedBy}).Error
	_ = database.DB.First(&item, item.ID).Error
	basehandler.OK(c, buildReviewResponse(item))
}

func (h *ContentProductionHandler) AgentConfigs(c *gin.Context) {
	names := []string{
		accountdiagnosis.Name,
		benchmark.Name,
		hotspottopic.Name,
		copywriting.Name,
		storyboard.Name,
		review.Name,
	}
	configs := make([]provider.AgentConfig, 0, len(names))
	for _, name := range names {
		configs = append(configs, provider.StepFunConfigForAgent(name))
	}
	basehandler.OK(c, configs)
}

func buildCopywritingInput(payload scriptGeneratePayload, operatorID uint64) (copywriting.Input, error) {
	var topic model.ContentTopic
	if err := database.DB.First(&topic, payload.TopicID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return copywriting.Input{}, errors.New("选题不存在")
		}
		return copywriting.Input{}, errors.New("读取选题失败")
	}
	var merchant model.Merchant
	if err := database.DB.First(&merchant, topic.MerchantID).Error; err != nil {
		return copywriting.Input{}, errors.New("读取商家失败")
	}
	var packages []model.MerchantPackage
	_ = database.DB.Where("merchant_id = ? AND status = ?", merchant.ID, model.MerchantPackageStatusEnabled).Find(&packages).Error
	products := make([]copywriting.ProductContext, 0, len(packages))
	for _, item := range packages {
		products = append(products, copywriting.ProductContext{
			Name:          item.Name,
			SellingPrice:  item.SellingPrice,
			OriginalPrice: item.OriginalPrice,
			TrafficLabel:  item.TrafficLabel,
			ProfitGuard:   item.ProfitGuard,
		})
	}
	return copywriting.Input{
		MerchantID:       merchant.ID,
		TopicID:          topic.ID,
		TopicTitle:       topic.Title,
		TopicHook:        topic.Hook,
		TopicAngle:       topic.Angle,
		ConversionTarget: topic.Target,
		Merchant: copywriting.MerchantContext{
			Name:          merchant.Name,
			Industry:      merchant.Industry,
			City:          merchant.City,
			SellingPoints: splitTextList(merchant.Remark),
		},
		Products:         products,
		ExtraRequirement: strings.TrimSpace(payload.ExtraRequirement),
		Options:          agent.RunOptions{OperatorID: operatorID, DryRun: true},
	}, nil
}

func buildStoryboardInput(payload storyboardGeneratePayload, operatorID uint64) (storyboard.Input, model.ContentScript, error) {
	var script model.ContentScript
	if err := database.DB.First(&script, payload.ScriptID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return storyboard.Input{}, model.ContentScript{}, errors.New("文案不存在")
		}
		return storyboard.Input{}, model.ContentScript{}, errors.New("读取文案失败")
	}
	locations := payload.Locations
	if len(locations) == 0 {
		locations = []string{"门店门头", "门店环境", "产品/套餐展示区", "收银/团购核销处"}
	}
	return storyboard.Input{
		MerchantID:  script.MerchantID,
		ScriptID:    script.ID,
		ScriptTitle: script.Title,
		ScriptText:  script.FullScript,
		Locations:   locations,
		Options:     agent.RunOptions{OperatorID: operatorID, DryRun: true},
	}, script, nil
}

func buildScheduleFromPayload(payload schedulePayload) (model.PublishSchedule, error) {
	title := strings.TrimSpace(payload.VideoTitle)
	if title == "" {
		return model.PublishSchedule{}, errors.New("请输入视频标题")
	}
	status := strings.TrimSpace(payload.Status)
	if status == "" {
		status = model.PublishScheduleStatusPending
	}
	if !isValidScheduleStatus(status) {
		return model.PublishSchedule{}, errors.New("排期状态不正确")
	}
	item := model.PublishSchedule{
		VideoTitle:     title,
		Owner:          strings.TrimSpace(payload.Owner),
		DouyinAccount:  strings.TrimSpace(payload.DouyinAccount),
		MaterialStatus: firstText(strings.TrimSpace(payload.MaterialStatus), "待拍摄"),
		Status:         status,
		Remark:         strings.TrimSpace(payload.Remark),
	}
	if payload.PublishTime != "" {
		parsed, err := time.ParseInLocation("2006-01-02 15:04:05", payload.PublishTime, time.Local)
		if err != nil {
			return model.PublishSchedule{}, errors.New("发布时间格式应为 2006-01-02 15:04:05")
		}
		item.PublishTime = &parsed
	}
	if payload.StoryboardID > 0 {
		var board model.ContentStoryboard
		if err := database.DB.First(&board, payload.StoryboardID).Error; err != nil {
			return model.PublishSchedule{}, errors.New("分镜不存在")
		}
		item.MerchantID = board.MerchantID
		item.MerchantName = board.MerchantName
		item.TopicID = board.TopicID
		item.TopicTitle = board.TopicTitle
		item.ScriptID = board.ScriptID
		item.ScriptTitle = board.ScriptTitle
		item.StoryboardID = board.ID
	} else if payload.MerchantID > 0 {
		var merchant model.Merchant
		if err := database.DB.First(&merchant, payload.MerchantID).Error; err != nil {
			return model.PublishSchedule{}, errors.New("商家不存在")
		}
		item.MerchantID = merchant.ID
		item.MerchantName = merchant.Name
		item.DouyinAccount = firstText(item.DouyinAccount, merchant.DouyinAccount)
	} else {
		return model.PublishSchedule{}, errors.New("请选择商家或分镜")
	}
	return item, nil
}

func buildShootingTaskFromPayload(payload shootingTaskPayload) (model.ShootingTask, error) {
	var board model.ContentStoryboard
	if err := database.DB.First(&board, payload.StoryboardID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ShootingTask{}, errors.New("分镜不存在")
		}
		return model.ShootingTask{}, errors.New("读取分镜失败")
	}
	var shots []storyboard.Shot
	_ = json.Unmarshal([]byte(board.ShotsJSON), &shots)
	title := strings.TrimSpace(payload.TaskTitle)
	if title == "" {
		title = board.MerchantName + "｜" + board.ScriptTitle + "｜拍摄任务"
	}
	status := strings.TrimSpace(payload.Status)
	if status == "" {
		status = model.ShootingTaskStatusPending
	}
	if !isValidShootingTaskStatus(status) {
		return model.ShootingTask{}, errors.New("拍摄任务状态不正确")
	}
	shootTime, err := parseOptionalDateTime(payload.ShootTime)
	if err != nil {
		return model.ShootingTask{}, errors.New("拍摄时间格式应为 2006-01-02 15:04:05")
	}
	deadline, err := parseOptionalDateTime(payload.Deadline)
	if err != nil {
		return model.ShootingTask{}, errors.New("截止时间格式应为 2006-01-02 15:04:05")
	}
	return model.ShootingTask{
		MerchantID:   board.MerchantID,
		MerchantName: board.MerchantName,
		TopicID:      board.TopicID,
		TopicTitle:   board.TopicTitle,
		ScriptID:     board.ScriptID,
		ScriptTitle:  board.ScriptTitle,
		StoryboardID: board.ID,
		TaskTitle:    title,
		ShotCount:    len(shots),
		ShotsJSON:    board.ShotsJSON,
		Assignee:     strings.TrimSpace(payload.Assignee),
		ShootTime:    shootTime,
		Deadline:     deadline,
		Status:       status,
		MaterialURL:  strings.TrimSpace(payload.MaterialURL),
		Remark:       strings.TrimSpace(payload.Remark),
	}, nil
}

func scheduleUpdateMap(item model.PublishSchedule, operatorID uint64) map[string]any {
	return map[string]any{
		"merchant_id":     item.MerchantID,
		"merchant_name":   item.MerchantName,
		"topic_id":        item.TopicID,
		"topic_title":     item.TopicTitle,
		"script_id":       item.ScriptID,
		"script_title":    item.ScriptTitle,
		"storyboard_id":   item.StoryboardID,
		"video_title":     item.VideoTitle,
		"publish_time":    item.PublishTime,
		"owner":           item.Owner,
		"douyin_account":  item.DouyinAccount,
		"material_status": item.MaterialStatus,
		"status":          item.Status,
		"remark":          item.Remark,
		"updated_by":      operatorID,
	}
}

func buildShootingTaskResponses(list []model.ShootingTask) []shootingTaskResponse {
	result := make([]shootingTaskResponse, 0, len(list))
	for _, item := range list {
		result = append(result, buildShootingTaskResponse(item))
	}
	return result
}

func buildShootingTaskResponse(item model.ShootingTask) shootingTaskResponse {
	var shots []storyboard.Shot
	_ = json.Unmarshal([]byte(item.ShotsJSON), &shots)
	return shootingTaskResponse{
		ShootingTask: item,
		Shots:        shots,
	}
}

func buildScriptResponses(list []model.ContentScript) []contentScriptResponse {
	result := make([]contentScriptResponse, 0, len(list))
	for _, item := range list {
		result = append(result, buildScriptResponse(item))
	}
	return result
}

func buildScriptResponse(item model.ContentScript) contentScriptResponse {
	var notes []string
	_ = json.Unmarshal([]byte(item.ShootingNotesJSON), &notes)
	return contentScriptResponse{
		ContentScript: item,
		ShootingNotes: notes,
		Input:         parseJSONField(item.InputSnapshot),
		Result:        parseJSONField(item.ResultJSON),
	}
}

func buildStoryboardResponses(list []model.ContentStoryboard) []contentStoryboardResponse {
	result := make([]contentStoryboardResponse, 0, len(list))
	for _, item := range list {
		result = append(result, buildStoryboardResponse(item))
	}
	return result
}

func buildStoryboardResponse(item model.ContentStoryboard) contentStoryboardResponse {
	var shots []storyboard.Shot
	_ = json.Unmarshal([]byte(item.ShotsJSON), &shots)
	return contentStoryboardResponse{
		ContentStoryboard: item,
		Shots:             shots,
		Input:             parseJSONField(item.InputSnapshot),
		Result:            parseJSONField(item.ResultJSON),
	}
}

func buildReviewResponse(item model.ContentReviewTask) reviewTaskResponse {
	return reviewTaskResponse{
		ContentReviewTask: item,
		Input:             parseJSONField(item.InputSnapshot),
		Result:            parseJSONField(item.ResultJSON),
	}
}

func isValidScriptStatus(status string) bool {
	return status == model.ContentScriptStatusDraft || status == model.ContentScriptStatusConfirmed || status == model.ContentScriptStatusDisabled
}

func isValidStoryboardStatus(status string) bool {
	return status == model.ContentStoryboardStatusDraft || status == model.ContentStoryboardStatusConfirmed
}

func isValidShootingTaskStatus(status string) bool {
	return status == model.ShootingTaskStatusPending ||
		status == model.ShootingTaskStatusShooting ||
		status == model.ShootingTaskStatusShot ||
		status == model.ShootingTaskStatusEdited ||
		status == model.ShootingTaskStatusDone
}

func isValidScheduleStatus(status string) bool {
	return status == model.PublishScheduleStatusPending || status == model.PublishScheduleStatusPublished || status == model.PublishScheduleStatusReviewed
}

func parseOptionalDateTime(raw string) (*time.Time, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(raw), time.Local)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseOptionalDate(raw string) (*time.Time, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	parsed, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(raw), time.Local)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseReviewPeriod(startRaw string, endRaw string) (*time.Time, *time.Time, error) {
	start, err := parseOptionalDate(startRaw)
	if err != nil {
		return nil, nil, errors.New("复盘开始日期格式应为 2006-01-02")
	}
	end, err := parseOptionalDate(endRaw)
	if err != nil {
		return nil, nil, errors.New("复盘结束日期格式应为 2006-01-02")
	}
	if start != nil && end != nil && end.Before(*start) {
		return nil, nil, errors.New("复盘结束日期不能早于开始日期")
	}
	return start, end, nil
}

func formatDatePtr(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("2006-01-02")
}

func splitTextList(raw string) []string {
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == '\n' || r == '，' || r == ',' || r == ';' || r == '；'
	})
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func firstText(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
