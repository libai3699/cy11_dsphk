package admin

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"cy11dsphk/server/internal/agent"
	benchmarkagent "cy11dsphk/server/internal/agent/benchmark"
	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BenchmarkHandler struct{}

func NewBenchmarkHandler() *BenchmarkHandler {
	return &BenchmarkHandler{}
}

type benchmarkAccountPayload struct {
	MerchantID     uint64  `json:"merchantId" binding:"required"`
	AccountName    string  `json:"accountName" binding:"required"`
	Platform       string  `json:"platform"`
	City           string  `json:"city"`
	Industry       string  `json:"industry"`
	AccountURL     string  `json:"accountUrl"`
	FollowerCount  float64 `json:"followerCount"`
	BestPlayCount  float64 `json:"bestPlayCount"`
	LatestHitTitle string  `json:"latestHitTitle"`
	Takeaway       string  `json:"takeaway"`
	Risk           string  `json:"risk"`
	Status         string  `json:"status"`
	Remark         string  `json:"remark"`
}

type benchmarkAnalyzeBatchPayload struct {
	MerchantID   uint64   `json:"merchantId" binding:"required"`
	BenchmarkIDs []uint64 `json:"benchmarkIds"`
}

type benchmarkAccountListResponse struct {
	List  []model.BenchmarkAccount `json:"list"`
	Total int64                    `json:"total"`
	Page  int                      `json:"page"`
	Size  int                      `json:"size"`
}

type benchmarkAnalysisListResponse struct {
	List  []benchmarkAnalysisResponse `json:"list"`
	Total int64                       `json:"total"`
	Page  int                         `json:"page"`
	Size  int                         `json:"size"`
}

type benchmarkAnalysisResponse struct {
	model.BenchmarkAnalysisTask
	Input  any `json:"input,omitempty"`
	Result any `json:"result,omitempty"`
}

func (h *BenchmarkHandler) List(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	merchantID := parseUint64Query(c.Query("merchantId"))
	status := strings.TrimSpace(c.Query("status"))

	query := database.DB.Model(&model.BenchmarkAccount{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"merchant_name LIKE ? OR account_name LIKE ? OR platform LIKE ? OR city LIKE ? OR industry LIKE ? OR latest_hit_title LIKE ? OR takeaway LIKE ? OR risk LIKE ?",
			like,
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
		basehandler.ServerError(c, "读取对标账号数量失败")
		return
	}

	var list []model.BenchmarkAccount
	if err := query.Order("id DESC").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取对标账号失败")
		return
	}

	basehandler.OK(c, benchmarkAccountListResponse{List: list, Total: total, Page: page, Size: size})
}

func (h *BenchmarkHandler) Get(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var item model.BenchmarkAccount
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "对标账号不存在")
			return
		}
		basehandler.ServerError(c, "读取对标账号失败")
		return
	}

	basehandler.OK(c, item)
}

func (h *BenchmarkHandler) Create(c *gin.Context) {
	var payload benchmarkAccountPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并输入对标账号")
		return
	}

	item, err := buildBenchmarkAccountFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	item.CreatedBy = c.GetUint64("admin_user_id")
	item.UpdatedBy = item.CreatedBy

	if err := database.DB.Create(&item).Error; err != nil {
		basehandler.ServerError(c, "创建对标账号失败")
		return
	}

	basehandler.OK(c, item)
}

func (h *BenchmarkHandler) Update(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var payload benchmarkAccountPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并输入对标账号")
		return
	}

	updates, err := buildBenchmarkAccountFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	updates.UpdatedBy = c.GetUint64("admin_user_id")

	var item model.BenchmarkAccount
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "对标账号不存在")
			return
		}
		basehandler.ServerError(c, "读取对标账号失败")
		return
	}

	if err := database.DB.Model(&item).Updates(benchmarkAccountUpdateMap(updates)).Error; err != nil {
		basehandler.ServerError(c, "更新对标账号失败")
		return
	}
	if err := database.DB.First(&item, id).Error; err != nil {
		basehandler.ServerError(c, "读取对标账号失败")
		return
	}

	basehandler.OK(c, item)
}

func (h *BenchmarkHandler) Delete(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	if err := database.DB.Delete(&model.BenchmarkAccount{}, id).Error; err != nil {
		basehandler.ServerError(c, "删除对标账号失败")
		return
	}
	basehandler.OK(c, gin.H{"id": id})
}

func (h *BenchmarkHandler) Analyze(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var item model.BenchmarkAccount
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "对标账号不存在")
			return
		}
		basehandler.ServerError(c, "读取对标账号失败")
		return
	}

	var merchant model.Merchant
	if err := database.DB.First(&merchant, item.MerchantID).Error; err != nil {
		basehandler.ServerError(c, "读取商家失败")
		return
	}

	task, runErr := runBenchmarkAnalysis(c, merchant, []model.BenchmarkAccount{item})
	if runErr != nil {
		basehandler.ServerError(c, "对标分析失败")
		return
	}
	basehandler.OK(c, buildBenchmarkAnalysisResponse(task))
}

func (h *BenchmarkHandler) AnalyzeBatch(c *gin.Context) {
	var payload benchmarkAnalyzeBatchPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家")
		return
	}

	var merchant model.Merchant
	if err := database.DB.First(&merchant, payload.MerchantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "商家不存在")
			return
		}
		basehandler.ServerError(c, "读取商家失败")
		return
	}

	var accounts []model.BenchmarkAccount
	query := database.DB.Where("merchant_id = ?", merchant.ID)
	if len(payload.BenchmarkIDs) > 0 {
		query = query.Where("id IN ?", payload.BenchmarkIDs)
	}
	if err := query.Order("best_play_count DESC, follower_count DESC, id DESC").Find(&accounts).Error; err != nil {
		basehandler.ServerError(c, "读取对标账号失败")
		return
	}
	if len(accounts) == 0 {
		basehandler.BadRequest(c, "该商家还没有可分析的对标账号")
		return
	}

	task, runErr := runBenchmarkAnalysis(c, merchant, accounts)
	if runErr != nil {
		basehandler.ServerError(c, "对标分析失败")
		return
	}
	basehandler.OK(c, buildBenchmarkAnalysisResponse(task))
}

func runBenchmarkAnalysis(c *gin.Context, merchant model.Merchant, accounts []model.BenchmarkAccount) (model.BenchmarkAnalysisTask, error) {
	operatorID := c.GetUint64("admin_user_id")
	input := benchmarkagent.Input{
		MerchantID: merchant.ID,
		Industry:   merchant.Industry,
		City:       merchant.City,
		Options:    agent.RunOptions{OperatorID: operatorID, DryRun: true},
	}
	for _, item := range accounts {
		input.BenchmarkAccounts = append(input.BenchmarkAccounts, benchmarkagent.BenchmarkAccount{
			Name:       item.AccountName,
			PlatformID: item.AccountURL,
			Reason:     item.Takeaway,
			Tags:       []string{item.City, item.Industry, item.LatestHitTitle},
		})
	}

	benchmarkID, benchmarkName := benchmarkIdentity(accounts)
	snapshot := map[string]any{
		"merchant":       merchant,
		"benchmarks":     accounts,
		"agentInput":     input,
		"snapshotTime":   time.Now().Format("2006-01-02 15:04:05"),
		"snapshotSource": "admin_benchmark_analysis",
	}
	inputSnapshot, err := json.Marshal(snapshot)
	if err != nil {
		return model.BenchmarkAnalysisTask{}, errors.New("输入快照生成失败")
	}

	task := model.BenchmarkAnalysisTask{
		MerchantID:         merchant.ID,
		MerchantName:       merchant.Name,
		BenchmarkAccountID: benchmarkID,
		BenchmarkName:      benchmarkName,
		Status:             model.BenchmarkAnalysisStatusRunning,
		InputSnapshot:      string(inputSnapshot),
		CreatedBy:          operatorID,
		UpdatedBy:          operatorID,
	}
	if err := database.DB.Create(&task).Error; err != nil {
		return task, err
	}

	output, runErr := benchmarkagent.Agent{}.Run(context.Background(), input)
	if runErr != nil {
		task.Status = model.BenchmarkAnalysisStatusFailed
		task.ErrorMessage = runErr.Error()
		_ = database.DB.Model(&task).Updates(map[string]any{
			"status":        task.Status,
			"error_message": task.ErrorMessage,
			"updated_by":    task.UpdatedBy,
		}).Error
		return task, runErr
	}

	resultJSON, err := json.Marshal(output)
	if err != nil {
		return task, err
	}
	task.Status = model.BenchmarkAnalysisStatusCompleted
	task.ResultJSON = string(resultJSON)
	if err := database.DB.Model(&task).Updates(map[string]any{
		"status":      task.Status,
		"result_json": task.ResultJSON,
		"updated_by":  task.UpdatedBy,
	}).Error; err != nil {
		return task, err
	}
	ids := make([]uint64, 0, len(accounts))
	for _, item := range accounts {
		ids = append(ids, item.ID)
	}
	_ = database.DB.Model(&model.BenchmarkAccount{}).Where("id IN ?", ids).Updates(map[string]any{
		"status":     model.BenchmarkAccountStatusAnalyzed,
		"updated_by": task.UpdatedBy,
	}).Error

	if err := database.DB.First(&task, task.ID).Error; err != nil {
		return task, err
	}
	return task, nil
}

func (h *BenchmarkHandler) AnalysisList(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	merchantID := parseUint64Query(c.Query("merchantId"))
	benchmarkID := parseUint64Query(c.Query("benchmarkAccountId"))

	query := database.DB.Model(&model.BenchmarkAnalysisTask{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if benchmarkID > 0 {
		query = query.Where("benchmark_account_id = ?", benchmarkID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		basehandler.ServerError(c, "读取分析任务数量失败")
		return
	}
	var list []model.BenchmarkAnalysisTask
	if err := query.Order("id DESC").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取分析任务失败")
		return
	}
	basehandler.OK(c, benchmarkAnalysisListResponse{
		List:  buildBenchmarkAnalysisResponses(list),
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func buildBenchmarkAccountFromPayload(payload benchmarkAccountPayload) (model.BenchmarkAccount, error) {
	name := strings.TrimSpace(payload.AccountName)
	if name == "" {
		return model.BenchmarkAccount{}, errors.New("请输入对标账号")
	}
	if payload.MerchantID == 0 {
		return model.BenchmarkAccount{}, errors.New("请选择商家")
	}
	var merchant model.Merchant
	if err := database.DB.First(&merchant, payload.MerchantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.BenchmarkAccount{}, errors.New("商家不存在")
		}
		return model.BenchmarkAccount{}, errors.New("读取商家失败")
	}

	status := strings.TrimSpace(payload.Status)
	if status == "" {
		status = model.BenchmarkAccountStatusPending
	}
	if !isValidBenchmarkStatus(status) {
		return model.BenchmarkAccount{}, errors.New("对标账号状态不正确")
	}
	platform := strings.TrimSpace(payload.Platform)
	if platform == "" {
		platform = "抖音"
	}

	city := strings.TrimSpace(payload.City)
	if city == "" {
		city = merchant.City
	}
	industry := strings.TrimSpace(payload.Industry)
	if industry == "" {
		industry = merchant.Industry
	}

	return model.BenchmarkAccount{
		MerchantID:     merchant.ID,
		MerchantName:   merchant.Name,
		AccountName:    name,
		Platform:       platform,
		City:           city,
		Industry:       industry,
		AccountURL:     strings.TrimSpace(payload.AccountURL),
		FollowerCount:  payload.FollowerCount,
		BestPlayCount:  payload.BestPlayCount,
		LatestHitTitle: strings.TrimSpace(payload.LatestHitTitle),
		Takeaway:       strings.TrimSpace(payload.Takeaway),
		Risk:           strings.TrimSpace(payload.Risk),
		Status:         status,
		Remark:         strings.TrimSpace(payload.Remark),
	}, nil
}

func benchmarkAccountUpdateMap(item model.BenchmarkAccount) map[string]any {
	return map[string]any{
		"merchant_id":      item.MerchantID,
		"merchant_name":    item.MerchantName,
		"account_name":     item.AccountName,
		"platform":         item.Platform,
		"city":             item.City,
		"industry":         item.Industry,
		"account_url":      item.AccountURL,
		"follower_count":   item.FollowerCount,
		"best_play_count":  item.BestPlayCount,
		"latest_hit_title": item.LatestHitTitle,
		"takeaway":         item.Takeaway,
		"risk":             item.Risk,
		"status":           item.Status,
		"remark":           item.Remark,
		"updated_by":       item.UpdatedBy,
	}
}

func isValidBenchmarkStatus(status string) bool {
	switch status {
	case model.BenchmarkAccountStatusPending, model.BenchmarkAccountStatusAnalyzed, model.BenchmarkAccountStatusDisabled:
		return true
	default:
		return false
	}
}

func benchmarkIdentity(accounts []model.BenchmarkAccount) (uint64, string) {
	if len(accounts) == 0 {
		return 0, ""
	}
	if len(accounts) == 1 {
		return accounts[0].ID, accounts[0].AccountName
	}
	names := make([]string, 0, len(accounts))
	for _, item := range accounts {
		names = append(names, item.AccountName)
	}
	return 0, strings.Join(names, "、")
}

func buildBenchmarkAnalysisResponses(list []model.BenchmarkAnalysisTask) []benchmarkAnalysisResponse {
	result := make([]benchmarkAnalysisResponse, 0, len(list))
	for _, item := range list {
		result = append(result, buildBenchmarkAnalysisResponse(item))
	}
	return result
}

func buildBenchmarkAnalysisResponse(task model.BenchmarkAnalysisTask) benchmarkAnalysisResponse {
	return benchmarkAnalysisResponse{
		BenchmarkAnalysisTask: task,
		Input:                 parseJSONField(task.InputSnapshot),
		Result:                parseJSONField(task.ResultJSON),
	}
}
