package topicplanner

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

//go:embed prompt.md input.schema.json output.schema.json template.hbs
var assets embed.FS

const Name = "topicplanner"

type Agent struct{}

// Input 对齐来源 A 的 schemas/a4.input.json。
type Input struct {
	MerchantID  uint64           `json:"merchantId"`
	ClientID    string           `json:"clientId,omitempty"`
	WeekLabel   string           `json:"weekLabel"`
	HotspotIDs  []string         `json:"hotspotIds"`
	ViralIDs    []string         `json:"viralIds"`
	Quota       A4Quota          `json:"quota"`
	Constraints A4Constraints    `json:"constraints"`
	Options     agent.RunOptions `json:"options,omitempty"`
}

type A4Quota struct {
	Total    int            `json:"total"`
	ByPillar map[string]int `json:"byPillar,omitempty"`
}

type A4Constraints struct {
	Manpower       string `json:"人力"`
	ShootCondition string `json:"拍摄条件"`
}

// Output 字段对齐来源 A 的 schemas/a4.output.json（嵌套 struct + json tag）。
type Output struct {
	agent.Result
	Week          string          `json:"week"`
	Topics        []A4Topic       `json:"topics"`
	PillarBalance A4PillarBalance `json:"pillarBalance"`
	Dropped       []A4Dropped     `json:"dropped"`
	Summary       string          `json:"summary"`
}

type A4Topic struct {
	Seq           int      `json:"seq"`
	Title         string   `json:"title"`
	Pillar        string   `json:"pillar"`
	LinkedProduct string   `json:"linkedProduct"`
	HookIdea      string   `json:"hookIdea"`
	SourceRefs    []string `json:"sourceRefs"`
	Priority      string   `json:"priority"`
	Effort        string   `json:"effort"`
	ExpectedRole  string   `json:"expectedRole"`
}

type A4PillarBalance struct {
	ByPillar map[string]float64 `json:"byPillar"`
}

type A4Dropped struct {
	Idea   string `json:"idea"`
	Reason string `json:"reason"`
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

	payload, err := json.MarshalIndent(input, "", "  ")
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
			Content: "请根据以下输入生成本周选题表 JSON：\n" + string(payload),
		},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Week          string          `json:"week"`
		Topics        []A4Topic       `json:"topics"`
		PillarBalance A4PillarBalance `json:"pillarBalance"`
		Dropped       []A4Dropped     `json:"dropped"`
		Summary       string          `json:"summary"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if parsed.Summary == "" || len(parsed.Topics) == 0 {
		return Output{}, fmt.Errorf("stepfun topicplanner response is empty")
	}

	return Output{
		Result: agent.Result{
			Agent:   Name,
			Version: "v1.0.0-stepfun",
			Status:  "completed",
			Summary: parsed.Summary,
			Artifacts: map[string]any{
				"provider":         "stepfun",
				"model":            resp.Model,
				"input":            input,
				"rawModelResponse": resp.Content,
			},
		},
		Week:          parsed.Week,
		Topics:        parsed.Topics,
		PillarBalance: parsed.PillarBalance,
		Dropped:       parsed.Dropped,
		Summary:       parsed.Summary,
	}, nil
}
