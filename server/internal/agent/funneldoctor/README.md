# funneldoctor（A7 漏斗医生）

- **角色**：本地生活餐饮赛道的漏斗诊断专家，还原门店六阶段漏斗、定位瓶颈、给出最多 3 条一周内可执行动作。
- **域**：餐饮获客（内容营销 Agent 域）。
- **用法**：`provider.NewStepFunForAgent(funneldoctor.Name)` 注入 StepFun 客户端后调用 `(funneldoctor.Agent{}).Run(ctx, input)`。
- **Name**：`"funneldoctor"`
- **对应源 agent**：`知识库/marketingskills/落地实测/餐饮获客/官网和工作台-app/server/ai/agents/a7-funnel-doctor.agent.yaml`
  - 提示词：`prompts/funnel-diagnosis.md`
  - 输入/输出契约：`schemas/a7.input.json`、`schemas/a7.output.json`
  - 渲染模板：`templates/funnel-report.md.hbs`

## 输入输出对齐

- `Input` 含 `MerchantID uint64`，外加来源 A 输入契约字段（`clientId`/`period`/`metrics`/`baseline`/`notes`）；`metrics` 嵌套六阶段指标（中文 key：`播放`/`5秒完播率`/`团购点击`/`成交`/`下单`/`核销`）。
- `Output` 严格对齐 `a7.output.json`：`period`、`funnel`(嵌套 `stage`/`value`/`rate`/`benchmark`/`health`)、`bottleneck`(嵌套 `stage`/`severity`/`evidence`)、`actions`(嵌套 5 字段)、`watchlist`，并内嵌 `agent.Result`。
