---
name: "yolo-team-skill"
description: "Yolo-Team CLI 完整操作指南，涵盖项目/任务/文档/执行人管理。当 AI 需要操作 yolo-cli 进行项目管理时使用。"
---

# Yolo-Team CLI 技能文档

## 概述

Yolo-Team CLI (`yolo`) 是一个基于 **AI 员工** 的工作流系统命令行工具。它以"项目（Project）+ 任务卡片（Issue）"为核心模型，让人与 AI 在同一套工作流体系中协作管理任务、流转状态、沉淀知识。

- **编程语言**: Go (Cobra + Gin + GORM)
- **数据库**: SQLite
- **单一二进制**: 前端嵌入在二进制中，`yolo serve` 同时提供 API 和前端页面

---

## 1. CLI 命令结构

### 1.1 全局选项

| 选项 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--json` | bool | `false` | 以 JSON 格式输出（所有命令均支持） |
| `--help`, `-h` | bool | — | 查看帮助 |

### 1.2 项目管理 (`yolo project`)

| 命令 | 说明 | 必需参数 | 可选参数 |
|------|------|---------|---------|
| `yolo project list` | 列出所有项目 | — | `--json` |
| `yolo project info <key>` | 查看项目详情 | `<key>` (如 `YOLO-PROJECT-1`) | `--json` |

项目规则：
- Key 自动生成，格式 `YOLO-PROJECT-{id}`
- 项目名称最长 32 字符
- 项目描述最长 128 字符

### 1.4 任务卡片管理 (`yolo issue`)

| 命令 | 说明 | 必需参数 | 可选参数 |
|------|------|---------|---------|
| `yolo issue list` | 列出卡片 | — | `-p, --project` / `-s, --status` / `-e, --executor` / `--parent` |
| `yolo issue create` | 创建卡片 | `-t, --title` / `-p, --project` | `-d, --description` / `--priority` / `--parent` / `--executor` / `--deadline` |
| `yolo issue info <key>` | 卡片详情 | `<key>` (如 `YOLO-ISSUE-1`) | `--json` |
| `yolo issue update <key>` | 更新卡片 | `<key>` | `-t, --title` / `-s, --status` / `--executor` / `--deadline` / `--repo-url` / `--repo-name` / `--branch` |
| `yolo issue todo` | 待办查询 | `-e, --executor` | `-n, --limit` (默认 20) |

#### Issue 状态流转规则

- 四种状态：`created` → `in_progress` → `done` → `archived`
- 不允许跳过状态（如 `in_progress` 不能直接到 `archived`）
- 不允许逆向流转

#### 优先级

- `low` / `medium`（默认）/ `high` / `critical`

#### 待办排序规则

`yolo issue todo` 按 **截至时间（最近优先）** + **优先级（critical > high > medium > low）** 排序。

#### 子卡片

- 仅顶级卡片（`parent_id` 为 null）可以作为父卡片
- 子卡片不支持再拥有子卡片（仅一级嵌套）

#### 仓库关联

每个 Issue 可关联一个代码仓库：

| 字段 | 说明 | 示例 |
|------|------|------|
| `--repo-url` | 仓库 URL | `https://github.com/example/repo` |
| `--repo-name` | 仓库名称 | `example/repo` |
| `--branch` | 分支名称 | `feature/login` |

### 1.5 执行人管理 (`yolo executor`)

| 命令 | 说明 | 必需参数 | 可选参数 |
|------|------|---------|---------|
| `yolo executor list` | 列出所有执行人 | — | `--json` |
| `yolo executor create` | 创建执行人（已隐藏） | `-n, --name` | `-s, --soul` |

执行人角色：
- `leader` — 项目负责人
- `architect` — 架构师
- `developer` — 开发人员（默认）
- `qa` — 测试人员

创建规则：
- 名称最长 64 字符，不能包含空格
- `soul` 字段为 AI 执行人的灵魂/系统提示词

### 1.6 文档管理 (`yolo doc`)

| 命令 | 说明 | 必需参数 | 可选参数 |
|------|------|---------|---------|
| `yolo doc list` | 列出项目文档 | `-p, --project` (项目 Key) | `--json` |
| `yolo doc create` | 创建文档 | `-p, --project` / `-t, --title` | `-c, --content` |
| `yolo doc info <key>` | 查看文档详情 | `<key>` (如 `YOLO-DOC-1`) | `--json` |

文档规则：
- 文档内容以 Markdown 文件存储在工作空间的 `documents/{project_id}/{title}.md`
- 创建后返回文档 Key 和本地文件路径，AI 可直接通过路径读写内容
- 标题最长 128 字符

### 1.7 HTTP 服务 (`yolo serve`)

启动 Web 服务，同时提供 REST API 和嵌入式前端。

```bash
yolo serve                    # 默认端口 8080
yolo serve -p 9090            # 自定义端口
```

---

## 3. AI 员工操作指南

### 3.1 初始化流程

作为 AI 员工，首次接入 yolo-team 时：

1. 检查是否已初始化：运行 `yolo project list`
2. 如果未初始化（提示 "Not initialized"），告知人类员工先运行 `yolo init`
3. 初始化完成后，确认工作空间路径

### 3.2 典型工作流程

#### 当人类员工说"帮我创建一个任务"

```bash
# 1. 查看项目列表
yolo project list

# 2. 创建任务（指定项目和标题）
yolo issue create -p <project_id> -t "任务标题" \
    -d "任务描述" \
    --priority high \
    --deadline "2026-07-01" \
    --executor <executor_id>
```

#### 当人类员工说"我今天有什么待办"

```bash
yolo issue todo -e <executor_id>
```

#### 推进卡片状态

```bash
yolo issue update YOLO-ISSUE-1 -s in_progress
yolo issue update YOLO-ISSUE-1 -s done
```

注意：状态必须遵循流转规则，不能跳过或逆向。

#### 创建项目文档

```bash
# 1. 创建文档记录
yolo doc create -p YOLO-PROJECT-1 -t "需求文档"

# 2. 通过返回的文件路径直接写入 Markdown 内容
# 文档路径：{workspace}/documents/1/需求文档.md
```

#### 关联代码仓库

```bash
yolo issue update YOLO-ISSUE-1 \
    --repo-url "https://github.com/example/repo" \
    --repo-name "example/repo" \
    --branch "feature/login"
```

### 3.3 最佳实践

1. **始终优先使用 `--json` 获取结构化数据**，便于程序化处理
2. **创建 Issue 时尽量填写完整信息**（优先级、截至时间、执行人），便于待办排序
3. **文档内容存放在本地文件系统中**，AI 可以直接使用文件读写工具进行操作
4. **卡片状态变更前，先确认当前状态**（使用 `yolo issue info <key>`），确保流转合法
5. **使用合适的执行人角色**：`leader`/`architect`/`developer`/`qa`，便于任务分派

---

## 4. 数据模型

### 4.1 Project（项目）

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | int64 | 自增主键 |
| Key | string | 唯一标识（如 `YOLO-PROJECT-1`） |
| Name | string | 名称（最长 32 字符） |
| Description | string | 描述（最长 128 字符） |

### 4.2 Executor（执行人）

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | int64 | 自增主键 |
| Name | string | 名称（最长 64 字符，不能含空格） |
| Role | string | 角色：`leader`/`architect`/`developer`/`qa` |
| Soul | string | AI 灵魂/系统提示词 |

### 4.3 Issue（任务卡片）

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | int64 | 自增主键 |
| Key | string | 唯一标识（如 `YOLO-ISSUE-1`） |
| ProjectID | int64 | 所属项目 |
| ParentID | *int64 | 父卡片 ID（null 为顶级卡片） |
| Status | string | `created`/`in_progress`/`done`/`archived` |
| Title | string | 标题（最长 200 字符） |
| Description | string | 描述（最长 512 字符） |
| Deadline | *time.Time | 截至时间 |
| Priority | string | `low`/`medium`/`high`/`critical` |
| ExecutorID | *int64 | 执行人 |
| RepoURL | string | 仓库 URL |
| RepoName | string | 仓库名称 |
| BranchName | string | 分支名称 |

### 4.4 Document（文档）

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | int64 | 自增主键 |
| Key | string | 唯一标识（如 `YOLO-DOC-1`） |
| ProjectID | int64 | 所属项目 |
| Title | string | 标题（最长 128 字符） |
| Creator | string | 创建者（默认 "人类"） |

---

## 5. 错误处理

| 情形 | 表现 | 处理方式 |
|------|------|---------|
| 未初始化 | `Not initialized. Please run: yolo init` | 先执行 `yolo init` |
| 缺少必需参数 | 错误提示缺少的参数 | 补全必需参数 |
| 状态流转不合法 | `Error: 状态流转不合法` | 检查当前状态和目标状态 |
| 参数不合法 | 对应校验错误信息 | 按提示修正参数 |
