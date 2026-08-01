package benchmark

import (
	"context"
	"fmt"
	"strings"

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
	accountName := "对标账号"
	reason := ""
	if len(input.BenchmarkAccounts) > 0 {
		accountName = input.BenchmarkAccounts[0].Name
		reason = input.BenchmarkAccounts[0].Reason
	}
	patterns := buildPatterns(input, reason)
	risks := buildRisks(reason)
	suggestions := buildSuggestions(input, accountName)

	return Output{
		Result: agent.Result{
			Agent:   Name,
			Version: "v0.1.0",
			Status:  "completed",
			Summary: fmt.Sprintf(
				"已完成 %s 的对标分析。建议只拆结构和转化路径，不直接搬运文案、镜头和套餐价格。",
				accountName,
			),
			Suggestions: suggestions,
			Artifacts: map[string]any{
				"city":              input.City,
				"industry":          input.Industry,
				"benchmarkAccounts": input.BenchmarkAccounts,
			},
		},
		Patterns: patterns,
		Risks:    risks,
	}, nil
}

func buildPatterns(input Input, reason string) []string {
	patterns := []string{
		"拆开场：先看前 3 秒是否用了价格、反差、老板出镜、排队或份量感。",
		"拆结构：记录门头、产品特写、人物口播、套餐展示、到店引导的顺序。",
		"拆转化：重点看有没有把流量落到团购套餐、评论私信或到店核销。",
	}
	if input.City != "" {
		patterns = append(patterns, "同城维度：优先保留本地地名、商圈、门店距离和真实场景。")
	}
	if input.Industry != "" {
		patterns = append(patterns, "行业维度：围绕"+input.Industry+"用户最关心的价格、品质、体验和场景做拆解。")
	}
	if strings.TrimSpace(reason) != "" {
		patterns = append(patterns, "人工可抄点："+strings.TrimSpace(reason))
	}
	return patterns
}

func buildRisks(reason string) []string {
	risks := []string{
		"不要直接照搬对标账号的文案、镜头顺序和口播表达。",
		"不要为了追求销量盲目压低商家套餐价格。",
		"不要复制与本商家客单价、位置、人设不匹配的内容。",
	}
	if strings.Contains(reason, "低价") || strings.Contains(reason, "便宜") {
		risks = append(risks, "如果对标账号靠低价爆量，要单独核算扣除提点后的利润。")
	}
	return risks
}

func buildSuggestions(input Input, accountName string) []string {
	return []string{
		"把 " + accountName + " 最近爆款拆成 3 个可复用选题方向。",
		"每个方向都要绑定一个主推套餐，避免只涨播放不成交。",
		"先生成 5 条选题做小批量测试，再决定是否放大。",
	}
}
