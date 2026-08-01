package review

import (
	"context"

	"cy11dsphk/server/internal/agent"
)

const Name = "review"

type Agent struct{}

type Input struct {
	MerchantID uint64           `json:"merchantId"`
	VideoID    uint64           `json:"videoId"`
	Title      string           `json:"title"`
	Metrics    Metrics          `json:"metrics"`
	Options    agent.RunOptions `json:"options,omitempty"`
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
	_ = ctx
	return Output{
		Result: agent.Result{
			Agent:   Name,
			Version: "v0.1.0",
			Status:  "draft",
			Summary: "已生成数据复盘草稿，等待接入真实数据和模型。",
			Suggestions: []string{
				"复盘要区分内容表现和成交表现。",
				"复盘只给建议，不直接影响结算。",
			},
		},
		Conclusion: "当前为占位复盘结论。",
		NextTopics: []string{
			"继续围绕高转化套餐做体验类选题。",
		},
		Optimizes: []string{
			"补充前 3 秒钩子测试。",
		},
	}, nil
}
