package testcases

import (
	"testing"

	"yolo-team/yolo-cli-test/suite"

	"github.com/stretchr/testify/assert"
)

// TestIssueCreateMissingRequired
//   场景：创建 issue 时缺少必填参数，应输出明确的错误提示
//
//   步骤：
//     1) 创建项目 demo 作为前置条件
//     2) 执行 issue create -p 1（缺少 -t title），断言输出包含 "title"
//     3) 执行 issue create -t "no project"（缺少 -p project），断言输出包含 "project"
func TestIssueCreateMissingRequired(t *testing.T) {
	suite.Run("project", "create", "-n", "demo")

	out := suite.RunFail("issue", "create", "-p", "1")
	assert.Contains(t, out, "title")

	out = suite.RunFail("issue", "create", "-t", "no project")
	assert.Contains(t, out, "project")
}

// TestIssueTodo
//   场景：查询指定执行人的待办列表，应返回所有非完成/非归档的 issue
//
//   步骤：
//     1) 创建项目 demo
//     2) 创建执行人 worker
//     3) 创建 2 个 issue（task A 和 task B），均指派给 executor 1
//     4) 执行 issue todo -e 1 --json
//     5) 断言返回 2 条记录（created 状态属于待办）
func TestIssueTodo(t *testing.T) {
	suite.Run("project", "create", "-n", "demo")
	suite.Run("executor", "create", "-n", "worker")

	suite.Run("issue", "create", "-p", "1", "-t", "task A", "--executor", "1")
	suite.Run("issue", "create", "-p", "1", "-t", "task B", "--executor", "1")

	out := suite.Run("--json", "issue", "todo", "-e", "1")
	items := suite.JSONDecode(t, out)
	assert.Len(t, items, 2)
}

// TestIssueFilterByStatus
//   场景：按状态过滤 issue 列表，验证状态筛选条件生效
//
//   步骤：
//     1) 创建项目 demo
//     2) 创建 2 个 issue（task X 和 task Y），默认状态为 created
//     3) 执行 issue list --json -p 1，断言返回 2 条
//     4) 执行 issue list --json -p 1 -s in_progress，断言返回 0 条
//     5) 将第一个 issue 更新为 in_progress 状态
//     6) 再次执行 -s in_progress 过滤，断言返回 1 条
func TestIssueFilterByStatus(t *testing.T) {
	suite.Run("project", "create", "-n", "demo")
	suite.Run("issue", "create", "-p", "1", "-t", "task X")
	suite.Run("issue", "create", "-p", "1", "-t", "task Y")

	out := suite.Run("--json", "issue", "list", "-p", "1")
	items := suite.JSONDecode(t, out)
	assert.Len(t, items, 2)

	out = suite.Run("--json", "issue", "list", "-p", "1", "-s", "in_progress")
	items = suite.JSONDecode(t, out)
	assert.Len(t, items, 0)

	out = suite.Run("--json", "issue", "list", "-p", "1")
	items = suite.JSONDecode(t, out)
	key := items[0]["key"].(string)
	suite.Run("issue", "update", key, "--status", "in_progress")

	out = suite.Run("--json", "issue", "list", "-p", "1", "-s", "in_progress")
	items = suite.JSONDecode(t, out)
	assert.Len(t, items, 1)
}
