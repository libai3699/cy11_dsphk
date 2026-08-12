package hotspotradar

import (
	"context"
	"encoding/json"
	"fmt"
	"embed"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

const Name = "hotspotradar"

//go:embed prompt.md input.schema.json output.schema.json template.hbs
var assets embed.FS

type Agent struct{}

type RawItemInput struct {
	Text       string `json:"text"`
	Source     string `json:"source,omitempty"`
	CapturedAt string `json:"capturedAt,omitempty"`
}

type Input struct {
	MerchantID uint64         `json:"merchantId"`
	Region     string         `json:"region,omitempty"`
	Industry   string         `json:"industry,omitempty"`
	RawItems   []RawItemInput `json:"rawItems"`
	WindowDays int            `json:"windowDays,omitempty"`
	Options    agent.RunOptions `json:"options,omitempty"`
}

type HotspotItem struct {
	Title          string `json:"title"`
	Category       string `json:"category"`
	HeatLevel      string `json:"heatLevel"`
	LocalRelevance string `json:"localRelevance"`
	DecayEstimate  string `json:"decayEstimate"`
	ExpiresAt      string `json:"expiresAt"`
	SourceRef      string `json:"sourceRef"`
}

type DiscardedItem struct {
	Item   string `json:"item"`
	Reason string `json:"reason"`
}

type Output struct {
	agent.Result
	Summary   string          `json:"summary"`
	Items     []HotspotItem   `json:"items"`
	TopPicks  []string        `json:"topPicks"`
	Discarded []DiscardedItem `json:"discarded"`
}

func (Agent) Run(ctx context.Context, input Input) (Output, error) {
	client, ok := provider.NewStepFunForAgent(Name)
	if !ok {
		return Output{}, agent.ErrProviderNotConfigured
	}

	promptBytes, err := assets.ReadFile("prompt.md")
	if err != nil {
		return Output{}, err
	}

	resp, err := client.Chat(ctx, []provider.Message{
		{
			Role:    "system",
			Content: string(promptBytes),
		},
		{
			Role:    "user",
			Content: "请根据以下输入生成热点清单 JSON：\n" + mustJSON(input),
		},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Summary   string          `json:"summary"`
		Items     []HotspotItem   `json:"items"`
		TopPicks  []string        `json:"topPicks"`
		Discarded []DiscardedItem `json:"discarded"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if len(parsed.Items) == 0 || len(parsed.TopPicks) == 0 {
		return Output{}, fmt.Errorf("stepfun hotspotradar response is empty")
	}

	return Output{
		Result: agent.Result{
			Agent:   Name,
			Version: "v0.2.0-stepfun",
			Status:  "completed",
			Summary: parsed.Summary,
			Artifacts: map[string]any{
				"provider":         "stepfun",
				"model":            resp.Model,
				"input":            input,
				"rawModelResponse": resp.Content,
			},
		},
		Summary:   parsed.Summary,
		Items:     parsed.Items,
		TopPicks:  parsed.TopPicks,
		Discarded: parsed.Discarded,
	}, nil
}

func mustJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}
