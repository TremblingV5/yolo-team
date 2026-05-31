# 01 — 数据库与模型

> 参考：[models.md](../docs/models.md)

## 子任务

### 1.1 项目初始化

- [x] 初始化 `yolo-cli/go.mod`（module `yolo-team/yolo-cli`）
- [x] 安装依赖：`cobra`、`gin`、`gorm`、`gorm.io/driver/sqlite`
- [x] 创建目录结构：`cmd/yolo/`、`internal/{handler,repository,model,query,router,config,common,utils}`
- [x] 实现变更为：使用 `gorm cli`（`gorm.io/cli/gorm`）替代 `gorm/gen` 生成查询代码，产物在 `internal/query/`
- [x] 实现变更为：删除了 `service` 层，handler 直接调用 repository
- [ ] ~~配置 GORM 支持切换 SQLite / MySQL~~（仅使用 SQLite）

### 1.2 数据库迁移（GORM AutoMigrate）

- [x] 在 `internal/db/db.go` 中调用 `DB.AutoMigrate()` 注册所有模型
- [x] 每次启动 `yolo serve` 时自动执行表结构同步，无需手动迁移

### 1.3 Go 模型定义

- [x] `internal/model/project.go` — Project 结构体
- [x] `internal/model/issue.go` — Issue 结构体（含 Children []Issue）
- [x] `internal/model/executor.go` — Executor 结构体
- [x] `internal/model/document.go` — Document 结构体
- [x] 所有结构体带 JSON tag（snake_case），time 类型用 `time.Time`

### 1.4 Key 生成

- [x] 实现 Key 生成函数在 `internal/utils/keyutils/`：
  - `keyutils.Project(id)` → `YOLO-PROJECT-{id}`
  - `keyutils.Issue(id)` → `YOLO-ISSUE-{id}`
  - `keyutils.Document(id)` → `YOLO-DOC-{id}`
- [x] Repository 层实现"先 insert → 取 id → update key"策略

### 1.5 Repository 层 + GORM CLI 查询

- [x] 使用 `gorm cli`（`gorm.io/cli/gorm`）从 model 生成 field helpers → `internal/query/`
- [x] `internal/repository/project_repo.go` — 项目数据访问（使用 field helpers + `gorm.G[T]()`）
- [x] `internal/repository/issue_repo.go` — 卡片数据访问
- [x] `internal/repository/executor_repo.go` — 执行人数据访问
- [x] `internal/repository/document_repo.go` — 文档数据访问
- [x] 所有 Repo 通过生成的 field helpers 操作数据库，不手写 SQL 字符串
- [x] 通用泛型快捷方法 `db.G[T]()` 封装在 `internal/db/db.go`

## 验收标准

- [x] `go build` 通过
- [x] 数据库文件可创建，4 张表结构正确
- [x] Key 生成格式符合 `YOLO-{TYPE}-{id}` 格式
- [x] UNIQUE 约束（projects.key, issues.key, executors.name, project_documents.key）生效
