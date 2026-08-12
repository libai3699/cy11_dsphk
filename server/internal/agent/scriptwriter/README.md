# scriptwriter（A5 脚本编剧）

- **角色**：本地生活餐饮赛道的脚本编剧，把已确认选题写成口播脚本 + 逐镜头分镜表 + 3 个备选标题 + 话题标签 + 自评首评。
- **域**：餐饮获客（内容营销 Agent 域）。
- **用法**：`provider.NewStepFunForAgent(scriptwriter.Name)` 注入 StepFun 客户端后调用 `(scriptwriter.Agent{}).Run(ctx, input)`。
- **Name**：`"scriptwriter"`
- **对应源 agent**：`知识库/marketingskills/落地实测/餐饮获客/官网和工作台-app/server/ai/agents/a5-script-writer.agent.yaml`
  - 提示词：`prompts/script-storyboard.md`
  - 输入/输出契约：`schemas/a5.input.json`、`schemas/a5.output.json`
  - 渲染模板：`templates/script-package.md.hbs`

## 输入输出对齐

- `Input` 含 `MerchantID uint64`，外加来源 A 输入契约字段（`clientId`/`topicId`/`type`/`durationSec`/`tone`/`mustInclude`）。
- `Output` 严格对齐 `a5.output.json`：`topicRef`、`oral`(嵌套 `hook`/`body`/`cta`/`wordCount`/`estDurationSec`)、`storyboard`(嵌套 `shotNo` 等)、`titles`、`hashtags`、`firstComment`、`complianceNotes`，并内嵌 `agent.Result`。
