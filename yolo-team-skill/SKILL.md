---
name: "yolo-team-skill"
description: "Yolo-Team CLI 完整操作指南，涵盖项目/任务/子任务/文档/执行人管理。当 AI 需要操作 yolo-cli 进行项目管理时使用。"
---

# Yolo-Team CLI 技能文档

## 概述

Yolo-Team CLI (`yolo`) 是一个基于 **AI 员工** 的工作流系统命令行工具。它以"项目（Project）+ 任务卡片（Issue）+ 子任务（Task）"为核心模型，让人与 AI 在同一套工作流体系中协作管理任务、流转状态、沉淀知识。

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

### 1.3 任务卡片管理 (`yolo issue`)

| 命令 | 说明 | 必需参数 | 可选参数 |
|------|------|---------|---------|
| `yolo issue list` | 列出卡片 | — | `-p, --project` / `-s, --status` / `-e, --executor` |
| `yolo issue create` | 创建卡片 | `-t, --title` / `-p, --project` | `-d, --description` / `--priority` / `--executor` / `--deadline` |
| `yolo issue info <key>` | 卡片详情 | `<key>` (如 `YOLO-ISSUE-1`) | `--json` |
| `yolo issue update <key>` | 更新卡片 | `<key>` | `-t, --title` / `-s, --status` / `--executor` / `--deadline` / `--repo-url` / `--repo-name` / `--branch` |
| `yolo issue todo` | 待办查询 | `-e, --executor` | `-n, --limit` (默认 20) |

#### Issue 状态流转规则

- 四种状态：`created` → `in_progress` → `done` → `archived`
- 所有 Issue 平级，不支持父子关系（子任务拆分请使用 `yolo task`）

**推进到 `done` 前必须满足**：该 Issue 下所有 Task 均已为 `done`，否则返回 `40001: all tasks must be completed before closing the issue`

#### 优先级

- `low` / `medium`（默认）/ `high` / `critical`

#### 待办排序规则

`yolo issue todo` 按 **截至时间（最近优先）** + **优先级（critical > high > medium > low）** 排序。

#### 仓库关联

每个 Issue 可关联一个代码仓库：

| 字段 | 说明 | 示例 |
|------|------|------|
| `--repo-url` | 仓库 URL | `https://github.com/example/repo` |
| `--repo-name` | 仓库名称 | `example/repo` |
| `--branch` | 分支名称 | `feature/login` |

### 1.4 子任务管理 (`yolo task`)

子任务（Task）归属于某个 Issue，用于将卡片拆分为更细粒度的执行单元。Key 格式 `YOLO-TASK-{id}`，所有交互以 Key 为主。

| 命令 | 说明 | 必需参数 | 可选参数 |
|------|------|---------|---------|
| `yolo task list` | 列出子任务 | `-i, --issue` (Issue Key) | `--json` |
| `yolo task create` | 创建子任务 | `-i, --issue` / `-t, --title` | `-d, --description` |
| `yolo task update` | 更新子任务 | `-k, --key` (Task Key) | `-t, --title` / `-d, --description` / `-s, --status` |
| `yolo task delete` | 删除子任务 | `-k, --key` (Task Key) | — |

#### Task 状态

三种状态：`created` → `in_progress` → `done`

#### 使用示例

```bash
# 列出某个 Issue 下的所有子任务
yolo task list -i YOLO-ISSUE-1

# 创建子任务
yolo task create -i YOLO-ISSUE-1 -t "实现登录接口" -d "包括账号密码校验和 Token 签发"

# 推进子任务
yolo task update -k YOLO-TASK-1 -s in_progress
yolo task update -k YOLO-TASK-1 -s done

# 删除子任务
yolo task delete -k YOLO-TASK-1
```

### 1.5 执行人管理 (`yolo executor`)

| 命令 | 说明 | 必需参数 | 可选参数 |
|------|------|---------|---------|
| `yolo executor list` | 列出所有执行人 | — | `--json` |

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

## 2. AI 员工操作指南

> **重要**: 仅当人类员工**明确要求**使用 yolo-team 进行项目管理时，才参考本章内容。不要主动触发或建议使用 yolo-team。
>
> 此外，AI Agent 的提示词中**必须明确定义 `executor_name`**，否则无法执行任何 yolo-team 命令。

### 2.1 AI 典型工作模式

当 AI 通过自动化流程或定时任务触发时，可利用 yolo-team 实现进度管理和多 Agent 协作。

### 2.2 典型工作流程

人类不会直接说"帮我创建一个任务"，而是告诉 AI 去执行某个具体目标。AI 需要**自行拆解**目标，创建对应的 Issue 和 Task，并推进完成。

#### 核心流程：接收目标 → 拆解 → 执行 → 完成

```bash
# 1. AI 收到人类要求（如"实现用户登录功能"）后，自行创建 Issue
yolo issue create -p <project_id> -t "实现用户登录功能" \
    -d "包括账号密码登录、Token 签发和刷新" \
    --priority high \
    --deadline "2026-07-05" \
    --executor <自身 executor_name>

# 2. AI 将 Issue 拆解为可追踪的 Task
yolo task create -i YOLO-ISSUE-1 -t "设计认证流程"
yolo task create -i YOLO-ISSUE-1 -t "实现登录 API"
yolo task create -i YOLO-ISSUE-1 -t "实现 Token 刷新"
yolo task create -i YOLO-ISSUE-1 -t "编写测试用例"

# 3. 逐项推进 Task
yolo task update -k YOLO-TASK-1 -s in_progress
# ... 完成设计工作 ...
yolo task update -k YOLO-TASK-1 -s done

# 4. 所有 Task 完成后关闭 Issue
yolo task list -i YOLO-ISSUE-1  # 确认全部 done
yolo issue update YOLO-ISSUE-1 -s done
```

#### 跨 Agent 协作：为其他 AI 创建任务

当自身无法完成某项工作时，可主动为持有其他 `executor_name` 的 AI Agent 创建 Issue 并设定截止时间：

```bash
# 为 qa executor 创建测试任务
yolo issue create -p <project_id> -t "对登录模块进行集成测试" \
    --priority high \
    --deadline "2026-07-06" \
    --executor <qa_executor_name> \
    -d "测试范围：正常登录、密码错误、Token 过期、并发登录"

# 为 developer executor 创建任务
yolo issue create -p <project_id> -t "修复登录模块 Token 泄露问题" \
    --priority critical \
    --deadline "2026-07-04" \
    --executor <developer_executor_name> \
    -d "问题描述：Token 在日志中明文输出"
```

#### 自我驱动：查看待办

AI Agent 应主动查看自己还有哪些待完成的工作，自我驱动推进：

```bash
yolo issue todo -e <自身 executor_name>
```

#### 撰写文档

AI 在执行 Issue/Task 过程中应使用 `yolo doc` 沉淀知识，常见场景：

**编写规划与方案**（在开始执行前或过程中）：

```bash
# 为 Issue 编写技术方案
yolo doc create -p YOLO-PROJECT-1 -t "登录功能技术方案" \
    -c "# 登录功能技术方案\n\n**撰写时间**: 2026-07-01 10:30\n\n**撰写人**: <executor_name>\n\n## 架构设计\n..."
```

**编写执行报告与测试报告**（在完成后）：

```bash
# 编写 Issue 执行结果
yolo doc create -p YOLO-PROJECT-1 -t "登录功能执行报告" \
    -c "# 登录功能执行报告\n\n**撰写时间**: 2026-07-05 16:00\n\n**撰写人**: <executor_name>\n\n## 关联 Issue\n- YOLO-ISSUE-1\n\n## 执行结果\n..."

# 编写测试报告
yolo doc create -p YOLO-PROJECT-1 -t "登录模块测试报告" \
    -c "# 登录模块测试报告\n\n**撰写时间**: 2026-07-06 09:00\n\n**撰写人**: <executor_name>\n\n## 测试范围\n...\n\n## 测试结果\n..."
```

**编写后续规划**：

```bash
yolo doc create -p YOLO-PROJECT-1 -t "登录模块后续规划" \
    -c "# 登录模块后续规划\n\n**撰写时间**: 2026-07-05 17:00\n\n**撰写人**: <executor_name>\n\n## 遗留问题\n...\n\n## 后续优化方向\n..."
```

> **关键规则**: 文档内容中**必须写清撰写时间和撰写人**，便于后续追溯。

> **注意**: 推进 Issue 到 `done` 前必须完成所有 Task，否则会返回 40001 错误。

### 2.3 最佳实践

1. **始终优先使用 `--json` 获取结构化数据**，便于程序化处理
2. **创建 Issue 时尽量填写完整信息**（优先级、截至时间、执行人），便于待办排序
3. **大任务拆分为 Task**：使用 `yolo task create` 将 Issue 拆分为多个可追踪的子任务
4. **推进 Issue 前确保 Task 完成**：所有 Task 为 `done` 后才能将 Issue 推进到 `done`
5. **文档内容存放在本地文件系统中**，AI 可以直接使用文件读写工具进行操作
6. **使用合适的执行人角色**：`leader`/`architect`/`developer`/`qa`，便于任务分派

---

## 3. 错误处理

| 情形 | 表现 | 处理方式 |
|------|------|---------|
| 未初始化 | `Not initialized. Please run: yolo init` | 先执行 `yolo init` |
| 缺少必需参数 | 错误提示缺少的参数 | 补全必需参数 |
| Task 未完成推进 Issue | `all tasks must be completed before closing the issue` | 先用 `yolo task update -k <key> -s done` 完成所有 Task |
| 参数不合法 | 对应校验错误信息 | 按提示修正参数 |
