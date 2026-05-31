# 03 — Issue API

> 参考：[api.md](../docs/api.md) 3 节、[models.md](../docs/models.md) 3 节

这是最复杂的 API 模块，涉及 8 状态流转和角色权限校验。

## 子任务

### 3.1 基础 CRUD

- [x] `GET /api/v1/issues` — 列表，支持 query 筛选：
  - `?project_id=` / `?status=` / `?parent_id=` / `?executor_id=`
  - 顶级卡片响应包含 `children` 子卡片数组
- [x] `POST /api/v1/issues` — 创建，body:
  - 必填：`project_id`（int）、`title`（≤200）
  - 可选：`description`（≤512）、`priority`（默认 medium）、`parent_id`、`executor_id`、`deadline`（ISO8601）
  - 新建 status 固定 `created`
  - 创建子卡片校验父卡片为顶级（40003）
- [x] `GET /api/v1/issues/:key` — 详情（顶级含 children）
- [x] `PUT /api/v1/issues/:key` — 更新任意字段（见 3.2）
- [x] `DELETE /api/v1/issues/:key` — 删除（级联子卡片）

### 3.2 状态流转

实现完整的 8 状态流转引擎：

```
created → design → review → implementation → qa → pending_review → done → archived
```

**流转规则表**：

| 转换 | 允许角色 | 说明 |
|------|----------|------|
| `created` → `design` | `leader` | 开始设计 |
| `design` → `review` | `leader`、`architect` | 提交评审 |
| `review` → `implementation` | `leader` | 评审通过 |
| `review` → `design` | `leader`、`qa` | 不通过 |
| `implementation` → `qa` | `developer` | 提交质检 |
| `qa` → `pending_review` | `qa` | QA 通过 |
| `qa` → `implementation` | `qa` | QA 不通过 |
| `pending_review` → `done` | `leader` | 审查通过 |
| `pending_review` → `design/impl/qa` | `leader` | 不通过退回到指定状态 |
| `done` → `archived` | 仅页面 | 归档 |

- [x] 实现 `TransitionValidator` 逻辑在 `model/issue.go`（`CanTransitionTo` 方法）
- [x] CLI 不可设置 `review`、`pending_review`、`archived`（通过 `X-Yolo-CLI` header 判断，返回 40301）
- [x] 角色权限不足返回 40301
- [x] 非法流转返回 40002
- [x] `review` 和 `pending_review` 必须在请求中识别来源（header 标记），CLI 来源拒绝

### 3.3 待办查询

- [x] `GET /api/v1/issues/todo?executor_id=&limit=`
  - 条件：`executor_id` 匹配 且 `status NOT IN ('done', 'archived')`
  - 排序：`deadline ASC` 优先，`priority DESC` 次要
  - limit 默认 20

### 3.4 代码仓库关联

- [x] Issue 更新支持 `repo_url`、`repo_name`、`branch_name` 字段
- [x] 响应中包含仓库信息

### 3.5 Handler 层

- [x] ~~Service 层已删除，Handler 直接调用 Repository~~
- [x] `IssueHandler.List(req)` / `Create(req)` / `Get(req)` / `Update(req)` / `Delete(req)`
- [x] `IssueHandler.Todo(req)`
- [x] 创建子卡片前校验父卡片的 `parent_id IS NULL`
- [x] 删除顶级卡片时级联删除子卡片
- [x] 所有 Handler 方法使用泛型 `Wrap` 包装器

## 验收标准

- [x] 所有 CRUD 端点可用
- [x] 8 状态流转规则全部生效
- [x] 角色权限校验正确（leader/architect/developer/qa 各有对应权限）
- [x] CLI 来源无法设置 review/pending_review/archived
- [x] 待办查询排序正确
- [x] 父子卡片约束生效（仅一级）
- [x] Key 格式 `YOLO-ISSUE-{id}` 正确生成
