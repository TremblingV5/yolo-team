package testcases

import (
	"testing"

	"yolo-team/yolo-cli-test/suite"

	"github.com/stretchr/testify/assert"
)

// TestExecutorCreateList
//   场景：验证执行人的创建和列表功能
//
//   步骤：
//     1) 创建执行人 Alice，soul 为 "frontend expert"
//     2) 创建执行人 Bob，soul 为 "backend expert"
//     3) 执行 executor list --json
//     4) 断言返回 2 条记录
func TestExecutorCreateList(t *testing.T) {
	suite.Run("executor", "create", "-n", "Alice", "-s", "frontend expert")
	suite.Run("executor", "create", "-n", "Bob", "-s", "backend expert")

	out := suite.Run("--json", "executor", "list")
	items := suite.JSONDecode(t, out)
	assert.Len(t, items, 2)
}
