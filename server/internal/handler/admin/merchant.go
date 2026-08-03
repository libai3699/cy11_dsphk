package admin

import (
	"errors"
	"strconv"
	"strings"

	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"
	"cy11dsphk/server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MerchantHandler struct{}

func NewMerchantHandler() *MerchantHandler {
	return &MerchantHandler{}
}

type merchantPayload struct {
	Name               string  `json:"name" binding:"required"`
	Industry           string  `json:"industry"`
	City               string  `json:"city"`
	Address            string  `json:"address"`
	ContactName        string  `json:"contactName"`
	ContactPhone       string  `json:"contactPhone"`
	DouyinAccount      string  `json:"douyinAccount"`
	DouyinLaikeAccount string  `json:"douyinLaikeAccount"`
	CooperationType    string  `json:"cooperationType"`
	CommissionRate     float64 `json:"commissionRate"`
	Stage              string  `json:"stage"`
	Status             *int    `json:"status"`
	Remark             string  `json:"remark"`
}

type merchantListResponse struct {
	List  []model.Merchant `json:"list"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

type merchantWorkspaceResponse struct {
	Merchant     model.Merchant              `json:"merchant"`
	Completeness int                         `json:"completeness"`
	Metrics      merchantWorkspaceMetrics    `json:"metrics"`
	Requirements []merchantRequirementStatus `json:"requirements"`
}

type merchantWorkspaceMetrics struct {
	PackageCount             int64 `json:"packageCount"`
	EnabledPackageCount      int64 `json:"enabledPackageCount"`
	AccountAuthCount         int64 `json:"accountAuthCount"`
	ActiveAccountAuthCount   int64 `json:"activeAccountAuthCount"`
	DiagnosisCount           int64 `json:"diagnosisCount"`
	CompletedDiagnosisCount  int64 `json:"completedDiagnosisCount"`
	BenchmarkCount           int64 `json:"benchmarkCount"`
	AnalyzedBenchmarkCount   int64 `json:"analyzedBenchmarkCount"`
	TopicCount               int64 `json:"topicCount"`
	AcceptedTopicCount       int64 `json:"acceptedTopicCount"`
	ScriptCount              int64 `json:"scriptCount"`
	ConfirmedScriptCount     int64 `json:"confirmedScriptCount"`
	StoryboardCount          int64 `json:"storyboardCount"`
	ConfirmedStoryboardCount int64 `json:"confirmedStoryboardCount"`
	ShootingTaskCount        int64 `json:"shootingTaskCount"`
	ReadyShootingTaskCount   int64 `json:"readyShootingTaskCount"`
	ScheduleCount            int64 `json:"scheduleCount"`
	PublishedScheduleCount   int64 `json:"publishedScheduleCount"`
	ReviewCount              int64 `json:"reviewCount"`
}

type merchantRequirementStatus struct {
	Key     string   `json:"key"`
	Title   string   `json:"title"`
	Done    bool     `json:"done"`
	Missing []string `json:"missing"`
}

func (h *MerchantHandler) List(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	keyword := strings.TrimSpace(c.Query("keyword"))

	query := database.DB.Model(&model.Merchant{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"name LIKE ? OR industry LIKE ? OR city LIKE ? OR contact_name LIKE ? OR contact_phone LIKE ?",
			like,
			like,
			like,
			like,
			like,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		basehandler.ServerError(c, "读取商家数量失败")
		return
	}

	var list []model.Merchant
	if err := query.
		Order("id DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取商家列表失败")
		return
	}

	basehandler.OK(c, merchantListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func (h *MerchantHandler) Get(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var merchant model.Merchant
	if err := database.DB.First(&merchant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "商家不存在")
			return
		}
		basehandler.ServerError(c, "读取商家失败")
		return
	}

	basehandler.OK(c, merchant)
}

func (h *MerchantHandler) Workspace(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var merchant model.Merchant
	if err := database.DB.First(&merchant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "商家不存在")
			return
		}
		basehandler.ServerError(c, "读取商家失败")
		return
	}

	metrics := buildMerchantWorkspaceMetrics(merchant.ID)
	requirements := buildMerchantRequirements(merchant, metrics)
	basehandler.OK(c, merchantWorkspaceResponse{
		Merchant:     merchant,
		Completeness: calculateCompleteness(requirements),
		Metrics:      metrics,
		Requirements: requirements,
	})
}

func (h *MerchantHandler) Create(c *gin.Context) {
	var payload merchantPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请输入商家名称")
		return
	}

	merchant, err := buildMerchantFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	merchant.CreatedBy = c.GetUint64("admin_user_id")
	merchant.UpdatedBy = merchant.CreatedBy

	if err := database.DB.Create(&merchant).Error; err != nil {
		basehandler.ServerError(c, "创建商家失败")
		return
	}

	basehandler.OK(c, merchant)
}

func buildMerchantWorkspaceMetrics(merchantID uint64) merchantWorkspaceMetrics {
	return merchantWorkspaceMetrics{
		PackageCount:             countByMerchant(&model.MerchantPackage{}, merchantID, ""),
		EnabledPackageCount:      countByMerchant(&model.MerchantPackage{}, merchantID, "status = ?", model.MerchantPackageStatusEnabled),
		AccountAuthCount:         countByMerchant(&model.MerchantAccountAuth{}, merchantID, ""),
		ActiveAccountAuthCount:   countByMerchant(&model.MerchantAccountAuth{}, merchantID, "auth_status = ?", model.MerchantAccountAuthStatusActive),
		DiagnosisCount:           countByMerchant(&model.AccountDiagnosisTask{}, merchantID, ""),
		CompletedDiagnosisCount:  countByMerchant(&model.AccountDiagnosisTask{}, merchantID, "status = ?", model.AccountDiagnosisStatusCompleted),
		BenchmarkCount:           countByMerchant(&model.BenchmarkAccount{}, merchantID, ""),
		AnalyzedBenchmarkCount:   countByMerchant(&model.BenchmarkAccount{}, merchantID, "status = ?", model.BenchmarkAccountStatusAnalyzed),
		TopicCount:               countByMerchant(&model.ContentTopic{}, merchantID, ""),
		AcceptedTopicCount:       countByMerchant(&model.ContentTopic{}, merchantID, "status = ?", model.ContentTopicStatusAccepted),
		ScriptCount:              countByMerchant(&model.ContentScript{}, merchantID, ""),
		ConfirmedScriptCount:     countByMerchant(&model.ContentScript{}, merchantID, "status = ?", model.ContentScriptStatusConfirmed),
		StoryboardCount:          countByMerchant(&model.ContentStoryboard{}, merchantID, ""),
		ConfirmedStoryboardCount: countByMerchant(&model.ContentStoryboard{}, merchantID, "status = ?", model.ContentStoryboardStatusConfirmed),
		ShootingTaskCount:        countByMerchant(&model.ShootingTask{}, merchantID, ""),
		ReadyShootingTaskCount:   countByMerchant(&model.ShootingTask{}, merchantID, "status IN ?", []string{model.ShootingTaskStatusEdited, model.ShootingTaskStatusDone}),
		ScheduleCount:            countByMerchant(&model.PublishSchedule{}, merchantID, ""),
		PublishedScheduleCount:   countByMerchant(&model.PublishSchedule{}, merchantID, "status IN ?", []string{model.PublishScheduleStatusPublished, model.PublishScheduleStatusReviewed}),
		ReviewCount:              countByMerchant(&model.ContentReviewTask{}, merchantID, ""),
	}
}

func countByMerchant(table any, merchantID uint64, condition string, args ...any) int64 {
	query := database.DB.Model(table).Where("merchant_id = ?", merchantID)
	if condition != "" {
		query = query.Where(condition, args...)
	}
	var total int64
	_ = query.Count(&total).Error
	return total
}

func buildMerchantRequirements(merchant model.Merchant, metrics merchantWorkspaceMetrics) []merchantRequirementStatus {
	return []merchantRequirementStatus{
		{
			Key:     "profile",
			Title:   "基础档案",
			Done:    merchant.Name != "" && merchant.Industry != "" && merchant.City != "" && merchant.ContactName != "" && merchant.ContactPhone != "",
			Missing: missingFields([]fieldCheck{{"商家名称", merchant.Name != ""}, {"行业", merchant.Industry != ""}, {"城市", merchant.City != ""}, {"联系人", merchant.ContactName != ""}, {"联系电话", merchant.ContactPhone != ""}}),
		},
		{
			Key:     "cooperation",
			Title:   "合作规则",
			Done:    merchant.CooperationType != "" && merchant.CommissionRate > 0,
			Missing: missingFields([]fieldCheck{{"合作方式", merchant.CooperationType != ""}, {"分成比例", merchant.CommissionRate > 0}}),
		},
		{
			Key:     "package",
			Title:   "团购套餐",
			Done:    metrics.EnabledPackageCount > 0,
			Missing: missingCount(metrics.EnabledPackageCount > 0, "至少 1 个启用套餐"),
		},
		{
			Key:     "account",
			Title:   "账号授权",
			Done:    merchant.DouyinAccount != "" && metrics.ActiveAccountAuthCount > 0,
			Missing: missingFields([]fieldCheck{{"抖音账号", merchant.DouyinAccount != ""}, {"至少 1 条已授权记录", metrics.ActiveAccountAuthCount > 0}}),
		},
	}
}

type fieldCheck struct {
	Label string
	Done  bool
}

func missingFields(fields []fieldCheck) []string {
	result := []string{}
	for _, field := range fields {
		if !field.Done {
			result = append(result, field.Label)
		}
	}
	return result
}

func missingCount(done bool, label string) []string {
	if done {
		return []string{}
	}
	return []string{label}
}

func calculateCompleteness(requirements []merchantRequirementStatus) int {
	if len(requirements) == 0 {
		return 0
	}
	done := 0
	for _, requirement := range requirements {
		if requirement.Done {
			done++
		}
	}
	return done * 100 / len(requirements)
}

func (h *MerchantHandler) Update(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var payload merchantPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请输入商家名称")
		return
	}

	updates, err := buildMerchantFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	updates.UpdatedBy = c.GetUint64("admin_user_id")

	var merchant model.Merchant
	if err := database.DB.First(&merchant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "商家不存在")
			return
		}
		basehandler.ServerError(c, "读取商家失败")
		return
	}

	if err := database.DB.Model(&merchant).Updates(merchantUpdateMap(updates)).Error; err != nil {
		basehandler.ServerError(c, "更新商家失败")
		return
	}
	if err := database.DB.First(&merchant, id).Error; err != nil {
		basehandler.ServerError(c, "读取商家失败")
		return
	}

	basehandler.OK(c, merchant)
}

func (h *MerchantHandler) Delete(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	if !currentAdminHasRole(c, "super_admin") {
		basehandler.Forbidden(c, "只有超级管理员可以删除商家")
		return
	}

	var merchant model.Merchant
	if err := database.DB.First(&merchant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "商家不存在")
			return
		}
		basehandler.ServerError(c, "读取商家失败")
		return
	}

	deleted, err := deleteMerchantCascade(id)
	if err != nil {
		basehandler.ServerError(c, "删除商家失败")
		return
	}

	basehandler.OK(c, gin.H{
		"id":      id,
		"name":    merchant.Name,
		"deleted": deleted,
	})
}

func deleteMerchantCascade(merchantID uint64) (map[string]int64, error) {
	deleted := map[string]int64{}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		steps := []struct {
			key   string
			model any
		}{
			{key: "content_reviews", model: &model.ContentReviewTask{}},
			{key: "publish_schedules", model: &model.PublishSchedule{}},
			{key: "shooting_tasks", model: &model.ShootingTask{}},
			{key: "storyboards", model: &model.ContentStoryboard{}},
			{key: "scripts", model: &model.ContentScript{}},
			{key: "topics", model: &model.ContentTopic{}},
			{key: "hotspot_topic_tasks", model: &model.HotspotTopicTask{}},
			{key: "benchmark_analysis_tasks", model: &model.BenchmarkAnalysisTask{}},
			{key: "benchmark_accounts", model: &model.BenchmarkAccount{}},
			{key: "account_diagnoses", model: &model.AccountDiagnosisTask{}},
			{key: "account_auths", model: &model.MerchantAccountAuth{}},
			{key: "packages", model: &model.MerchantPackage{}},
		}

		for _, step := range steps {
			result := tx.Where("merchant_id = ?", merchantID).Delete(step.model)
			if result.Error != nil {
				return result.Error
			}
			deleted[step.key] = result.RowsAffected
		}

		result := tx.Delete(&model.Merchant{}, merchantID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		deleted["merchants"] = result.RowsAffected
		return nil
	})
	return deleted, err
}

func currentAdminHasRole(c *gin.Context, roleCode string) bool {
	roleCodes, err := service.GetUserRoleCodes(c.GetUint64("admin_user_id"))
	if err != nil {
		return false
	}
	for _, code := range roleCodes {
		if code == roleCode {
			return true
		}
	}
	return false
}

func merchantUpdateMap(merchant model.Merchant) map[string]any {
	return map[string]any{
		"name":                 merchant.Name,
		"industry":             merchant.Industry,
		"city":                 merchant.City,
		"address":              merchant.Address,
		"contact_name":         merchant.ContactName,
		"contact_phone":        merchant.ContactPhone,
		"douyin_account":       merchant.DouyinAccount,
		"douyin_laike_account": merchant.DouyinLaikeAccount,
		"cooperation_type":     merchant.CooperationType,
		"commission_rate":      merchant.CommissionRate,
		"stage":                merchant.Stage,
		"status":               merchant.Status,
		"remark":               merchant.Remark,
		"updated_by":           merchant.UpdatedBy,
	}
}

func buildMerchantFromPayload(payload merchantPayload) (model.Merchant, error) {
	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return model.Merchant{}, errors.New("请输入商家名称")
	}

	commissionRate := payload.CommissionRate
	if commissionRate <= 0 {
		commissionRate = 10
	}
	if commissionRate > 100 {
		return model.Merchant{}, errors.New("分成比例不能超过 100")
	}

	status := model.MerchantStatusEnabled
	if payload.Status != nil {
		status = *payload.Status
	}
	if status != model.MerchantStatusEnabled && status != model.MerchantStatusDisabled {
		return model.Merchant{}, errors.New("商家状态不正确")
	}

	stage := strings.TrimSpace(payload.Stage)
	if stage == "" {
		stage = "已建档"
	}
	cooperationType := strings.TrimSpace(payload.CooperationType)
	if cooperationType == "" {
		cooperationType = "成交提点"
	}

	return model.Merchant{
		Name:               name,
		Industry:           strings.TrimSpace(payload.Industry),
		City:               strings.TrimSpace(payload.City),
		Address:            strings.TrimSpace(payload.Address),
		ContactName:        strings.TrimSpace(payload.ContactName),
		ContactPhone:       strings.TrimSpace(payload.ContactPhone),
		DouyinAccount:      strings.TrimSpace(payload.DouyinAccount),
		DouyinLaikeAccount: strings.TrimSpace(payload.DouyinLaikeAccount),
		CooperationType:    cooperationType,
		CommissionRate:     commissionRate,
		Stage:              stage,
		Status:             status,
		Remark:             strings.TrimSpace(payload.Remark),
	}, nil
}

func parsePositiveInt(raw string, fallback int) int {
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func parseUint64Param(c *gin.Context, name string) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		basehandler.BadRequest(c, "参数错误")
		return 0, false
	}
	return id, true
}
