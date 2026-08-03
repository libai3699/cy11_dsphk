package storyboard

import (
	"context"
	"encoding/json"
	"fmt"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

const Name = "storyboard"

type Agent struct{}

type Input struct {
	MerchantID  uint64           `json:"merchantId"`
	ScriptID    uint64           `json:"scriptId"`
	ScriptTitle string           `json:"scriptTitle,omitempty"`
	ScriptText  string           `json:"scriptText"`
	Locations   []string         `json:"locations,omitempty"`
	Options     agent.RunOptions `json:"options,omitempty"`
}

type Shot struct {
	Index    int    `json:"index"`
	Duration string `json:"duration"`
	Location string `json:"location"`
	Camera   string `json:"camera"`
	Content  string `json:"content"`
	Line     string `json:"line"`
	Note     string `json:"note"`
}

type Output struct {
	agent.Result
	Shots []Shot `json:"shots,omitempty"`
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
			Content: `你是本地生活短视频“分镜脚本”Agent。
你只负责把文案脚本拆成拍摄人员能执行的镜头表。
你必须输出严格 JSON 对象，不要输出 Markdown，不要输出解释，不要反问。
JSON 结构：
{
  "summary": "一句话说明分镜策略",
  "suggestions": ["拍摄执行建议"],
  "shots": [
    {
      "index": 1,
      "duration": "0-3s",
      "location": "拍摄地点",
      "camera": "镜头景别/运动",
      "content": "画面内容",
      "line": "镜头里说的话/字幕",
      "note": "执行注意事项"
    }
  ]
}
要求：
1. 输出 5-8 个镜头。
2. 每个镜头必须写清楚地点、画面、台词、注意事项。
3. 控制拍摄成本，优先门店现场可完成。
4. 避免拍到无授权路人正脸、敏感信息和夸大宣传。`,
		},
		{
			Role:    "user",
			Content: "请根据以下输入生成分镜脚本 JSON：\n" + string(payload),
		},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Summary     string   `json:"summary"`
		Suggestions []string `json:"suggestions"`
		Shots       []Shot   `json:"shots"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if len(parsed.Shots) == 0 {
		return Output{}, fmt.Errorf("stepfun storyboard response is empty")
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
		Shots: parsed.Shots,
	}, nil
}
