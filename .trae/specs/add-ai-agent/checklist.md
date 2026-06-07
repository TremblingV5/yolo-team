# Checklist

## AI 底座
- [x] AI 配置读取（环境变量 > settings.json.ai 段）代码实现与规范一致
- [x] `Settings` 结构体包含 `AI` 配置字段
- [x] langchaingo Agent 封装代码实现与规范一致（LLM 初始化、工具注册）
- [x] 所有 26 个工具注册代码实现与规范一致
- [x] 工具直接调用 Repository 层（而非子进程 CLI）
- [x] 会话管理器支持 session_id 隔离上下文
- [x] 会话持久化到 workspace/sessions/ 目录
- [x] sessions.json 索引 + .jsonl 消息记录格式正确
- [x] 会话 Key 格式 YOLO-CHAT-{YYYYMMDD}-{NNN}
- [x] 标签检测正则表达式 `/<\w+\s*[^>]*\/\s*>/g` 实现正确

## Chat API
- [x] POST /api/v1/ai/chat SSE 流式响应实现正确
- [x] GET /api/v1/ai/sessions 返回会话列表
- [x] DELETE /api/v1/ai/sessions/:session_id 删除会话
- [x] POST /api/v1/ai/sessions/reset 清除所有会话
- [x] 跨域支持（CORS）已配置（Gin default）

## 路由
- [x] AI 相关路由已全部注册
- [x] SSE 端点 Content-Type 配置为 `text/event-stream`

## TUI Chat (yolo chat)
- [x] `yolo chat` 命令可正常启动全屏 TUI（顶层命令）
- [x] 左侧会话列表 + 右侧消息列表 + 底部输入框布局正确
- [x] Enter 发送、Tab 切换焦点
- [x] Ctrl+C 退出、Ctrl+L 清屏、Ctrl+N 新建会话
- [x] `/sessions` 列出并选择会话、`/clear` 清空会话
- [x] 会话持久化（复用 internal/ai.SessionManager）

## 前端 AI Chat (Ant Design X)
- [x] `/ai` 页面可访问
- [x] 使用 `@ant-design/x` Bubble.List 组件显示聊天消息
- [x] SSE 流式接收和渲染正确
- [x] 侧边栏显示会话历史列表（Conversations 组件）
- [x] 导航栏有 AI Chat 入口（TopBar 添加 RobotOutlined 按钮）

## 构建与测试
- [x] `go build -o yolo.exe ./cmd/yolo` 编译无错误
- [x] `pnpm build` 前端构建无错误
- [x] `go vet ./internal/...` 通过
- [x] `yolo.exe` 二进制成功生成
