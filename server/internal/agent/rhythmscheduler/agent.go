package rhythmscheduler

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

const Name = "rhythmscheduler"

//go:embed prompt.md input.schema.json output.schema.json template.hbs
var assets embed.FS

type Agent struct{}

type Constraints struct {
	ShootablePerWeek int      `json:"每周可拍条数"`
	PublishSlots     []string `json:"可发布时段"`
	Holidays        []string `json:"节假日"`
}

type Input struct {
	MerchantID  uint64           `json:"merchantId"`
	ClientID    string           `json:"clientId"`
	Weeks       int              `json:"weeks"`
	Stage       string           `json:"stage"`
	ScriptIDs   []string         `json:"scriptIds,omitempty"`
	Constraints Constraints      `json:"constraints"`
	Options     agent.RunOptions `json:"options,omitempty"`
}

type Slot struct {
	Date       string `json:"date"`
	TimeSlot   string `json:"timeSlot"`
	ScriptRef  string `json:"scriptRef"`
	ContentType string `json:"contentType"`
	Rationale  string `json:"rationale"`
}

type WeeklyRhythm struct {
	PostsPerWeek  int      `json:"postsPerWeek"`
	PillarSequence []string `json:"pillarSequence"`
}

type Risk struct {
	Risk        string `json:"risk"`
	Mitigation  string `json:"mitigation"`
}

type Output struct {
	agent.Result
	Range        string        `json:"range"`
	Slots        []Slot        `json:"slots"`
	WeeklyRhythm WeeklyRhythm  `json:"weeklyRhythm"`
	BufferSlots  []string      `json:"bufferSlots,omitempty"`
	Risks        []Risk        `json:"risks"`
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
		{Role: "user", Content: "请根据以下输入生成节奏排期 JSON：\n" + string(payload)},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Range        string        `json:"range"`
		Slots        []Slot        `json:"slots"`
		WeeklyRhythm WeeklyRhythm  `json:"weeklyRhythm"`
		BufferSlots  []string      `json:"bufferSlots"`
		Risks        []Risk        `json:"risks"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if parsed.Range == "" || len(parsed.Slots) == 0 || len(parsed.Risks) < 2 {
		return Output{}, fmt.Errorf("stepfun rhythmscheduler response is incomplete")
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
		Range:        parsed.Range,
		Slots:        parsed.Slots,
		WeeklyRhythm: parsed.WeeklyRhythm,
		BufferSlots:  parsed.BufferSlots,
		Risks:        parsed.Risks,
	}, nil
}
