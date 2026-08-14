# viralanatomist（爆款解剖师 A3）

## 角色
本地生活餐饮赛道的短视频结构分析专家。把一条已爆视频拆成可复用结构件：结构分段、爆因归因、可复制性判断、与本店契合切入角度。不写脚本、不排期、不选题。

## 域
餐饮获客域（M3 第二组，来源 A 内容营销 Agent 迁移）。

## 用法
```go
out, err := viralanatomist.Agent{}.Run(ctx, viralanatomist.Input{
    MerchantID: 123,
    Transcript: "视频逐字稿（≥50字）",
    Metrics:    viralanatomist.A3Metrics{Views: 100000, Likes: 5000, Comments: 300, Shares: 200, Collects: 800},
})
```
- `Name` = `"viralanatomist"`
- 依赖 `provider.NewStepFunForAgent(Name)`（未配置返回 `agent.ErrProviderNotConfigured`）。
- `Output` 字段对齐 `output.schema.json`：`summary` / `structure[]` / `hookAnalysis` / `viralFactors[]` / `replicability` / `angles[]`，并内嵌 `agent.Result`。

## 对应源 agent 路径（只读）
- `marketingskills/落地实测/餐饮获客/官网和工作台-app/server/ai/agents/a3-viral-anatomist.agent.yaml`
- prompt：`prompts/viral-deconstruct.md`
- schema：`schemas/a3.input.json`、`schemas/a3.output.json`
- template：`templates/viral-report.md.hbs`
