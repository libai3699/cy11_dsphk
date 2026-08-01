package admin

import (
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/service"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct{}

func NewMenuHandler() *MenuHandler {
	return &MenuHandler{}
}

func (h *MenuHandler) List(c *gin.Context) {
	userID := c.GetUint64("admin_user_id")
	menus, err := service.GetUserMenus(userID)
	if err != nil {
		basehandler.ServerError(c, "读取菜单失败")
		return
	}
	basehandler.OK(c, service.BuildMenuRoutes(menus))
}
