# yolo-client

yolo-team 的 Web 前端，使用 React + Umi 构建。编译为静态文件后通过 Go `embed` 嵌入到 `yolo-cli`。

## 定位

`yolo-client` 是前端源码目录。它不独立部署，而是作为构建管线的一部分：

1. **开发时**：通过 `pnpm dev` 独立启动开发服务器（:8000），API 代理到后端
2. **构建时**：通过 `pnpm build` 产出静态文件到 `dist/`，然后复制到 `yolo-cli/client-dist/`
3. **运行时**：由 `yolo-cli` 通过 `go:embed` 嵌入并统一提供服务

## 技术栈

| 组件 | 选择 |
|------|------|
| 框架 | React 18 + Umi 4 |
| UI 组件库 | Ant Design |
| API 客户端 | restful-react（从 OpenAPI 生成 TS 类型） |
| 状态管理 | Umi 内置（基于 React Context） |
| 请求库 | umi-request / axios |
| 包管理 | pnpm |

## 核心页面

- `/` — 全局看板（主页面），按 8 种状态分列展示 Issue，顶栏项目筛选
- `/project/:key/manage` — 项目管理（编辑名称与描述）
- `/project/:key/docs` — 项目文档列表
- `/doc/:key` — 文档详情与编辑

> Issue 详情通过**左侧抽屉**展示，不占用独立路由。

## 看板视图

项目详情页以看板形式展示 Issue，**按八种固定状态分列**：

| 列 | 状态值 | 颜色 |
|---|--------|------|
| 已创建 | `created` | 灰 |
| 方案设计 | `design` | 蓝 |
| 方案评审 | `review` | 紫 |
| 代码实现 | `implementation` | 黄 |
| QA质检 | `qa` | 橙 |
| 待审查 | `pending_review` | 粉 |
| 已完成 | `done` | 绿 |
| 已归档 | `archived` | 浅灰 |

每个卡片支持展开显示子任务（仅一级）。

## 构建与嵌入流程

```
yolo-client/
    │
    │  pnpm build
    ▼
dist/                        # Umi 构建产物
├── index.html
├── umi.js
├── umi.css
└── ...
    │
    │  复制到 yolo-cli/client-dist/
    ▼
yolo-cli/client-dist/        # ★ 被 //go:embed client-dist/* 嵌入
```

## 开发

### 日常开发（热更新）

```bash
cd yolo-client

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev
# → http://localhost:8000
# API 请求自动代理到 http://localhost:8080（yolo-cli serve）
```

### 生产构建

```bash
cd yolo-client

# 构建
pnpm build
# → 产物在 dist/

# 复制到 embed 目录
cp -r dist ../yolo-cli/client-dist/
```

## 代理配置

开发模式下，API 请求代理到后端（由 `yolo serve` 提供），在 `.umirc.ts` 中配置：

```ts
export default {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
};
```

## 目录结构

```
yolo-client/
├── src/
│   ├── pages/              # 页面组件（项目列表、项目详情、文档列表等）
│   ├── components/         # 通用组件（卡片、列容器、代码关联表单等）
│   ├── services/           # API 调用封装
│   └── models/             # 类型定义
├── dist/                   # 构建产物 → 复制到 yolo-cli/client-dist/
├── package.json
└── .umirc.ts
```
