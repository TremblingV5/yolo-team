package testcases

import (
	"testing"

	"yolo-team/yolo-cli-test/suite"

	"github.com/stretchr/testify/assert"
)

// TestIssueCRUD
//
//	场景：端到端验证 Issue 完整生命周期 — 创建 → 查看 → 更新 → 删除
//	该测试串联了 project、executor、issue 三个命令组，作为冒烟测试的核心用例
//
//	步骤：
//	  1) 创建项目 demo
//	  2) 创建执行人 Charlie
//	  3) 执行 issue create -p 1 -t "fix login bug" --priority high --executor 1
//	  4) 执行 issue list --json -p 1，断言返回 1 条记录
//	  5) 断言 title="fix login bug", priority="high"
//	  6) 执行 issue info <key> --json，断言包含 key
//	  7) 执行 issue update <key> --status in_progress --title "fix login bug [urgent]"
//	  8) 执行 issue info <key> --json，断言包含 "in_progress" 和 "[urgent]"
//	  9) 执行 issue delete <key>
//	 10) 执行 issue list --json -p 1，断言返回 0 条记录
func TestIssueCRUD(t *testing.T) {
	suite.Run("project", "create", "-n", "demo")
	suite.Run("executor", "create", "-n", "Charlie")

	out := suite.Run("issue", "create", "-p", "1", "-t", "fix login bug",
		"-d", "users cannot login", "--priority", "high", "--executor", "1")
	assert.Contains(t, out, "Issue created")

	out = suite.Run("--json", "issue", "list", "-p", "1")
	items := suite.JSONDecode(t, out)
	assert.Len(t, items, 1)
	assert.Equal(t, "fix login bug", items[0]["title"])
	assert.Equal(t, "high", items[0]["priority"])

	key := items[0]["key"].(string)

	out = suite.Run("--json", "issue", "info", key)
	assert.Contains(t, out, key)

	out = suite.Run("issue", "update", key, "--status", "in_progress", "--title", "fix login bug [urgent]")
	assert.Contains(t, out, "Issue updated")

	out = suite.Run("--json", "issue", "info", key)
	assert.Contains(t, out, "in_progress")
	assert.Contains(t, out, "[urgent]")

	out = suite.Run("issue", "delete", key)
	assert.Contains(t, out, "Issue deleted")

	out = suite.Run("--json", "issue", "list", "-p", "1")
	items = suite.JSONDecode(t, out)
	assert.Len(t, items, 0)
}
