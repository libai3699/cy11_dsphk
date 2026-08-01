package admin

import (
	"strconv"
	"strings"

	"cy11dsphk/server/internal/config"
	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"
	"cy11dsphk/server/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	cfg config.Config
}

func NewAuthHandler(cfg config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type userInfoResponse struct {
	UserID   string   `json:"userId"`
	Username string   `json:"username"`
	RealName string   `json:"realName"`
	Avatar   string   `json:"avatar"`
	Desc     string   `json:"desc"`
	HomePath string   `json:"homePath"`
	Roles    []string `json:"roles"`
	Token    string   `json:"token"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		basehandler.BadRequest(c, "请输入账号和密码")
		return
	}

	username := strings.TrimSpace(req.Username)
	var user model.AdminUser
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		h.writeLoginLog(c, 0, username, "failed", "账号或密码错误")
		basehandler.Unauthorized(c, "账号或密码错误")
		return
	}
	if user.Status != model.AdminUserStatusEnabled {
		h.writeLoginLog(c, user.ID, user.Username, "failed", "账号已禁用")
		basehandler.Forbidden(c, "账号已禁用")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		h.writeLoginLog(c, user.ID, user.Username, "failed", "账号或密码错误")
		basehandler.Unauthorized(c, "账号或密码错误")
		return
	}

	token, err := h.issueToken(user)
	if err != nil {
		basehandler.ServerError(c, "登录失败")
		return
	}

	h.writeLoginLog(c, user.ID, user.Username, "success", "登录成功")
	basehandler.OK(c, tokenResponse{AccessToken: token})
}

func (h *AuthHandler) UserInfo(c *gin.Context) {
	userID := c.GetUint64("admin_user_id")
	var user model.AdminUser
	if err := database.DB.First(&user, userID).Error; err != nil {
		basehandler.Unauthorized(c, "登录状态已失效")
		return
	}

	roleCodes, err := service.GetUserRoleCodes(user.ID)
	if err != nil {
		basehandler.ServerError(c, "读取角色失败")
		return
	}

	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	basehandler.OK(c, userInfoResponse{
		UserID:   strconv.FormatUint(user.ID, 10),
		Username: user.Username,
		RealName: user.RealName,
		Avatar:   user.Avatar,
		Desc:     "本地商家抖音获客后台",
		HomePath: user.HomePath,
		Roles:    roleCodes,
		Token:    token,
	})
}

func (h *AuthHandler) AccessCodes(c *gin.Context) {
	userID := c.GetUint64("admin_user_id")
	codes, err := service.GetUserAccessCodes(userID)
	if err != nil {
		basehandler.ServerError(c, "读取权限失败")
		return
	}
	basehandler.OK(c, codes)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	claims, exists := c.Get("admin_claims")
	if !exists {
		basehandler.Unauthorized(c, "登录状态已失效")
		return
	}
	adminClaims, ok := claims.(*config.AdminClaims)
	if !ok {
		basehandler.Unauthorized(c, "登录状态已失效")
		return
	}
	token, err := h.cfg.GenerateAdminToken(adminClaims.UserID, adminClaims.Username, adminClaims.Roles)
	if err != nil {
		basehandler.ServerError(c, "刷新 token 失败")
		return
	}
	basehandler.OK(c, token)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	basehandler.OK(c, true)
}

func (h *AuthHandler) issueToken(user model.AdminUser) (string, error) {
	roleCodes, err := service.GetUserRoleCodes(user.ID)
	if err != nil {
		return "", err
	}
	return h.cfg.GenerateAdminToken(user.ID, user.Username, roleCodes)
}

func (h *AuthHandler) writeLoginLog(c *gin.Context, userID uint64, username string, result string, message string) {
	_ = database.DB.Create(&model.AdminLoginLog{
		UserID:    userID,
		Username:  username,
		IP:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Result:    result,
		Message:   message,
	}).Error
}
