# AI Agent 智能助手 Spec

## Why

当前 yolo-team 的交互方式局限于：
- 人类用户在 Web UI 手动拖拽看板、填写表单来管理项目和任务
- AI 员工通过 CLI（yolo-team-skill）以子进程方式逐条执行命令

缺少一个**自然语言交互层**，让人类用户可以直接用对话方式了解工作状态、分配任务、推进流转，也让 AI Agent 能够以更智能的方式（而非机械的命令拼接）来操作 yolo-team 系统。

## What Changes

- **新增** `yolo chat` CLI 子命令 — TUI 交互模式，使用 charmbracelet/bubbles 和 lipgloss
- **新增** 前端 AI Chat 页面 — 使用 Ant Design X 的聊天组件
- **新增** AI 底座模块 — 基于 langchaingo 的智能体引擎，提供：
  - LLM 集成（OpenAI-compatible API）
  - 工具注册（映射到 yolo-team 的所有核心操作）
  - 会话管理（上下文记忆 + 文件系统持久化）
  - 标签渲染框架（双端共享标签模板）
- **新增** REST API `POST /api/v1/ai/chat` — SSE 流式聊天接口
- **修改** `yolo serve` — 增加 SSE 支持，注册 AI 路由
- **修改** `settings.json` 结构 — 增加 `ai` 配置段
- **新增** 标签模板系统 — 双端统一解析 `<TagName>` 标签实现特殊展示
- **新增** 会话文件系统管理体系 — workspace 下 `sessions/` 目录管理所有会话记录
- 不破坏现有 CLI、API、前端看板逻辑

## Impact

- Affected specs: 项目/任务/子任务/文档/执行人管理
- Affected code:
  - `yolo-cli/cmd/` — 新增 `chat.go` CLI 命令（`yolo chat`）
  - `yolo-cli/internal/ai/` — 新增 AI 底座模块（langchaingo 集成、工具定义、会话管理、标签渲染）
  - `yolo-cli/internal/handler/` — 新增 `ai_handler.go`（Chat API）
  - `yolo-cli/internal/router/router.go` — 注册 AI 路由
  - `yolo-cli/internal/config/config.go` — `Settings` 增加 `AI` 配置段
  - `yolo-client/` — 新增 AI Chat 页面
  - `go.mod` / `go.sum` — 新增 langchaingo、bubbles、lipgloss 等依赖
  - `package.json` — 新增 `@ant-design/x` 依赖

---

## Requirements

### Requirement: AI 底座（langchaingo）

AI 底座是核心引擎，封装在 `yolo-cli/internal/ai/` 包中，双端共享。

#### Scenario: LLM 集成

- **WHEN** AI 底座启动
- **THEN** 它 SHALL 支持 OpenAI-compatible API 的 LLM 调用
- **AND** 它 SHALL 可通过环境变量 `OPENAI_BASE_URL`、`OPENAI_API_KEY`、`OPENAI_MODEL` 配置
- **AND** 它 SHALL 默认使用 `gpt-4o` 模型
- **AND** 它 SHALL 从 `settings.json` 的 `ai` 配置段读取默认值（环境变量优先级更高）

#### Scenario: 工具注册

- **WHEN** AI 底座初始化
- **THEN** 它 SHALL 注册以下工具到 langchaingo：

| 工具名 | 对应操作 | 说明 |
|--------|---------|------|
| `list_projects` | 列出所有项目 | — |
| `get_project` | 查看单个项目 | 参数: key |
| `create_project` | 创建项目 | 参数: name, description |
| `update_project` | 更新项目 | 参数: key, name?, description? |
| `delete_project` | 删除项目 | 参数: key |
| `list_executors` | 列出所有执行人 | — |
| `create_executor` | 创建执行人 | 参数: name, soul? |
| `update_executor` | 更新执行人 | 参数: name, role?, soul? |
| `delete_executor` | 删除执行人 | 参数: name |
| `list_issues` | 列出任务卡片 | 参数: project_id?, status?, executor_name? |
| `get_issue` | 查看单个卡片 | 参数: key |
| `create_issue` | 创建卡片 | 参数: project_id, title, description?, priority?, executor_name?, deadline? |
| `update_issue` | 更新卡片 | 参数: key, title?, description?, status?, priority?, executor_name?, deadline? |
| `delete_issue` | 删除卡片 | 参数: key |
| `todo_issues` | 待办查询 | 参数: executor_name, limit? |
| `create_task` | 创建子任务 | 参数: issue_key, title, description? |
| `list_tasks` | 列出子任务 | 参数: issue_key |
| `update_task` | 更新子任务 | 参数: task_key, title?, description?, status? |
| `delete_task` | 删除子任务 | 参数: task_key |
| `list_documents` | 列出项目文档 | 参数: project_key |
| `get_document` | 查看文档 | 参数: key |
| `create_document` | 创建文档 | 参数: project_key, title, content, creator? |
| `update_document` | 更新文档 | 参数: key, title?, content? |
| `delete_document` | 删除文档 | 参数: key |
| `list_sessions` | 列出所有会话 | — |
| `get_session` | 查看单条会话 | 参数: session_id |
| `delete_session` | 删除会话 | 参数: session_id |

- **AND** 所有工具 SHALL 直接调用对应 repository 层的方法，而非子进程调用 CLI

#### Scenario: 会话管理与持久化

- **WHEN** 用户通过 Chat API 发送消息
- **THEN** AI 底座 SHALL 维护会话上下文（消息历史）
- **AND** 支持通过 `session_id` 区分不同会话
- **AND** 支持设置系统提示词（System Prompt），包含当前 `executor_name` 等信息
- **AND** 会话记录 SHALL 持久化到 workspace 的 `sessions/` 目录中
- **AND** 持久化格式 SHALL 使用 JSON Lines（.jsonl），每行一条消息

#### Scenario: 会话目录结构

- **WHEN** 系统创建会话记录
- **THEN** 目录结构 SHALL 为：

```
{workspace}/
└── sessions/
    ├── sessions.json          # 会话索引：{session_id: {title, created_at, updated_at, message_count}}
    ├── YOLO-CHAT-20260605-001.jsonl  # 单个会话消息记录
    ├── YOLO-CHAT-20260605-002.jsonl
    └── ...
```

- **AND** 会话 Key 格式：`YOLO-CHAT-{YYYYMMDD}-{NNN}`（如 `YOLO-CHAT-20260605-001`）
- **AND** `sessions.json` SHALL 维护所有会话的元数据索引
- **AND** 每个 `.jsonl` 文件 SHALL 包含该会话的完整消息历史

#### Scenario: 会话管理体系

- **WHEN** 用户发送消息创建新会话
- **THEN** 系统 SHALL 在 `sessions.json` 中添加条目
- **AND** 创建对应的 `.jsonl` 文件
- **AND** `sessions.json` 元数据包括：`session_id`、`title`（自动从首条消息生成）、`created_at`、`updated_at`、`message_count`

- **WHEN** 用户删除一个会话
- **THEN** 系统 SHALL 从 `sessions.json` 中移除该条目
- **AND** 删除对应的 `.jsonl` 文件

- **WHEN** 用户列出所有会话
- **THEN** 系统 SHALL 展示 `sessions.json` 中按 `updated_at` 倒序排列的会话列表
- **AND** TUI 模式下通过 `yolo session list` 或 TUI 内 `/sessions` 命令查看

#### Scenario: 多步推理

- **WHEN** 用户请求一个复杂操作（如"帮我创建一个高优先级的 Issue，然后分配给张三"）
- **THEN** AI 底座 SHALL 自主调用多个工具按序完成
- **AND** 每一步的结果反馈给 LLM，用于决定下一步操作

---

### Requirement: 标签模板系统

双端 UI 通过标签模板统一渲染结构化的数据展示。AI 输出的文本中嵌入标签，两端各按自己的方式渲染。

#### Scenario: 标签定义

系统 SHALL 支持以下标签：

| 标签 | 属性 | 前端渲染 | TUI 渲染 |
|------|------|---------|----------|
| `<ProjectList />` | — | Ant Design Table | 表格列表 |
| `<ProjectCard key="..." />` | key | Ant Design Card | 带边框卡片 |
| `<IssueList project="..." status="..." />` | project, status, executor | Ant Design List | 项目符号列表 |
| `<IssueCard key="..." />` | key | Ant Design Card | 带边框卡片 |
| `<TaskList issue="..." />` | issue | Ant Design Checkbox 列表 | 复选框列表 |
| `<ExecutorList />` | — | Ant Design Table | 表格列表 |
| `<KanbanBoard project="..." />` | project | 看板视图（四列：created/in_progress/done/archived） | 四列布局 |
| `<StatusBadge status="..." />` | status | 彩色标签 | 彩色文字 |
| `<PriorityBadge priority="..." />` | priority | 彩色标签 | 彩色文字 |
| `<DocumentList project="..." />` | project | Ant Design List | 项目符号列表 |

#### Scenario: 标签解析

- **WHEN** AI 底座生成响应
- **THEN** 响应内容中 SHALL 包含纯文本 + 标签的混合内容
- **AND** 响应 SHALL 通过 SSE 流式传输
- **AND** 前端/TUI SHALL 解析标签并替换为对应 UI 组件
- **AND** 双端 SHALL 使用相同的正则表达式 `/<\w+\s*[^>]*\/\s*>/g` 匹配自闭合标签

---

### Requirement: Chat API

#### Scenario: SSE 流式聊天

- **WHEN** 用户发送 `POST /api/v1/ai/chat` 请求
- **THEN** 服务端 SHALL 返回 `text/event-stream` 格式的 SSE 响应
- **AND** 请求体格式：`{"session_id": "...", "message": "..."}`
- **AND** 响应体 SHALL 包含 AI 思考过程的流式文本（含标签）
- **AND** 支持跨域请求（CORS）

#### Scenario: 会话管理 API

- **WHEN** 用户发送 `GET /api/v1/ai/sessions` 请求
- **THEN** 服务端 SHALL 返回所有会话的元数据列表

- **WHEN** 用户发送 `DELETE /api/v1/ai/sessions/{session_id}` 请求
- **THEN** 服务端 SHALL 删除指定会话的文件和索引

- **WHEN** 用户发送 `POST /api/v1/ai/sessions/reset` 请求
- **THEN** 服务端 SHALL 清除所有会话文件

---

### Requirement: 前端 AI Chat（Ant Design X）

#### Scenario: Chat 页面

- **WHEN** 用户访问 `/ai` 路径
- **THEN** 前端 SHALL 显示一个 Ant Design X 聊天界面
- **AND** 发送消息 SHALL 调用 `POST /api/v1/ai/chat` SSE API
- **AND** 接收到的标签 SHALL 被解析并渲染为对应 Ant Design 组件
- **AND** 消息气泡 SHALL 区分用户消息和 AI 消息
- **AND** 侧边栏 SHALL 显示会话历史列表

#### Scenario: 标签渲染

- **WHEN** AI 响应中包含 `<ProjectList />` 等标签
- **THEN** 前端 SHALL 解析标签，并发起对应数据请求填充内容
- **AND** 渲染为对应的 Ant Design 组件

---

### Requirement: TUI Chat（bubbles + lipgloss）

#### Scenario: TUI 启动

- **WHEN** 用户执行 `yolo chat`
- **THEN** TUI SHALL 启动一个全屏终端聊天界面（基于 bubbles）
- **AND** TUI SHALL 使用 lipgloss 进行样式美化
- **AND** 底部的文本输入框 SHALL 支持消息输入（`Enter` 发送、`Tab` 切换焦点）
- **AND** 左侧边栏 SHALL 显示会话列表（可折叠）

#### Scenario: 标签渲染

- **WHEN** AI 响应中包含标签
- **THEN** TUI SHALL 解析标签并渲染为对应的终端 UI 元素（颜色、表格、边框等）
- **AND** 对于列表类标签（如 `<ProjectList />`），TUI SHALL 使用 lipgloss 表格样式渲染
- **AND** 对于卡片类标签（如 `<IssueCard />`），TUI SHALL 使用 lipgloss 带边框样式渲染
- **AND** 对于看板标签（如 `<KanbanBoard />`），TUI SHALL 使用四列布局分屏显示

#### Scenario: TUI 命令

- **WHEN** 用户在输入框输入 `/sessions`
- **THEN** TUI SHALL 显示会话列表供选择和切换
- **WHEN** 用户在输入框输入 `/clear`
- **THEN** TUI SHALL 新建会话

#### Scenario: TUI 快捷键

- **WHEN** 用户在 TUI 中按下 `Ctrl+C`
- **THEN** TUI SHALL 退出程序
- **WHEN** 用户在 TUI 中按下 `Ctrl+L`
- **THEN** TUI SHALL 清屏
- **WHEN** 用户在 TUI 中按下 `Ctrl+N`
- **THEN** TUI SHALL 新建会话

---

### Requirement: 配置

#### Scenario: AI 配置段

- **WHEN** `yolo chat` 或 `yolo serve` 启动
- **THEN** 系统 SHALL 从以下来源读取 AI 配置（优先级从高到低）：
  1. 环境变量：`OPENAI_BASE_URL`、`OPENAI_API_KEY`、`OPENAI_MODEL`
  2. `settings.json` 中的 `ai` 配置段

- **AND** `settings.json` 扩展后的格式：
  ```json
  {
    "workspace": "...",
    "ai": {
      "base_url": "https://api.openai.com/v1",
      "api_key": "",
      "model": "gpt-4o",
      "system_prompt": "你是一个 yolo-team AI 助手..."
    }
  }
  ```

---

## REMOVED Requirements

N/A
