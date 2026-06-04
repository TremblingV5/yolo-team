# 文档管理技能 (Document Management)

## 技能描述

管理 yolo-team 工作区中的文档，包括创建文档、列出项目文档和查看文档详情。文档以 Markdown 格式存储在工作区中。

## 相关命令

| 命令 | 用途 |
|------|------|
| `yolo doc list -p <项目 KEY>` | 列出项目中的文档 |
| `yolo doc create -p <项目 KEY> -t <标题>` | 创建新文档 |
| `yolo doc info <文档 KEY>` | 查看文档详情 |

## 命令使用指南

### 1. 列出项目文档

```bash
yolo doc list -p <项目 KEY>
```

#### 参数说明
- `-p, --project` (必需): 项目 KEY

#### 示例
```bash
# 列出 YT001 项目的所有文档
yolo doc list -p YT001
```

#### 示例输出
```
KEY                TITLE                CREATOR    UPDATED
YT001-DOC001       项目需求文档         CLI        2024-01-15 10:30:00
YT001-DOC002       API设计文档          CLI        2024-01-14 15:20:00
```

### 2. 创建新文档

```bash
yolo doc create -p <项目 KEY> -t <标题> [选项]
```

#### 参数说明
- `-p, --project` (必需): 项目 KEY
- `-t, --title` (必需): 文档标题
- `-c, --content`: 文档内容

#### 示例
```bash
# 创建基本文档
yolo doc create -p YT001 -t "项目需求文档"

# 创建带内容的文档
yolo doc create -p YT001 -t "API设计文档" -c "# API 设计\n\n这是 API 设计文档的内容..."
```

#### 创建成功输出
```
Document created: YT001-DOC001
Path: /path/to/workspace/documents/YT001-DOC001.md
```

### 3. 查看文档详情

```bash
yolo doc info <文档 KEY>
```

#### 示例
```bash
yolo doc info YT001-DOC001
```

#### 示例输出
```
File: /path/to/workspace/documents/YT001-DOC001.md

# 项目需求文档

这是项目需求文档的内容...
```

## 文档存储

- 文档以 Markdown 格式存储在工作区的 `documents` 目录中
- 每个文档有唯一的 KEY 作为标识符
- 文档文件以 KEY 命名，后缀为 `.md`

## 输入要求

- 必须先初始化工作区（运行 `yolo init`）
- 创建文档时，项目 KEY 和标题是必填项
- 文档内容可以是空的，但建议包含有意义的内容

## 输出格式

- 列表格式：显示文档 KEY、标题、创建者和更新时间
- 详情格式：显示文件路径和完整文档内容
- JSON 格式：使用 `--json` 选项获取结构化数据

## 质量标准

- 文档标题应准确反映文档内容
- 使用 Markdown 格式保持文档的可读性和结构
- 定期更新文档以保持信息的及时性
- 文档应包含足够的细节，便于团队成员理解
- 合理组织文档内容，使用标题、列表等元素增强可读性
