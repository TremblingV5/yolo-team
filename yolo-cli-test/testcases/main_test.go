package testcases

import (
	"os"
	"testing"

	"yolo-team/yolo-cli-test/suite"
)

func TestMain(m *testing.M) {
	suite.Setup()
	code := m.Run()
	suite.Teardown()
	os.Exit(code)
}
