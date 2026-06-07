# Tasks

- [x] Task 1: 新增 Go 依赖 & 扩展配置 — 添加 langchaingo、bubbles、lipgloss 依赖，扩展 Settings 结构
  - [x] 运行 `go get github.com/tmc/langchaingo`
  - [x] 运行 `go get github.com/charmbracelet/bubbles`
  - [x] 运行 `go get github.com/charmbracelet/lipgloss`
  - [x] 修改 `internal/config/config.go` — `Settings` 增加 `AI` 配置段（base_url, api_key, model, system_prompt）

- [x] Task 2: 创建 AI 底座模块 — `yolo-cli/internal/ai/` 包
  - [x] 创建 `internal/ai/config.go` — AI 配置读取（环境变量 > settings.json）
  - [x] 创建 `internal/ai/agent.go` — langchaingo Agent 封装（LLM 初始化、工具注册）
  - [x] 创建 `internal/ai/tools.go` — 注册所有 yolo-team 操作工具（调用 repository 层），含 23 + 3 个工具
  - [x] 创建 `internal/ai/session.go` — 会话管理器（文件系统持久化：sessions.json + .jsonl）
  - [x] 创建 `internal/ai/tags.go` — 标签检测和渲染定义（正则匹配、标签结构体定义）

- [x] Task 3: 创建 Chat API Handler — `yolo-cli/internal/handler/ai_handler.go`
  - [x] 实现 SSE 流式响应处理（`gin.Context.Stream`）
  - [x] 处理 `POST /api/v1/ai/chat` 请求（流式聊天）
  - [x] 处理 `GET /api/v1/ai/sessions` 请求（会话列表）
  - [x] 处理 `DELETE /api/v1/ai/sessions/:session_id` 请求（删除会话）
  - [x] 处理 `POST /api/v1/ai/sessions/reset` 请求（重置所有会话）

- [x] Task 4: 注册 AI 路由 — 修改 `internal/router/router.go`
  - [x] 添加 `POST /api/v1/ai/chat` SSE 路由
  - [x] 添加 `GET /api/v1/ai/sessions` 路由
  - [x] 添加 `DELETE /api/v1/ai/sessions/:session_id` 路由
  - [x] 添加 `POST /api/v1/ai/sessions/reset` 路由

- [x] Task 5: 创建 TUI Chat 命令 — `yolo-cli/cmd/chat.go` + `internal/tui/chat.go`
  - [x] 创建 `yolo chat` 子命令（顶层 Cobra 命令）
  - [x] 实现 TUI 主界面（bubbles: 左侧会话列表 + 右侧消息列表 + 底部输入框）
  - [x] 实现与 AI 底座的内联调用
  - [x] 实现快捷键（Ctrl+C 退出、Ctrl+L 清屏、Ctrl+N 新建会话）
  - [x] 实现斜杠命令（`/sessions` 列出会话、`/clear` 清空会话）

- [x] Task 6: 创建前端 AI Chat 页面 — `yolo-client/src/pages/AiChat.tsx`
  - [x] 安装 `@ant-design/x` 依赖
  - [x] 基于 `@ant-design/x` 的 Bubble.List/Sender/Conversations 构建聊天界面
  - [x] 实现 SSE 流式接收和渲染
  - [x] 实现会话侧边栏（列出历史会话，支持删除）
  - [x] 在路由中注册 `/ai` 路径
  - [x] 在导航栏添加 AI Chat 入口

- [x] Task 7: 集成测试和验证
  - [x] 编译 `yolo.exe` 确认无编译错误
  - [x] `pnpm build` 前端构建成功
  - [x] `go vet ./internal/...` 代码静态检查通过
  - [x] 后端 AI 底座模块 Go 代码编译通过
  - [x] 验证所有文件完整性
