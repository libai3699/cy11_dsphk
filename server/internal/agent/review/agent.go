package review

import (
	"context"
	"encoding/json"
	"fmt"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

const Name = "review"

type Agent struct{}

type Input struct {
	MerchantID  uint64           `json:"merchantId"`
	VideoID     uint64           `json:"videoId"`
	Title       string           `json:"title"`
	PeriodStart string           `json:"periodStart,omitempty"`
	PeriodEnd   string           `json:"periodEnd,omitempty"`
	Metrics     Metrics          `json:"metrics"`
	Options     agent.RunOptions `json:"options,omitempty"`
}

type Metrics struct {
	PlayCount      int64   `json:"playCount"`
	LikeCount      int64   `json:"likeCount"`
	CommentCount   int64   `json:"commentCount"`
	ShareCount     int64   `json:"shareCount"`
	DealCount      int64   `json:"dealCount"`
	WriteOffAmount float64 `json:"writeOffAmount"`
}

type Output struct {
	agent.Result
	Conclusion string   `json:"conclusion"`
	NextTopics []string `json:"nextTopics,omitempty"`
	Optimizes  []string `json:"optimizes,omitempty"`
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
			Content: `你是本地生活商家短视频“数据复盘”Agent。
你只负责根据视频数据给运营复盘，不做财务结算，不承诺下一条一定爆。
你必须输出严格 JSON 对象，不要输出 Markdown，不要输出解释，不要反问。
JSON 结构：
{
  "summary": "一句话总结表现",
  "suggestions": ["后续动作建议"],
  "conclusion": "复盘结论",
  "nextTopics": ["下一轮可测试选题"],
  "optimizes": ["具体优化点"]
}
要求：
1. 区分播放、互动、成交/核销三类表现。
2. 给出下一轮选题方向和拍摄优化动作。
3. 不把低价促销作为唯一建议。
4. 如果数据缺失，要说明缺什么数据。`,
		},
		{
			Role:    "user",
			Content: "请根据以下输入生成复盘 JSON：\n" + string(payload),
		},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Summary     string   `json:"summary"`
		Suggestions []string `json:"suggestions"`
		Conclusion  string   `json:"conclusion"`
		NextTopics  []string `json:"nextTopics"`
		Optimizes   []string `json:"optimizes"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if parsed.Conclusion == "" {
		return Output{}, fmt.Errorf("stepfun review response is empty")
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
		Conclusion: parsed.Conclusion,
		NextTopics: parsed.NextTopics,
		Optimizes:  parsed.Optimizes,
	}, nil
}
