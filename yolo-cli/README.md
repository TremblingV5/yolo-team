# yolo-cli

yolo-team 的 **唯一代码目录与编译入口**，使用 Go + Cobra 构建。

## 职责

`yolo-cli` 包含项目的全部 Go 代码，编译为唯一产物 `yolo.exe`，集所有能力于一身：

1. **CLI 工具** — 项目、任务卡片、文档的命令行操作
2. **HTTP 服务** — 内置 REST API（handler/service/repository 分层）
3. **静态资源托管** — 通过 Go `embed` 将前端构建产物嵌入二进制，统一端口提供服务

## 技术栈

| 组件 | 选择 |
|------|------|
| 语言 | Go 1.21+ |
| CLI 框架 | [Cobra](https://github.com/spf13/cobra) |
| Web 框架 | Gin |
| API 文档 | swaggo/swag（OpenAPI 3.0） |
| 数据库 | SQLite / MySQL（GORM 切换） |
| ORM | GORM + gorm/gen |
| 前端嵌入 | `embed` (Go 1.16+) 标准库 |
| HTTP 客户端 | `net/http` + 自定义封装 |

## 构建管线

```
yolo-client/dist/ ──► 复制到 client-dist/ ──► go:embed ──► go build ──► yolo.exe
```

## 命令概览

```bash
# === 服务管理 ===
yolo serve                       # 启动 HTTP 服务（API + 前端页面）

# === 项目管理 ===
yolo project list                # 查看所有项目
yolo project info <key>          # 查看项目详情

# === 任务卡片（Issue） ===
yolo issue list -p <project_id>          # 查看项目下任务列表
yolo issue list -p <pid> -e <executor>   # 按执行人筛选
yolo issue create -t <title> -p <pid>    # 创建顶级任务
yolo issue create -t <title> -p <pid> --executor <id>  # 创建并指定执行人
yolo issue create -t <title> -p <pid> --parent <id>  # 创建子任务
yolo issue create -t <title> -p <pid> --deadline "2026-06-30"  # 创建带截至时间的任务
yolo issue info <id>                      # 查看任务详情
yolo issue update <id> -s <status>        # 更新任务状态
yolo issue update <id> --executor <id>    # 更新执行人
yolo issue update <id> --deadline "2026-06-30"  # 更新截至时间
yolo issue update <id> --repo-url <url> --repo-name <n> --branch <b>  # 关联代码仓库
yolo issue todo -e <executor>             # 查看待办列表
yolo issue todo -e <executor> -n 10       # 查看待办列表（指定数量）

# === 执行人 ===
yolo executor list                        # 查看所有执行人

# === 项目文档 ===
yolo doc list -p <project_key>            # 查看项目文档列表
yolo doc create -p <project_key> -t <title>  # 创建项目文档（返回本地路径）
yolo doc info <key>                       # 查看文档详情（含本地路径）

# === 全局选项 ===
--server <url>                   # 指定服务端地址（默认 http://localhost:8080）
```

## 卡片状态

任务卡片固定八种状态：

```
已创建 → 方案设计 → 方案评审 → 代码实现 → QA质检 → 待审查 → 已完成 → 已归档
```

| 状态值 | 显示名 | 说明 |
|--------|--------|------|
| `created` | 已创建 | — |
| `design` | 方案设计 | — |
| `review` | 方案评审 | CLI 不可设置 |
| `implementation` | 代码实现 | — |
| `qa` | QA质检 | — |
| `pending_review` | 待审查 | CLI 不可设置 |
| `done` | 已完成 | 仅 `leader` 可推进至此 |
| `archived` | 已归档 | 仅页面操作，CLI 不可设置 |

## `yolo serve` 详细说明

启动后，一个端口提供所有服务：

```
http://localhost:8080/
├── /api/v1/*    → REST API（内置 handler）
├── /            → 嵌入式前端 SPA（来源于 client-dist/）
├── /*           → 兜底到 index.html（SPA 路由）
```

支持参数：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--port` | `8080` | 监听端口 |
| `--db` | `{workspace}/yolo.db` | SQLite 数据库文件路径 |

## 开发

### 纯 CLI 开发（无需前端）

```bash
cd yolo-cli
go run . project list --server http://localhost:8080
```

### 带前端的生产构建

```bash
# 1. 先构建前端
cd ../yolo-client && pnpm build

# 2. 复制产物到 embed 目录
cp -r dist ../yolo-cli/client-dist/

# 3. 编译
cd ../yolo-cli && go build -o yolo.exe .

# 4. 启动
./yolo.exe serve
```

### 开发模式（前后端分离）

```bash
# 终端 1
cd yolo-cli && go run . serve

# 终端 2
cd yolo-client && pnpm dev
# 访问 http://localhost:8000（热更新，API 代理到 :8080）
```

## 目录结构

```
yolo-cli/
├── main.go              # 程序入口
├── embed.go             # //go:embed client-dist/*
├── client-dist/         # 前端构建产物（gitignore，构建时生成）
├── cmd/                 # Cobra 命令定义
│   ├── root.go
│   ├── serve.go         # yolo serve
│   ├── project.go
│   ├── issue.go
│   ├── executor.go
│   └── doc.go
├── internal/
│   ├── handler/         # HTTP 处理器（禁止直接操作 DB）
│   ├── service/         # 业务逻辑、校验、状态流转
│   ├── repository/      # 数据访问（使用 gorm/gen DAO）
│   ├── model/           # 数据模型 (Project, Issue, Document, Executor)
│   ├── query/           # gorm/gen 自动生成的查询代码
│   ├── router/          # 路由注册
│   └── config/          # 服务配置
├── docs/                # swaggo 生成的 OpenAPI 文档
├── go.mod
```
