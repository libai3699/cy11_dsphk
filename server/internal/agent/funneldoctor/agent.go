package funneldoctor

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

const Name = "funneldoctor"

//go:embed prompt.md input.schema.json output.schema.json template.hbs
var assets embed.FS

type Agent struct{}

type Metrics struct {
	Plays        float64 `json:"播放"`
	Stay5sRate   float64 `json:"5秒完播率"`
	GroupBuyClick float64 `json:"团购点击"`
	Deals        float64 `json:"成交"`
	Orders       float64 `json:"下单"`
	Redemptions  float64 `json:"核销"`
}

type Input struct {
	MerchantID uint64           `json:"merchantId"`
	ClientID    string           `json:"clientId"`
	Period      string           `json:"period"`
	Metrics     Metrics          `json:"metrics"`
	Baseline    map[string]any   `json:"baseline,omitempty"`
	Notes       string           `json:"notes,omitempty"`
	Options     agent.RunOptions `json:"options,omitempty"`
}

type FunnelStage struct {
	Stage      string  `json:"stage"`
	Value      float64 `json:"value"`
	Rate       float64 `json:"rate"`
	Benchmark  string  `json:"benchmark"`
	Health     string  `json:"health"`
}

type Bottleneck struct {
	Stage    string `json:"stage"`
	Severity string `json:"severity"`
	Evidence string `json:"evidence"`
}

type Action struct {
	Action       string `json:"action"`
	TargetStage  string `json:"targetStage"`
	ExpectedLift string `json:"expectedLift"`
	Effort       string `json:"effort"`
	Owner        string `json:"owner"`
}

type Output struct {
	agent.Result
	Period     string        `json:"period"`
	Funnel     []FunnelStage `json:"funnel"`
	Bottleneck Bottleneck    `json:"bottleneck"`
	Actions    []Action      `json:"actions"`
	Watchlist  []string      `json:"watchlist,omitempty"`
}

func (Agent) Run(ctx context.Context, input Input) (Output, error) {
	client, ok := provider.NewStepFunForAgent(Name)
	if !ok {
		return Output{}, agent.ErrProviderNotConfigured
	}

	promptBytes, err := assets.ReadFile("prompt.md")
	if err != nil {
		return Output{}, fmt.Errorf("read prompt: %w", err)
	}
	systemPrompt := strings.ReplaceAll(string(promptBytes), "{{guardrails}}", "")

	payload, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return Output{}, err
	}

	resp, err := client.Chat(ctx, []provider.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: "请根据以下输入生成漏斗诊断 JSON：\n" + string(payload)},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Period     string        `json:"period"`
		Funnel     []FunnelStage `json:"funnel"`
		Bottleneck Bottleneck    `json:"bottleneck"`
		Actions    []Action      `json:"actions"`
		Watchlist  []string      `json:"watchlist"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if parsed.Period == "" || len(parsed.Funnel) == 0 || len(parsed.Actions) == 0 || parsed.Bottleneck.Stage == "" {
		return Output{}, fmt.Errorf("stepfun funneldoctor response is incomplete")
	}

	return Output{
		Result: agent.Result{
			Agent:       Name,
			Version:     "v1.0.0-stepfun",
			Status:      "completed",
			Artifacts: map[string]any{
				"provider":         "stepfun",
				"model":            resp.Model,
				"input":            input,
				"rawModelResponse": resp.Content,
			},
		},
		Period:     parsed.Period,
		Funnel:     parsed.Funnel,
		Bottleneck: parsed.Bottleneck,
		Actions:    parsed.Actions,
		Watchlist:  parsed.Watchlist,
	}, nil
}
