# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Yolo-Team 是基于 **AI 员工** 的工作流系统，以"项目（Project）+ 任务卡片（Issue）"为核心模型，让人类员工和 AI 员工在同一套工作流体系中协作。整个项目编译为**单一二进制文件** `yolo.exe`，前端通过 `go:embed` 嵌入。

## 常用命令

### 后端（yolo-cli）

```bash
# 开发模式启动服务
cd yolo-cli && go run ./cmd/yolo serve

# 编译二进制
cd yolo-cli && go build -o yolo.exe ./cmd/yolo

# 生成 gorm/gen 类型安全查询代码（修改 model 后必须执行）
cd yolo-cli && make gen

# 生成 Swagger/OpenAPI 文档
cd yolo-cli && make swag

# 完整构建（前端+后端）
cd yolo-cli && make build
```

### 前端（yolo-client）

```bash
# 开发模式（端口 8000，API 代理到 8080）
cd yolo-client && pnpm dev

# 构建（产物直接输出到 yolo-cli/internal/frontend/client-dist/）
cd yolo-client && pnpm build

# 从 OpenAPI 生成 TypeScript API 客户端
cd yolo-client && pnpm run generate-api
```

### 测试（yolo-cli-test）

```bash
# 运行全部测试（会自动编译二进制，启动独立工作空间）
cd yolo-cli-test && go test ./testcases/ -v

# 运行单个测试
cd yolo-cli-test && go test ./testcases/ -v -run TestProjectCreate

# 测试框架：编译 yolo 二进制 → 创建临时工作空间 → 通过 CLI 子进程执行命令并校验输出
```

## 架构

### 单一二进制架构

```
yolo-client (React) ──pnpm build──► client-dist/ ──go:embed──► yolo-cli (Go) ──go build──► yolo.exe
```

运行模式：
- **服务模式** `yolo serve` — 启动 HTTP 服务（默认 :8080），同时提供 REST API 和嵌入式前端 SPA
- **CLI 模式** `yolo project list` / `yolo issue create` 等 — 通过内置 HTTP 客户端操作

### 后端三层架构（严格分层）

```
handler/ → repository/ → db (GORM/gen)
```

注意：**当前代码中无独立 service 层**，业务逻辑（校验、状态流转）分散在 handler 和 model 中。Handler 直接调用 Repository。

- **handler/** — HTTP 处理器，使用泛型 `Wrap[R, T]` 函数统一处理请求绑定和响应。Handler 方法签名：`func(ctx context.Context, req R) (T, error)`
- **repository/** — 数据访问层，使用 `db.G[Model]()` + gorm/gen 生成的类型安全查询（`internal/query/`）
- **model/** — 数据模型，包含 `New*()` 构造函数和 `Validate()` 校验方法
- **common/** — 统一响应格式 `Response{Code, Data, Message}`，`AppError` 错误类型（Code: 40001=参数错误, 40401=未找到, 40901=冲突, 50001=内部错误）
- **config/** — 配置管理，工作空间路径存储在 `~/.yolo-team/settings.json`
- **db/** — SQLite 初始化和 GORM 注册，AutoMigrate 自动建表

### 前端路由

| 路径 | 页面 |
|------|------|
| `/` | KanbanPage（看板） |
| `/project/:key/docs` | DocumentList（项目文档） |
| `/doc/:key` | DocumentDetail（文档详情） |
| `/executors` | ExecutorManage（执行人管理） |

### 数据存储

- **数据库**: SQLite，位于 `{workspace}/yolo.db`
- **文档内容**: 文件系统，位于 `{workspace}/documents/{project_id}/{title}.md`
- **工作空间配置**: `~/.yolo-team/settings.json`（通过 `yolo init` 初始化）

### 核心数据模型

- **Project** — Key 格式 `YOLO-PROJECT-{id}`
- **Issue** — Key 格式 `YOLO-ISSUE-{id}`，四种状态：`created` → `in_progress` → `done` → `archived`。所有 Issue 平级，不支持父子关系
- **Task** — Key 格式 `YOLO-TASK-{id}`，Issue 下的子任务，三种状态：`created` → `in_progress` → `done`，通过 `issue_id` 关联 Issue。所有交互以 Key 为主
- **Executor** — 角色：`leader` / `architect` / `developer` / `qa`
- **Document** — Key 格式 `YOLO-DOC-{id}`，元数据在 DB，正文在文件系统

### API 路由

所有 API 在 `/api/v1` 下：`/projects`、`/issues`、`/executors`、`/documents`。路由定义在 `internal/router/router.go`。

## 关键约定

- 修改 `internal/model/` 后必须运行 `make gen` 重新生成 `internal/query/` 的 DAO 代码
- 前端构建产物输出到 `yolo-cli/internal/frontend/client-dist/`（由 vite.config.ts 配置）
- 前端 API 客户端代码 `src/generated.tsx` 由 `restful-react` 从 swagger.json 生成，修改 API 后需运行 `pnpm run generate-api`
- Handler 使用 Swagger 注释（`@Summary`、`@Router` 等）生成 OpenAPI 文档，修改 API 后需运行 `make swag`
- 测试套件（`yolo-cli-test/`）通过编译二进制并以子进程方式执行 CLI 命令来验证功能，每个测试函数有独立的工作空间隔离

## REST API 约定

- **不使用 DELETE 方法**。`restful-react` 对返回 `struct{}`（空响应体）的 DELETE 端点存在代码生成缺陷：路径参数会被错误生成为请求体

  示例：handler 定义 `DELETE /api/v1/documents/{key}` → restful-react 生成 `POST /api/v1/documents` + body: `key`

  改为使用 POST：
  ```go
  // ❌ 错误
  documents.DELETE("/:key", handler.Wrap(docH.Delete))
  // ✅ 正确
  documents.POST("/:key", handler.Wrap(docH.Delete))
  ```

- `fetch` 只能出现在 `yolo-cli-test/`（测试套件）中，前端代码**禁止**使用裸 `await fetch()`，必须使用 `generated.tsx` 的 hook 或 `restful-react` 的 `useMutate`/`useGet`

- `restful-react` 的 `useMutate` 返回的 `mutate` 函数签名是 `(body, options?)`，**不是** `(body, queryParams, pathParams)`。覆盖 pathParams 的正确写法：
  ```ts
  // ❌ 错误
  mutate(body, undefined, { key: 'val' })
  // ✅ 正确
  mutate(body, { pathParams: { key: 'val' } })
  ```
