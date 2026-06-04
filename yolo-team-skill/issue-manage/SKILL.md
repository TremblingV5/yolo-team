# 任务管理技能 (Issue Management)

## 技能描述

管理 yolo-team 工作区中的任务（Issue），包括创建、查看、更新、删除任务，以及查看执行者的待办列表。

## 相关命令

| 命令 | 用途 |
|------|------|
| `yolo issue list` | 列出任务 |
| `yolo issue create` | 创建新任务 |
| `yolo issue info <key>` | 查看任务详情 |
| `yolo issue update <key>` | 更新任务 |
| `yolo issue delete <key>` | 删除任务 |
| `yolo issue todo` | 查看执行者待办列表 |

## 命令使用指南

### 1. 列出任务

```bash
yolo issue list [选项]
```

#### 选项
- `-p, --project <id>`: 按项目 ID 筛选
- `-s, --status <状态>`: 按状态筛选
- `-e, --executor <id>`: 按执行者 ID 筛选
- `--parent <id>`: 按父任务 ID 筛选（0 表示顶层任务）

#### 示例
```bash
# 列出所有任务
yolo issue list

# 列出项目 1 的所有任务
yolo issue list -p 1

# 列出状态为 "todo" 的任务
yolo issue list -s todo
```

#### 示例输出
```
KEY                STATUS         PRIORITY TITLE                   DEADLINE     REPO
YT001-001          todo           high     实现用户登录功能        2024-02-01   -
YT001-002          in-progress    medium   优化数据库查询          -            -
```

### 2. 创建新任务

```bash
yolo issue create -t <标题> -p <项目ID> [选项]
```

#### 参数说明
- `-t, --title` (必需): 任务标题
- `-p, --project` (必需): 项目 ID
- `-d, --description`: 任务描述
- `--priority`: 优先级
- `--parent`: 父任务 ID
- `--executor`: 执行者 ID
- `--deadline`: 截止日期 (YYYY-MM-DD)

#### 示例
```bash
# 创建基本任务
yolo issue create -t "修复登录bug" -p 1

# 创建完整任务
yolo issue create -t "实现API文档" -p 1 -d "为所有API接口生成文档" --priority high --deadline 2024-02-15
```

### 3. 查看任务详情

```bash
yolo issue info <任务 KEY>
```

#### 示例
```bash
yolo issue info YT001-001
```

#### 示例输出
```
Key:         YT001-001
Status:      todo
Title:       实现用户登录功能
Priority:    high
Description: 实现基于JWT的用户认证系统
Deadline:    2024-02-01
Executor:    Alice (developer)
Repo:        -
```

### 4. 更新任务

```bash
yolo issue update <任务 KEY> [选项]
```

#### 选项
- `-s, --status <状态>`: 更新状态
- `-t, --title <标题>`: 更新标题
- `--executor <id>`: 更新执行者 ID
- `--deadline <日期>`: 更新截止日期
- `--repo-url <url>`: 更新仓库 URL
- `--repo-name <名称>`: 更新仓库名称
- `--branch <分支名>`: 更新分支名称

#### 示例
```bash
# 更新任务状态为 in-progress
yolo issue update YT001-001 -s in-progress

# 更新任务截止日期
yolo issue update YT001-001 --deadline 2024-02-10
```

### 5. 删除任务

```bash
yolo issue delete <任务 KEY>
```

#### 示例
```bash
yolo issue delete YT001-001
```

### 6. 查看执行者待办列表

```bash
yolo issue todo -e <执行者 ID> [选项]
```

#### 选项
- `-e, --executor <id>` (必需): 执行者 ID
- `-n, --limit <数量>`: 最大结果数 (默认 20)

#### 示例
```bash
# 查看执行者 1 的待办任务
yolo issue todo -e 1

# 查看执行者 2 的前 10 个待办任务
yolo issue todo -e 2 -n 10
```

## 输入要求

- 必须先初始化工作区（运行 `yolo init`）
- 创建任务时，标题和项目 ID 是必填项
- 日期格式必须为 YYYY-MM-DD

## 任务状态

常见状态包括：
- `todo`: 待处理
- `in-progress`: 进行中
- `done`: 已完成
- `blocked`: 被阻塞

## 优先级

常见优先级包括：
- `low`: 低
- `medium`: 中
- `high`: 高
- `critical`: 紧急

## 质量标准

- 任务标题应简洁清晰地描述任务内容
- 任务描述应详细说明任务的目标和验收标准
- 合理设置截止日期和优先级
- 及时更新任务状态以反映当前进度
- 为任务分配合适的执行者
