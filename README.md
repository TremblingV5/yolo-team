# Yolo-Team

基于 **AI 员工** 的工作流系统。

## 核心理念

让人与 AI 不再是工具与使用者的关系，而是在同一套工作流体系中的 **协作伙伴**。yolo-team 以"项目（Project）+ 任务卡片（Issue）"为核心模型，让人类员工和 AI 员工共同管理任务、流转状态、沉淀知识。

## 单一二进制架构

除 `yolo-team-skill`（纯文档）外，整个项目仅编译产出 **一个二进制文件**：

```
yolo-client (React + Umi)  ──pnpm build──►  dist/
                                                 │
                                                 │  go:embed
                                                 ▼
                                        yolo-cli (Go + Cobra)
                                                 │
                                           go build
                                                 ▼
                                            yolo.exe  ★ 唯一产物
```

- `yolo.exe serve` — 启动 HTTP 服务，同时提供 API 和前端页面
- `yolo.exe project list` / `yolo.exe issue create` … — CLI 操作项目和任务

## 项目结构

| 目录 | 类型 | 说明 |
|------|------|------|
| [yolo-cli](./yolo-cli/) | Go | **唯一代码目录**，包含 CLI + HTTP 服务 + 数据访问等全部后端代码 |
| [yolo-client](./yolo-client/) | React + Umi 源码 | 前端源码，编译为静态文件后嵌入到 yolo-cli |
| [yolo-team-skill](./yolo-team-skill/) | Markdown / 配置 | AI 员工技能定义与行为配置 |
| [docs](./docs/) | — | 项目设计文档 |

## 核心功能

- **项目管理** — 创建项目、每个项目有唯一 Key（`YOLO-PROJECT-{id}`）
- **任务卡片（Issue）** — 固定五种状态流转，支持截至时间、优先级设置
- **执行人（Executor）** — 管理系统中的执行人，含角色（leader/architect/developer/qa），卡片可分配执行人
- **待办查询** — 按执行人查询待办卡片，按截至时间和优先级排序
- **父子卡片** — 支持一级子任务拆分
- **代码仓库关联** — 每个卡片可关联代码仓库 URL、名称、分支
- **项目文档** — 每个项目下可创建和管理文档，内容以 Markdown 文件存储，响应含本地路径供 AI 直接读写
- **CLI 操作** — 终端内完成所有操作，适合 AI Agent 调用
- **AI 技能体系** — 定义 AI 员工可执行的工作技能

## 快速开始

### 构建（生产模式）

```bash
# 1. 构建前端
cd yolo-client && pnpm build

# 2. 将前端产物复制到 yolo-cli 的 embed 目录
cp -r dist ../yolo-cli/client-dist/

# 3. 编译唯一二进制
cd ../yolo-cli && go build -o yolo.exe .

# 4. 启动一切
./yolo.exe serve
# → API:      http://localhost:8080/api/v1/
# → 前端页面: http://localhost:8080/
```

### 开发模式

```bash
# 终端 1：启动服务端
cd yolo-cli && go run . serve

# 终端 2：启动前端开发服务器（热更新）
cd yolo-client && pnpm dev
# → http://localhost:8000，API 自动代理到 :8080

# CLI 直接操作
yolo issue create -t "完成登录模块" -p 1 --executor 1 --deadline "2026-06-30"
yolo issue update YOLO-ISSUE-1 -s design
yolo issue todo -e 1
yolo doc create -p 1 -t "需求文档"
```

## 技术选型

- **数据库**：SQLite / MySQL（通过 GORM + gorm/gen 切换，类型安全查询）
- **后端**：Go + Gin + GORM，swaggo/swag 生成 OpenAPI 文档
- **前端**：React + Umi + Ant Design，restful-react 从 OpenAPI 生成 TS 类型
- **CLI**：Go + Cobra
- **前端嵌入**：Go `embed` 标准库

## 设计文档

详见 [docs/architecture.md](./docs/architecture.md)

## In Progress

1. agent tui，用于便捷管理
2. 页面上的agent
