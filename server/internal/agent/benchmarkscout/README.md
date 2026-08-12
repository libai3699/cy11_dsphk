# benchmarkscout（对标猎手 / a1-benchmark-scout）

## 角色
本地生活餐饮赛道的**对标分析专家**。从人工提供的候选账号与粘贴内容中，筛出真正可对标的账号，归纳内容类别、发布规律、可借鉴点。不做热点判断、不出选题、不评价单条视频爆因。

## 域
餐饮获客域 · 内容营销 Agent（来源 A 迁移，里程碑 M3 第一组）。

## 用法
```go
import "cy11dsphk/server/internal/agent/benchmarkscout"

out, err := benchmarkscout.Agent{}.Run(ctx, benchmarkscout.Input{
    MerchantID:    123,
    SeedAccounts:  []string{"账号A", "账号B"},
    PastedSamples: []benchmarkscout.PastedSampleInput{{Author: "x", Text: "..."}},
    Options:       agent.RunOptions{OperatorID: 1},
})
```
- `Name` = `"benchmarkscout"`
- 走 `provider.NewStepFunForAgent(Name)`，未配置 `STEP_API_KEY` 时返回 `agent.ErrProviderNotConfigured`。
- 资源经 `//go:embed` 内联：`prompt.md`、`input.schema.json`、`output.schema.json`、`template.hbs`。

## 对应源 agent 路径
- 源 spec：`餐饮获客/官网和工作台-app/server/ai/agents/a1-benchmark-scout.agent.yaml`
- 提示词：`.../server/ai/prompts/benchmark.md`
- 输入 schema：`.../server/ai/schemas/a1.input.json`
- 输出 schema：`.../server/ai/schemas/a1.output.json`
- 模板：`.../server/ai/templates/benchmark-report.md.hbs`
