# hotspotradar（热点雷达 / a2-hotspot-radar）

## 角色
本地生活餐饮赛道的**热点雷达**。把用户杂乱粘贴进来的热搜词、话题、链接，整理成带时效与相关度评级的热点清单（逐条标注分类、热度、本地相关度、时效衰减与过期时间）。不分析爆款结构、不出切入角度、不写脚本、不选题。

## 域
餐饮获客域 · 内容营销 Agent（来源 A 迁移，里程碑 M3 第一组）。

## 用法
```go
import "cy11dsphk/server/internal/agent/hotspotradar"

out, err := hotspotradar.Agent{}.Run(ctx, hotspotradar.Input{
    MerchantID: 123,
    Region:     "长沙",
    Industry:   "湘菜/本地生活",
    RawItems:   []hotspotradar.RawItemInput{{Text: "..."}},
    WindowDays: 7,
    Options:    agent.RunOptions{OperatorID: 1},
})
```
- `Name` = `"hotspotradar"`
- 走 `provider.NewStepFunForAgent(Name)`，未配置 `STEP_API_KEY` 时返回 `agent.ErrProviderNotConfigured`。
- 资源经 `//go:embed` 内联：`prompt.md`、`input.schema.json`、`output.schema.json`、`template.hbs`。

## 对应源 agent 路径
- 源 spec：`餐饮获客/官网和工作台-app/server/ai/agents/a2-hotspot-radar.agent.yaml`
- 提示词：`.../server/ai/prompts/hotspot.md`
- 输入 schema：`.../server/ai/schemas/a2.input.json`
- 输出 schema：`.../server/ai/schemas/a2.output.json`
- 模板：`.../server/ai/templates/hotspot-list.md.hbs`
