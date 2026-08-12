# M3：来源 A 7 个内容营销 Agent 领域资产迁移（餐饮获客域）

> 目标仓库：server/internal/agent/（Go + StepFun）
> 来源：D:\李峰荣\知识库\marketingskills\落地实测\餐饮获客\官网和工作台-app\server\ai\
> 策略：只搬内容/资产不搬 Node 框架；沿用目标既有 agent 模式；不改动现有 6 agent、不动 router/handler（接线留 M4/M6）。

## 一、映射（来源 A → 新包名）
| 源 Agent | 新包名 server/internal/agent/ | 角色 |
|---|---|---|
| a1-benchmark-scout | benchmarkscout | 对标猎手（账号/对标分析）|
| a2-hotspot-radar | hotspotradar | 热点雷达（热点选题）|
| a3-viral-anatomist | viralanatomist | 爆款解剖师（爆款拆解）|
| a4-topic-planner | topicplanner | 选题策划官（选题表）|
| a5-script-writer | scriptwriter | 脚本作家（分镜脚本）|
| a6-rhythm-scheduler | rhythmscheduler | 节奏排期官（发布排期）|
| a7-funnel-doctor | funneldoctor | 漏斗医生（数据诊断/复盘）|

## 二、每个新包的产物
- `agent.go`：沿用目标模式
  - `const Name = "<新包名>"`
  - `type Agent struct{}`
  - `type Input struct { MerchantID uint64; /* 领域参数 */; Options agent.RunOptions }`
  - `type Output struct { agent.Result; /* 字段对齐源 output.schema.json */ }`
  - `func (Agent) Run(ctx, input) (Output, error)`：内嵌 system prompt → `provider.NewStepFunForAgent(Name)` → `.Chat(ctx, messages)`（json_object）→ 解析进 Output；未配置返回 `agent.ErrProviderNotConfigured`
  - `//go:embed prompt.md input.schema.json output.schema.json template.hbs` + `var assets embed.FS`
- `prompt.md`：来源 A 对应 prompt 原文（或 agent.yaml 的 role 提炼）
- `input.schema.json` / `output.schema.json`：来源 A 对应 schema
- `template.hbs`：来源 A 对应 Handlebars 模板
- `README.md`：角色、域、用法、`Name`、对应源 agent

## 三、参考（只读）
- `server/internal/agent/types.go`（Result / RunOptions / ErrProviderNotConfigured）
- `server/internal/agent/provider/provider.go`（NewStepFunForAgent / Chat）
- `server/internal/agent/copywriting/agent.go`（模式范本，照抄 import 路径与结构）

## 四、子 agent 分片
- SubA：benchmarkscout + hotspotradar（a1, a2）
- SubB：viralanatomist + topicplanner（a3, a4）
- SubC：scriptwriter + rhythmscheduler + funneldoctor（a5, a6, a7）

## 五、硬约束
- 只读：上述参考文件 + 来源 A 的对应 agent 文件
- 只写：server/internal/agent/<新包名>/ 目录
- 禁止改现有 6 agent / router / handler / provider；禁止删；不跑 go build
- 报告 ≤400 字/组：列出创建文件、字段映射、假设
