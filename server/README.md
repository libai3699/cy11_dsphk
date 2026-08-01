# CY11DSPHK Go 后端

当前阶段只做后台基础能力：

- 登录、退出
- JWT 鉴权
- 用户、角色、菜单权限
- 根据角色返回后端菜单
- MySQL 自动建库、自动迁移、基础数据 seed

## 本地启动

```powershell
cd server
copy .env.example .env
go run ./cmd
```

默认数据库名：`cy11_dsphk`

默认后台账号：

- 账号：`admin`
- 密码：`admin123456`

前端开发环境已经代理 `/api` 到 `http://127.0.0.1:8989`。
