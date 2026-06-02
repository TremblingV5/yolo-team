# 07 — 嵌入与构建

> 参考：[architecture.md](../docs/architecture.md) 3 节

## 子任务

### 7.1 前端构建嵌入

- [x] `yolo-cli/client-dist/` 目录（含 .gitkeep 占位）
- [x] `embed.go`（`//go:embed client-dist/*`）使用 `fs.Sub` 挂载为静态文件系统
- [x] `vite.config.ts` 配置 `outDir: '../yolo-cli/client-dist'`，前端构建直接输出到嵌入目录

### 7.2 `yolo serve` 命令

- [x] `cmd/serve.go` 启动 HTTP 服务 + 数据库初始化
- [x] API 路由（`internal/router`）
- [x] 嵌入式前端静态文件挂载到 `/`（`r.StaticFS("/", http.FS(staticFS))`）
- [x] SPA 路由回退：`NoRoute` 处理器将非 `/api` 路径回退到 `index.html`
- [x] 支持参数：`--port`（默认 8080）、`--db`（默认 `yolo.db`）

### 7.3 开发模式

- [x] 后端开发：`cd yolo-cli && go run ./cmd/yolo serve`
- [x] 前端开发：`cd yolo-client && npm run dev`（Vite dev server :8000，proxy → :8080）
- [x] 纯 CLI 开发：直接运行 `go run ./cmd/yolo <command>`

### 7.4 生产构建

- [x] Makefile：`make build` 一键构建（先 `build-frontend` 输出到 `client-dist`，再 `build-backend` 编译 Go 二进制）
- [x] 最终产物 `yolo.exe` 可独立运行，前端已嵌入二进制

### 7.5 SPA 路由处理

- [x] `r.StaticFS("/", http.FS(staticFS))` 服务静态资源（JS/CSS）
- [x] `NoRoute` 中间件：非 `/api` 路径返回 `client-dist/index.html`（SPA 支持）

## 验收标准

- [x] `yolo serve` 启动后可通过 `localhost:8080` 访问前端
- [x] `/api/v1/*` 路径正常响应 JSON
- [x] SPA 路由刷新不 404
- [x] `yolo.exe` 可脱离源码目录独立运行
- [x] 开发模式前后端分离正常工作
