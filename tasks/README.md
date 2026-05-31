# 任务拆解

## 依赖关系

```
01-database-and-models
    │
    ├──► 02-api-projects-executors
    │       │
    │       └──► 03-api-issues ──► 04-api-documents
    │                                       │
    │                                       ▼
    └──────────────────────────────► 05-cli
                                        │
                                        ▼
                                    06-frontend
                                        │
                                        ▼
                                    07-embed-and-serve
```

## 任务索引

| 序号 | 任务 | 说明 | 预估复杂度 | 状态 |
|------|------|------|------------|------|
| [01](./01-database-and-models.md) | 数据库与模型 | SQLite 建表、Go 模型定义、迁移、Key 生成 | ⭐⭐ | ✅ 完成 |
| [02](./02-api-projects-executors.md) | Project & Executor API | 项目 CRUD、执行人管理（role+name+soul） | ⭐⭐ | ✅ 完成 |
| [03](./03-api-issues.md) | Issue API | 卡片 CRUD、8 状态流转、角色权限校验、待办查询 | ⭐⭐⭐⭐ | ✅ 完成 |
| [04](./04-api-documents.md) | Document API | 文档 CRUD、文件系统存储、file_path 响应 | ⭐⭐ | ✅ 完成 |
| [05](./05-cli.md) | CLI 命令 | 5 组命令行工具（project/issue/executor/doc/serve） | ⭐⭐⭐ | ✅ 大部分完成 |
| [06](./06-frontend.md) | 前端页面 | 全局看板、抽屉详情、项目管理、文档页面 | ⭐⭐⭐⭐ | 🟡 核心完成 |
| [07](./07-embed-and-serve.md) | 嵌入与构建 | 前端 embed、yolo serve、构建管线 | ⭐⭐ | 🔴 未开始 |

## 实现变更记录

以下是在开发过程中相比原始设计的重要变更：

| 模块 | 原始设计 | 实际实现 |
|------|----------|----------|
| 查询生成 | `gorm/gen` 库 | `gorm cli`（`gorm.io/cli/gorm`）生成 field helpers，配合 `gorm.G[T]()` 泛型 API |
| 架构分层 | handler → service → repo | handler → repo（service 层已删除，校验逻辑移至 model） |
| 前端框架 | Umi 4 | React 18 + Vite + Ant Design 5 |
| API 客户端 | `restful-react` 自动生成 | 手写 `api.ts` |
| 公共类型 | `handler/response.go` | `internal/common/common.go`（Response、AppError） |
| 工具函数 | 散落在各 handler/repo 中 | `internal/utils/{timeutils,fileutils,keyutils,dbutils}/` |
| 数据库 | 支持 SQLite / MySQL 切换 | 仅 SQLite |
| Swagger | swaggo/swag 注解生成 API 文档 | 未实现 |

## 推荐执行顺序

1. **01 → 02 → 03 → 04**（后端 API 按依赖链）
2. **05**（CLI 可与 02-04 并行，共用 HTTP 客户端封装）
3. **06**（前端需等 API 就绪后开始）
4. **07**（最后集成）
