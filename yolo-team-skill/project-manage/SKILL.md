# 项目管理技能 (Project Management)

## 技能描述

管理 yolo-team 工作区中的项目，包括创建项目、查看项目列表和查看项目详情。

## 相关命令

| 命令 | 用途 |
|------|------|
| `yolo project list` | 列出所有项目 |
| `yolo project create` | 创建新项目 |
| `yolo project info <key>` | 查看项目详情 |

## 命令使用指南

### 1. 列出所有项目

```bash
yolo project list
```

#### 说明
- 显示所有项目的列表，包含项目 ID、KEY、名称和更新时间
- 支持 `--json` 选项以 JSON 格式输出

#### 示例输出
```
ID   KEY                  NAME                             UPDATED
1    YT001                Yolo Team Project                2024-01-15 10:30:00
2    TSK002               Task System                      2024-01-14 15:20:00
```

### 2. 创建新项目

```bash
yolo project create -n <项目名称> -d <项目描述>
```

#### 参数说明
- `-n, --name` (必需): 项目名称
- `-d, --description` (可选): 项目描述

#### 示例
```bash
yolo project create -n "New Project" -d "这是一个新的项目"
```

### 3. 查看项目详情

```bash
yolo project info <项目 KEY>
```

#### 参数说明
- `<项目 KEY>` (必需): 项目的唯一标识符

#### 示例
```bash
yolo project info YT001
```

#### 示例输出
```
Key:         YT001
Name:        Yolo Team Project
Description: 团队协作项目
Created:     2024-01-01 09:00:00
Updated:     2024-01-15 10:30:00
```

## 输入要求

- 必须先初始化工作区（运行 `yolo init`）
- 项目名称不能为空
- 项目 KEY 是系统自动生成的唯一标识符

## 输出格式

- 文本格式：人类可读的表格或键值对
- JSON 格式：使用 `--json` 选项获取结构化数据

## 质量标准

- 确保项目名称清晰明确
- 项目描述应简洁概括项目的目标和范围
- 定期更新项目信息以保持准确性
