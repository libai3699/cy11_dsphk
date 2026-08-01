package accountdiagnosis

import (
	"context"

	"cy11dsphk/server/internal/agent"
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
	_ = ctx
	return Output{
		Result: agent.Result{
			Agent:   Name,
			Version: "v0.1.0",
			Status:  "draft",
			Summary: "已生成账号诊断草稿，等待接入真实模型后输出完整分析。",
			Suggestions: []string{
				"先补齐商家定位、主推套餐、近 10 条视频数据。",
				"诊断结果只进入待确认状态，不自动改账号。",
			},
		},
		Problems:     []string{"账号诊断 Agent 尚未接入模型。"},
		NextActions:  []string{"接入模型 provider", "保存诊断报告", "增加人工确认节点"},
		AccountScore: 60,
		ContentScore: 60,
		ConvertScore: 60,
	}, nil
}
