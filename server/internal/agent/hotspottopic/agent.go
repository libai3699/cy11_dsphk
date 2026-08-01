package hotspottopic

import (
	"context"

	"cy11dsphk/server/internal/agent"
)

const Name = "hotspot_topic"

type Agent struct{}

type Input struct {
	MerchantID uint64           `json:"merchantId"`
	Industry   string           `json:"industry"`
	City       string           `json:"city"`
	Products   []Product        `json:"products,omitempty"`
	Hotspots   []Hotspot        `json:"hotspots,omitempty"`
	Options    agent.RunOptions `json:"options,omitempty"`
}

type Product struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	GrossMargin float64 `json:"grossMargin"`
}

type Hotspot struct {
	Title  string `json:"title"`
	Source string `json:"source"`
	Scope  string `json:"scope"`
}

type Topic struct {
	Title           string   `json:"title"`
	Hook            string   `json:"hook"`
	Angle           string   `json:"angle"`
	Target          string   `json:"target"`
	RiskLevel       string   `json:"riskLevel"`
	RecommendReason string   `json:"recommendReason"`
	Tags            []string `json:"tags,omitempty"`
}

type Output struct {
	agent.Result
	Topics []Topic `json:"topics,omitempty"`
}

func (Agent) Run(ctx context.Context, input Input) (Output, error) {
	_ = ctx
	return Output{
		Result: agent.Result{
			Agent:   Name,
			Version: "v0.1.0",
			Status:  "draft",
			Summary: "已生成选题池草稿，等待接入热点来源和模型。",
			Suggestions: []string{
				"选题必须绑定商家主推套餐或门店卖点。",
				"同城热点和全国热点要分开标记。",
			},
		},
		Topics: []Topic{
			{
				Title:           "门店主推套餐体验",
				Hook:            "本地人最近都在问这家店值不值",
				Angle:           "真实体验 + 套餐承接",
				Target:          "团购转化",
				RiskLevel:       "low",
				RecommendReason: "适合作为冷启动标准选题。",
			},
		},
	}, nil
}
