package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	serverURL string
	useJSON   bool
)

var rootCmd = &cobra.Command{
	Use:   "yolo",
	Short: "Yolo-Team CLI 鈥?project & task management for human + AI teams",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "http://localhost:8080", "server URL")
	rootCmd.PersistentFlags().BoolVar(&useJSON, "json", false, "output as JSON")
}
