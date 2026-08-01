package storyboard

import (
	"context"

	"cy11dsphk/server/internal/agent"
)

const Name = "storyboard"

type Agent struct{}

type Input struct {
	MerchantID uint64           `json:"merchantId"`
	ScriptID   uint64           `json:"scriptId"`
	ScriptText string           `json:"scriptText"`
	Locations  []string         `json:"locations,omitempty"`
	Options    agent.RunOptions `json:"options,omitempty"`
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
	_ = ctx
	return Output{
		Result: agent.Result{
			Agent:   Name,
			Version: "v0.1.0",
			Status:  "draft",
			Summary: "已生成分镜草稿，等待接入真实模型。",
			Suggestions: []string{
				"每条分镜必须能被拍摄人员直接执行。",
				"分镜需要保留人工调整入口。",
			},
		},
		Shots: []Shot{
			{
				Index:    1,
				Duration: "0-3s",
				Location: "门店门头",
				Camera:   "中近景",
				Content:  "展示门店环境和招牌",
				Line:     "开头钩子待填充",
				Note:     "避免路人正脸入镜。",
			},
		},
	}, nil
}
