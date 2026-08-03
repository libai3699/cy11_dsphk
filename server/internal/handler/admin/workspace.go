package admin

import (
	"fmt"
	"math"
	"time"

	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"

	"github.com/gin-gonic/gin"
)

type WorkspaceHandler struct{}

func NewWorkspaceHandler() *WorkspaceHandler {
	return &WorkspaceHandler{}
}

type workspaceOverviewResponse struct {
	Metrics       []workspaceMetricRow       `json:"metrics"`
	Merchants     []workspaceMerchantRow     `json:"merchants"`
	ShootingTasks []workspaceShootingTaskRow `json:"shootingTasks"`
	Topics        []workspaceTopicRow        `json:"topics"`
	Reviews       []workspaceReviewRow       `json:"reviews"`
}

type workspaceMetricRow struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
	Hint  string `json:"hint"`
	Path  string `json:"path"`
}

type workspaceMerchantRow struct {
	ID                   uint64  `json:"id"`
	Name                 string  `json:"name"`
	Industry             string  `json:"industry"`
	City                 string  `json:"city"`
	Stage                string  `json:"stage"`
	CommissionRate       float64 `json:"commissionRate"`
	RecentWriteOffAmount float64 `json:"recentWriteOffAmount"`
	EstimatedCommission  float64 `json:"estimatedCommission"`
	NextAction           string  `json:"nextAction"`
}

type workspaceShootingTaskRow struct {
	ID           uint64     `json:"id"`
	MerchantID   uint64     `json:"merchantId"`
	MerchantName string     `json:"merchantName"`
	TaskTitle    string     `json:"taskTitle"`
	Assignee     string     `json:"assignee"`
	Deadline     *time.Time `json:"deadline"`
	Status       string     `json:"status"`
	ShotCount    int        `json:"shotCount"`
}

type workspaceTopicRow struct {
	ID            uint64 `json:"id"`
	MerchantID    uint64 `json:"merchantId"`
	MerchantName  string `json:"merchantName"`
	Title         string `json:"title"`
	Hook          string `json:"hook"`
	PublishWindow string `json:"publishWindow"`
	Status        string `json:"status"`
}

type workspaceReviewRow struct {
	ID             uint64  `json:"id"`
	MerchantID     uint64  `json:"merchantId"`
	MerchantName   string  `json:"merchantName"`
	VideoTitle     string  `json:"videoTitle"`
	PlayCount      int64   `json:"playCount"`
	DealCount      int64   `json:"dealCount"`
	WriteOffAmount float64 `json:"writeOffAmount"`
	Status         string  `json:"status"`
}

func (h *WorkspaceHandler) Overview(c *gin.Context) {
	now := time.Now()
	recentStart := now.AddDate(0, 0, -7)

	var recentReviews []model.ContentReviewTask
	if err := database.DB.Where("created_at >= ?", recentStart).Find(&recentReviews).Error; err != nil {
		basehandler.ServerError(c, "读取近 7 天复盘数据失败")
		return
	}

	merchantRates, err := loadMerchantCommissionRates()
	if err != nil {
		basehandler.ServerError(c, "读取商家分成比例失败")
		return
	}

	recentWriteOffAmount, estimatedCommission := calculateRecentAmounts(recentReviews, merchantRates)

	merchants, err := buildWorkspaceMerchants(recentStart)
	if err != nil {
		basehandler.ServerError(c, "读取重点商家失败")
		return
	}
	shootingTasks, err := buildWorkspaceShootingTasks()
	if err != nil {
		basehandler.ServerError(c, "读取拍摄任务失败")
		return
	}
	topics, err := buildWorkspaceTopics()
	if err != nil {
		basehandler.ServerError(c, "读取选题池失败")
		return
	}
	reviews, err := buildWorkspaceReviews()
	if err != nil {
		basehandler.ServerError(c, "读取复盘数据失败")
		return
	}

	basehandler.OK(c, workspaceOverviewResponse{
		Metrics:       buildWorkspaceMetrics(recentWriteOffAmount, estimatedCommission),
		Merchants:     merchants,
		ShootingTasks: shootingTasks,
		Topics:        topics,
		Reviews:       reviews,
	})
}

func buildWorkspaceMetrics(recentWriteOffAmount float64, estimatedCommission float64) []workspaceMetricRow {
	merchantTotal := countAll(&model.Merchant{})
	activeMerchantTotal := countWhere(&model.Merchant{}, "status = ?", model.MerchantStatusEnabled)
	shootingPendingTotal := countWhere(&model.ShootingTask{}, "status IN ?", []string{
		model.ShootingTaskStatusPending,
		model.ShootingTaskStatusShooting,
		model.ShootingTaskStatusShot,
	})
	pendingScheduleTotal := countWhere(&model.PublishSchedule{}, "status = ?", model.PublishScheduleStatusPending)
	enabledPackageTotal := countWhere(&model.MerchantPackage{}, "status = ?", model.MerchantPackageStatusEnabled)
	reviewTotal := countAll(&model.ContentReviewTask{})

	return []workspaceMetricRow{
		{
			Key:   "merchants",
			Label: "在运营商家",
			Value: fmt.Sprintf("%d", activeMerchantTotal),
			Hint:  fmt.Sprintf("总商家 %d 家，启用套餐 %d 个", merchantTotal, enabledPackageTotal),
			Path:  "/users/list",
		},
		{
			Key:   "shooting",
			Label: "待拍/待剪任务",
			Value: fmt.Sprintf("%d", shootingPendingTotal),
			Hint:  fmt.Sprintf("待发布视频 %d 条", pendingScheduleTotal),
			Path:  "/logs/user",
		},
		{
			Key:   "writeoff",
			Label: "近 7 天核销额",
			Value: formatMoney(recentWriteOffAmount),
			Hint:  fmt.Sprintf("已复盘 %d 条视频", reviewTotal),
			Path:  "/logs/admin",
		},
		{
			Key:   "commission",
			Label: "预估分成",
			Value: formatMoney(estimatedCommission),
			Hint:  "按各商家分成比例估算",
			Path:  "/plans/orders",
		},
	}
}

func buildWorkspaceMerchants(recentStart time.Time) ([]workspaceMerchantRow, error) {
	var list []model.Merchant
	if err := database.DB.Order("updated_at DESC, id DESC").Limit(6).Find(&list).Error; err != nil {
		return nil, err
	}
	result := make([]workspaceMerchantRow, 0, len(list))
	for _, item := range list {
		var recentWriteOffAmount float64
		_ = database.DB.Model(&model.ContentReviewTask{}).
			Where("merchant_id = ? AND created_at >= ?", item.ID, recentStart).
			Select("COALESCE(SUM(write_off_amount), 0)").
			Scan(&recentWriteOffAmount).Error
		result = append(result, workspaceMerchantRow{
			ID:                   item.ID,
			Name:                 item.Name,
			Industry:             item.Industry,
			City:                 item.City,
			Stage:                item.Stage,
			CommissionRate:       item.CommissionRate,
			RecentWriteOffAmount: recentWriteOffAmount,
			EstimatedCommission:  roundMoney(recentWriteOffAmount * item.CommissionRate / 100),
			NextAction:           buildMerchantNextAction(item),
		})
	}
	return result, nil
}

func buildMerchantNextAction(merchant model.Merchant) string {
	metrics := buildMerchantWorkspaceMetrics(merchant.ID)
	requirements := buildMerchantRequirements(merchant, metrics)
	for _, item := range requirements {
		if !item.Done {
			return item.Title
		}
	}
	if metrics.AcceptedTopicCount == 0 {
		return "生成并采用选题"
	}
	if metrics.ConfirmedScriptCount == 0 {
		return "生成并确认文案"
	}
	if metrics.ConfirmedStoryboardCount == 0 {
		return "生成并确认分镜"
	}
	if metrics.ShootingTaskCount == 0 {
		return "创建拍摄任务"
	}
	if metrics.PublishedScheduleCount == 0 {
		return "安排发布排期"
	}
	return "复盘数据并迭代选题"
}

func buildWorkspaceShootingTasks() ([]workspaceShootingTaskRow, error) {
	var list []model.ShootingTask
	query := database.DB.Where("status IN ?", []string{
		model.ShootingTaskStatusPending,
		model.ShootingTaskStatusShooting,
		model.ShootingTaskStatusShot,
		model.ShootingTaskStatusEdited,
	})
	if err := query.Order("deadline IS NULL, deadline ASC, id DESC").Limit(6).Find(&list).Error; err != nil {
		return nil, err
	}
	result := make([]workspaceShootingTaskRow, 0, len(list))
	for _, item := range list {
		result = append(result, workspaceShootingTaskRow{
			ID:           item.ID,
			MerchantID:   item.MerchantID,
			MerchantName: item.MerchantName,
			TaskTitle:    item.TaskTitle,
			Assignee:     item.Assignee,
			Deadline:     item.Deadline,
			Status:       item.Status,
			ShotCount:    item.ShotCount,
		})
	}
	return result, nil
}

func buildWorkspaceTopics() ([]workspaceTopicRow, error) {
	var list []model.ContentTopic
	if err := database.DB.
		Where("status <> ?", model.ContentTopicStatusDisabled).
		Order("id DESC").
		Limit(6).
		Find(&list).Error; err != nil {
		return nil, err
	}
	result := make([]workspaceTopicRow, 0, len(list))
	for _, item := range list {
		result = append(result, workspaceTopicRow{
			ID:            item.ID,
			MerchantID:    item.MerchantID,
			MerchantName:  item.MerchantName,
			Title:         item.Title,
			Hook:          item.Hook,
			PublishWindow: item.PublishWindow,
			Status:        item.Status,
		})
	}
	return result, nil
}

func buildWorkspaceReviews() ([]workspaceReviewRow, error) {
	var list []model.ContentReviewTask
	if err := database.DB.Order("id DESC").Limit(6).Find(&list).Error; err != nil {
		return nil, err
	}
	result := make([]workspaceReviewRow, 0, len(list))
	for _, item := range list {
		result = append(result, workspaceReviewRow{
			ID:             item.ID,
			MerchantID:     item.MerchantID,
			MerchantName:   item.MerchantName,
			VideoTitle:     item.VideoTitle,
			PlayCount:      item.PlayCount,
			DealCount:      item.DealCount,
			WriteOffAmount: item.WriteOffAmount,
			Status:         item.Status,
		})
	}
	return result, nil
}

func loadMerchantCommissionRates() (map[uint64]float64, error) {
	var merchants []model.Merchant
	if err := database.DB.Select("id", "commission_rate").Find(&merchants).Error; err != nil {
		return nil, err
	}
	result := make(map[uint64]float64, len(merchants))
	for _, merchant := range merchants {
		result[merchant.ID] = merchant.CommissionRate
	}
	return result, nil
}

func calculateRecentAmounts(reviews []model.ContentReviewTask, merchantRates map[uint64]float64) (float64, float64) {
	var writeOffAmount float64
	var commission float64
	for _, item := range reviews {
		writeOffAmount += item.WriteOffAmount
		rate := merchantRates[item.MerchantID]
		if rate <= 0 {
			rate = 10
		}
		commission += item.WriteOffAmount * rate / 100
	}
	return roundMoney(writeOffAmount), roundMoney(commission)
}

func countAll(table any) int64 {
	var total int64
	_ = database.DB.Model(table).Count(&total).Error
	return total
}

func countWhere(table any, condition string, args ...any) int64 {
	var total int64
	_ = database.DB.Model(table).Where(condition, args...).Count(&total).Error
	return total
}

func formatMoney(value float64) string {
	return fmt.Sprintf("¥%.2f", roundMoney(value))
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}
