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
