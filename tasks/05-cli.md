# 05 — CLI 命令

> 参考：[cli.md](../docs/cli.md)

## 子任务
### 5.1 框架搭建

- [x] `cmd/root.go` — Cobra 根命令，`--server` 全局选项（默认 `http://localhost:8080`），`--json` 输出切换
- [x] 实现变更为：HTTP 客户端直接封装在 `cmd/client.go`（`apiGet`/`apiPost`/`apiPut`/`apiDelete`），未单独创建 `internal/client/http.go`
- [x] 自动处理 JSON 序列化、错误响应解析

### 5.2 初始化命令
- [x] `yolo init --workspace <path>`

### 5.3 服务命令

- [x] `cmd/serve.go` — `yolo serve`（启动全栈 HTTP 服务，含数据库初始化）

### 5.4 项目命令（`cmd/project.go`）
- [x] `yolo project list` — 表格输出 id/key/name/description/updated
- [x] `yolo project info <key>` — 详情表格

> 不含 `yolo project create`（项目仅页面创建）

### 5.5 Issue 命令（`cmd/issue.go`）
- [x] `yolo issue list [-p project_id] [-s status] [-e executor_id] [--parent 0]`
  - 表格输出：KEY / STATUS / PRIORITY / TITLE / DEADLINE / REPO
- [x] `yolo issue create -t title -p project_id [--description] [--priority] [--parent] [--executor] [--deadline]`
  - 输出：`Issue created: YOLO-ISSUE-X (标题)`
- [x] `yolo issue info <key>` — 详情表格
  - 包含：Key/Status/Title/Description/Priority/Deadline/Repo
- [x] `yolo issue update <key> [-s status] [-t title] [--executor] [--deadline] [--repo-url/name/branch]`
  - CLI 更新带 `X-Yolo-CLI` header 自动触发权限校验
  - 状态受流转规则约束，CLI 无法设置 review/pending_review/archived
- [x] `yolo issue delete <key>` — 删除
- [x] `yolo issue todo -e executor_id [-n limit]` — 待办查询
  - 按 deadline ASC + priority DESC 排序

> 不含 `yolo issue done`（使用 `update -s done` 替代）

### 5.6 执行人命令（`cmd/executor.go`）
- [x] `yolo executor list` — 表格输出 id/name/role/created

> 不含 `yolo executor create`（仅页面创建）

### 5.7 文档命令（`cmd/executor.go` 中定义 `docCmd`）
- [x] `yolo doc list -p <project_key>` — 表格输出 key/title/updated
- [x] `yolo doc create -p <project_key> -t <title>`
  - 创建空文档（仅标题），输出：`Document created: YOLO-DOC-X` + `Path: ...`
- [x] `yolo doc info <key>` — 输出文件路径 + Markdown 内容
- [x] 实现变更为：doc 命令定义在 `cmd/executor.go` 中（未独立为 `cmd/doc.go`）

> 不含 `yolo doc update`（AI 通过文件路径直接操作）

### 5.8 输出格式

- [x] 默认：可读表格
- [x] `--json`：原始 JSON
- [ ] 错误输出到 stderr，非零退出码

## 验收标准

- [x] 所有命令可正常调用，与 API 交互正确
- [x] 表格输出对齐，信息完整（含仓库列）
- [x] `--json` 输出合法 JSON
- [x] 错误场景正确提示（参数错误、API 错误、权限不足）
- [x] 文档 key 参数使用正确（project key、doc key）
