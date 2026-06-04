# 工作区设置与服务技能 (Workspace Setup & Server)

## 技能描述

管理 yolo-team 工作区的初始化和 HTTP 服务器的启动，这是使用 yolo-cli 的第一步。

## 相关命令

| 命令 | 用途 |
|------|------|
| `yolo init` | 初始化 yolo-team 工作区 |
| `yolo serve` | 启动 HTTP 服务器 |

## 命令使用指南

### 1. 初始化工作区

```bash
yolo init [选项]
```

#### 说明
- 第一次使用 yolo-cli 时必须先初始化工作区
- 可以交互式选择工作区位置，也可以直接指定
- 初始化后会在工作区创建 `documents` 目录用于存储文档

#### 选项
- `-w, --workspace <路径>`: 直接指定工作区路径（跳过交互式）

#### 示例

##### 交互式初始化
```bash
yolo init
```

运行后会看到交互式向导：
```
  ╔══════════════════════════════════╗
  ║     Yolo-Team Setup Wizard       ║
  ╚══════════════════════════════════╝

? Choose workspace location:
  ▸ Default: /path/to/default/workspace
    Custom path (enter manually)
```

选择后确认即可完成初始化。

##### 直接指定路径
```bash
yolo init -w /path/to/my/workspace
```

#### 初始化成功输出
```
Initialized: /path/to/my/workspace
```

### 2. 启动 HTTP 服务器

```bash
yolo serve [选项]
```

#### 说明
- 启动一个 HTTP 服务器，提供 Web 界面和 API 访问
- 服务器默认监听 8080 端口
- 可以通过浏览器访问 http://localhost:8080 使用 Web 界面

#### 选项
- `-p, --port <端口>`: 指定监听端口（默认 8080）

#### 示例
```bash
# 使用默认端口启动
yolo serve

# 使用指定端口启动
yolo serve -p 3000
```

#### 启动成功输出
```
Server started at http://localhost:8080
```

## 工作区结构

初始化后的工作区包含以下内容：

```
workspace/
├── documents/          # 文档存储目录
└── (其他系统文件)
```

## 输入要求

### 初始化工作区
- 工作区路径必须是有效的文件系统路径
- 如果目录不存在，会自动创建
- 需要有足够的文件系统权限

### 启动服务器
- 端口必须未被占用
- 需要有网络监听权限
- 工作区必须已初始化

## 全局选项

所有 yolo 命令都支持以下全局选项：

- `--json`: 以 JSON 格式输出结果
- `-h, --help`: 显示帮助信息

## 使用流程

1. 首次使用：运行 `yolo init` 初始化工作区
2. 创建项目：使用 `yolo project` 命令管理项目
3. 创建任务：使用 `yolo issue` 命令管理任务
4. （可选）启动 Web 界面：运行 `yolo serve`

## 质量标准

- 工作区路径应选择合适的位置，便于访问和备份
- 确保工作区目录有足够的磁盘空间
- 定期备份工作区数据
- 服务器端口选择避免与其他服务冲突
- 不要在生产环境使用默认端口
