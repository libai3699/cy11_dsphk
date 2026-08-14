package hotspottopic

import (
	"context"
	"encoding/json"
	"fmt"

	"cy11dsphk/server/internal/agent"
	"cy11dsphk/server/internal/agent/provider"
)

const Name = "hotspot_topic"

type Agent struct{}

type Input struct {
	MerchantID        uint64             `json:"merchantId"`
	MerchantName      string             `json:"merchantName"`
	Industry          string             `json:"industry"`
	City              string             `json:"city"`
	Products          []Product          `json:"products,omitempty"`
	Hotspots          []Hotspot          `json:"hotspots,omitempty"`
	BenchmarkAccount  string             `json:"benchmarkAccount,omitempty"`
	BenchmarkSummary  string             `json:"benchmarkSummary,omitempty"`
	BenchmarkAccounts []BenchmarkAccount `json:"benchmarkAccounts,omitempty"`
	ExtraRequirement  string             `json:"extraRequirement,omitempty"`
	Options           agent.RunOptions   `json:"options,omitempty"`
}

type BenchmarkAccount struct {
	Name           string  `json:"name"`
	Platform       string  `json:"platform,omitempty"`
	City           string  `json:"city,omitempty"`
	Industry       string  `json:"industry,omitempty"`
	FollowerCount  float64 `json:"followerCount,omitempty"`
	BestPlayCount  float64 `json:"bestPlayCount,omitempty"`
	LatestHitTitle string  `json:"latestHitTitle,omitempty"`
	Takeaway       string  `json:"takeaway,omitempty"`
	Risk           string  `json:"risk,omitempty"`
}

type Product struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	GrossMargin float64 `json:"grossMargin"`
}

type Hotspot struct {
	Title  string `json:"title"`
	Source string `json:"source"`
	Scope  string `json:"scope"`
}

type Topic struct {
	Title           string   `json:"title"`
	Hook            string   `json:"hook"`
	Angle           string   `json:"angle"`
	Target          string   `json:"target"`
	RiskLevel       string   `json:"riskLevel"`
	RecommendReason string   `json:"recommendReason"`
	Tags            []string `json:"tags,omitempty"`
}

type Output struct {
	agent.Result
	Topics []Topic `json:"topics,omitempty"`
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
			Content: `你是本地生活商家短视频“找爆款/选题”Agent。
你只负责找爆款方向和生成可执行选题，不写完整文案、不写分镜。
你必须输出严格 JSON 对象，不要输出 Markdown，不要输出解释，不要反问。
如果信息缺失，也必须基于已有信息保守输出，并在 riskLevel/recommendReason 里说明。
JSON 结构：
{
  "summary": "一句话说明本轮找爆款依据",
  "suggestions": ["建议动作"],
  "topics": [
    {
      "title": "选题标题",
      "hook": "前三秒钩子",
      "angle": "内容角度",
      "target": "转化目标",
      "riskLevel": "low/medium/high",
      "recommendReason": "推荐理由",
      "tags": ["标签"]
    }
  ]
}
要求：
1. topics 输出 5 条。
2. 每条必须绑定商家、套餐、同城/行业热点或对标账号结构。
3. 不能承诺一定爆，不输出虚假夸大。
4. 避免低价伤商家利润。
5. 适合抖音本地生活、团购、到店核销。`,
		},
		{
			Role:    "user",
			Content: "请根据以下输入找爆款方向并生成选题 JSON：\n" + string(payload),
		},
	})
	if err != nil {
		return Output{}, err
	}

	var parsed struct {
		Summary     string   `json:"summary"`
		Suggestions []string `json:"suggestions"`
		Topics      []Topic  `json:"topics"`
	}
	if err := json.Unmarshal([]byte(provider.ExtractJSONObject(resp.Content)), &parsed); err != nil {
		return Output{}, err
	}
	if parsed.Summary == "" || len(parsed.Topics) == 0 {
		return Output{}, fmt.Errorf("stepfun hotspot response is empty")
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
		Topics: parsed.Topics,
	}, nil
}
