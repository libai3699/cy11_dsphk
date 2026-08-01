package admin

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/accountdiagnosis"
	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountDiagnosisHandler struct{}

func NewAccountDiagnosisHandler() *AccountDiagnosisHandler {
	return &AccountDiagnosisHandler{}
}

type accountDiagnosisPayload struct {
	MerchantID       uint64   `json:"merchantId" binding:"required"`
	AccountAuthID    uint64   `json:"accountAuthId"`
	FollowerCount    float64  `json:"followerCount"`
	AvgPlayCount     float64  `json:"avgPlayCount"`
	BestVideoTitle   string   `json:"bestVideoTitle"`
	BestVideoPlay    float64  `json:"bestVideoPlay"`
	RecentVideoCount float64  `json:"recentVideoCount"`
	OwnerAppearance  string   `json:"ownerAppearance"`
	CurrentProblems  string   `json:"currentProblems"`
	TargetPackage    string   `json:"targetPackage"`
	OperatorGoal     string   `json:"operatorGoal"`
	RecentVideos     []string `json:"recentVideos"`
	Remark           string   `json:"remark"`
}

type accountDiagnosisListResponse struct {
	List  []accountDiagnosisResponse `json:"list"`
	Total int64                      `json:"total"`
	Page  int                        `json:"page"`
	Size  int                        `json:"size"`
}

type accountDiagnosisResponse struct {
	model.AccountDiagnosisTask
	Input  any `json:"input,omitempty"`
	Result any `json:"result,omitempty"`
}

func (h *AccountDiagnosisHandler) List(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("size"), 10)
	if size > 100 {
		size = 100
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	merchantID := parseUint64Query(c.Query("merchantId"))
	status := strings.TrimSpace(c.Query("status"))

	query := database.DB.Model(&model.AccountDiagnosisTask{})
	if merchantID > 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"merchant_name LIKE ? OR account_name LIKE ? OR status LIKE ? OR error_message LIKE ?",
			like,
			like,
			like,
			like,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		basehandler.ServerError(c, "读取诊断任务数量失败")
		return
	}

	var list []model.AccountDiagnosisTask
	if err := query.
		Order("id DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&list).Error; err != nil {
		basehandler.ServerError(c, "读取诊断任务失败")
		return
	}

	basehandler.OK(c, accountDiagnosisListResponse{
		List:  buildAccountDiagnosisResponses(list),
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func (h *AccountDiagnosisHandler) Get(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}

	var task model.AccountDiagnosisTask
	if err := database.DB.First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "诊断任务不存在")
			return
		}
		basehandler.ServerError(c, "读取诊断任务失败")
		return
	}

	basehandler.OK(c, buildAccountDiagnosisResponse(task))
}

func (h *AccountDiagnosisHandler) Create(c *gin.Context) {
	var payload accountDiagnosisPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家")
		return
	}

	task, agentInput, err := buildAccountDiagnosisTask(payload, c.GetUint64("admin_user_id"))
	if err != nil {
		basehandler.BadRequest(c, err.Error())
		return
	}

	if err := database.DB.Create(&task).Error; err != nil {
		basehandler.ServerError(c, "创建诊断任务失败")
		return
	}

	task.Status = model.AccountDiagnosisStatusRunning
	task.UpdatedBy = task.CreatedBy
	if err := database.DB.Model(&task).Updates(map[string]any{
		"status":     task.Status,
		"updated_by": task.UpdatedBy,
	}).Error; err != nil {
		basehandler.ServerError(c, "更新诊断任务失败")
		return
	}

	output, runErr := accountdiagnosis.Agent{}.Run(context.Background(), agentInput)
	if runErr != nil {
		task.Status = model.AccountDiagnosisStatusFailed
		task.ErrorMessage = runErr.Error()
		_ = database.DB.Model(&task).Updates(map[string]any{
			"status":        task.Status,
			"error_message": task.ErrorMessage,
			"updated_by":    task.UpdatedBy,
		}).Error
		basehandler.ServerError(c, "账号诊断执行失败")
		return
	}

	resultJSON, err := json.Marshal(output)
	if err != nil {
		basehandler.ServerError(c, "诊断结果序列化失败")
		return
	}
	task.Status = model.AccountDiagnosisStatusCompleted
	task.ResultJSON = string(resultJSON)
	task.ErrorMessage = ""
	if err := database.DB.Model(&task).Updates(map[string]any{
		"status":        task.Status,
		"result_json":   task.ResultJSON,
		"error_message": task.ErrorMessage,
		"updated_by":    task.UpdatedBy,
	}).Error; err != nil {
		basehandler.ServerError(c, "保存诊断结果失败")
		return
	}

	if err := database.DB.First(&task, task.ID).Error; err != nil {
		basehandler.ServerError(c, "读取诊断任务失败")
		return
	}
	basehandler.OK(c, buildAccountDiagnosisResponse(task))
}

func buildAccountDiagnosisTask(payload accountDiagnosisPayload, operatorID uint64) (model.AccountDiagnosisTask, accountdiagnosis.Input, error) {
	if payload.MerchantID == 0 {
		return model.AccountDiagnosisTask{}, accountdiagnosis.Input{}, errors.New("请选择商家")
	}

	var merchant model.Merchant
	if err := database.DB.First(&merchant, payload.MerchantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.AccountDiagnosisTask{}, accountdiagnosis.Input{}, errors.New("商家不存在")
		}
		return model.AccountDiagnosisTask{}, accountdiagnosis.Input{}, errors.New("读取商家失败")
	}

	var auth model.MerchantAccountAuth
	if payload.AccountAuthID > 0 {
		if err := database.DB.First(&auth, payload.AccountAuthID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return model.AccountDiagnosisTask{}, accountdiagnosis.Input{}, errors.New("授权记录不存在")
			}
			return model.AccountDiagnosisTask{}, accountdiagnosis.Input{}, errors.New("读取授权记录失败")
		}
	}

	var packages []model.MerchantPackage
	_ = database.DB.
		Where("merchant_id = ? AND status = ?", merchant.ID, model.MerchantPackageStatusEnabled).
		Order("id DESC").
		Find(&packages).Error

	accountName := auth.AccountName
	if accountName == "" {
		accountName = merchant.DouyinAccount
	}
	if accountName == "" {
		accountName = merchant.Name
	}

	snapshot := map[string]any{
		"merchant":       merchant,
		"accountAuth":    auth,
		"packages":       packages,
		"manualInput":    payload,
		"snapshotTime":   time.Now().Format("2006-01-02 15:04:05"),
		"snapshotSource": "admin_manual_input",
	}
	inputSnapshot, err := json.Marshal(snapshot)
	if err != nil {
		return model.AccountDiagnosisTask{}, accountdiagnosis.Input{}, errors.New("输入快照生成失败")
	}

	profileParts := []string{
		"商家：" + merchant.Name,
		"行业：" + merchant.Industry,
		"城市：" + merchant.City,
		"合作方式：" + merchant.CooperationType,
		"老板诉求：" + strings.TrimSpace(payload.OperatorGoal),
		"当前问题：" + strings.TrimSpace(payload.CurrentProblems),
		"主推套餐：" + strings.TrimSpace(payload.TargetPackage),
		"是否愿意出镜：" + strings.TrimSpace(payload.OwnerAppearance),
		"备注：" + strings.TrimSpace(payload.Remark),
	}

	agentInput := accountdiagnosis.Input{
		MerchantID:  merchant.ID,
		AccountName: accountName,
		Industry:    merchant.Industry,
		City:        merchant.City,
		Profile:     strings.Join(profileParts, "\n"),
		Metrics: map[string]float64{
			"followerCount":    payload.FollowerCount,
			"avgPlayCount":     payload.AvgPlayCount,
			"bestVideoPlay":    payload.BestVideoPlay,
			"recentVideoCount": payload.RecentVideoCount,
		},
		Options: agent.RunOptions{
			OperatorID: operatorID,
			DryRun:     true,
		},
	}
	for _, title := range payload.RecentVideos {
		title = strings.TrimSpace(title)
		if title == "" {
			continue
		}
		agentInput.RecentVideos = append(agentInput.RecentVideos, accountdiagnosis.VideoSnapshot{
			Title: title,
		})
	}
	if payload.BestVideoTitle != "" {
		agentInput.RecentVideos = append(agentInput.RecentVideos, accountdiagnosis.VideoSnapshot{
			Title:     strings.TrimSpace(payload.BestVideoTitle),
			PlayCount: int64(payload.BestVideoPlay),
		})
	}

	task := model.AccountDiagnosisTask{
		MerchantID:    merchant.ID,
		MerchantName:  merchant.Name,
		AccountAuthID: auth.ID,
		AccountName:   accountName,
		Status:        model.AccountDiagnosisStatusPending,
		InputSnapshot: string(inputSnapshot),
		CreatedBy:     operatorID,
		UpdatedBy:     operatorID,
	}
	return task, agentInput, nil
}

func buildAccountDiagnosisResponses(list []model.AccountDiagnosisTask) []accountDiagnosisResponse {
	result := make([]accountDiagnosisResponse, 0, len(list))
	for _, item := range list {
		result = append(result, buildAccountDiagnosisResponse(item))
	}
	return result
}

func buildAccountDiagnosisResponse(task model.AccountDiagnosisTask) accountDiagnosisResponse {
	return accountDiagnosisResponse{
		AccountDiagnosisTask: task,
		Input:                parseJSONField(task.InputSnapshot),
		Result:               parseJSONField(task.ResultJSON),
	}
}

func parseJSONField(raw string) any {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var value any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return nil
	}
	return value
}
