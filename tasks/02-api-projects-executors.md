# 02 — Project & Executor API

> 参考：[api.md](../docs/api.md) 1、2 节

## 2A. 项目 API

### 2A.1 路由与 Handler

- [ ] 为每个 handler 添加 swaggo 注释（`@Summary`、`@Tags`、`@Param`、`@Success`、`@Failure`）
- [x] `GET /api/v1/projects` — 列表，返回 id/key/name/description/timestamps
- [x] `POST /api/v1/projects` — 创建，body: {name(必填,≤32), description(≤128)}
  - 创建后调用 Key 生成策略写入 `key`
- [x] `GET /api/v1/projects/:key` — 详情，key 路由参数
- [x] `PUT /api/v1/projects/:key` — 更新 name/description
- [x] `DELETE /api/v1/projects/:key` — 删除

### 2A.2 Handler 层

- [x] ~~Service 层已删除，Handler 直接调用 Repository~~
- [x] `ProjectHandler.List()` / `Create()` / `Get()` / `Update()` / `Delete()`
- [x] 参数校验：name ≤ 32、description ≤ 128（通过 model.Validate()）

### 2A.3 通用约定

- [x] 统一响应结构 `{code, data, message}` 定义在 `internal/common/common.go`
- [x] 错误码：0 成功、40001 参数错误、40401 不存在
- [x] 所有 Handler 方法使用泛型 `Wrap` 包装器：入参 `(ctx, req)` → 出参 `(resp, error)`

---

## 2B. 执行人 API

### 2B.1 路由与 Handler

- [x] `GET /api/v1/executors` — 列表，返回 id/name/role/soul/timestamps
- [x] `POST /api/v1/executors` — 创建，body: {name(必填,≤64,不含空格,唯一), soul(可选)}
  - `role` 不通过 API 设置，默认 `developer`
- [x] `DELETE /api/v1/executors/:name` — 按 name 删除

### 2B.2 Handler 层

- [x] ~~Service 层已删除，Handler 直接调用 Repository~~
- [x] `ExecutorHandler.List()` / `Create()` / `Delete()`
- [x] 参数校验：name ≤ 64、不含空格、唯一（40901）
- [x] 不支持通过 API 更新 role（仅页面可改）

### 2B.3 Router 注册

- [x] `internal/router/router.go` — 注册所有路由分组
- [x] `/api/v1/projects` → projectHandler
- [x] `/api/v1/executors` → executorHandler

## 验收标准

- [x] 项目 CRUD 全部可用，通过 `:key` 访问
- [x] 创建项目后返回正确的 `YOLO-PROJECT-{id}` key
- [x] 创建执行人成功，重名返回 40901
- [x] 执行人不含 role 设置接口
- [x] 统一 JSON 响应格式正确
