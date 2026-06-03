package testcases

import (
	"testing"

	"yolo-team/yolo-cli-test/suite"

	"github.com/stretchr/testify/assert"
)

// TestInitOverwrite
//   场景：重复执行 yolo init 应始终成功，不会因已有配置而报错
//
//   步骤：
//     1) 执行 init -w .，断言输出包含 "Initialized"
//     2) 再次执行 init -w .（覆盖已有配置）
//     3) 断言输出仍然包含 "Initialized"
func TestInitOverwrite(t *testing.T) {
	out := suite.Run("init", "-w", ".")
	assert.Contains(t, out, "Initialized")

	out = suite.Run("init", "-w", ".")
	assert.Contains(t, out, "Initialized")
}
