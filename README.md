# Blog-Backend

<div style="text-align: center;">

![Go](https://img.shields.io/badge/Go-1.24.0-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17.6-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-8.0.5-DC382D?logo=redis&logoColor=white)
![Nginx](https://img.shields.io/badge/Nginx-1.26.3-009639?logo=nginx&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.11.0-00ADD8?logo=go&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-1.31.1-00ADD8?logo=go&logoColor=white)
![release](https://img.shields.io/github/v/release/xianbo-deep/BLOG-BACKEND?label=release&include_prereleases)

</div>





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


## 部署说明

本项目使用Github Action执行自动化脚本，将推送后的代码自动编译、部署到服务器