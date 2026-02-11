# Blog-Backend

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.24.0-00ADD8?logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/PostgreSQL-17.6-4169E1?logo=postgresql&logoColor=white" />
  <img src="https://img.shields.io/badge/Redis-8.0.5-DC382D?logo=redis&logoColor=white" />
  <img src="https://img.shields.io/badge/Nginx-1.26.3-009639?logo=nginx&logoColor=white" />
  <img src="https://img.shields.io/badge/Gin-1.11.0-00ADD8?logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/GORM-1.31.1-00ADD8?logo=go&logoColor=white" />
  <img src="https://img.shields.io/github/v/release/xianbo-deep/BLOG-BACKEND?label=release&include_prereleases" />

这是一个基于 Go 语言和 Gin 框架开发的博客后台服务系统。主要用于博客站点的流量统计、性能监控以及后台管理功能。
</div>

---

## 项目目的

本项目旨在为博客提供一个轻量级、高性能的后端服务，核心目标包括：

1.  **流量统计**：收集并分析访客数据，包括 PV/UV、地理位置、设备信息等。
2.  **性能监控**：监控页面加载延迟，识别慢加载页面，辅助优化前端性能。
3.  **后台管理**：提供可视化的数据看板，展示访问趋势、访客地图及详细访问日志。

## 技术栈

-   **编程语言**: Go
-   **Web 框架**: Gin
-   **数据库**: PostgreSQL
-   **ORM**: GORM
-   **缓存**: Redis
-   **认证**: JWT (JSON Web Token)
-   **定时任务**: Robfig Cron
-   **IP 地理位置**: GeoIP2


## 项目结构

```text
├── api/              # Serverless 函数入口点（如用于 Vercel 部署）
├── bootstrap/        # 应用启动和依赖注入
├── consts/           # 全局常量（环境变量键名、业务状态码）
├── core/             # 核心基础设施（数据库、Redis、GeoIP 初始化）
├── dto/              # 数据传输对象（请求/响应数据结构）
├── internal/         # 内部应用逻辑（外部不可直接访问）
│   ├── controller/   # HTTP 处理器（公共接口和管理后台）
│   ├── dao/          # 数据访问对象（数据库操作层）
│   ├── job/          # 后台任务（数据同步、死链检查等）
│   ├── notify/       # 通知服务（邮件通知等）
│   ├── service/      # 业务逻辑服务层
│   ├── task/         # 定时任务调度器
│   └── ws/           # WebSocket 中心（实时推送）
├── middleware/       # Gin 中间件（认证、跨域、超时控制等）
├── model/           # GORM 数据库模型
├── router/          # 路由定义
├── test/            # 单元测试和集成测试
├── thirdparty/       # 第三方 API 客户端（GitHub GraphQL 等）
└── utils/           # 共享工具函数（JWT、GeoIP 助手等）
```


## 提交类型说明

**格式**

- `type`:提交类型
- `scope`:模块名
- `message`:描述信息
```text
<type>(scope): <message>
```

**提交类型**

|    类型    |  含义   |
|:--------:|:-----:|
|   feat   |  新功能  |
|   fix    | 修复bug |
| refactor |  重构   |
|   docs   | 文档变更  |
|  chore   |  杂项   |
|    ci    | CI/CD |


## 部署说明

本项目使用Github Action执行自动化脚本，将推送后的代码自动编译、部署到服务器，并使用了Nginx进行反向代理。

## 快速开始

1. 克隆本项目

```shell
git clone https://github.com/xianbo-deep/BLOG-BACKEND.git
```


2. 查看项目依赖

```shell
cd BLOG-BACKEND
cat go.mod
```

3. 下载依赖

```shell
go mod download
```

4. 运行

- 直接运行

```shell
go run main.go
```

- 编译后运行

```shell
go build -o blog-backend .
./blog-backend
```

**注意**

- 执行上述命令前请先下载好 Go 编译器
- 请基于Bash执行上述命令
- 需要使用其它数据库请预先下载好对应的驱动库
- 请将环境变量更替为你自己的值
- 您拉取的分支可能并不是最新的，如遇报错请提交Issue
- 请根据您博客的部署情况配置您自己的可信代理以获取真实客户端IP

**环境变量**

```go
const (
	EnvPgURI               = "PG_URI"                   // 数据库地址
	EnvRedisURL            = "REDIS_URL"                // Redis地址
	EnvJWTSecret           = "JWT_SECRET"               // JWT密钥
	EnvAdminUser           = "ADMIN_USER"               // 统计后台用户名
	EnvAdminPwd            = "ADMIN_PASSWORD"           // 统计后台用户密码
	EnvGeoDBPath           = "GEODB_PATH"               // Geo数据库路径
	EnvBaseURL             = "BASE_URL"                 // 你的博客地址
	EnvAdminURL            = "ADMIN_URL"                // 统计后台地址
	EnvPort                = "PORT"                     // Gin监听的端口
	EnvDiscussionToken     = "DISCUSSION_TOKEN"         // Github Discussion密钥
	EnvGithubWebhookSecret = "GITHUB_WEBHOOK_SECRET"    // Github Webhook密钥
)
```

**可信代理**

请在[router.go](https://github.com/xianbo-deep/BLOG-BACKEND/blob/main/router/router.go)中配置您自己的可信代理
