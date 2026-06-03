package cmd

import (
	"fmt"

	"yolo-team/yolo-cli/internal/router"

	"github.com/spf13/cobra"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server (optional)",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := router.Setup(workspace())
		addr := fmt.Sprintf(":%d", servePort)
		fmt.Printf("Server started at http://localhost%s\n", addr)
		return r.Run(addr)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "listen port")
}
