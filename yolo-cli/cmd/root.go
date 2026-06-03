package cmd

import (
	"context"
	"fmt"
	"os"

	"yolo-team/yolo-cli/internal/config"
	"yolo-team/yolo-cli/internal/db"

	"github.com/spf13/cobra"
)

var (
	settings *config.Settings
	useJSON  bool
)

var rootCmd = &cobra.Command{
	Use:   "yolo",
	Short: "Yolo-Team CLI — project & task management for human + AI teams",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "init" || cmd.Name() == "help" {
			return nil
		}
		help, _ := cmd.Flags().GetBool("help")
		if help {
			return nil
		}

		s, err := config.LoadSettings()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Not initialized. Please run: yolo init")
			os.Exit(1)
		}
		settings = s

		return db.Init(settings.Workspace)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func workspace() string {
	if settings != nil {
		return settings.Workspace
	}
	return "."
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&useJSON, "json", false, "output as JSON")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

var ctx = context.Background()
