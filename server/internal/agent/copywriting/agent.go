package copywriting

import (
	"context"
	"encoding/json"
	"fmt"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

const Name = "copywriting"

type Agent struct{}

type Input struct {
	MerchantID       uint64           `json:"merchantId"`
	TopicID          uint64           `json:"topicId"`
	TopicTitle       string           `json:"topicTitle"`
	TopicHook        string           `json:"topicHook,omitempty"`
	TopicAngle       string           `json:"topicAngle,omitempty"`
	ConversionTarget string           `json:"conversionTarget,omitempty"`
	Merchant         MerchantContext  `json:"merchant"`
	Products         []ProductContext `json:"products,omitempty"`
	ExtraRequirement string           `json:"extraRequirement,omitempty"`
	Options          agent.RunOptions `json:"options,omitempty"`
}

type MerchantContext struct {
	Name          string   `json:"name"`
	Industry      string   `json:"industry"`
	City          string   `json:"city"`
	SellingPoints []string `json:"sellingPoints,omitempty"`
}

type ProductContext struct {
	Name          string  `json:"name"`
	SellingPrice  float64 `json:"sellingPrice"`
	OriginalPrice float64 `json:"originalPrice"`
	TrafficLabel  string  `json:"trafficLabel,omitempty"`
	ProfitGuard   string  `json:"profitGuard,omitempty"`
}

type Output struct {
	agent.Result
	Title         string   `json:"title"`
	Opening       string   `json:"opening"`
	Body          string   `json:"body"`
	Ending        string   `json:"ending"`
	CTA           string   `json:"cta"`
	FullScript    string   `json:"fullScript"`
	ShootingNotes []string `json:"shootingNotes,omitempty"`
}

func (Agent) Run(ctx context.Context, input Input) (Output, error) {
	client, ok := provider.NewStepFunForAgent(Name)
	if !ok {
		return Output{}, agent.ErrProviderNotConfigured
	}

	payload, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return Output{}, err
	}

	resp, err := client.Chat(ctx, []provider.Message{
		{
			Role: "system",
			Content: `你是本地生活商家短视频“文案脚本”Agent。
你只负责把已确认选题写成可拍摄、可口播、可转化的短视频脚本。
你必须输出严格 JSON 对象，不要输出 Markdown，不要输出解释，不要反问。
JSON 结构：
{
  "summary": "一句话说明脚本策略",
  "suggestions": ["运营注意事项"],
  "title": "视频标题",
  "opening": "0-3 秒开场钩子",
  "body": "主体口播/旁白",
  "ending": "结尾收束",
  "cta": "到店/团购/评论/私信转化引导",
  "fullScript": "完整脚本文案",
  "shootingNotes": ["拍摄注意事项"]
}
要求：
1. 用普通人能听懂的话，不写空泛营销词。
2. 不承诺绝对效果，不夸大疗效、收益、食品安全等。
3. 优先突出商家真实卖点、套餐价值、到店理由。
4. 适合 30-60 秒抖音本地生活视频。`,
		},
		{
			Role:    "user",
			Content: "请根据以下输入生成文案脚本 JSON：\n" + string(payload),
		},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Summary       string   `json:"summary"`
		Suggestions   []string `json:"suggestions"`
		Title         string   `json:"title"`
		Opening       string   `json:"opening"`
		Body          string   `json:"body"`
		Ending        string   `json:"ending"`
		CTA           string   `json:"cta"`
		FullScript    string   `json:"fullScript"`
		ShootingNotes []string `json:"shootingNotes"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if parsed.Title == "" || parsed.Opening == "" || parsed.FullScript == "" {
		return Output{}, fmt.Errorf("stepfun copywriting response is empty")
	}

	return Output{
		Result: agent.Result{
			Agent:       Name,
			Version:     "v0.2.0-stepfun",
			Status:      "completed",
			Summary:     parsed.Summary,
			Suggestions: parsed.Suggestions,
			Artifacts: map[string]any{
				"provider":         "stepfun",
				"model":            resp.Model,
				"input":            input,
				"rawModelResponse": resp.Content,
			},
		},
		Title:         parsed.Title,
		Opening:       parsed.Opening,
		Body:          parsed.Body,
		Ending:        parsed.Ending,
		CTA:           parsed.CTA,
		FullScript:    parsed.FullScript,
		ShootingNotes: parsed.ShootingNotes,
	}, nil
}
