package router

import (
	"cy11dsphk/server/internal/config"
	adminhandler "cy11dsphk/server/internal/handler/admin"
	"cy11dsphk/server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func New(cfg config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), middleware.CORS())

	adminAuth := adminhandler.NewAuthHandler(cfg)
	adminMenu := adminhandler.NewMenuHandler()

	api := r.Group("/api/admin")
	{
		api.POST("/auth/login", adminAuth.Login)
		api.POST("/auth/logout", adminAuth.Logout)

		protected := api.Group("")
		protected.Use(middleware.AdminAuthRequired(cfg))
		{
			protected.POST("/auth/refresh", adminAuth.RefreshToken)
			protected.GET("/auth/codes", adminAuth.AccessCodes)
			protected.GET("/user/info", adminAuth.UserInfo)
			protected.GET("/menus", adminMenu.List)
		}
	}

	return r
}
