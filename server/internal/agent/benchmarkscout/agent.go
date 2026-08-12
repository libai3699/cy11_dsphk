package benchmarkscout

import (
	"context"
	"encoding/json"
	"fmt"
	"embed"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

const Name = "benchmarkscout"

//go:embed prompt.md input.schema.json output.schema.json template.hbs
var assets embed.FS

type Agent struct{}

type PastedSampleInput struct {
	Author  string         `json:"author"`
	Text    string         `json:"text"`
	Metrics map[string]any `json:"metrics,omitempty"`
}

type Input struct {
	MerchantID       uint64              `json:"merchantId"`
	Platform         string              `json:"platform,omitempty"`
	CityKeyword      string              `json:"cityKeyword,omitempty"`
	SeedAccounts     []string            `json:"seedAccounts"`
	ExcludeAccounts  []string            `json:"excludeAccounts,omitempty"`
	PastedSamples    []PastedSampleInput `json:"pastedSamples"`
	Options          agent.RunOptions    `json:"options,omitempty"`
}

type RecommendedAccount struct {
	Account    string `json:"account"`
	MatchScore int    `json:"matchScore"`
	Reason     string `json:"reason"`
	Tier       string `json:"tier"`
}

type ContentCategory struct {
	Name     string   `json:"name"`
	Share    float64  `json:"share"`
	Examples []string `json:"examples"`
}

type PostingPattern struct {
	Frequency string   `json:"frequency"`
	TimeSlots []string `json:"timeSlots"`
	FormatMix string   `json:"formatMix"`
}

type Borrowable struct {
	Point         string `json:"point"`
	Evidence      string `json:"evidence"`
	Applicability string `json:"applicability"`
}

type Output struct {
	agent.Result
	Summary          string              `json:"summary"`
	Recommended      []RecommendedAccount `json:"recommended"`
	ContentCategories []ContentCategory  `json:"contentCategories"`
	PostingPattern   PostingPattern      `json:"postingPattern"`
	Borrowables      []Borrowable        `json:"borrowables"`
	Gaps             []string            `json:"gaps"`
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
			Content: "请根据以下输入生成对标分析报告 JSON：\n" + mustJSON(input),
		},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Summary           string                `json:"summary"`
		Recommended       []RecommendedAccount  `json:"recommended"`
		ContentCategories []ContentCategory     `json:"contentCategories"`
		PostingPattern    PostingPattern        `json:"postingPattern"`
		Borrowables       []Borrowable          `json:"borrowables"`
		Gaps              []string              `json:"gaps"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if len(parsed.Recommended) == 0 || len(parsed.Borrowables) == 0 {
		return Output{}, fmt.Errorf("stepfun benchmarkscout response is empty")
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
		Summary:           parsed.Summary,
		Recommended:       parsed.Recommended,
		ContentCategories: parsed.ContentCategories,
		PostingPattern:    parsed.PostingPattern,
		Borrowables:       parsed.Borrowables,
		Gaps:              parsed.Gaps,
	}, nil
}

func mustJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}
