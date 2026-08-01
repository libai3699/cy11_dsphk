package middleware

import (
	"strings"

	"cy11dsphk/server/internal/config"
	basehandler "cy11dsphk/server/internal/handler"

	"github.com/gin-gonic/gin"
)

func AdminAuthRequired(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := bearerToken(c.GetHeader("Authorization"))
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			basehandler.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		claims, err := cfg.ParseAdminToken(token)
		if err != nil {
			basehandler.Unauthorized(c, "登录状态已失效")
			c.Abort()
			return
		}

		c.Set("admin_claims", claims)
		c.Set("admin_user_id", claims.UserID)
		c.Set("admin_username", claims.Username)
		c.Set("admin_roles", claims.Roles)
		c.Next()
	}
}

func AdminRequireRole(roleCodes ...string) gin.HandlerFunc {
	allowedRoles := map[string]struct{}{}
	for _, roleCode := range roleCodes {
		allowedRoles[roleCode] = struct{}{}
	}

	return func(c *gin.Context) {
		value, exists := c.Get("admin_roles")
		if !exists {
			basehandler.Forbidden(c, "没有权限")
			c.Abort()
			return
		}

		roles, ok := value.([]string)
		if !ok {
			basehandler.Forbidden(c, "没有权限")
			c.Abort()
			return
		}

		for _, role := range roles {
			if _, ok := allowedRoles[role]; ok {
				c.Next()
				return
			}
		}

		basehandler.Forbidden(c, "没有权限")
		c.Abort()
	}
}

func bearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
