package cmd

import (
	"fmt"
	"os"

	"yolo-team/yolo-cli/internal/config"
	"yolo-team/yolo-cli/internal/db"
	"yolo-team/yolo-cli/internal/tui"

	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start AI chat TUI",
	Long: `启动 AI 智能助手对话界面（TUI 模式）

直接输入消息与 AI 助手交互，支持项目管理、任务分配等操作。
内置会话管理，自动保存对话历史。

快捷键:
  Ctrl+N  新建会话
  Tab     切换焦点（会话列表 / 输入框）
  Ctrl+L  清屏
  Ctrl+C  退出

斜杠命令:
  /sessions  查看会话列表
  /clear     清空当前会话
  /delete <ID>  删除指定会话`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := config.LoadSettings()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Not initialized. Please run: yolo init")
			os.Exit(1)
		}
		settings = s

		if err := db.Init(settings.Workspace); err != nil {
			return fmt.Errorf("init db: %w", err)
		}

		return tui.Run(settings.Workspace)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
