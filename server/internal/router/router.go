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
	adminMerchant := adminhandler.NewMerchantHandler()
	adminMerchantPackage := adminhandler.NewMerchantPackageHandler()
	adminMerchantAccountAuth := adminhandler.NewMerchantAccountAuthHandler()
	adminAccountDiagnosis := adminhandler.NewAccountDiagnosisHandler()
	adminBenchmark := adminhandler.NewBenchmarkHandler()
	adminContentTopic := adminhandler.NewContentTopicHandler()
	adminContentProduction := adminhandler.NewContentProductionHandler()
	adminWorkspace := adminhandler.NewWorkspaceHandler()

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
			protected.GET("/workspace/overview", adminWorkspace.Overview)
			protected.GET("/merchants", adminMerchant.List)
			protected.POST("/merchants", adminMerchant.Create)
			protected.GET("/merchants/:id/workspace", adminMerchant.Workspace)
			protected.GET("/merchants/:id", adminMerchant.Get)
			protected.PUT("/merchants/:id", adminMerchant.Update)
			protected.DELETE("/merchants/:id", middleware.AdminRequireRole("super_admin"), adminMerchant.Delete)
			protected.GET("/packages", adminMerchantPackage.List)
			protected.POST("/packages", adminMerchantPackage.Create)
			protected.GET("/packages/:id", adminMerchantPackage.Get)
			protected.PUT("/packages/:id", adminMerchantPackage.Update)
			protected.DELETE("/packages/:id", adminMerchantPackage.Delete)
			protected.GET("/account-auths", adminMerchantAccountAuth.List)
			protected.POST("/account-auths", adminMerchantAccountAuth.Create)
			protected.GET("/account-auths/:id", adminMerchantAccountAuth.Get)
			protected.PUT("/account-auths/:id", adminMerchantAccountAuth.Update)
			protected.DELETE("/account-auths/:id", adminMerchantAccountAuth.Delete)
			protected.GET("/account-diagnoses", adminAccountDiagnosis.List)
			protected.POST("/account-diagnoses", adminAccountDiagnosis.Create)
			protected.GET("/account-diagnoses/:id", adminAccountDiagnosis.Get)
			protected.GET("/benchmarks", adminBenchmark.List)
			protected.POST("/benchmarks", adminBenchmark.Create)
			protected.GET("/benchmarks/:id", adminBenchmark.Get)
			protected.PUT("/benchmarks/:id", adminBenchmark.Update)
			protected.DELETE("/benchmarks/:id", adminBenchmark.Delete)
			protected.POST("/benchmarks/:id/analyze", adminBenchmark.Analyze)
			protected.GET("/benchmark-analyses", adminBenchmark.AnalysisList)
			protected.GET("/topics", adminContentTopic.List)
			protected.POST("/topics/generate", adminContentTopic.Generate)
			protected.PUT("/topics/:id/status", adminContentTopic.UpdateStatus)
			protected.GET("/scripts", adminContentProduction.ListScripts)
			protected.POST("/scripts/generate", adminContentProduction.GenerateScript)
			protected.PUT("/scripts/:id/status", adminContentProduction.UpdateScriptStatus)
			protected.GET("/storyboards", adminContentProduction.ListStoryboards)
			protected.POST("/storyboards/generate", adminContentProduction.GenerateStoryboard)
			protected.PUT("/storyboards/:id/status", adminContentProduction.UpdateStoryboardStatus)
			protected.GET("/shooting-tasks", adminContentProduction.ListShootingTasks)
			protected.POST("/shooting-tasks", adminContentProduction.CreateShootingTask)
			protected.PUT("/shooting-tasks/:id/status", adminContentProduction.UpdateShootingTaskStatus)
			protected.GET("/schedules", adminContentProduction.ListSchedules)
			protected.POST("/schedules", adminContentProduction.CreateSchedule)
			protected.PUT("/schedules/:id", adminContentProduction.UpdateSchedule)
			protected.PUT("/schedules/:id/status", adminContentProduction.UpdateScheduleStatus)
			protected.GET("/reviews", adminContentProduction.ListReviews)
			protected.POST("/reviews/generate", adminContentProduction.GenerateReview)
			protected.GET("/agent-configs", adminContentProduction.AgentConfigs)
		}
	}

	return r
}
