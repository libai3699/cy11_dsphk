# Agent 目录约定

每一个 Agent 都必须有自己的独立目录，不把所有提示词、输入输出、工具调用混在一个大文件里。

当前第一版按内容运营流程拆 6 个 Agent：

- `accountdiagnosis`：商家账号诊断 Agent
- `benchmark`：对标账号分析 Agent
- `hotspottopic`：热点 / 选题 Agent
- `copywriting`：文案脚本 Agent
- `storyboard`：分镜脚本 Agent
- `review`：数据复盘 Agent

目录原则：

- 每个 Agent 自己维护 `Input`、`Output`、`Agent.Run`
- 公共返回结构放在 `server/internal/agent/types.go`
- Agent 只生成建议、内容和分析，不直接发布、不结算、不删除数据
- 后端 handler/service 负责拿数据库上下文、调用 Agent、保存结果、进入人工确认流程

后续接大模型时，每个目录可以继续拆：

```text
agent.go      # Go 调用入口
prompt.md     # 系统提示词
tools.go      # 该 Agent 能用的工具
schema.go     # 输入输出结构
```
