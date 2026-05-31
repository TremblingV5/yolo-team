# 04 — Document API

> 参考：[api.md](../docs/api.md) 4 节、[models.md](../docs/models.md) 4 节

文档采用"DB 存元数据 + 文件系统存内容"方案。

## 子任务

### 4.1 文档存储路径

- [x] `{workspace}/documents/{project-id}/` 目录自动创建
- [x] 文件名 = `{title}.md`
- [x] 实现变更为：workspace 通过 Handler 构造函数传入，使用 `fileutils.BuildDocPath(workspace, projectID, title)` 构建路径

### 4.2 路由与 Handler

- [x] `GET /api/v1/projects/:key/documents` — 文档列表（仅元数据，不含正文）
  - 响应含 `key`（YOLO-DOC-{id}）、`title`、`timestamps`
- [x] `POST /api/v1/projects/:key/documents` — 创建文档
  - body: {title(必填,≤128,作为文件名), content(必填,Markdown)}
  - 写入 `{workspace}/documents/{project-id}/{title}.md`
  - 响应含 `file_path` 字段（本地绝对路径）
- [x] `GET /api/v1/documents/:key` — 文档详情
  - 从文件系统读取正文，响应含 `content` + `file_path`
- [x] `PUT /api/v1/documents/:key` — 更新文档
  - 可更新 title（重命名文件）和 content（覆写文件）
- [x] `DELETE /api/v1/documents/:key` — 删除文档
  - 删除 DB 记录 + 对应 `.md` 文件

### 4.3 Key 生成

- [x] 文档 key 格式 `YOLO-DOC-{id}`，使用 `keyutils.Document(id)`，"先 insert → 回写 key"策略

### 4.4 Handler 层

- [x] ~~Service 层已删除，Handler 直接调用 Repository~~
- [x] `DocumentHandler.List(req)` — 查 DB 元数据
- [x] `DocumentHandler.Create(req)` — 写 DB + 写文件
- [x] `DocumentHandler.Get(req)` — 查 DB 元数据 + 读文件内容
- [x] `DocumentHandler.Update(req)` — 更新 DB + 文件（title 变更时重命名）
- [x] `DocumentHandler.Delete(req)` — 删 DB 记录 + 删文件
- [x] 所有 Handler 方法使用泛型 `Wrap` 包装器

### 4.5 文件操作

- [x] 封装文件操作工具函数在 `internal/utils/fileutils/`：`fileutils.Write(path, content)` / `fileutils.Read(path)` / `fileutils.BuildDocPath(...)`
- [ ] 路径安全校验：防止路径穿越

## 验收标准

- [x] 文档 CRUD 全部可用
- [x] 文件实际写入 `{workspace}/documents/{id}/{title}.md`
- [x] 详情响应含 `file_path`，可直接用于本地文件读写
- [x] 更新 title 时文件名同步修改
- [x] 删除时 DB 和文件同步清除
- [x] Key 格式 `YOLO-DOC-{id}` 正确
