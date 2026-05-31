# Yolo-Team CLI 命令定义

## 约定

| 项目 | 说明 |
|------|------|
| 全局参数 | `--server` 指定服务端地址，默认 `http://localhost:8080` |
| 输出格式 | 所有命令输出为可读的表格或 JSON（通过 `--json` 切换） |
| 错误处理 | 网络错误、API 错误均以非零退出码 + stderr 输出错误信息 |

### 通用选项

| 选项 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--server` | string | `http://localhost:8080` | 服务端地址 |
| `--json` | bool | `false` | 以 JSON 格式输出 |

---

## 1. 服务管理

### `yolo serve`

启动 HTTP 服务，同时提供 API 和嵌入式前端。

```
yolo serve [flags]
```

**选项**：

| 选项 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--port`, `-p` | int | `8080` | 监听端口 |
| `--db` | string | `{workspace}/yolo.db` | SQLite 数据库文件路径 |

**示例**：

```bash
yolo serve
# → 输出：Server started at http://localhost:8080

yolo serve --port 9090
# → 输出：Server started at http://localhost:9090
```

---

## 2. 项目管理（`yolo project`）

### 2.1 项目列表

```
yolo project list [flags]
```

**入参**：无

**输出**（表格）：

```
ID  KEY                NAME       DESCRIPTION  UPDATED
1   YOLO-PROJECT-1     某某项目    项目描述      2026-05-20 08:30
2   YOLO-PROJECT-2     另一个项目                2026-05-19 12:00
```

**输出**（`--json`）：

```json
[
  {
    "id": 1,
    "key": "YOLO-PROJECT-1",
    "name": "某某项目",
    "description": "项目描述",
    "created_at": "2026-05-01T10:00:00Z",
    "updated_at": "2026-05-20T08:30:00Z"
  }
]
```

### 2.2 项目详情

```
yolo project info <key> [flags]
```

**入参**：

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `<key>` | string | 是 | 项目 Key，如 `YOLO-PROJECT-1` |

**示例**：

```bash
yolo project info YOLO-PROJECT-1
```

**输出**（表格）：

```
Key:        YOLO-PROJECT-1
Name:       某某项目
Description: 项目描述
Created:    2026-05-01 10:00
Updated:    2026-05-20 08:30
```

---

## 3. 任务卡片管理（`yolo issue`）

### 3.1 卡片列表

```
yolo issue list [flags]
```

**选项**：

| 选项 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `--project`, `-p` | int | 否 | 按项目 ID 筛选 |
| `--status`, `-s` | string | 否 | 按状态筛选（`created`/`design`/`review`/`implementation`/`qa`/`pending_review`/`done`/`archived`） |
| `--executor`, `-e` | int | 否 | 按执行人 ID 筛选 |
| `--parent` | int | 否 | 按父卡片筛选（`0` 只看顶级卡片） |

**示例**：

```bash
# 查看项目 1 下所有顶级卡片
yolo issue list -p 1 --parent 0

# 查看某执行人的所有卡片
yolo issue list -e 1
```

**输出**（表格）：

```
KEY               STATUS          PRIORITY  TITLE           DEADLINE     REPO
YOLO-ISSUE-1      implementation  high      完成登录模块      2026-06-30   example/repo (feature/login)
  └── YOLO-ISSUE-3  design        medium    设计登录UI       -            -
YOLO-ISSUE-2      created         medium    编写测试用例      -            -
```

> 子卡片以缩进 + `└──` 显示在父卡片下方。无仓库关联时 REPO 列显示 `-`。

**输出**（`--json`）：

```json
[
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
    "repo_url": "",
    "repo_name": "",
    "branch_name": "",
    "sort_order": 0,
    "created_at": "2026-05-31T10:00:00Z",
    "updated_at": "2026-05-31T10:00:00Z",
    "children": [ ... ]
  }
]
```

### 3.2 创建卡片

```
yolo issue create [flags]
```

**选项**：

| 选项 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `--title`, `-t` | string | 是 | 任务标题，最长 200 字符 |
| `--project`, `-p` | int | 是 | 所属项目 ID |
| `--description`, `-d` | string | 否 | 任务描述，最长 512 字符 |
| `--priority` | string | 否 | 优先级，默认 `medium`（`low`/`medium`/`high`/`critical`） |
| `--parent` | int | 否 | 父卡片 ID（创建子卡片） |
| `--executor` | int | 否 | 执行人 ID |
| `--deadline` | string | 否 | 截至时间，格式 `YYYY-MM-DD` |

**示例**：

```bash
# 创建顶级任务
yolo issue create -t "完成登录模块" -p 1 --priority high --deadline "2026-06-30"

# 创建子任务
yolo issue create -t "设计登录UI" -p 1 --parent 1

# 创建并指定执行人
yolo issue create -t "编写测试用例" -p 1 --executor 1
```

**输出**：`Issue created: YOLO-ISSUE-3 (设计登录UI)`

### 3.3 卡片详情

```
yolo issue info <key> [flags]
```

**入参**：

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `<key>` | string | 是 | 卡片 Key，如 `YOLO-ISSUE-1` |

**示例**：

```bash
yolo issue info YOLO-ISSUE-1
```

**输出**（表格）：

```
Key:         YOLO-ISSUE-1
Project:     某某项目
Parent:      -
Status:      implementation
Title:       完成登录模块
Description: 实现用户名密码登录功能
Priority:    high
Executor:    zhangsan
Deadline:    2026-06-30
Repo:        example/repo (feature/login)
Children:
  └── YOLO-ISSUE-3  design  设计登录UI
```

### 3.4 更新卡片

```
yolo issue update <key> [flags]
```

**入参**：

| 参数/选项 | 类型 | 必需 | 说明 |
|-----------|------|------|------|
| `<key>` | string | 是 | 卡片 Key |
| `--title`, `-t` | string | 否 | 任务标题 |
| `--description`, `-d` | string | 否 | 任务描述 |
| `--status`, `-s` | string | 否 | 状态值（须符合流转规则） |
| `--priority` | string | 否 | 优先级 |
| `--executor` | int | 否 | 执行人 ID |
| `--deadline` | string | 否 | 截至时间 |
| `--repo-url` | string | 否 | 仓库 URL |
| `--repo-name` | string | 否 | 仓库名称 |
| `--branch` | string | 否 | 分支名称 |

**示例**：

```bash
# 推进状态
yolo issue update YOLO-ISSUE-1 -s design

# 指派执行人
yolo issue update YOLO-ISSUE-1 --executor 2

# 关联代码仓库
yolo issue update YOLO-ISSUE-1 --repo-url "https://github.com/example/repo" --repo-name "example/repo" --branch "feature/login"
```

**错误示例**：

```bash
yolo issue update YOLO-ISSUE-1 -s done
# → Error: 状态流转不合法：implementation 不能直接变更为 done
```

> 状态流转规则见 [models.md](./models.md) 3.2 节。`review`、`pending_review`、`archived` 状态不可通过 CLI 设置。

### 3.5 删除卡片

```
yolo issue delete <key> [flags]
```

**入参**：

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `<key>` | string | 是 | 卡片 Key |

**示例**：

```bash
yolo issue delete YOLO-ISSUE-3
# → Issue deleted: YOLO-ISSUE-3
#   → 1 children also deleted (如果有子卡片)
```

### 3.6 待办查询

```
yolo issue todo [flags]
```

**选项**：

| 选项 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `--executor`, `-e` | int | 是 | 执行人 ID |
| `--limit`, `-n` | int | 否 | 返回数量上限，默认 20 |

**示例**：

```bash
yolo issue todo -e 1
yolo issue todo -e 1 -n 5
```

**输出**（表格，按截至时间 + 优先级排序）：

```
KEY               STATUS    PRIORITY  TITLE           DEADLINE
YOLO-ISSUE-5      created   critical  修复安全问题      2026-06-01
YOLO-ISSUE-1      design    high      完成登录模块      2026-06-30
YOLO-ISSUE-8      qa        medium    更新API文档      2026-07-15
```

---

## 4. 执行人管理（`yolo executor`）

### 4.1 执行人列表

```
yolo executor list [flags]
```

**入参**：无

**输出**（表格）：

```
ID  NAME        CREATED
1   zhangsan    2026-05-01 10:00
2   lisi        2026-05-15 14:30
```

---

## 5. 项目文档管理（`yolo doc`）

### 5.1 文档列表

```
yolo doc list [flags]
```

**选项**：

| 选项 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `--project`, `-p` | string | 是 | 项目 Key，如 `YOLO-PROJECT-1` |

**示例**：

```bash
yolo doc list -p YOLO-PROJECT-1
```

**输出**（表格）：

```
KEY              TITLE       UPDATED
YOLO-DOC-1       需求文档      2026-05-31 10:00
YOLO-DOC-2       技术方案      2026-05-30 15:00
```

### 5.2 创建文档

```
yolo doc create [flags]
```

**选项**：

| 选项 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `--project`, `-p` | string | 是 | 项目 Key，如 `YOLO-PROJECT-1` |
| `--title`, `-t` | string | 是 | 文档标题（同时作为文件名），最长 128 字符 |

**示例**：

```bash
yolo doc create -p YOLO-PROJECT-1 -t "需求文档"
# → Document created: YOLO-DOC-3
# → Path: {workspace}/documents/1/需求文档.md
```

> 创建后返回文档的 Key 和本地文件路径，AI 可直接通过路径写入内容。

### 5.3 文档详情

```
yolo doc info <key> [flags]
```

**入参**：

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `<key>` | string | 是 | 文档 Key，如 `YOLO-DOC-1` |

**输出**：以 Markdown 渲染形式（终端内）显示文档正文，并标注本地文件路径。

**示例**：

```bash
yolo doc info YOLO-DOC-1
# → File: {workspace}/documents/1/需求文档.md
# → 在终端内渲染显示 Markdown 内容
```

---

## 退出码

| 退出码 | 含义 |
|--------|------|
| `0` | 成功 |
| `1` | 通用错误（网络、API 错误等） |
| `2` | 参数错误（缺少必需参数、值不合法） |
