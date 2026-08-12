# rhythmscheduler（A6 节奏调度师）

- **角色**：本地生活餐饮赛道的节奏调度师，依据门店阶段、产能/时段/节假日约束排出周/月发布节奏与时段。
- **域**：餐饮获客（内容营销 Agent 域）。
- **用法**：`provider.NewStepFunForAgent(rhythmscheduler.Name)` 注入 StepFun 客户端后调用 `(rhythmscheduler.Agent{}).Run(ctx, input)`。
- **Name**：`"rhythmscheduler"`
- **对应源 agent**：`知识库/marketingskills/落地实测/餐饮获客/官网和工作台-app/server/ai/agents/a6-rhythm-scheduler.agent.yaml`
  - 提示词：`prompts/content-rhythm.md`
  - 输入/输出契约：`schemas/a6.input.json`、`schemas/a6.output.json`
  - 渲染模板：`templates/calendar-plan.md.hbs`

## 输入输出对齐

- `Input` 含 `MerchantID uint64`，外加来源 A 输入契约字段（`clientId`/`weeks`/`stage`/`scriptIds`/`constraints`）；`constraints` 嵌套 `每周可拍条数`/`可发布时段`/`节假日`。
- `Output` 严格对齐 `a6.output.json`：`range`、`slots`(嵌套 `date`/`timeSlot`/`scriptRef`/`contentType`/`rationale`)、`weeklyRhythm`(嵌套 `postsPerWeek`/`pillarSequence`)、`bufferSlots`、`risks`(嵌套 `risk`/`mitigation`)，并内嵌 `agent.Result`。
