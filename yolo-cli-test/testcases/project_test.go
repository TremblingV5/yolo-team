package testcases

import (
	"testing"

	"yolo-team/yolo-cli-test/suite"

	"github.com/stretchr/testify/assert"
)

// TestProjectCreateList
//   场景：验证项目创建、列表和详情功能
//
//   步骤：
//     1) 创建项目 alpha，带描述 "first project"
//     2) 执行 project list --json，断言返回 1 条记录且 key 以 YOLO-PROJECT 开头
//     3) 执行 project info <key>，断言返回的 JSON 中包含 name="alpha" 和 description="first project"
//     4) 创建第二个项目 beta
//     5) 再次执行 project list --json，断言返回 2 条记录
func TestProjectCreateList(t *testing.T) {
	out := suite.Run("project", "create", "-n", "alpha", "-d", "first project")
	assert.Contains(t, out, "Project created")

	out = suite.Run("--json", "project", "list")
	items := suite.JSONDecode(t, out)
	assert.Len(t, items, 1)
	assert.Contains(t, items[0]["key"], "YOLO-PROJECT")

	key := items[0]["key"].(string)
	out = suite.Run("--json", "project", "info", key)
	assert.Contains(t, out, "alpha")
	assert.Contains(t, out, "first project")

	suite.Run("project", "create", "-n", "beta")
	out = suite.Run("--json", "project", "list")
	items = suite.JSONDecode(t, out)
	assert.Len(t, items, 2)
}

// TestProjectInfoNonExistent
//   场景：查询不存在的项目 key 时应返回明确的错误信息
//
//   步骤：
//     1) 执行 project info YOLO-PROJECT-99（不存在的 key）
//     2) 断言命令以非零退出码结束
//     3) 断言错误输出中包含 "not found"
func TestProjectInfoNonExistent(t *testing.T) {
	out := suite.RunFail("project", "info", "YOLO-PROJECT-99")
	assert.Contains(t, out, "not found")
}
