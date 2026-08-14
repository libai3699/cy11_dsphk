# topicplanner（选题策划官 A4）

## 角色
本地生活餐饮赛道的内容选题策划专家。综合门店档案、已确认热点清单、已确认爆款拆解（及可选漏斗诊断），产出本周选题表：每条选题含支柱类型、关联真实团购套餐、切入钩子、来源引用、优先级、拍摄难度与预期出镜角色。不写脚本文案、不排期。

## 域
餐饮获客域（M3 第二组，来源 A 内容营销 Agent 迁移）。

## 用法
```go
out, err := topicplanner.Agent{}.Run(ctx, topicplanner.Input{
    MerchantID: 123,
    WeekLabel:  "2026-W33",
    HotspotIDs: []string{"hotspot-1"},
    ViralIDs:   []string{"viral-1"},
    Quota:      topicplanner.A4Quota{Total: 5},
    Constraints: topicplanner.A4Constraints{Manpower: "店长+后厨每周最多4条", ShootCondition: "仅店内实拍"},
})
```
- `Name` = `"topicplanner"`
- 依赖 `provider.NewStepFunForAgent(Name)`（未配置返回 `agent.ErrProviderNotConfigured`）。
- `Output` 字段对齐 `output.schema.json`：`week` / `topics[]` / `pillarBalance.byPillar` / `dropped[]` / `summary`，并内嵌 `agent.Result`。

## 对应源 agent 路径（只读）
- `marketingskills/落地实测/餐饮获客/官网和工作台-app/server/ai/agents/a4-topic-planner.agent.yaml`
- prompt：`prompts/topic-planning.md`
- schema：`schemas/a4.input.json`、`schemas/a4.output.json`
- template：`templates/topic-plan.md.hbs`
