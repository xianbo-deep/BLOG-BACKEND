# Blog-Backend

<div align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.24.0-00ADD8?logo=go&logoColor=white" />
  <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-17.6-4169E1?logo=postgresql&logoColor=white" />
  <img alt="Redis" src="https://img.shields.io/badge/Redis-8.0.5-DC382D?logo=redis&logoColor=white" />
  <img alt="Nginx" src="https://img.shields.io/badge/Nginx-1.26.3-009639?logo=nginx&logoColor=white" />
  <img alt="Gin" src="https://img.shields.io/badge/Gin-1.11.0-00ADD8?logo=go&logoColor=white" />
  <img alt="GORM" src="https://img.shields.io/badge/GORM-1.31.1-00ADD8?logo=go&logoColor=white" />
  <img alt="Release" src="https://img.shields.io/github/v/release/xianbo-deep/BLOG-BACKEND?label=release&include_prereleases" />

一个基于 Go + Gin 的博客后台服务。主要提供：统计采集、数据看板、性能监控、订阅邮件通知、GitHub Webhook 与实时 WebSocket 推送。
</div>

---

## 目录

- [功能特性](#功能特性)
- [启动流程](#启动流程)
- [项目结构](#项目结构)
- [快速开始](#快速开始)
  - [前置依赖](#前置依赖)
  - [配置环境变量](#配置环境变量)
  - [运行（Windows PowerShell）](#运行windows-powershell)
  - [运行（Linux/macOS Bash）](#运行linuxmacos-bash)
- [鉴权与安全](#鉴权与安全)
- [API 路由概览](#api-路由概览)
- [定时任务](#定时任务)
- [测试与质量](#测试与质量)
- [部署说明](#部署说明)
- [常见问题](#常见问题)

## 功能特性

- **博客统计采集**：采集 PV/UV、来源（Referer）、UA/设备、IP 地理位置等（`/blog/collect`）
- **订阅与邮件**：订阅邮箱、验证码校验、邮件模板渲染（`/blog/subscribe`、`/blog/verify`，模板见 `internal/notify/email/template/`）
- **管理后台**：登录获取 JWT，查询看板/趋势/访问日志/性能数据（`/admin/*`）
- **页面分析**：路径排行、来源分析、细分维度（device/source 等）（`/admin/analysis/*`）
- **访客地图**：世界地图/中国地图维度统计（`/admin/visitormap/*`）
- **GitHub 集成**：支持 GitHub Webhook 推送、Discussion 数据（`/webhook/*`、`/admin/discussionmap/*`）
- **WebSocket 实时推送**：管理端实时数据连接（`/admin/ws`）
- **定时任务**：数据同步、死链检查、Discussion 报告（见 `internal/task`、`internal/job`）

## 启动流程

入口为 `main.go`，启动顺序（与代码一致）：

1. `core.Init()`：初始化 PostgreSQL / Redis / GeoIP / IP2Region（`core/db.go`）
2. `bootstrap.InitComponet()`：组装 DAO / Service / Controller / Mailer / WS Hub 等（`bootstrap/component.go`）
3. `task.InitCron(comps)`：注册并启动定时任务（`internal/task/cron.go`）
4. `router.SetupRouter(comps)`：初始化 Gin 路由（`router/router.go`）

## 项目结构

```text
├── api/              # Serverless 函数入口点（如 Vercel 部署）
├── bootstrap/        # 应用启动与依赖注入
├── consts/           # 全局常量（环境变量键名、错误、业务常量）
├── core/             # 核心基础设施初始化（PostgreSQL、Redis、GeoIP、IP2Region）
├── dto/              # 请求/响应 DTO
├── internal/         # 内部逻辑
│   ├── controller/   # HTTP Handler（admin/public/github）
│   ├── dao/          # DAO（含 cache 子模块）
│   ├── job/          # 任务实现（sync/deadlink/discReport 等）
│   ├── notify/       # 通知（邮件）
│   ├── service/      # 业务服务层
│   ├── task/         # Cron 调度
│   └── ws/           # WebSocket hub/client
├── middleware/       # Gin 中间件（CORS、超时、鉴权、Webhook 校验等）
├── model/            # GORM model
├── router/           # 路由定义（README 的 API 清单以此为准）
├── test/             # 单测（geo/ua/referer 等）
├── thirdparty/       # 第三方 API（GitHub GraphQL 等）
└── utils/            # 公共工具（JWT、Geo、分页等）
```

## 快速开始

### 前置依赖

- Go（与 `go.mod` 声明版本一致）
- PostgreSQL（需要创建数据库并提供 DSN）
- Redis
- GeoIP2 数据文件（mmdb，例如 MaxMind GeoLite2 City）
- IP2Region 数据文件（xdb，IPv4/IPv6 各一份）

### 配置环境变量

项目**不会自动加载 `.env` 文件**（代码中未引入 dotenv），需要在运行环境/容器/CI 中注入环境变量。

#### 必需环境变量（来自 `consts/env.go`）

| Key | 说明 |
|---|---|
| `PG_URI` | PostgreSQL DSN |
| `REDIS_URL` | Redis URL（形如 `redis://:password@host:6379/0`） |
| `JWT_SECRET` | JWT 签名密钥 |
| `ADMIN_USER` | 管理后台用户名 |
| `ADMIN_PASSWORD` | 管理后台密码 |
| `GEODB_PATH` | GeoIP2 mmdb 文件路径 |
| `BASE_URL` | 博客站点 URL（用于 WS Origin 白名单等） |
| `ADMIN_URL` | 管理后台前端 URL（用于 WS Origin 白名单等） |
| `PORT` | 服务端口（为空默认 `8080`） |
| `DISCUSSION_TOKEN` | GitHub Token（用于 Discussion 数据） |
| `GITHUB_WEBHOOK_SECRET` | `/webhook/github` 的签名校验密钥 |
| `GITHUB_NOTIFY_SECRET` | `/webhook/notify` 的签名校验密钥 |

#### 额外环境变量（来自 `core/db.go`）

| Key | 说明 |
|---|---|
| `IP2REGION_V4_PATH` | IP2Region IPv4 xdb 文件路径 |
| `IP2REGION_V6_PATH` | IP2Region IPv6 xdb 文件路径 |

#### 示例（仅供本地开发参考）

> 注意：请不要把真实密钥提交到仓库。

```dotenv
PG_URI="host=127.0.0.1 user=postgres password=postgres dbname=blog port=5432 sslmode=disable TimeZone=Asia/Shanghai"
REDIS_URL="redis://:password@127.0.0.1:6379/0"

JWT_SECRET="change_me_to_a_random_string"
ADMIN_USER="admin"
ADMIN_PASSWORD="admin"

GEODB_PATH="E:\\data\\GeoLite2-City.mmdb"
IP2REGION_V4_PATH="E:\\data\\ip2region_v4.xdb"
IP2REGION_V6_PATH="E:\\data\\ip2region_v6.xdb"

BASE_URL="https://your-blog.example.com"
ADMIN_URL="https://your-admin.example.com"

DISCUSSION_TOKEN="ghp_xxx"
GITHUB_WEBHOOK_SECRET="change_me"
GITHUB_NOTIFY_SECRET="change_me"

PORT="8080"
```

### 运行（Windows PowerShell）

> 下面示例只演示最小必需变量，你也可以把所有变量都设置完整再启动。


### 运行（Linux/macOS Bash）

```bash
export PG_URI='host=127.0.0.1 user=postgres password=postgres dbname=blog port=5432 sslmode=disable TimeZone=Asia/Shanghai'
export REDIS_URL='redis://:password@127.0.0.1:6379/0'
export JWT_SECRET='change_me'
export ADMIN_USER='admin'
export ADMIN_PASSWORD='admin'
export GEODB_PATH='/data/GeoLite2-City.mmdb'
export IP2REGION_V4_PATH='/data/ip2region_v4.xdb'
export IP2REGION_V6_PATH='/data/ip2region_v6.xdb'
export BASE_URL='https://your-blog.example.com'
export ADMIN_URL='https://your-admin.example.com'
export DISCUSSION_TOKEN='ghp_xxx'
export GITHUB_WEBHOOK_SECRET='change_me'
export GITHUB_NOTIFY_SECRET='change_me'
export PORT='8080'

go run .
```

## 鉴权与安全

- **Admin API（HTTP）**：使用 `Authorization: Bearer <token>`（见 `middleware/auth.go` 与 `router/router.go`）
- **WebSocket**：使用 query 参数 `token`：`GET /admin/ws?token=<token>`（见 `middleware/auth.go` 与 `router/router.go`）
- **GitHub Webhook**：必须携带 `X-Hub-Signature-256`，并使用对应 secret 校验：
  - `/webhook/github` 使用 `GITHUB_WEBHOOK_SECRET`
  - `/webhook/notify` 使用 `GITHUB_NOTIFY_SECRET`
- **可信代理**：当前仅信任 `127.0.0.1`（`router/router.go` 的 `SetTrustedProxies`），生产环境请按你的反代拓扑调整

## API 路由概览

> 以 `router/router.go` 为准。

### Public（`/blog`）

| Method | Path | 说明 |
|---|---|---|
| `ANY` | `/blog/collect` | 采集统计数据 |
| `GET` | `/blog/subscribe` | 订阅 |
| `GET` | `/blog/verify` | 获取/校验验证码 |

### Admin（`/admin`）

| Method | Path | Auth |
|---|---|---|
| `POST` | `/admin/login` | 否 |
| `GET` | `/admin/dashboard/summary` | 是 |
| `GET` | `/admin/dashboard/trend` | 是 |
| `GET` | `/admin/dashboard/insights` | 是 |
| `GET` | `/admin/accesslog/logs` | 是 |
| `GET` | `/admin/performance/averageDelay` | 是 |
| `GET` | `/admin/performance/slowPages` | 是 |
| `GET` | `/admin/analysis/metrics` | 是 |
| `GET` | `/admin/analysis/trend` | 是 |
| `GET` | `/admin/analysis/rank` | 是 |
| `GET` | `/admin/analysis/path` | 是 |
| `GET` | `/admin/analysis/source` | 是 |
| `GET` | `/admin/analysis/querypath` | 是 |
| `GET` | `/admin/analysis/pathDetail/trend` | 是 |
| `GET` | `/admin/analysis/pathDetail/metric` | 是 |
| `GET` | `/admin/analysis/pathDetail/source` | 是 |
| `GET` | `/admin/analysis/pathDetail/device` | 是 |
| `GET` | `/admin/visitormap/map` | 是 |
| `GET` | `/admin/visitormap/chineseMap` | 是 |
| `GET` | `/admin/discussionmap/metric` | 是 |
| `GET` | `/admin/discussionmap/trend` | 是 |
| `GET` | `/admin/discussionmap/activeuser` | 是 |
| `GET` | `/admin/discussionmap/feed` | 是 |

### WebSocket（`/admin/ws`）

| Method | Path | Auth |
|---|---|---|
| `GET` | `/admin/ws?token=<jwt>` | 是（query token） |

### Webhook（`/webhook`）

| Method | Path | 说明 |
|---|---|---|
| `POST` | `/webhook/github` | GitHub webhook（discussion/discussion_comment） |
| `POST` | `/webhook/notify` | GitHub webhook（push 通知订阅用户） |

## 定时任务

定时任务在启动时由 `task.InitCron(comps)` 注册并启动（见 `internal/task/cron.go`）。具体任务实现位于 `internal/job/*`（如 sync / deadlink / discReport 等）。

## 测试与质量

```bash
go test ./...
```

## 部署说明

本项目可配合 GitHub Actions 构建并部署到服务器，通过 Nginx 反向代理对外提供服务。

## 常见问题

- 启动时报数据库/Redis 未配置：请检查 `PG_URI` / `REDIS_URL`
- GeoIP 或 IP2Region 文件路径错误：请检查 `GEODB_PATH`、`IP2REGION_V4_PATH`、`IP2REGION_V6_PATH`
- Webhook 返回 401/403：请确认 `X-Hub-Signature-256` 与对应 secret 是否一致
- WebSocket 连接失败：请确认 `BASE_URL` / `ADMIN_URL` 是否加入了 Origin 白名单，以及 `token` 是否有效
