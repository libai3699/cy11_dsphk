package copywriting

import (
	"context"

	"cy11dsphk/server/internal/agent"
)

const Name = "copywriting"

type Agent struct{}

type Input struct {
	MerchantID uint64           `json:"merchantId"`
	TopicID    uint64           `json:"topicId"`
	TopicTitle string           `json:"topicTitle"`
	Merchant   MerchantContext  `json:"merchant"`
	Options    agent.RunOptions `json:"options,omitempty"`
}

type MerchantContext struct {
	Name          string   `json:"name"`
	Industry      string   `json:"industry"`
	City          string   `json:"city"`
	SellingPoints []string `json:"sellingPoints,omitempty"`
}

type Output struct {
	agent.Result
	Opening string `json:"opening"`
	Body    string `json:"body"`
	Ending  string `json:"ending"`
	CTA     string `json:"cta"`
}

func (Agent) Run(ctx context.Context, input Input) (Output, error) {
	_ = ctx
	return Output{
		Result: agent.Result{
			Agent:   Name,
			Version: "v0.1.0",
			Status:  "draft",
			Summary: "已生成文案脚本草稿，等待接入真实模型。",
			Suggestions: []string{
				"文案要围绕选题、套餐、门店卖点生成。",
				"生成后必须由运营人员确认。",
			},
		},
		Opening: "开头 3 秒钩子待模型生成。",
		Body:    "主体表达待模型生成。",
		Ending:  "结尾引导待模型生成。",
		CTA:     "引导到店或团购下单。",
	}, nil
}
