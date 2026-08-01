package benchmark

import (
	"context"

	"cy11dsphk/server/internal/agent"
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
	_ = ctx
	return Output{
		Result: agent.Result{
			Agent:   Name,
			Version: "v0.1.0",
			Status:  "draft",
			Summary: "已生成对标分析草稿，等待接入真实账号数据和模型。",
			Suggestions: []string{
				"对标账号要按同城、同行、全国优秀账号分组。",
				"输出必须包含可借鉴点和风险提醒。",
			},
		},
		Patterns: []string{"爆款结构、开头钩子、团购承接方式待模型分析。"},
		Risks:    []string{"禁止直接搬运文案和镜头。"},
	}, nil
}
