package cmd

import (
	"fmt"

	"yolo-team/yolo-cli/internal/db"
	"yolo-team/yolo-cli/internal/router"

	"github.com/spf13/cobra"
)

var servePort int
var serveDB string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace := "."
		if err := db.Init(workspace); err != nil {
			return fmt.Errorf("database init: %w", err)
		}

		r := router.Setup(workspace)
		addr := fmt.Sprintf(":%d", servePort)
		fmt.Printf("Server started at http://localhost%s\n", addr)
		return r.Run(addr)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "listen port")
	serveCmd.Flags().StringVar(&serveDB, "db", "yolo.db", "database path")
}
