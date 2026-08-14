package admin

import (
	"errors"
	"fmt"
	"strings"

	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FollowUpHandler struct{}

func NewFollowUpHandler() *FollowUpHandler {
	return &FollowUpHandler{}
}

type followUpPayload struct {
	MerchantID     uint64 `json:"merchantId" binding:"required"`
	Stage          string `json:"stage"`
	LatestTalk     string `json:"latestTalk" binding:"required"`
	Objection      string `json:"objection"`
	NextStep       string `json:"nextStep"`
	Owner          string `json:"owner"`
	FollowTime     string `json:"followTime"`
	NextFollowTime string `json:"nextFollowTime"`
}

type followUpSuggestionResponse struct {
	TalkScript string   `json:"talkScript"`
	Actions    []string `json:"actions"`
}

func (h *FollowUpHandler) List(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	merchantID := parseUint64Query(c.Query("merchantId"))
	stage := strings.TrimSpace(c.Query("stage"))

	query := database.DB.Model(&model.MerchantFollowUpLog{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if stage != "" {
		query = query.Where("stage = ?", stage)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"merchant_name LIKE ? OR stage LIKE ? OR latest_talk LIKE ? OR objection LIKE ? OR next_step LIKE ? OR owner LIKE ?",
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
		basehandler.ServerError(c, "读取跟进记录数量失败")
		return
	}

	var list []model.MerchantFollowUpLog
	if err := query.Order("id DESC").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取跟进记录失败")
		return
	}

	basehandler.OK(c, listResponse[model.MerchantFollowUpLog]{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func (h *FollowUpHandler) Create(c *gin.Context) {
	var payload followUpPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并填写最近沟通")
		return
	}

	item, err := buildFollowUpFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}
	item.CreatedBy = c.GetUint64("admin_user_id")
	item.UpdatedBy = item.CreatedBy

	if err := database.DB.Create(&item).Error; err != nil {
		basehandler.ServerError(c, "创建跟进记录失败")
		return
	}
	basehandler.OK(c, item)
}

func (h *FollowUpHandler) Update(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var payload followUpPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并填写最近沟通")
		return
	}

	updates, err := buildFollowUpFromPayload(payload)
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}

	var item model.MerchantFollowUpLog
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "跟进记录不存在")
			return
		}
		basehandler.ServerError(c, "读取跟进记录失败")
		return
	}

	if err := database.DB.Model(&item).Updates(followUpUpdateMap(updates, c.GetUint64("admin_user_id"))).Error; err != nil {
		basehandler.ServerError(c, "更新跟进记录失败")
		return
	}
	_ = database.DB.First(&item, id).Error
	basehandler.OK(c, item)
}

func (h *FollowUpHandler) Delete(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	if err := database.DB.Delete(&model.MerchantFollowUpLog{}, id).Error; err != nil {
		basehandler.ServerError(c, "删除跟进记录失败")
		return
	}
	basehandler.OK(c, gin.H{"id": id})
}

func (h *FollowUpHandler) Suggestion(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var item model.MerchantFollowUpLog
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "跟进记录不存在")
			return
		}
		basehandler.ServerError(c, "读取跟进记录失败")
		return
	}

	talkScript := buildFollowUpTalkScript(item)
	actions := buildFollowUpActions(item)
	basehandler.OK(c, followUpSuggestionResponse{TalkScript: talkScript, Actions: actions})
}

func buildFollowUpFromPayload(payload followUpPayload) (model.MerchantFollowUpLog, error) {
	if payload.MerchantID == 0 {
		return model.MerchantFollowUpLog{}, errors.New("请选择商家")
	}
	latestTalk := strings.TrimSpace(payload.LatestTalk)
	if latestTalk == "" {
		return model.MerchantFollowUpLog{}, errors.New("请输入最近沟通")
	}

	var merchant model.Merchant
	if err := database.DB.First(&merchant, payload.MerchantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.MerchantFollowUpLog{}, errors.New("商家不存在")
		}
		return model.MerchantFollowUpLog{}, errors.New("读取商家失败")
	}

	stage := strings.TrimSpace(payload.Stage)
	if stage == "" {
		stage = model.MerchantFollowUpStageContacting
	}
	if !isValidFollowUpStage(stage) {
		return model.MerchantFollowUpLog{}, errors.New("跟进阶段不正确")
	}

	followTime, err := parseOptionalDateTime(payload.FollowTime)
	if err != nil {
		return model.MerchantFollowUpLog{}, errors.New("沟通时间格式应为 2006-01-02 15:04:05")
	}
	nextFollowTime, err := parseOptionalDateTime(payload.NextFollowTime)
	if err != nil {
		return model.MerchantFollowUpLog{}, errors.New("下次跟进时间格式应为 2006-01-02 15:04:05")
	}

	return model.MerchantFollowUpLog{
		MerchantID:     merchant.ID,
		MerchantName:   merchant.Name,
		Stage:          stage,
		LatestTalk:     latestTalk,
		Objection:      strings.TrimSpace(payload.Objection),
		NextStep:       strings.TrimSpace(payload.NextStep),
		Owner:          strings.TrimSpace(payload.Owner),
		FollowTime:     followTime,
		NextFollowTime: nextFollowTime,
	}, nil
}

func followUpUpdateMap(item model.MerchantFollowUpLog, operatorID uint64) map[string]any {
	return map[string]any{
		"merchant_id":      item.MerchantID,
		"merchant_name":    item.MerchantName,
		"stage":            item.Stage,
		"latest_talk":      item.LatestTalk,
		"objection":        item.Objection,
		"next_step":        item.NextStep,
		"owner":            item.Owner,
		"follow_time":      item.FollowTime,
		"next_follow_time": item.NextFollowTime,
		"updated_by":       operatorID,
	}
}

func isValidFollowUpStage(stage string) bool {
	return stage == model.MerchantFollowUpStageContacting ||
		stage == model.MerchantFollowUpStageNegotiating ||
		stage == model.MerchantFollowUpStageContracted ||
		stage == model.MerchantFollowUpStagePaused ||
		stage == model.MerchantFollowUpStageLost
}

func buildFollowUpTalkScript(item model.MerchantFollowUpLog) string {
	objection := firstText(item.Objection, "老板还没有明确卡点")
	nextStep := firstText(item.NextStep, "确认是否愿意先跑一轮低风险测试")
	return fmt.Sprintf(
		"老板，我这边不是先收你钱做运营，而是先围绕门店真实团购转化跑内容。上次你提到的核心顾虑是「%s」，所以这次我们只确认三件事：第一，主推套餐利润边界；第二，账号和来客后台能不能看数据；第三，先用 3-5 条视频验证有没有自然成交。只要数据不成立，我们就调整打法，不让你盲目降价亏钱。今天我们可以先把「%s」定下来。",
		objection,
		nextStep,
	)
}

func buildFollowUpActions(item model.MerchantFollowUpLog) []string {
	actions := []string{
		"补齐商家套餐售价、成本和提点比例",
		"确认抖音账号和抖音来客后台是否可查看",
		"约定一轮最小测试：3-5 条视频、一个主推套餐、一个统计周期",
	}
	if item.Stage == model.MerchantFollowUpStageContracted {
		actions = []string{
			"进入账号授权和账号诊断",
			"补齐主推套餐并核算利润",
			"建立对标账号库，启动平台调研",
		}
	}
	return actions
}
