# Yolo-Team API 定义

## 约定

| 项目 | 说明 |
|------|------|
| 基础路径 | `http://localhost:8080/api/v1` |
| 请求格式 | `application/json`（GET/DELETE 无 body） |
| 响应格式 | `application/json` |
| 时间格式 | ISO 8601（`2026-06-30T15:04:05Z`） |
| 资源定位 | Project、Issue、Document 使用 `key`（如 `YOLO-PROJECT-1`），Executor 使用 `name`（不含空格） |

### 通用响应结构

```json
// 成功（单条）
{
  "code": 0,
  "data": { ... }
}

// 成功（列表）
{
  "code": 0,
  "data": [ ... ]
}

// 错误
{
  "code": 40001,
  "message": "错误描述"
}
```

### 错误码

| code | 含义 |
|------|------|
| `0` | 成功 |
| `40001` | 请求参数校验失败 |
| `40002` | 状态流转不合法 |
| `40003` | 父子卡片嵌套层级超限 |
| `40301` | 权限不足（如非 leader 推进至 done） |
| `40401` | 资源不存在 |
| `40901` | 资源冲突（如执行人名称重复） |
| `50001` | 服务内部错误 |

---

## 1. 项目（Project）

### 1.1 项目列表

```
GET /api/v1/projects
```

**入参**：无

**响应**：

```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "key": "YOLO-PROJECT-1",
      "name": "某某项目",
      "description": "项目描述",
      "created_at": "2026-05-01T10:00:00Z",
      "updated_at": "2026-05-20T08:30:00Z"
    }
  ]
}
```

### 1.2 创建项目

```
POST /api/v1/projects
```

**入参**（JSON Body）：

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 项目名称，最长 32 字符 |
| `description` | string | 否 | 项目描述，最长 128 字符 |

```json
{
  "name": "某某项目",
  "description": "这是项目描述"
}
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "key": "YOLO-PROJECT-1",
    "name": "某某项目",
    "description": "这是项目描述",
    "created_at": "2026-05-31T10:00:00Z",
    "updated_at": "2026-05-31T10:00:00Z"
  }
}
```

### 1.3 项目详情

```
GET /api/v1/projects/:key
```

**入参**：

| 字段 | 类型 | 位置 | 必需 | 说明 |
|------|------|------|------|------|
| `key` | string | path | 是 | 项目 Key，如 `YOLO-PROJECT-1` |

**响应**：同 1.2 创建项目的 data 结构

### 1.4 更新项目

```
PUT /api/v1/projects/:key
```

**入参**：

| 字段 | 类型 | 位置 | 必需 | 说明 |
|------|------|------|------|------|
| `key` | string | path | 是 | 项目 Key |
| `name` | string | body | 否 | 项目名称，最长 32 字符 |
| `description` | string | body | 否 | 项目描述，最长 128 字符 |

```json
{
  "name": "新名称"
}
```

**响应**：同 1.2 创建项目的 data 结构

### 1.5 删除项目

```
DELETE /api/v1/projects/:key
```

**入参**：

| 字段 | 类型 | 位置 | 必需 | 说明 |
|------|------|------|------|------|
| `key` | string | path | 是 | 项目 Key |

**响应**：

```json
{
  "code": 0,
  "data": null
}
```

---

## 2. 执行人（Executor）

### 2.1 执行人列表

```
GET /api/v1/executors
```

**入参**：无

**响应**：

```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "name": "zhangsan",
      "role": "developer",
      "soul": "你是一个资深后端工程师...",
      "created_at": "2026-05-01T10:00:00Z",
      "updated_at": "2026-05-01T10:00:00Z"
    }
  ]
}
```

### 2.2 创建执行人

```
POST /api/v1/executors
```

**入参**（JSON Body）：

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 执行人名称，最长 64 字符，唯一且不含空格 |
| `soul` | string | 否 | 角色定义，描述执行人的能力与行为准则 |

> `role` 字段不通过 API 设置，创建时默认为 `developer`，仅可通过页面修改。

```json
{
  "name": "zhangsan",
  "soul": "你是一个资深后端工程师，擅长 Go 和数据库设计..."
}
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "zhangsan",
    "role": "developer",
    "soul": "你是一个资深后端工程师，擅长 Go 和数据库设计...",
    "created_at": "2026-05-31T10:00:00Z",
    "updated_at": "2026-05-31T10:00:00Z"
  }
}
```

### 2.3 删除执行人

```
DELETE /api/v1/executors/:name
```

**入参**：

| 字段 | 类型 | 位置 | 必需 | 说明 |
|------|------|------|------|------|
| `name` | string | path | 是 | 执行人名称 |

**响应**：

```json
{
  "code": 0,
  "data": null
}
```

---

## 3. 任务卡片（Issue）

### 3.1 卡片列表

```
GET /api/v1/issues
```

**入参**（Query）：

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `project_id` | int | 否 | 按项目 ID 筛选 |
| `status` | string | 否 | 按状态筛选（`created`/`design`/`review`/`implementation`/`qa`/`pending_review`/`done`/`archived`） |
| `parent_id` | int | 否 | 按父卡片筛选（`0` 只返回顶级卡片） |
| `executor_id` | int | 否 | 按执行人 ID 筛选 |

**响应**：

```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "key": "YOLO-ISSUE-1",
      "project_id": 1,
      "parent_id": null,
      "status": "implementation",
      "title": "完成登录模块",
      "description": "实现用户名密码登录功能",
      "deadline": "2026-06-30T00:00:00Z",
      "priority": "high",
      "executor_id": 1,
      "repo_url": "https://github.com/example/repo",
      "repo_name": "example/repo",
      "branch_name": "feature/login",
      "sort_order": 0,
      "created_at": "2026-05-31T10:00:00Z",
      "updated_at": "2026-05-31T10:00:00Z",
      "children": []
    }
  ]
}
```

> `children` 为子卡片数组，仅顶级卡片返回。子卡片的 `children` 始终为空。

### 3.2 创建卡片

```
POST /api/v1/issues
```

**入参**（JSON Body）：

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `project_id` | int | 是 | 所属项目 ID |
| `title` | string | 是 | 任务标题，最长 200 字符 |
| `description` | string | 否 | 任务描述，最长 512 字符 |
| `priority` | string | 否 | 优先级，默认 `medium`（`low`/`medium`/`high`/`critical`） |
| `parent_id` | int | 否 | 父卡片 ID（创建子卡片时使用，父卡片必须为顶级卡片） |
| `executor_id` | int | 否 | 执行人 ID |
| `deadline` | string | 否 | 截至时间，ISO 8601 格式 |

```json
{
  "project_id": 1,
  "title": "完成登录模块",
  "description": "实现用户名密码登录功能",
  "priority": "high",
  "executor_id": 1,
  "deadline": "2026-06-30T00:00:00Z"
}
```

**响应**：同 3.1 列表中的单条卡片结构

**校验规则**：
- `parent_id` 不为空时，校验父卡片 `parent_id IS NULL`，否则返回 `40003`
- 新建卡片 `status` 固定为 `created`

### 3.3 卡片详情

```
GET /api/v1/issues/:key
```

**入参**：

| 字段 | 类型 | 位置 | 必需 | 说明 |
|------|------|------|------|------|
| `key` | string | path | 是 | 卡片 Key，如 `YOLO-ISSUE-1` |

**响应**：同 3.1 列表中的单条卡片结构（顶级卡片含 `children`）

### 3.4 更新卡片

```
PUT /api/v1/issues/:key
```

**入参**：

| 字段 | 类型 | 位置 | 必需 | 说明 |
|------|------|------|------|------|
| `key` | string | path | 是 | 卡片 Key |
| `title` | string | body | 否 | 任务标题 |
| `description` | string | body | 否 | 任务描述 |
| `status` | string | body | 否 | 状态值（须符合流转规则） |
| `priority` | string | body | 否 | 优先级 |
| `executor_id` | int | body | 否 | 执行人 ID |
| `deadline` | string | body | 否 | 截至时间 |
| `repo_url` | string | body | 否 | 仓库 URL |
| `repo_name` | string | body | 否 | 仓库名称 |
| `branch_name` | string | body | 否 | 分支名称 |

```json
{
  "status": "design",
  "deadline": "2026-07-15T00:00:00Z"
}
```

**状态流转规则**：

| 当前状态 | 允许变更为 | 限制 |
|----------|-----------|------|
| `created` | `design` | `leader` |
| `design` | `review` | `leader`/`architect` |
| `review` | `implementation` | `leader`（通过）；CLI 不可设置 |
| `review` | `design` | `leader`/`qa`（不通过）；CLI 不可设置 |
| `implementation` | `qa` | `developer` |
| `qa` | `pending_review` | `qa`（通过） |
| `qa` | `implementation` | `qa`（不通过） |
| `pending_review` | `done` | `leader`（通过）；CLI 不可设置 |
| `pending_review` | `design`/`implementation`/`qa` | `leader`（不通过）；CLI 不可设置 |
| `done` | `archived` | 仅页面操作 |
| `archived` | 不可变更 | — |

非法流转返回 `40002`，角色权限不足返回 `40301`。

**响应**：同 3.1 列表中的单条卡片结构

### 3.5 删除卡片

```
DELETE /api/v1/issues/:key
```

**入参**：

| 字段 | 类型 | 位置 | 必需 | 说明 |
|------|------|------|------|------|
| `key` | string | path | 是 | 卡片 Key |

**响应**：

```json
{
  "code": 0,
  "data": null
}
```

> 删除顶级卡片时，其下所有子卡片一并删除。

### 3.6 待办查询

```
GET /api/v1/issues/todo
```

**入参**（Query）：

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `executor_id` | int | 是 | 执行人 ID |
| `limit` | int | 否 | 返回数量上限，默认 20 |

**查询逻辑**：
- `executor_id` 匹配
- `status NOT IN ('done', 'archived')`
- 排序：`deadline ASC`（优先），`priority DESC`（`critical > high > medium > low`）

**响应**：同 3.1 卡片列表结构

---

## 4. 项目文档（Document）

### 4.1 文档列表

```
GET /api/v1/projects/:key/documents
```

**入参**：

| 字段 | 类型 | 位置 | 必需 | 说明 |
|------|------|------|------|------|
| `key` | string | path | 是 | 项目 Key |

**响应**：

```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "key": "YOLO-DOC-1",
      "project_id": 1,
      "title": "需求文档",
      "sort_order": 0,
      "created_at": "2026-05-31T10:00:00Z",
      "updated_at": "2026-05-31T10:00:00Z"
    }
  ]
}
```

> 文档列表**不返回正文内容**，正文需通过 4.3 详情接口获取。

### 4.2 创建文档

```
POST /api/v1/projects/:key/documents
```

**入参**：

| 字段 | 类型 | 位置 | 必需 | 说明 |
|------|------|------|------|------|
| `key` | string | path | 是 | 项目 Key |
| `title` | string | body | 是 | 文档标题（同时作为文件名），最长 128 字符 |
| `content` | string | body | 是 | 文档正文（Markdown） |

```json
{
  "title": "需求文档",
  "content": "# 需求文档\n\n..."
}
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "key": "YOLO-DOC-1",
    "project_id": 1,
    "title": "需求文档",
    "content": "# 需求文档\n\n...",
    "file_path": "{workspace}/documents/1/需求文档.md",
    "sort_order": 0,
    "created_at": "2026-05-31T10:00:00Z",
    "updated_at": "2026-05-31T10:00:00Z"
  }
}
```

> 响应中 `file_path` 为文档在本地文件系统中的绝对路径。AI 可通过此路径直接读写文件，避免大量内容经 HTTP 传输。同时文件也会写入 `{workspace}/documents/{id}/{title}.md`。

### 4.3 文档详情

```
GET /api/v1/documents/:key
```

**入参**：

| 字段 | 类型 | 位置 | 必需 | 说明 |
|------|------|------|------|------|
| `key` | string | path | 是 | 文档 Key，如 `YOLO-DOC-1` |

**响应**：同 4.2 创建文档的响应（含 `content` 正文和 `file_path` 本地路径）

### 4.4 更新文档

```
PUT /api/v1/documents/:key
```

**入参**：

| 字段 | 类型 | 位置 | 必需 | 说明 |
|------|------|------|------|------|
| `key` | string | path | 是 | 文档 Key |
| `title` | string | body | 否 | 文档标题（修改标题将重命名文件） |
| `content` | string | body | 否 | 文档正文 |

```json
{
  "title": "需求文档 v2",
  "content": "# 需求文档 v2\n\n..."
}
```

**响应**：同 4.2 创建文档的响应

### 4.5 删除文档

```
DELETE /api/v1/documents/:key
```

**入参**：

| 字段 | 类型 | 位置 | 必需 | 说明 |
|------|------|------|------|------|
| `key` | string | path | 是 | 文档 Key |

**响应**：

```json
{
  "code": 0,
  "data": null
}
```

> 删除时同步删除 DB 记录和 `{workspace}/documents/{id}/{title}.md` 文件。
