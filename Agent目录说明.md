# Agent 目录说明

后端里所有 Agent 都放在 `server/internal/agent` 下，并且每个 Agent 一个独立目录。

当前目录：

```text
server/internal/agent/
  types.go                 # Agent 公共返回结构
  accountdiagnosis/        # 商家账号诊断 Agent
  benchmark/               # 对标账号分析 Agent
  hotspottopic/            # 热点 / 选题 Agent
  copywriting/             # 文案脚本 Agent
  storyboard/              # 分镜脚本 Agent
  review/                  # 数据复盘 Agent
```

这样拆的原因：

- 每个 Agent 的输入、输出、提示词、工具调用都不一样。
- 后续哪个 Agent 复杂，就只扩哪个目录，不影响其他 Agent。
- 后端仍然是总控层，Agent 只是内容工人层。
- Agent 输出的是草稿、建议、分析，不直接发布、不直接结算、不直接删数据。

后续每个 Agent 目录可以继续扩成：

```text
agent.go      # Go 调用入口
prompt.md     # 系统提示词
tools.go      # 该 Agent 可用工具
schema.go     # 输入输出结构
provider.go   # 模型供应商适配
```
