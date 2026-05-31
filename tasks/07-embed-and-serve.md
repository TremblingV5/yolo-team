# 07 — 嵌入与构建

> 参考：[architecture.md](../docs/architecture.md) 3 节

## 子任务

### 7.1 前端构建嵌入

- [ ] `yolo-cli/client-dist/` 目录（未创建）
- [ ] `embed.go`（//go:embed 指令）（未实现）
- [ ] 前端构建脚本：`yolo-client/` 下 `pnpm build` 后复制 `dist/` → `../yolo-cli/client-dist/`

### 7.2 `yolo serve` 命令

- [x] `cmd/serve.go` 已实现基本 HTTP 服务启动
- [x] 启动时初始化数据库（调用 `db.Init`）
- [x] 注册 API 路由（`internal/router`）
- [ ] ~~挂载嵌入式前端静态文件到 `/`~~（前端嵌入未实现，前后端分离运行）
- [ ] ~~非 `/api` 路径回退到 `index.html`（SPA 支持）~~
- [x] 支持参数：`--port`（默认 8080）、`--db`（默认 `yolo.db`）

### 7.3 开发模式

- [x] 后端开发：`cd yolo-cli && go run . serve`
- [x] 纯 CLI 开发：直接运行 `go run . <command>`，无需启动 serve
- [ ] 前端开发：需要独立启动 `pnpm dev`（Vite dev server）

### 7.4 生产构建

- [ ] 构建脚本或 Makefile：前后端统一构建流程
- [ ] 最终产物 `yolo.exe` 可独立运行（含嵌入前端）

### 7.5 SPA 路由处理

- [ ] Gin 中间件：对非 `/api` 路径返回 `client-dist/index.html`
- [ ] 静态资源（JS/CSS）正确响应 MIME type

## 验收标准

- [x] `yolo serve` 启动后可通过 `localhost:8080` 访问后端 API
- [x] `/api/v1/*` 路径正常响应 JSON
- [ ] SPA 路由刷新不 404
- [ ] `yolo.exe` 可脱离源码目录独立运行
- [ ] 开发模式前后端分离正常工作
