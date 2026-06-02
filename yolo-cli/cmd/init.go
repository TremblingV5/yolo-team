package cmd

import (
	"fmt"
	"os"

	"yolo-team/yolo-cli/internal/config"

	"github.com/spf13/cobra"
)

var initWorkspace string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize yolo-team workspace",
	Long:  `Set up the workspace directory and save settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := os.Stat(initWorkspace)
		if err != nil {
			return fmt.Errorf("workspace path error: %w", err)
		}

		if !info.IsDir() {
			return fmt.Errorf("workspace path is not a directory: %s", initWorkspace)
		}

		settings := &config.Settings{Workspace: initWorkspace}
		if err := config.SaveSettings(settings); err != nil {
			return err
		}

		os.MkdirAll(config.DocumentsDir(initWorkspace), 0755)

		fmt.Printf("Workspace initialized: %s\n", initWorkspace)
		fmt.Printf("Settings saved to ~/.yolo-team/settings.json\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&initWorkspace, "workspace", "w", "", "workspace directory (required)")
}
