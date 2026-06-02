# 06 — 前端页面

> 参考：[use-cases.md](../docs/use-cases.md)

实现变更为：使用 React 18 + Vite + Ant Design 5 构建（非 Umi 4）。

## 子任务

### 6.1 项目初始化

- [x] 创建 `yolo-client` 项目（React 18 + Vite + TypeScript）
- [x] 安装依赖：`antd`、`@ant-design/icons`、`react-router-dom`、`dayjs`、`restful-react`、`react-markdown`、`@dnd-kit/core`
- [x] 使用 `restful-react` 从 swagger.json 自动生成 `src/generated.tsx`
- [x] Vite proxy 配置 `/api` → `http://localhost:8080`

### 6.2 全局看板（`/`）

- [x] 顶栏：项目下拉菜单、新建项目、项目管理、文档按钮
- [x] 8 列看板：每列按 status 分组展示 Issue 卡片
- [x] 卡片内容：标题、优先级标签、截至时间、仓库标记
- [x] 拖拽：卡片跨列拖拽 → `PUT /api/v1/issues/:key {status}` + 乐观更新 + 失败回滚
- [x] 快速指派：卡片上执行人下拉 → `PUT /api/v1/issues/:key {executor_id}` + 乐观更新
- [x] 归档操作：done 卡片"归档"按钮 → `PUT /api/v1/issues/:key {status: 'archived'}` + 乐观更新
- [x] 新建 Issue："已创建"列 `+` → 抽屉表单

### 6.3 Issue 详情抽屉

- [x] 点击卡片 → 左侧 Ant Design Drawer
- [x] 编辑区：标题、描述、状态下拉、优先级、执行人、仓库信息
- [x] 底部：保存 / 删除按钮

### 6.4 项目管理（`/project/:key/manage`）

- [x] 表单编辑项目名称和描述
- [x] 保存调用 `PUT /api/v1/projects/:key`

### 6.5 文档列表（`/project/:key/docs`）

- [x] 表格展示文档 key、title、更新时间
- [x] 新建文档
- [x] 删除文档（Popconfirm 确认）

### 6.6 文档详情（`/doc/:key`）

- [x] Markdown 渲染（react-markdown）
- [x] 编辑模式切换
- [x] 保存 → `PUT /api/v1/documents/:key`

### 6.7 API 封装

- [x] `restful-react` 自动生成类型和 hooks
- [x] 所有组件通过生成的 hooks 调用 API

### 6.8 组件化

- [x] 组件：`TopBar`、`IssueCard`、`KanbanBoard`（含 `DraggableIssueCard`、`KanbanColumn`）、`CreateIssueForm`、`IssueDrawer`、`CreateProjectModal`
- [x] 页面：`KanbanPage`、`ProjectManage`、`DocumentList`、`DocumentDetail`
- [x] 常量：`constants.ts`

## 验收标准

- [x] 全局看板 8 列展示 + 项目筛选
- [x] 拖拽卡片可改变状态，网络错误自动回滚
- [x] 快速指派执行人，乐观更新
- [x] 归档操作
- [x] 文档 CRUD 完整，Markdown 渲染
- [x] API 代理正常工作
