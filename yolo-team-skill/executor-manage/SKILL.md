# 执行者管理技能 (Executor Management)

## 技能描述

管理 yolo-team 工作区中的执行者（Executor），包括创建新执行者和查看所有执行者列表。执行者可以是人类或 AI 员工。

## 相关命令

| 命令 | 用途 |
|------|------|
| `yolo executor list` | 列出所有执行者 |
| `yolo executor create` | 创建新执行者 |

## 命令使用指南

### 1. 列出所有执行者

```bash
yolo executor list
```

#### 说明
- 显示所有执行者的列表，包含 ID、姓名、角色和创建时间
- 支持 `--json` 选项以 JSON 格式输出

#### 示例输出
```
ID   NAME                 ROLE         CREATED
1    Alice                developer    2024-01-01 09:00:00
2    Bob                  designer     2024-01-02 14:30:00
3    AI-Assistant         assistant    2024-01-03 10:15:00
```

### 2. 创建新执行者

```bash
yolo executor create -n <名称> [选项]
```

#### 参数说明
- `-n, --name` (必需): 执行者名称
- `-s, --soul` (可选): 执行者的灵魂/提示词（用于 AI 执行者）

#### 示例
```bash
# 创建基本执行者
yolo executor create -n "Charlie"

# 创建带灵魂配置的 AI 执行者
yolo executor create -n "AI-Coder" -s "你是一个专业的 Go 语言开发者，擅长编写高质量代码。"
```

## 执行者角色

系统会根据执行者的特性自动分配角色，常见角色包括：
- `developer`: 开发者
- `designer`: 设计师
- `manager`: 经理
- `assistant`: 助手
- `reviewer`: 审核者

## 输入要求

- 必须先初始化工作区（运行 `yolo init`）
- 执行者名称不能为空
- 对于 AI 执行者，`soul` 参数可以定义其行为特性

## 输出格式

- 文本格式：人类可读的表格或简单反馈
- JSON 格式：使用 `--json` 选项获取结构化数据

## 质量标准

- 执行者名称应清晰可辨
- AI 执行者的 `soul` 配置应明确其角色和行为规范
- 定期维护执行者列表，确保信息的准确性
- 合理分配任务给合适的执行者
