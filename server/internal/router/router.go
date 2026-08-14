package router

import (
	"cy11dsphk/server/internal/config"
	adminhandler "cy11dsphk/server/internal/handler/admin"
	"cy11dsphk/server/internal/middleware"
	"cy11dsphk/server/internal/model"
	"cy11dsphk/server/internal/service"

	"github.com/gin-gonic/gin"
)

func New(cfg config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), middleware.CORS())

	adminAuth := adminhandler.NewAuthHandler(cfg)
	adminMenu := adminhandler.NewMenuHandler()
	adminMerchant := adminhandler.NewMerchantHandler()
	adminMerchantPackage := adminhandler.NewMerchantPackageHandler()
	adminFollowUp := adminhandler.NewFollowUpHandler()
	adminSettlementOrder := adminhandler.NewSettlementOrderHandler()
	adminMerchantAccountAuth := adminhandler.NewMerchantAccountAuthHandler()
	adminAccountDiagnosis := adminhandler.NewAccountDiagnosisHandler()
	adminBenchmark := adminhandler.NewBenchmarkHandler()
	adminPlatformResearch := adminhandler.NewPlatformResearchHandler()
	adminContentTopic := adminhandler.NewContentTopicHandler()
	adminContentProduction := adminhandler.NewContentProductionHandler()
	adminWorkspace := adminhandler.NewWorkspaceHandler()
	adminCateringAgent := adminhandler.NewCateringAgentHandler()

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
			protected.GET("/follow-up-logs", adminFollowUp.List)
			protected.POST("/follow-up-logs", adminFollowUp.Create)
			protected.PUT("/follow-up-logs/:id", adminFollowUp.Update)
			protected.DELETE("/follow-up-logs/:id", adminFollowUp.Delete)
			protected.POST("/follow-up-logs/:id/suggestion", adminFollowUp.Suggestion)
			protected.GET("/settlement-orders", adminSettlementOrder.List)
			protected.POST("/settlement-orders/generate", adminSettlementOrder.Generate)
			protected.POST("/settlement-orders", adminSettlementOrder.Create)
			protected.PUT("/settlement-orders/:id", adminSettlementOrder.Update)
			protected.PUT("/settlement-orders/:id/status", adminSettlementOrder.UpdateStatus)
			protected.DELETE("/settlement-orders/:id", adminSettlementOrder.Delete)
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
			protected.POST("/benchmarks/analyze-batch", adminBenchmark.AnalyzeBatch)
			protected.POST("/benchmarks/:id/analyze", adminBenchmark.Analyze)
			protected.GET("/benchmark-analyses", adminBenchmark.AnalysisList)
			protected.GET("/platform-researches", adminPlatformResearch.List)
			protected.POST("/platform-researches/generate", adminPlatformResearch.Generate)
			protected.GET("/platform-researches/:id", adminPlatformResearch.Get)
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

			// 餐饮获客 Agent（M4'）：对标侦察 / 热点雷达 / 爆款拆解 / 选题规划
			protected.POST("/agents/benchmarkscout/run", adminCateringAgent.RunBenchmarkScout)
			protected.POST("/agents/hotspotradar/run", adminCateringAgent.RunHotspotRadar)
			protected.POST("/agents/viralanatomist/run", adminCateringAgent.RunViralAnatomist)
			protected.POST("/agents/topicplanner/run", adminCateringAgent.RunTopicPlanner)
			protected.POST("/agents/scriptwriter/run", adminCateringAgent.RunScriptWriter)
			protected.POST("/agents/rhythmscheduler/run", adminCateringAgent.RunRhythmScheduler)
			protected.POST("/agents/funneldoctor/run", adminCateringAgent.RunFunnelDoctor)

			// 运营知识库（M1）：痛点/案例/账号画像/平台规则/内容模板
			protected.GET("/knowledge/pain-points", knowledgeList[model.PainPoint](service.ListPainPoints))
			protected.GET("/knowledge/pain-points/:id", knowledgeGet[model.PainPoint](service.GetPainPoint))
			protected.POST("/knowledge/pain-points", knowledgeCreate[model.PainPoint](service.CreatePainPoint))
			protected.PUT("/knowledge/pain-points/:id", knowledgeUpdate[model.PainPoint](service.UpdatePainPoint))
			protected.DELETE("/knowledge/pain-points/:id", knowledgeDelete(service.DeletePainPoint))

			protected.GET("/knowledge/case-studies", knowledgeList[model.CaseStudy](service.ListCaseStudies))
			protected.GET("/knowledge/case-studies/:id", knowledgeGet[model.CaseStudy](service.GetCaseStudy))
			protected.POST("/knowledge/case-studies", knowledgeCreate[model.CaseStudy](service.CreateCaseStudy))
			protected.PUT("/knowledge/case-studies/:id", knowledgeUpdate[model.CaseStudy](service.UpdateCaseStudy))
			protected.DELETE("/knowledge/case-studies/:id", knowledgeDelete(service.DeleteCaseStudy))

			protected.GET("/knowledge/account-profiles", knowledgeList[model.AccountProfile](service.ListAccountProfiles))
			protected.GET("/knowledge/account-profiles/:id", knowledgeGet[model.AccountProfile](service.GetAccountProfile))
			protected.POST("/knowledge/account-profiles", knowledgeCreate[model.AccountProfile](service.CreateAccountProfile))
			protected.PUT("/knowledge/account-profiles/:id", knowledgeUpdate[model.AccountProfile](service.UpdateAccountProfile))
			protected.DELETE("/knowledge/account-profiles/:id", knowledgeDelete(service.DeleteAccountProfile))

			protected.GET("/knowledge/platform-rules", knowledgeList[model.PlatformRule](service.ListPlatformRules))
			protected.GET("/knowledge/platform-rules/:id", knowledgeGet[model.PlatformRule](service.GetPlatformRule))
			protected.POST("/knowledge/platform-rules", knowledgeCreate[model.PlatformRule](service.CreatePlatformRule))
			protected.PUT("/knowledge/platform-rules/:id", knowledgeUpdate[model.PlatformRule](service.UpdatePlatformRule))
			protected.DELETE("/knowledge/platform-rules/:id", knowledgeDelete(service.DeletePlatformRule))

			protected.GET("/knowledge/content-templates", knowledgeList[model.ContentTemplate](service.ListContentTemplates))
			protected.GET("/knowledge/content-templates/:id", knowledgeGet[model.ContentTemplate](service.GetContentTemplate))
			protected.POST("/knowledge/content-templates", knowledgeCreate[model.ContentTemplate](service.CreateContentTemplate))
			protected.PUT("/knowledge/content-templates/:id", knowledgeUpdate[model.ContentTemplate](service.UpdateContentTemplate))
			protected.DELETE("/knowledge/content-templates/:id", knowledgeDelete(service.DeleteContentTemplate))
		}
	}

	return r
}
