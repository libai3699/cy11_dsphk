package admin

import (
	"context"
	"errors"
	"time"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/benchmarkscout"
	"cy11dsphk/server/internal/agent/funneldoctor"
	"cy11dsphk/server/internal/agent/hotspotradar"
	"cy11dsphk/server/internal/agent/rhythmscheduler"
	"cy11dsphk/server/internal/agent/scriptwriter"
	"cy11dsphk/server/internal/agent/topicplanner"
	"cy11dsphk/server/internal/agent/viralanatomist"
	"cy11dsphk/server/internal/database"
	basehandler "cy11dsphk/server/internal/handler"
	"cy11dsphk/server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CateringAgentHandler 承接餐饮获客 Agent（M4'）：对标侦察 / 热点雷达 / 爆款拆解 / 选题规划。
type CateringAgentHandler struct{}

func NewCateringAgentHandler() *CateringAgentHandler {
	return &CateringAgentHandler{}
}

type benchmarkScoutRunPayload struct {
	     MerchantID                             uint64 `json:"merchantId" binding:"required"`
	       Platform                             string `json:"platform"`
	    CityKeyword                             string `json:"cityKeyword"`
	   SeedAccounts                           []string `json:"seedAccounts"`
	ExcludeAccounts                           []string `json:"excludeAccounts"`
	  PastedSamples []benchmarkscout.PastedSampleInput `json:"pastedSamples"`
}

type hotspotRadarRunPayload struct {
	MerchantID                      uint64 `json:"merchantId" binding:"required"`
	    Region                      string `json:"region"`
	  Industry                      string `json:"industry"`
	  RawItems []hotspotradar.RawItemInput `json:"rawItems"`
	WindowDays                         int `json:"windowDays"`
}

type viralAnatomistRunPayload struct {
	 MerchantID                   uint64 `json:"merchantId" binding:"required"`
	   ClientID                   string `json:"clientId"`
	   VideoURL                   string `json:"videoUrl"`
	 Transcript                   string `json:"transcript"`
	    Metrics viralanatomist.A3Metrics `json:"metrics"`
	ProductHook                   string `json:"productHook"`
	PublishedAt                   string `json:"publishedAt"`
}

type topicPlannerRunPayload struct {
	 MerchantID                     uint64 `json:"merchantId" binding:"required"`
	   ClientID                     string `json:"clientId"`
	  WeekLabel                     string `json:"weekLabel"`
	 HotspotIDs                   []string `json:"hotspotIds"`
	   ViralIDs                   []string `json:"viralIds"`
	      Quota       topicplanner.A4Quota `json:"quota"`
	Constraints topicplanner.A4Constraints `json:"constraints"`
}

type scriptWriterRunPayload struct {
	MerchantID  uint64   `json:"merchantId" binding:"required"`
	ClientID    string   `json:"clientId"`
	TopicID     string   `json:"topicId"`
	Type        string   `json:"type"`
	DurationSec int      `json:"durationSec"`
	Tone        string   `json:"tone"`
	MustInclude []string `json:"mustInclude"`
}

type rhythmSchedulerRunPayload struct {
	MerchantID  uint64                     `json:"merchantId" binding:"required"`
	ClientID    string                     `json:"clientId"`
	Weeks       int                        `json:"weeks"`
	Stage       string                     `json:"stage"`
	ScriptIDs   []string                   `json:"scriptIds"`
	Constraints rhythmscheduler.Constraints `json:"constraints"`
}

type funnelDoctorRunPayload struct {
	MerchantID uint64              `json:"merchantId" binding:"required"`
	ClientID   string              `json:"clientId"`
	Period     string              `json:"period"`
	Metrics    funneldoctor.Metrics `json:"metrics"`
	Baseline   map[string]any      `json:"baseline"`
	Notes      string              `json:"notes"`
}

// RunBenchmarkScout POST /api/admin/agents/benchmarkscout/run
func (h *CateringAgentHandler) RunBenchmarkScout(c *gin.Context) {
	var payload benchmarkScoutRunPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并填写对标侦察输入")
		return
	}
	if _, ok := cateringLoadMerchant(c, payload.MerchantID); !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()
	output, runErr := benchmarkscout.Agent{}.Run(ctx, benchmarkscout.Input{
		MerchantID:      payload.MerchantID,
		Platform:        payload.Platform,
		CityKeyword:     payload.CityKeyword,
		SeedAccounts:    payload.SeedAccounts,
		ExcludeAccounts: payload.ExcludeAccounts,
		PastedSamples:   payload.PastedSamples,
		Options:         agent.RunOptions{OperatorID: c.GetUint64("admin_user_id"), DryRun: true},
	})
	if runErr != nil {
		basehandler.ServerError(c, "对标侦察 Agent 执行失败："+runErr.Error())
		return
	}
	basehandler.OK(c, output)
}

// RunHotspotRadar POST /api/admin/agents/hotspotradar/run
func (h *CateringAgentHandler) RunHotspotRadar(c *gin.Context) {
	var payload hotspotRadarRunPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并填写热点雷达输入")
		return
	}
	merchant, ok := cateringLoadMerchant(c, payload.MerchantID)
	if !ok {
		return
	}
	industry := payload.Industry
	if industry == "" {
		industry = merchant.Industry
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()
	output, runErr := hotspotradar.Agent{}.Run(ctx, hotspotradar.Input{
		MerchantID: payload.MerchantID,
		Region:     payload.Region,
		Industry:   industry,
		RawItems:   payload.RawItems,
		WindowDays: payload.WindowDays,
		Options:    agent.RunOptions{OperatorID: c.GetUint64("admin_user_id"), DryRun: true},
	})
	if runErr != nil {
		basehandler.ServerError(c, "热点雷达 Agent 执行失败："+runErr.Error())
		return
	}
	basehandler.OK(c, output)
}

// RunViralAnatomist POST /api/admin/agents/viralanatomist/run
func (h *CateringAgentHandler) RunViralAnatomist(c *gin.Context) {
	var payload viralAnatomistRunPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并填写爆款拆解输入")
		return
	}
	if _, ok := cateringLoadMerchant(c, payload.MerchantID); !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()
	output, runErr := viralanatomist.Agent{}.Run(ctx, viralanatomist.Input{
		MerchantID:  payload.MerchantID,
		ClientID:    payload.ClientID,
		VideoURL:    payload.VideoURL,
		Transcript:  payload.Transcript,
		Metrics:     payload.Metrics,
		ProductHook: payload.ProductHook,
		PublishedAt: payload.PublishedAt,
		Options:     agent.RunOptions{OperatorID: c.GetUint64("admin_user_id"), DryRun: true},
	})
	if runErr != nil {
		basehandler.ServerError(c, "爆款拆解 Agent 执行失败："+runErr.Error())
		return
	}
	basehandler.OK(c, output)
}

// RunTopicPlanner POST /api/admin/agents/topicplanner/run
func (h *CateringAgentHandler) RunTopicPlanner(c *gin.Context) {
	var payload topicPlannerRunPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并填写选题规划输入")
		return
	}
	if _, ok := cateringLoadMerchant(c, payload.MerchantID); !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()
	output, runErr := topicplanner.Agent{}.Run(ctx, topicplanner.Input{
		MerchantID:  payload.MerchantID,
		ClientID:    payload.ClientID,
		WeekLabel:   payload.WeekLabel,
		HotspotIDs:  payload.HotspotIDs,
		ViralIDs:    payload.ViralIDs,
		Quota:       payload.Quota,
		Constraints: payload.Constraints,
		Options:     agent.RunOptions{OperatorID: c.GetUint64("admin_user_id"), DryRun: true},
	})
	if runErr != nil {
		basehandler.ServerError(c, "选题规划 Agent 执行失败："+runErr.Error())
		return
	}
	basehandler.OK(c, output)
}

// RunScriptWriter POST /api/admin/agents/scriptwriter/run
func (h *CateringAgentHandler) RunScriptWriter(c *gin.Context) {
	var payload scriptWriterRunPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并填写脚本作家输入")
		return
	}
	if _, ok := cateringLoadMerchant(c, payload.MerchantID); !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()
	output, runErr := scriptwriter.Agent{}.Run(ctx, scriptwriter.Input{
		MerchantID:  payload.MerchantID,
		ClientID:    payload.ClientID,
		TopicID:     payload.TopicID,
		Type:        payload.Type,
		DurationSec: payload.DurationSec,
		Tone:        payload.Tone,
		MustInclude: payload.MustInclude,
		Options:     agent.RunOptions{OperatorID: c.GetUint64("admin_user_id"), DryRun: true},
	})
	if runErr != nil {
		basehandler.ServerError(c, "脚本作家 Agent 执行失败："+runErr.Error())
		return
	}
	basehandler.OK(c, output)
}

// RunRhythmScheduler POST /api/admin/agents/rhythmscheduler/run
func (h *CateringAgentHandler) RunRhythmScheduler(c *gin.Context) {
	var payload rhythmSchedulerRunPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并填写节奏排期输入")
		return
	}
	if _, ok := cateringLoadMerchant(c, payload.MerchantID); !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()
	output, runErr := rhythmscheduler.Agent{}.Run(ctx, rhythmscheduler.Input{
		MerchantID:  payload.MerchantID,
		ClientID:    payload.ClientID,
		Weeks:       payload.Weeks,
		Stage:       payload.Stage,
		ScriptIDs:   payload.ScriptIDs,
		Constraints: payload.Constraints,
		Options:     agent.RunOptions{OperatorID: c.GetUint64("admin_user_id"), DryRun: true},
	})
	if runErr != nil {
		basehandler.ServerError(c, "节奏排期 Agent 执行失败："+runErr.Error())
		return
	}
	basehandler.OK(c, output)
}

// RunFunnelDoctor POST /api/admin/agents/funneldoctor/run
func (h *CateringAgentHandler) RunFunnelDoctor(c *gin.Context) {
	var payload funnelDoctorRunPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		basehandler.BadRequest(c, "请选择商家并填写漏斗诊断输入")
		return
	}
	if _, ok := cateringLoadMerchant(c, payload.MerchantID); !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()
	output, runErr := funneldoctor.Agent{}.Run(ctx, funneldoctor.Input{
		MerchantID: payload.MerchantID,
		ClientID:   payload.ClientID,
		Period:     payload.Period,
		Metrics:    payload.Metrics,
		Baseline:   payload.Baseline,
		Notes:      payload.Notes,
		Options:    agent.RunOptions{OperatorID: c.GetUint64("admin_user_id"), DryRun: true},
	})
	if runErr != nil {
		basehandler.ServerError(c, "漏斗诊断 Agent 执行失败："+runErr.Error())
		return
	}
	basehandler.OK(c, output)
}

// cateringLoadMerchant 校验商家存在性，对齐 content_topic.Generate 的“请选择商家/商家不存在”语义。
func cateringLoadMerchant(c *gin.Context, merchantID uint64) (*model.Merchant, bool) {
	var merchant model.Merchant
	if err := database.DB.First(&merchant, merchantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			basehandler.BadRequest(c, "商家不存在")
			return nil, false
		}
		basehandler.ServerError(c, "读取商家失败")
		return nil, false
	}
	return &merchant, true
}
