package scriptwriter

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

const Name = "scriptwriter"

//go:embed prompt.md input.schema.json output.schema.json template.hbs
var assets embed.FS

type Agent struct{}

type Input struct {
	MerchantID  uint64          `json:"merchantId"`
	ClientID    string          `json:"clientId"`
	TopicID     string          `json:"topicId"`
	Type        string          `json:"type"`
	DurationSec int             `json:"durationSec"`
	Tone        string          `json:"tone"`
	MustInclude []string        `json:"mustInclude,omitempty"`
	Options     agent.RunOptions `json:"options,omitempty"`
}

type Oral struct {
	Hook          string   `json:"hook"`
	Body          []string `json:"body"`
	CTA           string   `json:"cta"`
	WordCount     int      `json:"wordCount"`
	EstDurationSec int     `json:"estDurationSec"`
}

type Shot struct {
	ShotNo       int    `json:"shotNo"`
	TimeRange    string `json:"timeRange"`
	Visual       string `json:"visual"`
	Voiceover    string `json:"voiceover"`
	OnScreenText string `json:"onScreenText"`
	Note         string `json:"note"`
}

type Output struct {
	agent.Result
	TopicRef        string   `json:"topicRef"`
	Oral            Oral     `json:"oral"`
	Storyboard      []Shot   `json:"storyboard"`
	Titles          []string `json:"titles"`
	Hashtags        []string `json:"hashtags"`
	FirstComment    string   `json:"firstComment"`
	ComplianceNotes []string `json:"complianceNotes,omitempty"`
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
		{Role: "user", Content: "请根据以下输入生成脚本编剧 JSON：\n" + string(payload)},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		TopicRef        string   `json:"topicRef"`
		Oral            Oral     `json:"oral"`
		Storyboard      []Shot   `json:"storyboard"`
		Titles          []string `json:"titles"`
		Hashtags        []string `json:"hashtags"`
		FirstComment    string   `json:"firstComment"`
		ComplianceNotes []string `json:"complianceNotes"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if parsed.TopicRef == "" || parsed.Oral.Hook == "" || len(parsed.Storyboard) == 0 || len(parsed.Titles) < 3 {
		return Output{}, fmt.Errorf("stepfun scriptwriter response is incomplete")
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
		TopicRef:        parsed.TopicRef,
		Oral:            parsed.Oral,
		Storyboard:      parsed.Storyboard,
		Titles:          parsed.Titles,
		Hashtags:        parsed.Hashtags,
		FirstComment:    parsed.FirstComment,
		ComplianceNotes: parsed.ComplianceNotes,
	}, nil
}
