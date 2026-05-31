# Yolo-Team 系统架构设计

## 1. 项目概述

yolo-team 是一套基于"AI 员工"的工作流系统。核心理念是：**让人类员工与 AI 员工在同一套工作流体系中高效协作**。系统以"项目（Project）+ 任务卡片（Issue）"为核心模型，提供任务管理、状态流转、项目文档、代码关联等能力。

## 2. 系统架构总览

### 2.1 单一二进制架构

yolo-team（除 skill 外）最终仅编译产出 **一个二进制文件**：`yolo.exe`（或 `yolo`）。所有功能——CLI 操作、HTTP 服务、Web 前端——均由同一个二进制提供。

```
                         ┌──────────────────────┐
                         │     yolo.exe          │
                         │   (编译后的产物)       │
                         ├──────────────────────┤
                         │                      │
                         │   CLI 模式            │
                         │   yolo project list   │
                         │   yolo issue create   │
                         │   yolo issue update   │
                         │                      │
                         ├──────────────────────┤
                         │                      │
                         │   服务模式             │
                         │   yolo serve          │
                         │   ├─ REST API (:8080) │
                         │   └─ 静态资源 (前端)   │
                         │                      │
                         └──────────┬───────────┘
                                    │
                            ┌───────┴───────┐
                            │   SQLite DB   │
                            └───────────────┘
```

### 2.2 构建管线

```
yolo-client/               yolo-team-skill/
  (React + Umi)              (Markdown/配置)
       │                          │
       │ pnpm build               │
       ▼                          │
  dist/ (静态文件)                │
       │                          │
       │    复制到                │
       ▼                          │
  yolo-cli/client-dist/           │
       │                          │
       │ go:embed                 │
       ▼                          │
  yolo-cli/ (Go + Cobra)          │
  └─ go build → yolo.exe         │
                                  │
                           独立的文档体系
```

### 2.3 各子系统职责

| 子系统 | 类型 | 职责 |
|--------|------|------|
| `yolo-cli` | Go | **唯一代码目录**。包含 CLI、HTTP 服务、数据访问等全部后端代码，编译为唯一二进制 |
| `yolo-client` | React + Umi 源码 | Web 前端源码。编译为静态文件后，通过 `go:embed` 嵌入到 `yolo-cli` |
| `yolo-team-skill` | Markdown/配置 | AI 员工技能定义与行为配置，独立文档 |

### 2.4 双模式运行

| 模式 | 命令 | 说明 |
|------|------|------|
| CLI 模式 | `yolo project list` / `yolo issue create` 等 | 命令行操作项目和任务，内置 HTTP 客户端直连本地 server 或远程 server |
| 服务模式 | `yolo serve` | 启动嵌入式 HTTP 服务，同时提供 API 接口和前端页面 |

## 3. 前端嵌入方案

### 3.1 原理

使用 Go 1.16+ 的 `embed` 包，在编译时将 `yolo-client` 的构建产物嵌入到二进制中。

```go
//go:embed client-dist
var clientFS embed.FS
```

### 3.2 开发模式 vs 生产模式

| 阶段 | 前端运行方式 | 说明 |
|------|-------------|------|
| 开发 | `pnpm dev` 独立启动（:8000），代理 API 到 :8080 | 热更新，开发体验好 |
| 生产构建 | 前端编译为静态文件 → `go:embed` 嵌入 → `go build` | 单文件部署，无需额外前端服务 |

### 3.3 `yolo serve` 路由设计

```
:8080/
├── /api/v1/*         → REST API
├── /swagger/*         → OpenAPI 文档（swaggo）
├── /                 → 嵌入式前端 SPA (index.html)
├── /umi.js           → 嵌入式前端 JS bundle
├── /umi.css          → 嵌入式前端 CSS
└── /*               → 兜底到 index.html（SPA 路由）
```

## 4. 数据模型

> 完整的模型定义（含属性表、状态枚举、数据库 DDL、Key 生成策略）见 [models.md](./models.md)。

### 4.1 实体关系

```
Project ──1:N──► Issue ──1:N──► Issue（子卡片，仅一级）
    │                 │
    │                 └──► Executor（执行人）
    │
    └──1:N──► Document（元数据存 DB，正文存文件系统）
```

### 4.2 核心模型摘要

| 模型 | 表名 | 唯一标识 | 说明 |
|------|------|----------|------|
| Project | `projects` | `YOLO-PROJECT-{id}` | 项目 |
| Issue | `issues` | `YOLO-ISSUE-{id}` | 任务卡片，8 种固定状态 |
| Executor | `executors` | 自增 id，name 唯一 | 执行人，含 role（leader/architect/developer/qa） |
| Document | `project_documents` | `YOLO-DOC-{id}` | 项目文档（正文存文件系统） |

### 4.3 Issue 状态流转

```
已创建 → 方案设计 → 方案评审 → 代码实现 → QA质检 → 待审查 → 已完成 → 已归档
```

关键限制：
- `review`（方案评审）和 `pending_review`（待审查）CLI 不可设置
- 推进至 `done` 仅允许 `leader` 角色
- `archived` 仅人类通过页面操作
- 完整状态转换权限表见 [models.md](./models.md) 3.2 节

### 4.4 数据库

使用 **GORM** 作为 ORM，底层支持 **SQLite** 和 **MySQL** 两种存储后端（通过配置切换）。默认使用 SQLite。

数据库迁移使用 **GORM AutoMigrate**，每次 `yolo serve` 启动时自动同步表结构。

查询层使用 **gorm/gen** 自动生成类型安全的 DAO 代码，存放于 `internal/query/`。

### 4.5 分层设计

严格三层架构，**禁止越层调用**：

```
handler ──► service ──► repository (GORM/gen)
  │           │            │
  │           │            └── 数据库访问，使用 gorm/gen 生成的 DAO
  │           └── 业务逻辑、校验、状态流转规则
  └── HTTP 请求/响应处理、参数绑定、调用 service
```

**硬约束**：
- Handler 层**禁止**直接操作数据库
- Repository 层**禁止**包含业务逻辑
- Service 层是唯一的业务逻辑入口

### 4.6 工作目录

首次运行时，用户必须通过 `yolo init` 设定工作目录。工作目录路径存储在 `~/.yolo-team/settings.json` 中：

```json
{
  "workspace": "/home/user/yolo-workspace"
}
```

所有数据均保存在工作目录下：

```
{workspace}/
├── yolo.db                     # SQLite 数据库
└── documents/
    └── {project-id}/
        └── {title}.md          # 文档正文
```

## 6. 接口与命令

### 6.1 API 定义

> 完整的接口定义（含入参、响应结构、错误码）见 [api.md](./api.md)。

主要资源端点：

- `projects` — 项目 CRUD
- `issues` — 任务卡片 CRUD，含待办查询
- `executors` — 执行人管理
- `documents` — 项目文档管理（内容存储于文件系统）

### 6.2 CLI 命令

> 完整的命令定义（含参数、选项、示例输出）见 [cli.md](./cli.md)。

命令分组：

- `yolo init` — 初始化工作目录（首次运行必须）
- `yolo serve` — 启动 HTTP 服务
- `yolo project` — 项目管理
- `yolo issue` — 任务卡片管理
- `yolo executor` — 执行人管理
- `yolo doc` — 项目文档管理

### 6.3 状态变更规则

卡片状态变更需遵循流转顺序，只允许**向前流转**（或回到已创建）：

| 当前状态 | 允许变更为 |
|----------|-----------|
| `created` | `design` |
| `design` | `review` |
| `review` | `implementation`、`design` |
| `implementation` | `qa` |
| `qa` | `pending_review`、`implementation` |
| `pending_review` | `done`、`design`、`implementation`、`qa` |
| `done` | `archived` |
| `archived` | （不可变更） |

**角色限制**：完整状态转换权限表见 [models.md](./models.md) 3.2 节。`review` 和 `pending_review` CLI 不可设置。

## 7. 技术选型

| 层级 | 技术 | 选型理由 |
|------|------|----------|
| CLI 框架 | Cobra | Go 生态最成熟的 CLI 框架 |
| Web 框架 | Gin | 轻量高性能，配合 swaggo/swag 生成 OpenAPI |
| API 文档 | swaggo/swag | 通过注释自动生成 OpenAPI 3.0 规范 |
| 前端嵌入 | `embed` (Go 1.16+) | 标准库原生支持，零依赖 |
| 数据库 | SQLite / MySQL（通过 GORM 切换） | GORM + gorm/gen 自动生成查询代码 |
| ORM | GORM + gorm/gen | 类型安全查询，减少手写 SQL |
| 前端框架 | React + Umi | 企业级前端框架，开箱即用 |
| 前端 UI | Ant Design | 丰富的企业级组件库 |
| 前端 API | restful-react | 从 OpenAPI 自动生成 TypeScript 类型和请求函数 |

## 8. 目录结构规划

```
yolo-team/
├── README.md
├── docs/
│   ├── architecture.md              # 本文件（系统架构）
│   ├── models.md                     # 数据模型定义
│   ├── api.md                        # API 接口定义
│   ├── cli.md                        # CLI 命令定义
│   └── use-cases.md                  # 前端页面与用例
├── yolo-cli/                        # ★ 唯一代码目录
│   ├── main.go                      # 程序入口
│   ├── cmd/                         # Cobra 命令 (project, issue, doc, executor, serve)
│   │   ├── root.go
│   │   ├── serve.go                 # yolo serve 命令
│   │   ├── project.go
│   │   ├── issue.go
│   │   ├── executor.go
│   │   └── doc.go
│   ├── internal/
│   │   ├── handler/                 # HTTP 处理器（禁止直接操作 DB）
│   │   ├── service/                 # 业务逻辑、校验、状态流转
│   │   ├── repository/              # 数据访问（调用 gen DAO）
│   │   ├── model/                   # 数据模型
│   │   ├── query/                   # gorm/gen 生成的查询代码
│   │   ├── router/                  # 路由定义
│   │   └── config/                  # 服务配置
│   ├── client-dist/                 # ★ 前端构建产物（通过 embed 嵌入）
│   ├── embed.go                     # //go:embed 声明
│   └── go.mod
├── yolo-client/                     # 前端源码（开发用）
│   ├── src/
│   │   ├── pages/
│   │   ├── components/
│   │   ├── services/
│   │   └── models/
│   ├── dist/                        # 构建产物 → 复制到 yolo-cli/client-dist/
│   ├── package.json
│   └── .umirc.ts
└── yolo-team-skill/
    └── README.md
```
