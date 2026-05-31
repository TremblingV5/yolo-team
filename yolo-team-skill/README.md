# yolo-team-skill

AI 员工技能定义与行为配置。

## 定位

yolo-team-skill 是 yolo-team 的 **AI 技能体系** 模块。它定义了 AI 员工（AI Employee）在 yolo-team 工作流系统中所能执行的各类技能任务，以及对应的行为规范与配置。

## 核心理念

在 yolo-team 中，AI 不仅是工具，更是"员工"。和人类员工一样，AI 员工通过技能（Skill）来参与工作：

- 每个 Skill 定义一个 AI 员工可以独立完成的工作类型（如：代码审查、文档生成、测试用例编写等）
- Skill 包含明确的输入规范、输出格式和质量标准
- AI 员工通过 CLI（`yolo-cli`）接收任务、反馈进度、提交成果

## 技能示例（规划中）

| 技能名称 | 用途 |
|----------|------|
| `code-review` | 对 PR/代码变更执行自动化审查 |
| `doc-generate` | 根据代码变更自动生成/更新技术文档 |
| `test-write` | 根据功能描述自动生成测试用例 |
| `bug-analyze` | 分析 Bug 报告并给出修复建议 |
| `release-note` | 根据 Git 历史自动生成发布说明 |

## 目录结构（规划中）

```
yolo-team-skill/
├── README.md              # 本文件
├── code-review/
│   ├── SKILL.md          # 技能定义
│   └── prompt.md         # LLM 提示词模板
├── doc-generate/
│   ├── SKILL.md
│   └── prompt.md
└── ...
```

## 状态

> 本模块目前处于规划阶段，尚未实现具体技能。待后续版本迭代中完善。
