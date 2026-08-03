package benchmark

import (
	"context"
	"encoding/json"
	"fmt"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

const Name = "benchmark"

type Agent struct{}

type Input struct {
	MerchantID        uint64             `json:"merchantId"`
	Industry          string             `json:"industry"`
	City              string             `json:"city"`
	BenchmarkAccounts []BenchmarkAccount `json:"benchmarkAccounts,omitempty"`
	Options           agent.RunOptions   `json:"options,omitempty"`
}

type BenchmarkAccount struct {
	Name       string   `json:"name"`
	PlatformID string   `json:"platformId,omitempty"`
	Reason     string   `json:"reason,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

type Output struct {
	agent.Result
	Patterns []string `json:"patterns,omitempty"`
	Risks    []string `json:"risks,omitempty"`
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
			Content: `你是本地生活商家短视频“对标分析”Agent。
你只负责拆解对标账号可学习的内容结构、转化路径和风险，不直接搬运文案、镜头、套餐价格。
你必须输出严格 JSON 对象，不要输出 Markdown，不要输出解释，不要反问。
JSON 结构：
{
  "summary": "一句话总结对标拆解",
  "suggestions": ["后续运营动作"],
  "patterns": ["可学习结构"],
  "risks": ["不能照抄或需要注意的风险"]
}
要求：
1. patterns 输出 3-6 条，必须围绕开场、内容结构、转化路径。
2. risks 输出 2-5 条，必须提醒不要低价伤利润、不要侵权照搬。
3. 建议要能进入后续选题 Agent 使用。`,
		},
		{
			Role:    "user",
			Content: "请根据以下输入生成对标分析 JSON：\n" + string(payload),
		},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Summary     string   `json:"summary"`
		Suggestions []string `json:"suggestions"`
		Patterns    []string `json:"patterns"`
		Risks       []string `json:"risks"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if parsed.Summary == "" || len(parsed.Patterns) == 0 {
		return Output{}, fmt.Errorf("stepfun benchmark response is empty")
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
		Patterns: parsed.Patterns,
		Risks:    parsed.Risks,
	}, nil
}
