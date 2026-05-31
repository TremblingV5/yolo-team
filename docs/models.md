# Yolo-Team 数据模型

## 1. 模型总览

```
Project（项目）
  ├── Issue（任务卡片）
  │     ├── 子 Issue（仅一级）
  │     ├── Executor（执行人）
  │     └── 代码仓库关联
  │           ├── repo_url
  │           ├── repo_name
  │           └── branch_name
  └── Document（项目文档）

Executor（执行人）
  └── 系统中可分配任务的执行人
```

---

## 2. 项目（Project）

项目的核心组织单元，代表一个正在进行的项目。

| 属性 | 类型 | 说明 |
|------|------|------|
| `id` | INTEGER | 自增主键 |
| `key` | VARCHAR(64) | 唯一标识，格式 `YOLO-PROJECT-{id}` |
| `name` | VARCHAR(32) | 项目名称 |
| `description` | VARCHAR(128) | 项目描述 |
| `created_at` | DATETIME | 创建时间 |
| `updated_at` | DATETIME | 更新时间 |

---

## 3. 任务卡片（Issue）

任务卡片是项目中的工作单元。卡片状态固定为八种。

### 3.1 固定状态流转

```
  已创建 ──► 方案设计 ──► 方案评审 ──► 代码实现 ──► QA质检 ──► 待审查 ──► 已完成 ──► 已归档
(created)  (design)    (review)    (implementation) (qa)   (pending_review) (done)  (archived)
              ▲           │              ▲           │            │
              │  评审不通过 │              │  质检不通过 │            │ 审查不通过
              └───────────┘              └───────────┘            ├──► design
                                                                  ├──► implementation
                                                                  └──► qa
```

| 状态值 | 显示名 | 说明 | 操作限制 |
|--------|--------|------|----------|
| `created` | 已创建 | 任务新建，待规划 | — |
| `design` | 方案设计 | 正在设计方案 | — |
| `review` | 方案评审 | 评审设计方案 | CLl 不可设置 |
| `implementation` | 代码实现 | 正在编码实现 | — |
| `qa` | QA质检 | 测试验证中 | — |
| `pending_review` | 待审查 | 工作完成，待 leader 审查 | CLl 不可设置 |
| `done` | 已完成 | 任务完成 | 仅 `leader` 角色可将任务推进至此 |
| `archived` | 已归档 | 任务归档 | 仅人类通过页面操作，CLI 不可设置 |

### 3.2 状态转换权限

每种状态转换所允许的执行人角色：

| 转换 | 允许角色 | 说明 |
|------|----------|------|
| `created` → `design` | `leader` | 开始方案设计 |
| `design` → `review` | `leader`、`architect` | 提交方案评审 |
| `review` → `implementation` | `leader` | 评审通过，进入开发 |
| `review` → `design` | `leader`、`qa` | 评审不通过，重新设计 |
| `implementation` → `qa` | `developer` | 提交 QA 质检 |
| `qa` → `pending_review` | `qa` | QA 通过，提交审查 |
| `qa` → `implementation` | `qa` | QA 不通过，返工 |
| `pending_review` → `done` | `leader` | 审查通过 |
| `pending_review` → `design` | `leader` | 审查不通过，重新设计 |
| `pending_review` → `implementation` | `leader` | 审查不通过，重新实现 |
| `pending_review` → `qa` | `leader` | 审查不通过，重新 QA |
| `done` → `archived` | 仅页面操作 | 任务归档 |

### 3.3 属性

| 属性 | 类型 | 说明 |
|------|------|------|
| `id` | INTEGER | 自增主键 |
| `key` | VARCHAR(64) | 唯一标识，格式 `YOLO-ISSUE-{id}` |
| `project_id` | INTEGER | 所属项目 |
| `parent_id` | INTEGER\|NULL | 父卡片 ID。NULL 为顶级卡片，仅允许一级嵌套 |
| `status` | VARCHAR(16) | 当前状态，枚举值见 3.1，默认 `created` |
| `title` | VARCHAR(200) | 任务标题 |
| `description` | VARCHAR(512) | 任务描述（简明扼要） |
| `deadline` | DATETIME\|NULL | 截至时间 |
| `priority` | VARCHAR(16) | 优先级：`low` / `medium` / `high` / `critical` |
| `executor_id` | INTEGER\|NULL | 执行人 ID，关联 `executors` 表 |
| `repo_url` | VARCHAR(500) | 关联仓库 URL |
| `repo_name` | VARCHAR(255) | 关联仓库名称 |
| `branch_name` | VARCHAR(255) | 关联分支名称 |
| `sort_order` | INTEGER | 排序序号 |
| `created_at` | DATETIME | 创建时间 |
| `updated_at` | DATETIME | 更新时间 |

### 3.4 父子卡片约束

- `parent_id` 为 NULL 表示顶级卡片
- `parent_id` 不为 NULL 表示子卡片
- **已经为子卡片的卡片，不允许再拥有子卡片**（即嵌套深度最多 1 层）
- 创建子卡片时需校验父卡片 `parent_id IS NULL`

---

## 4. 项目文档（Document）

每个项目下可创建多个文档，用于承载需求说明、技术方案、会议记录等项目文档。

文档采用 **数据库记录元数据 + 文件系统存储内容** 的混合方案：

- 数据库 `project_documents` 表仅存储元数据（标题 `title` 即文件名）
- 实际内容以 `.md` 文件形式存储在 `{workspace}/documents/{project-id}/{title}.md`

| 属性 | 类型 | 说明 |
|------|------|------|
| `id` | INTEGER | 自增主键 |
| `key` | VARCHAR(64) | 唯一标识，格式 `YOLO-DOC-{id}` |
| `project_id` | INTEGER | 所属项目 |
| `title` | VARCHAR(128) | 文档标题，同时作为文件名 |
| `sort_order` | INTEGER | 排序序号 |
| `created_at` | DATETIME | 创建时间 |
| `updated_at` | DATETIME | 更新时间 |

### 4.1 文档存储路径

```
{workspace}/
├── yolo.db
└── documents/
    └── {project-id}/
        ├── 需求文档.md
        ├── 技术方案.md
        └── 会议记录.md
```

---

## 5. 执行人（Executor）

系统中可分配任务的执行人。`name` 用作接口标识（替代 id），必须唯一且不允许包含空格。

`role` 字段定义执行人的角色类型，**仅能通过页面修改，CLI 不可设置**。

| 属性 | 类型 | 说明 |
|------|------|------|
| `id` | INTEGER | 自增主键 |
| `name` | VARCHAR(64) | 执行人名称，唯一、不含空格 |
| `role` | VARCHAR(16) | 角色类型：`leader` / `architect` / `developer` / `qa`，默认 `developer` |
| `soul` | TEXT | 角色定义，描述执行人的能力与行为准则 |
| `created_at` | DATETIME | 创建时间 |
| `updated_at` | DATETIME | 更新时间 |

### 5.1 角色权限

| 角色 | 值 | 权限说明 |
|------|------|----------|
| Leader | `leader` | 完全控制，唯一可将任务推进至 `done` 的角色 |
| Architect | `architect` | 架构设计相关 |
| Developer | `developer` | 编码实现相关 |
| QA | `qa` | 测试质检相关 |

---

## 6. 数据库设计

使用 **GORM** 作为 ORM，底层支持 **SQLite** 和 **MySQL** 两种存储后端。

数据库迁移使用 **GORM AutoMigrate**，每次启动时自动同步表结构，无需手动执行迁移命令。查询使用 **gorm/gen** 生成类型安全 DAO。

**设计原则**：尽量不使用 TEXT 类型，改用 VARCHAR 并设定合理长度限制。

### 6.1 表结构

```sql
-- 项目表
CREATE TABLE projects (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    key         VARCHAR(64)  NOT NULL UNIQUE,
    name        VARCHAR(32)  NOT NULL,
    description VARCHAR(128) DEFAULT '',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 任务卡片表
CREATE TABLE issues (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    key           VARCHAR(64)  NOT NULL UNIQUE,
    project_id    INTEGER NOT NULL,
    parent_id     INTEGER DEFAULT NULL,
    status        VARCHAR(16)  NOT NULL DEFAULT 'created',
    title         VARCHAR(200) NOT NULL,
    description   VARCHAR(512) DEFAULT '',
    deadline      DATETIME DEFAULT NULL,
    priority      VARCHAR(16)  DEFAULT 'medium',
    executor_id   INTEGER DEFAULT NULL,
    repo_url      VARCHAR(500) DEFAULT '',
    repo_name     VARCHAR(255) DEFAULT '',
    branch_name   VARCHAR(255) DEFAULT '',
    sort_order    INTEGER DEFAULT 0,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (parent_id) REFERENCES issues(id),
    FOREIGN KEY (executor_id) REFERENCES executors(id)
);

-- 执行人表
CREATE TABLE executors (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        VARCHAR(64) NOT NULL UNIQUE,
    role        VARCHAR(16) NOT NULL DEFAULT 'developer',
    soul        TEXT DEFAULT '',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 项目文档表（仅元数据，内容存文件系统）
CREATE TABLE project_documents (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    key         VARCHAR(64)  NOT NULL UNIQUE,
    project_id  INTEGER NOT NULL,
    title       VARCHAR(128) NOT NULL,
    sort_order  INTEGER DEFAULT 0,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);
```

### 6.2 Key 生成策略

`key` 字段依赖于自增 `id`，因此采用"先插入、再回写"的策略：

```go
// 伪代码
func CreateProject(name string) *Project {
    p := &Project{Name: name}
    db.Insert(p)                           // INSERT → 获得 id
    p.Key = fmt.Sprintf("YOLO-PROJECT-%d", p.Id)
    db.Update(p, "key")                    // UPDATE key
    return p
}
```

同理，`YOLO-ISSUE-{id}` 和 `YOLO-DOC-{id}` 也采用相同方式生成。
