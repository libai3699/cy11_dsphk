package admin

import (
	"errors"
	"strconv"
	"strings"

	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MerchantPackageHandler struct{}

func NewMerchantPackageHandler() *MerchantPackageHandler {
	return &MerchantPackageHandler{}
}

type merchantPackagePayload struct {
	MerchantID     uint64  `json:"merchantId" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	OriginalPrice  float64 `json:"originalPrice"`
	SellingPrice   float64 `json:"sellingPrice" binding:"required"`
	CostPrice      float64 `json:"costPrice"`
	CommissionRate float64 `json:"commissionRate"`
	TrafficLabel   string  `json:"trafficLabel"`
	ProfitGuard    string  `json:"profitGuard"`
	Status         *int    `json:"status"`
	Remark         string  `json:"remark"`
}

type merchantPackageListResponse struct {
	List  []merchantPackageResponse `json:"list"`
	Total int64                     `json:"total"`
	Page  int                       `json:"page"`
	Size  int                       `json:"size"`
}

type merchantPackageResponse struct {
	model.MerchantPackage
	GrossProfit         float64 `json:"grossProfit"`
	MarginRate          float64 `json:"marginRate"`
	EstimatedCommission float64 `json:"estimatedCommission"`
	NetAfterCommission  float64 `json:"netAfterCommission"`
}

func (h *MerchantPackageHandler) List(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	merchantID := parseUint64Query(c.Query("merchantId"))

	query := database.DB.Model(&model.MerchantPackage{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"merchant_name LIKE ? OR name LIKE ? OR traffic_label LIKE ? OR profit_guard LIKE ? OR remark LIKE ?",
			like,
			like,
			like,
			like,
			like,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		basehandler.ServerError(c, "读取套餐数量失败")
		return
	}

	var list []model.MerchantPackage
	if err := query.
		Order("id DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取套餐列表失败")
		return
	}

	basehandler.OK(c, merchantPackageListResponse{
		List:  buildMerchantPackageResponses(list),
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func (h *MerchantPackageHandler) Get(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var item model.MerchantPackage
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "套餐不存在")
			return
		}
		basehandler.ServerError(c, "读取套餐失败")
		return
	}

	basehandler.OK(c, buildMerchantPackageResponse(item))
}

func (h *MerchantPackageHandler) Create(c *gin.Context) {
	var payload merchantPackagePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并输入套餐名称和售价")
		return
	}

	item, err := buildMerchantPackageFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	item.CreatedBy = c.GetUint64("admin_user_id")
	item.UpdatedBy = item.CreatedBy

	if err := database.DB.Create(&item).Error; err != nil {
		basehandler.ServerError(c, "创建套餐失败")
		return
	}

	basehandler.OK(c, buildMerchantPackageResponse(item))
}

func (h *MerchantPackageHandler) Update(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var payload merchantPackagePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并输入套餐名称和售价")
		return
	}

	updates, err := buildMerchantPackageFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	updates.UpdatedBy = c.GetUint64("admin_user_id")

	var item model.MerchantPackage
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "套餐不存在")
			return
		}
		basehandler.ServerError(c, "读取套餐失败")
		return
	}

	if err := database.DB.Model(&item).Updates(merchantPackageUpdateMap(updates)).Error; err != nil {
		basehandler.ServerError(c, "更新套餐失败")
		return
	}
	if err := database.DB.First(&item, id).Error; err != nil {
		basehandler.ServerError(c, "读取套餐失败")
		return
	}

	basehandler.OK(c, buildMerchantPackageResponse(item))
}

func (h *MerchantPackageHandler) Delete(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	if err := database.DB.Delete(&model.MerchantPackage{}, id).Error; err != nil {
		basehandler.ServerError(c, "删除套餐失败")
		return
	}
	basehandler.OK(c, gin.H{"id": id})
}

func buildMerchantPackageFromPayload(payload merchantPackagePayload) (model.MerchantPackage, error) {
	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return model.MerchantPackage{}, errors.New("请输入套餐名称")
	}
	if payload.MerchantID == 0 {
		return model.MerchantPackage{}, errors.New("请选择商家")
	}
	if payload.SellingPrice <= 0 {
		return model.MerchantPackage{}, errors.New("售价必须大于 0")
	}
	if payload.CostPrice < 0 || payload.OriginalPrice < 0 {
		return model.MerchantPackage{}, errors.New("价格不能小于 0")
	}
	if payload.CostPrice > payload.SellingPrice {
		return model.MerchantPackage{}, errors.New("成本不能高于售价")
	}

	var merchant model.Merchant
	if err := database.DB.First(&merchant, payload.MerchantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.MerchantPackage{}, errors.New("商家不存在")
		}
		return model.MerchantPackage{}, errors.New("读取商家失败")
	}

	commissionRate := payload.CommissionRate
	if commissionRate <= 0 {
		commissionRate = merchant.CommissionRate
	}
	if commissionRate <= 0 {
		commissionRate = 10
	}
	if commissionRate > 100 {
		return model.MerchantPackage{}, errors.New("提点比例不能超过 100")
	}

	status := model.MerchantPackageStatusEnabled
	if payload.Status != nil {
		status = *payload.Status
	}
	if status != model.MerchantPackageStatusEnabled && status != model.MerchantPackageStatusDisabled {
		return model.MerchantPackage{}, errors.New("套餐状态不正确")
	}

	profitGuard := strings.TrimSpace(payload.ProfitGuard)
	if profitGuard == "" {
		profitGuard = buildProfitGuard(payload.SellingPrice, payload.CostPrice, commissionRate)
	}

	return model.MerchantPackage{
		MerchantID:     merchant.ID,
		MerchantName:   merchant.Name,
		Name:           name,
		OriginalPrice:  payload.OriginalPrice,
		SellingPrice:   payload.SellingPrice,
		CostPrice:      payload.CostPrice,
		CommissionRate: commissionRate,
		TrafficLabel:   strings.TrimSpace(payload.TrafficLabel),
		ProfitGuard:    profitGuard,
		Status:         status,
		Remark:         strings.TrimSpace(payload.Remark),
	}, nil
}

func merchantPackageUpdateMap(item model.MerchantPackage) map[string]any {
	return map[string]any{
		"merchant_id":     item.MerchantID,
		"merchant_name":   item.MerchantName,
		"name":            item.Name,
		"original_price":  item.OriginalPrice,
		"selling_price":   item.SellingPrice,
		"cost_price":      item.CostPrice,
		"commission_rate": item.CommissionRate,
		"traffic_label":   item.TrafficLabel,
		"profit_guard":    item.ProfitGuard,
		"status":          item.Status,
		"remark":          item.Remark,
		"updated_by":      item.UpdatedBy,
	}
}

func buildMerchantPackageResponses(list []model.MerchantPackage) []merchantPackageResponse {
	result := make([]merchantPackageResponse, 0, len(list))
	for _, item := range list {
		result = append(result, buildMerchantPackageResponse(item))
	}
	return result
}

func buildMerchantPackageResponse(item model.MerchantPackage) merchantPackageResponse {
	grossProfit := item.SellingPrice - item.CostPrice
	marginRate := 0.0
	if item.SellingPrice > 0 {
		marginRate = grossProfit / item.SellingPrice * 100
	}
	estimatedCommission := item.SellingPrice * item.CommissionRate / 100
	return merchantPackageResponse{
		MerchantPackage:     item,
		GrossProfit:         grossProfit,
		MarginRate:          marginRate,
		EstimatedCommission: estimatedCommission,
		NetAfterCommission:  grossProfit - estimatedCommission,
	}
}

func buildProfitGuard(sellingPrice float64, costPrice float64, commissionRate float64) string {
	grossProfit := sellingPrice - costPrice
	estimatedCommission := sellingPrice * commissionRate / 100
	net := grossProfit - estimatedCommission
	if sellingPrice <= 0 {
		return ""
	}
	netRate := net / sellingPrice * 100
	switch {
	case net <= 0:
		return "扣除提点后亏损，不能作为主推套餐"
	case netRate < 15:
		return "扣除提点后利润偏低，只适合限量引流"
	case netRate < 30:
		return "利润空间一般，建议控制赠品和投放强度"
	default:
		return "利润空间健康，可以作为内容主推"
	}
}

func parseUint64Query(raw string) uint64 {
	if strings.TrimSpace(raw) == "" {
		return 0
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0
	}
	return value
}
