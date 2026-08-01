package accountdiagnosis

import (
	"context"
	"encoding/json"
	"fmt"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

const Name = "account_diagnosis"

type Agent struct{}

type Input struct {
	MerchantID   uint64             `json:"merchantId"`
	AccountName  string             `json:"accountName"`
	Industry     string             `json:"industry"`
	City         string             `json:"city"`
	Profile      string             `json:"profile"`
	RecentVideos []VideoSnapshot    `json:"recentVideos,omitempty"`
	Metrics      map[string]float64 `json:"metrics,omitempty"`
	Options      agent.RunOptions   `json:"options,omitempty"`
}

type VideoSnapshot struct {
	Title        string `json:"title"`
	PublishAt    string `json:"publishAt"`
	PlayCount    int64  `json:"playCount"`
	LikeCount    int64  `json:"likeCount"`
	CommentCount int64  `json:"commentCount"`
	DealCount    int64  `json:"dealCount"`
}

type Output struct {
	agent.Result
	Problems     []string `json:"problems,omitempty"`
	NextActions  []string `json:"nextActions,omitempty"`
	AccountScore int      `json:"accountScore"`
	ContentScore int      `json:"contentScore"`
	ConvertScore int      `json:"convertScore"`
}

func (Agent) Run(ctx context.Context, input Input) (Output, error) {
	client, ok := provider.NewStepFunFromEnv()
	if !ok {
		return Output{}, agent.ErrProviderNotConfigured
	}
	return runWithStepFun(ctx, client, input)
}

func runWithStepFun(ctx context.Context, client *provider.Client, input Input) (Output, error) {
	payload, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return Output{}, err
	}

	resp, err := client.Chat(ctx, []provider.Message{
		{
			Role: "system",
			Content: `你是一个只输出 JSON 的本地商家抖音账号诊断 Agent。
你只做账号诊断，不写选题、不写完整文案。
你必须输出严格 JSON 对象，不要输出 Markdown，不要输出解释，不要输出反问，不要要求补充信息。
如果输入缺少信息，就在 problems 里指出缺失，并基于已有信息给保守建议。
JSON 结构：
{
  "summary": "一句话总结",
  "accountScore": 0-100,
  "contentScore": 0-100,
  "convertScore": 0-100,
  "problems": ["问题1"],
  "nextActions": ["动作1"],
  "positioning": "账号定位建议",
  "contentDirection": "内容方向建议",
  "conversionAdvice": "套餐转化建议"
}
评分要保守，动作要具体，必须围绕本地商家、抖音团购、到店核销。`,
		},
		{
			Role:    "user",
			Content: "请根据以下输入生成账号诊断 JSON：\n" + string(payload),
		},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Summary          string   `json:"summary"`
		AccountScore     int      `json:"accountScore"`
		ContentScore     int      `json:"contentScore"`
		ConvertScore     int      `json:"convertScore"`
		Problems         []string `json:"problems"`
		NextActions      []string `json:"nextActions"`
		Positioning      string   `json:"positioning"`
		ContentDirection string   `json:"contentDirection"`
		ConversionAdvice string   `json:"conversionAdvice"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if parsed.Summary == "" {
		return Output{}, fmt.Errorf("stepfun response summary is empty")
	}

	return Output{
		Result: agent.Result{
			Agent:       Name,
			Version:     "v0.2.0-stepfun",
			Status:      "completed",
			Summary:     parsed.Summary,
			Suggestions: parsed.NextActions,
			Artifacts: map[string]any{
				"provider":         "stepfun",
				"model":            resp.Model,
				"profile":          input.Profile,
				"metrics":          input.Metrics,
				"recentVideos":     input.RecentVideos,
				"positioning":      parsed.Positioning,
				"contentDirection": parsed.ContentDirection,
				"conversionAdvice": parsed.ConversionAdvice,
				"rawModelResponse": resp.Content,
			},
		},
		Problems:     parsed.Problems,
		NextActions:  parsed.NextActions,
		AccountScore: clampScore(parsed.AccountScore),
		ContentScore: clampScore(parsed.ContentScore),
		ConvertScore: clampScore(parsed.ConvertScore),
	}, nil
}

func clampScore(score int) int {
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}
