package admin

import (
	"errors"
	"strings"

	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MerchantAccountAuthHandler struct{}

func NewMerchantAccountAuthHandler() *MerchantAccountAuthHandler {
	return &MerchantAccountAuthHandler{}
}

type merchantAccountAuthPayload struct {
	MerchantID  uint64 `json:"merchantId" binding:"required"`
	Platform    string `json:"platform"`
	AuthMethod  string `json:"authMethod"`
	AccountName string `json:"accountName"`
	AccountUID  string `json:"accountUid"`
	AuthStatus  string `json:"authStatus"`
	LastLoginAt string `json:"lastLoginAt"`
	Remark      string `json:"remark"`
}

type merchantAccountAuthListResponse struct {
	List  []model.MerchantAccountAuth `json:"list"`
	Total int64                       `json:"total"`
	Page  int                         `json:"page"`
	Size  int                         `json:"size"`
}

func (h *MerchantAccountAuthHandler) List(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	merchantID := parseUint64Query(c.Query("merchantId"))
	status := strings.TrimSpace(c.Query("status"))

	query := database.DB.Model(&model.MerchantAccountAuth{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if status != "" {
		query = query.Where("auth_status = ?", status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"merchant_name LIKE ? OR platform LIKE ? OR auth_method LIKE ? OR account_name LIKE ? OR account_uid LIKE ? OR auth_status LIKE ? OR remark LIKE ?",
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
		basehandler.ServerError(c, "读取授权数量失败")
		return
	}

	var list []model.MerchantAccountAuth
	if err := query.
		Order("id DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取授权列表失败")
		return
	}

	basehandler.OK(c, merchantAccountAuthListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func (h *MerchantAccountAuthHandler) Get(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var item model.MerchantAccountAuth
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "授权记录不存在")
			return
		}
		basehandler.ServerError(c, "读取授权记录失败")
		return
	}

	basehandler.OK(c, item)
}

func (h *MerchantAccountAuthHandler) Create(c *gin.Context) {
	var payload merchantAccountAuthPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家")
		return
	}

	item, err := buildMerchantAccountAuthFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	item.CreatedBy = c.GetUint64("admin_user_id")
	item.UpdatedBy = item.CreatedBy

	if err := database.DB.Create(&item).Error; err != nil {
		basehandler.ServerError(c, "创建授权记录失败")
		return
	}
	if err := syncMerchantAccountFields(item); err != nil {
		basehandler.ServerError(c, "同步商家账号信息失败")
		return
	}

	basehandler.OK(c, item)
}

func (h *MerchantAccountAuthHandler) Update(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var payload merchantAccountAuthPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家")
		return
	}

	updates, err := buildMerchantAccountAuthFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	updates.UpdatedBy = c.GetUint64("admin_user_id")

	var item model.MerchantAccountAuth
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "授权记录不存在")
			return
		}
		basehandler.ServerError(c, "读取授权记录失败")
		return
	}

	if err := database.DB.Model(&item).Updates(merchantAccountAuthUpdateMap(updates)).Error; err != nil {
		basehandler.ServerError(c, "更新授权记录失败")
		return
	}
	if err := database.DB.First(&item, id).Error; err != nil {
		basehandler.ServerError(c, "读取授权记录失败")
		return
	}
	if err := syncMerchantAccountFields(item); err != nil {
		basehandler.ServerError(c, "同步商家账号信息失败")
		return
	}

	basehandler.OK(c, item)
}

func (h *MerchantAccountAuthHandler) Delete(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	if err := database.DB.Delete(&model.MerchantAccountAuth{}, id).Error; err != nil {
		basehandler.ServerError(c, "删除授权记录失败")
		return
	}
	basehandler.OK(c, gin.H{"id": id})
}

func buildMerchantAccountAuthFromPayload(payload merchantAccountAuthPayload) (model.MerchantAccountAuth, error) {
	if payload.MerchantID == 0 {
		return model.MerchantAccountAuth{}, errors.New("请选择商家")
	}

	var merchant model.Merchant
	if err := database.DB.First(&merchant, payload.MerchantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.MerchantAccountAuth{}, errors.New("商家不存在")
		}
		return model.MerchantAccountAuth{}, errors.New("读取商家失败")
	}

	platform := strings.TrimSpace(payload.Platform)
	if platform == "" {
		platform = "抖音号"
	}
	authMethod := strings.TrimSpace(payload.AuthMethod)
	if authMethod == "" {
		authMethod = "验证码代登"
	}
	authStatus := strings.TrimSpace(payload.AuthStatus)
	if authStatus == "" {
		authStatus = model.MerchantAccountAuthStatusPending
	}
	if !isValidAccountAuthStatus(authStatus) {
		return model.MerchantAccountAuth{}, errors.New("授权状态不正确")
	}

	return model.MerchantAccountAuth{
		MerchantID:   merchant.ID,
		MerchantName: merchant.Name,
		Platform:     platform,
		AuthMethod:   authMethod,
		AccountName:  strings.TrimSpace(payload.AccountName),
		AccountUID:   strings.TrimSpace(payload.AccountUID),
		AuthStatus:   authStatus,
		LastLoginAt:  strings.TrimSpace(payload.LastLoginAt),
		Remark:       strings.TrimSpace(payload.Remark),
	}, nil
}

func merchantAccountAuthUpdateMap(item model.MerchantAccountAuth) map[string]any {
	return map[string]any{
		"merchant_id":   item.MerchantID,
		"merchant_name": item.MerchantName,
		"platform":      item.Platform,
		"auth_method":   item.AuthMethod,
		"account_name":  item.AccountName,
		"account_uid":   item.AccountUID,
		"auth_status":   item.AuthStatus,
		"last_login_at": item.LastLoginAt,
		"remark":        item.Remark,
		"updated_by":    item.UpdatedBy,
	}
}

func isValidAccountAuthStatus(status string) bool {
	switch status {
	case model.MerchantAccountAuthStatusPending,
		model.MerchantAccountAuthStatusActive,
		model.MerchantAccountAuthStatusRenewal,
		model.MerchantAccountAuthStatusExpired:
		return true
	default:
		return false
	}
}

func syncMerchantAccountFields(item model.MerchantAccountAuth) error {
	accountValue := strings.TrimSpace(item.AccountName)
	if accountValue == "" {
		accountValue = strings.TrimSpace(item.AccountUID)
	}
	if accountValue == "" {
		return nil
	}

	updates := map[string]any{}
	switch item.Platform {
	case "抖音号":
		updates["douyin_account"] = accountValue
	case "抖音来客":
		updates["douyin_laike_account"] = accountValue
	default:
		return nil
	}
	if len(updates) == 0 {
		return nil
	}
	return database.DB.
		Model(&model.Merchant{}).
		Where("id = ?", item.MerchantID).
		Updates(updates).Error
}
