# Blog-Backend

![go](https://img.shields.io/badge/Go-1.24.0?logo=go&logoColor=white)
![postgresql](https://img.shields.io/badge/PostgreSQL-17.6?logo=postgresql&logoColor=white)
![redis](https://img.shields.io/badge/Redis-8.0.5?logo=redis&logoColor=white)
![nginx](https://img.shields.io/badge/Nginx-1.26.3?logo=nginx&logoColor=white)
![gin](https://img.shields.io/badge/Gin-1.11.0)
![gorm](https://img.shields.io/badge/Gorm-1.31.1)
![release](https://img.shields.io/github/v/release/xianbo-deep/BLOG-BACKEND?label=release)


这是一个基于 Go 语言和 Gin 框架开发的博客后台服务系统。主要用于博客站点的流量统计、性能监控以及后台管理功能。

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

## 功能模块

### 1. 公共接口 (Public)
-   **数据采集 (`/blog/collect`)**: 接收前端上报的访问数据，包括时间戳、路径、延迟、IP、User-Agent 等信息。

### 2. 管理后台 (Admin)
后台接口需要 JWT 认证。

-   **仪表盘 (`/admin/dashboard`)**:
    -   数据概览 (Summary): 获取总访问量、今日访问量等概览数据。
    -   访问趋势 (Trend): 展示近期的访问趋势图表数据。
    -   洞察分析 (Insights): 提供基于数据的深度分析。
-   **访问日志 (`/admin/accesslog`)**:
    -   日志查询 (`/logs`): 分页查询详细的访客记录。
-   **性能监控 (`/admin/performance`)**:
    -   平均延迟 (`/averageDelay`): 统计页面的平均加载延迟。
    -   慢页面分析 (`/slowPages`): 识别加载速度较慢的页面。
-   **页面分析 (`/admin/analysis`)**:
    -   全站统计 (`/total`): 全站页面的访问统计。
    -   今日统计 (`/today`): 今日页面的访问统计。
-   **访客地图 (`/admin/visitormap`)**:
    -   世界地图 (`/map`): 基于 IP 的全球访客分布。
    -   中国地图 (`/chineseMap`): 基于 IP 的中国访客分布。

### 3. 定时任务
-   **数据同步**: 每日凌晨自动将 Redis 中的缓存数据同步到 PostgreSQL 数据库，确保数据持久化。

## 项目结构

```text
Blog-Backend/
├── api/            # Serverless 入口 (如 Vercel)
├── consts/         # 常量定义 (环境配置, 业务常量)
├── core/           # 核心组件初始化 (数据库, Redis, GeoIP)
├── dto/            # 数据传输对象 (Request/Response 定义)
├── internal/       # 内部业务逻辑
│   ├── controller/ # 控制器层 (处理 HTTP 请求)
│   ├── dao/        # 数据访问层 (数据库操作)
│   ├── service/    # 业务逻辑层
│   └── task/       # 定时任务 (数据同步)
├── middleware/     # 中间件 (认证, CORS)
├── model/          # 数据库模型定义
├── router/         # 路由配置
├── thirdparty/     # 第三方服务集成 (GitHub)
└── utils/          # 工具函数 (GeoIP, JWT, 分页)
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