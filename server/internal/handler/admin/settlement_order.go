package admin

import (
	"errors"
	"strings"
	"time"

	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SettlementOrderHandler struct{}

func NewSettlementOrderHandler() *SettlementOrderHandler {
	return &SettlementOrderHandler{}
}

type settlementOrderPayload struct {
	MerchantID     uint64  `json:"merchantId" binding:"required"`
	ScheduleID     uint64  `json:"scheduleId"`
	VideoTitle     string  `json:"videoTitle"`
	SourceVideo    string  `json:"sourceVideo"`
	OrderWindow    string  `json:"orderWindow"`
	PeriodStart    string  `json:"periodStart"`
	PeriodEnd      string  `json:"periodEnd"`
	RedeemedAmount float64 `json:"redeemedAmount"`
	CommissionRate float64 `json:"commissionRate"`
	Status         string  `json:"status"`
	Remark         string  `json:"remark"`
}

type settlementOrderStatusPayload struct {
	Status string `json:"status" binding:"required"`
	Remark string `json:"remark"`
}

type settlementGeneratePayload struct {
	MerchantID uint64 `json:"merchantId"`
}

type settlementGenerateResponse struct {
	Created int                     `json:"created"`
	List    []model.SettlementOrder `json:"list"`
}

func (h *SettlementOrderHandler) List(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	merchantID := parseUint64Query(c.Query("merchantId"))
	status := strings.TrimSpace(c.Query("status"))

	query := database.DB.Model(&model.SettlementOrder{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"merchant_name LIKE ? OR video_title LIKE ? OR source_video LIKE ? OR order_window LIKE ? OR remark LIKE ?",
			like,
			like,
			like,
			like,
			like,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		basehandler.ServerError(c, "读取分成订单数量失败")
		return
	}

	var list []model.SettlementOrder
	if err := query.Order("id DESC").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取分成订单失败")
		return
	}
	basehandler.OK(c, listResponse[model.SettlementOrder]{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func (h *SettlementOrderHandler) Create(c *gin.Context) {
	var payload settlementOrderPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并填写订单金额")
		return
	}

	item, err := buildSettlementOrderFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	item.CreatedBy = c.GetUint64("admin_user_id")
	item.UpdatedBy = item.CreatedBy
	if err := database.DB.Create(&item).Error; err != nil {
		basehandler.ServerError(c, "创建分成订单失败")
		return
	}
	basehandler.OK(c, item)
}

func (h *SettlementOrderHandler) Update(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	var payload settlementOrderPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并填写订单金额")
		return
	}

	updates, err := buildSettlementOrderFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}

	var item model.SettlementOrder
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "分成订单不存在")
			return
		}
		basehandler.ServerError(c, "读取分成订单失败")
		return
	}

	if err := database.DB.Model(&item).Updates(settlementOrderUpdateMap(updates, c.GetUint64("admin_user_id"))).Error; err != nil {
		basehandler.ServerError(c, "更新分成订单失败")
		return
	}
	_ = database.DB.First(&item, id).Error
	basehandler.OK(c, item)
}

func (h *SettlementOrderHandler) UpdateStatus(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	var payload settlementOrderStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil || !isValidSettlementStatus(strings.TrimSpace(payload.Status)) {
		basehandler.BadRequest(c, "分成订单状态不正确")
		return
	}

	var item model.SettlementOrder
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "分成订单不存在")
			return
		}
		basehandler.ServerError(c, "读取分成订单失败")
		return
	}

	updates := map[string]any{
		"status":     strings.TrimSpace(payload.Status),
		"updated_by": c.GetUint64("admin_user_id"),
	}
	if strings.TrimSpace(payload.Remark) != "" {
		updates["remark"] = strings.TrimSpace(payload.Remark)
	}
	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		basehandler.ServerError(c, "更新分成订单状态失败")
		return
	}
	_ = database.DB.First(&item, id).Error
	basehandler.OK(c, item)
}

func (h *SettlementOrderHandler) Delete(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	if err := database.DB.Delete(&model.SettlementOrder{}, id).Error; err != nil {
		basehandler.ServerError(c, "删除分成订单失败")
		return
	}
	basehandler.OK(c, gin.H{"id": id})
}

func (h *SettlementOrderHandler) Generate(c *gin.Context) {
	var payload settlementGeneratePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "生成参数不正确")
		return
	}

	query := database.DB.Model(&model.PublishSchedule{}).Where("status IN ?", []string{
		model.PublishScheduleStatusPublished,
		model.PublishScheduleStatusReviewed,
	})
	if payload.MerchantID > 0 {
		query = query.Where("merchant_id = ?", payload.MerchantID)
	}

	var schedules []model.PublishSchedule
	if err := query.Order("id DESC").Limit(100).Find(&schedules).Error; err != nil {
		basehandler.ServerError(c, "读取发布排期失败")
		return
	}

	created := make([]model.SettlementOrder, 0)
	operatorID := c.GetUint64("admin_user_id")
	for _, schedule := range schedules {
		var existing int64
		_ = database.DB.Model(&model.SettlementOrder{}).Where("schedule_id = ?", schedule.ID).Count(&existing).Error
		if existing > 0 {
			continue
		}
		item := buildSettlementOrderFromSchedule(schedule, operatorID)
		if err := database.DB.Create(&item).Error; err != nil {
			basehandler.ServerError(c, "生成分成订单失败")
			return
		}
		created = append(created, item)
	}
	basehandler.OK(c, settlementGenerateResponse{Created: len(created), List: created})
}

func buildSettlementOrderFromPayload(payload settlementOrderPayload) (model.SettlementOrder, error) {
	if payload.MerchantID == 0 {
		return model.SettlementOrder{}, errors.New("请选择商家")
	}
	if payload.RedeemedAmount < 0 {
		return model.SettlementOrder{}, errors.New("核销额不能小于 0")
	}

	var merchant model.Merchant
	if err := database.DB.First(&merchant, payload.MerchantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.SettlementOrder{}, errors.New("商家不存在")
		}
		return model.SettlementOrder{}, errors.New("读取商家失败")
	}

	item := model.SettlementOrder{
		MerchantID:     merchant.ID,
		MerchantName:   merchant.Name,
		ScheduleID:     payload.ScheduleID,
		VideoTitle:     strings.TrimSpace(payload.VideoTitle),
		SourceVideo:    strings.TrimSpace(payload.SourceVideo),
		OrderWindow:    strings.TrimSpace(payload.OrderWindow),
		RedeemedAmount: payload.RedeemedAmount,
		CommissionRate: payload.CommissionRate,
		Status:         strings.TrimSpace(payload.Status),
		Remark:         strings.TrimSpace(payload.Remark),
	}
	if item.CommissionRate <= 0 {
		item.CommissionRate = merchant.CommissionRate
	}
	if item.CommissionRate <= 0 {
		item.CommissionRate = 10
	}
	if item.CommissionRate > 100 {
		return model.SettlementOrder{}, errors.New("分成比例不能超过 100")
	}
	if item.Status == "" {
		item.Status = model.SettlementOrderStatusPending
	}
	if !isValidSettlementStatus(item.Status) {
		return model.SettlementOrder{}, errors.New("分成订单状态不正确")
	}
	periodStart, err := parseOptionalDate(payload.PeriodStart)
	if err != nil {
		return model.SettlementOrder{}, errors.New("统计开始日期格式应为 2006-01-02")
	}
	periodEnd, err := parseOptionalDate(payload.PeriodEnd)
	if err != nil {
		return model.SettlementOrder{}, errors.New("统计结束日期格式应为 2006-01-02")
	}
	if periodStart != nil && periodEnd != nil && periodEnd.Before(*periodStart) {
		return model.SettlementOrder{}, errors.New("统计结束日期不能早于开始日期")
	}
	item.PeriodStart = periodStart
	item.PeriodEnd = periodEnd

	if item.ScheduleID > 0 {
		var schedule model.PublishSchedule
		if err := database.DB.First(&schedule, item.ScheduleID).Error; err != nil {
			return model.SettlementOrder{}, errors.New("发布排期不存在")
		}
		if schedule.MerchantID != item.MerchantID {
			return model.SettlementOrder{}, errors.New("发布排期不属于当前商家")
		}
		item.VideoTitle = firstText(item.VideoTitle, schedule.VideoTitle)
		item.SourceVideo = firstText(item.SourceVideo, schedule.VideoTitle)
	}
	if item.SourceVideo == "" {
		item.SourceVideo = item.VideoTitle
	}
	if item.OrderWindow == "" {
		item.OrderWindow = buildOrderWindow(item.PeriodStart, item.PeriodEnd)
	}
	item.Commission = roundMoney(item.RedeemedAmount * item.CommissionRate / 100)
	return item, nil
}

func buildSettlementOrderFromSchedule(schedule model.PublishSchedule, operatorID uint64) model.SettlementOrder {
	var merchant model.Merchant
	_ = database.DB.First(&merchant, schedule.MerchantID).Error

	var review model.ContentReviewTask
	_ = database.DB.Where("schedule_id = ?", schedule.ID).Order("id DESC").First(&review).Error

	periodStart := review.PeriodStart
	periodEnd := review.PeriodEnd
	if periodStart == nil && schedule.PublishTime != nil {
		start := time.Date(schedule.PublishTime.Year(), schedule.PublishTime.Month(), schedule.PublishTime.Day(), 0, 0, 0, 0, schedule.PublishTime.Location())
		periodStart = &start
	}
	if periodEnd == nil && periodStart != nil {
		end := periodStart.AddDate(0, 0, 6)
		periodEnd = &end
	}

	commissionRate := merchant.CommissionRate
	if commissionRate <= 0 {
		commissionRate = 10
	}
	redeemedAmount := review.WriteOffAmount
	commission := roundMoney(redeemedAmount * commissionRate / 100)

	return model.SettlementOrder{
		MerchantID:     schedule.MerchantID,
		MerchantName:   schedule.MerchantName,
		ScheduleID:     schedule.ID,
		VideoTitle:     schedule.VideoTitle,
		SourceVideo:    schedule.VideoTitle,
		OrderWindow:    buildOrderWindow(periodStart, periodEnd),
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		RedeemedAmount: redeemedAmount,
		CommissionRate: commissionRate,
		Commission:     commission,
		Status:         model.SettlementOrderStatusPending,
		Remark:         "由已发布/已复盘视频自动生成，核销额可在编辑中校正",
		CreatedBy:      operatorID,
		UpdatedBy:      operatorID,
	}
}

func settlementOrderUpdateMap(item model.SettlementOrder, operatorID uint64) map[string]any {
	return map[string]any{
		"merchant_id":     item.MerchantID,
		"merchant_name":   item.MerchantName,
		"schedule_id":     item.ScheduleID,
		"video_title":     item.VideoTitle,
		"source_video":    item.SourceVideo,
		"order_window":    item.OrderWindow,
		"period_start":    item.PeriodStart,
		"period_end":      item.PeriodEnd,
		"redeemed_amount": item.RedeemedAmount,
		"commission_rate": item.CommissionRate,
		"commission":      item.Commission,
		"status":          item.Status,
		"remark":          item.Remark,
		"updated_by":      operatorID,
	}
}

func isValidSettlementStatus(status string) bool {
	return status == model.SettlementOrderStatusPending ||
		status == model.SettlementOrderStatusConfirmed ||
		status == model.SettlementOrderStatusPaid
}

func buildOrderWindow(start *time.Time, end *time.Time) string {
	if start == nil && end == nil {
		return ""
	}
	if start != nil && end != nil {
		return start.Format("2006-01-02") + " 至 " + end.Format("2006-01-02")
	}
	if start != nil {
		return start.Format("2006-01-02") + " 起"
	}
	return "截至 " + end.Format("2006-01-02")
}
