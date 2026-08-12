package viralanatomist

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

const Name = "viralanatomist"

type Agent struct{}

// Input 对齐来源 A 的 schemas/a3.input.json。
type Input struct {
	MerchantID  uint64           `json:"merchantId"`
	ClientID    string           `json:"clientId,omitempty"`
	VideoURL    string           `json:"videoUrl,omitempty"`
	Transcript  string           `json:"transcript"`
	Metrics     A3Metrics        `json:"metrics"`
	ProductHook string           `json:"productHook,omitempty"`
	PublishedAt string           `json:"publishedAt,omitempty"`
	Options     agent.RunOptions `json:"options,omitempty"`
}

type A3Metrics struct {
	Views    int64 `json:"views"`
	Likes    int64 `json:"likes"`
	Comments int64 `json:"comments"`
	Shares   int64 `json:"shares"`
	Collects int64 `json:"collects"`
}

// Output 字段对齐来源 A 的 schemas/a3.output.json（嵌套 struct + json tag）。
type Output struct {
	agent.Result
	Summary       string          `json:"summary"`
	Structure     []A3Segment     `json:"structure"`
	HookAnalysis  A3HookAnalysis  `json:"hookAnalysis"`
	ViralFactors  []A3ViralFactor `json:"viralFactors"`
	Replicability A3Replicability `json:"replicability"`
	Angles        []A3Angle       `json:"angles"`
}

type A3Segment struct {
	Segment  string  `json:"segment"`
	StartSec float64 `json:"startSec"`
	EndSec   float64 `json:"endSec"`
	Purpose  string  `json:"purpose"`
	Technique string `json:"technique"`
}

type A3HookAnalysis struct {
	Type            string `json:"type"`
	First3sec       string `json:"first3sec"`
	RetentionDriver string `json:"retentionDriver"`
}

type A3ViralFactor struct {
	Factor   string  `json:"factor"`
	Weight   float64 `json:"weight"`
	Evidence string  `json:"evidence"`
}

type A3Replicability struct {
	Score  int    `json:"score"`
	Reason string `json:"reason"`
}

type A3Angle struct {
	Angle        string `json:"angle"`
	ProfileBasis string `json:"profileBasis"`
	Difficulty   string `json:"difficulty"`
	ExpectedFit  string `json:"expectedFit"`
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
			Content: "请根据以下输入生成爆款拆解 JSON：\n" + string(payload),
		},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Summary       string          `json:"summary"`
		Structure     []A3Segment     `json:"structure"`
		HookAnalysis  A3HookAnalysis  `json:"hookAnalysis"`
		ViralFactors  []A3ViralFactor `json:"viralFactors"`
		Replicability A3Replicability `json:"replicability"`
		Angles        []A3Angle       `json:"angles"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if parsed.Summary == "" || len(parsed.Structure) == 0 || len(parsed.Angles) == 0 {
		return Output{}, fmt.Errorf("stepfun viralanatomist response is empty")
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
		Summary:       parsed.Summary,
		Structure:     parsed.Structure,
		HookAnalysis:  parsed.HookAnalysis,
		ViralFactors:  parsed.ViralFactors,
		Replicability: parsed.Replicability,
		Angles:        parsed.Angles,
	}, nil
}
