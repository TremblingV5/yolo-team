package testcases

import (
	"testing"

	"yolo-team/yolo-cli-test/suite"

	"github.com/stretchr/testify/assert"
)

// TestDocCRUD
//   场景：验证文档的创建、列表和详情功能，包含 Markdown 内容的正确写入
//
//   步骤：
//     1) 创建项目 docs-project
//     2) 通过 project list --json 获取项目 key
//     3) 执行 doc create -p <key> -t readme -c "# Hello\n\nWelcome to yolo-team"
//     4) 断言输出包含 "Document created"
//     5) 执行 doc list --json -p <key>，断言 1 条记录
//     6) 断言文档 title 为 "readme"，creator 为 CLI
//     7) 执行 doc info <docKey> --json
//     8) 断言返回的 JSON 中包含 "Hello" 和 "yolo-team"
func TestDocCRUD(t *testing.T) {
	suite.Run("project", "create", "-n", "docs-project")

	out := suite.Run("--json", "project", "list")
	items := suite.JSONDecode(t, out)
	projKey := items[0]["key"].(string)

	out = suite.Run("doc", "create", "-p", projKey, "-t", "readme", "-c", "# Hello\n\nWelcome to yolo-team")
	assert.Contains(t, out, "Document created")

	out = suite.Run("--json", "doc", "list", "-p", projKey)
	docs := suite.JSONDecode(t, out)
	assert.Len(t, docs, 1)
	assert.Equal(t, "readme", docs[0]["title"])
	assert.Equal(t, "CLI", docs[0]["creator"])

	docKey := docs[0]["key"].(string)

	out = suite.Run("--json", "doc", "info", docKey)
	assert.Contains(t, out, "Hello")
	assert.Contains(t, out, "yolo-team")
}

// TestDocCreateEmpty
//   场景：创建不传 -c 的空文档，应正常创建且列表可见
//
//   步骤：
//     1) 创建项目 empty-doc-proj
//     2) 获取项目 key
//     3) 执行 doc create -p <key> -t empty（不传 -c）
//     4) 执行 doc list --json -p <key>，断言返回 1 条记录
func TestDocCreateEmpty(t *testing.T) {
	suite.Run("project", "create", "-n", "empty-doc-proj")

	out := suite.Run("--json", "project", "list")
	items := suite.JSONDecode(t, out)
	projKey := items[0]["key"].(string)

	suite.Run("doc", "create", "-p", projKey, "-t", "empty")

	out = suite.Run("--json", "doc", "list", "-p", projKey)
	docs := suite.JSONDecode(t, out)
	assert.Len(t, docs, 1)
}
