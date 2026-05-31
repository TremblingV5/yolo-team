// @title           Yolo Team API
// @version         1.0
// @description     Yolo Team 项目管理 API
// @host            localhost:8080
// @BasePath        /api/v1

package main

import (
	_ "yolo-team/yolo-cli/docs"
	"yolo-team/yolo-cli/cmd"
)

func main() {
	cmd.Execute()
}
