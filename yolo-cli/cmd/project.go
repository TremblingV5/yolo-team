package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{Use: "project", Short: "Manage projects"}

var projectListCmd = &cobra.Command{
	Use: "list", Short: "List all projects",
	Run: func(cmd *cobra.Command, args []string) {
		var resp apiResponse
		must(apiGet("/projects", &resp))
		if useJSON {
			printJSON(resp)
			return
		}

		var projects []map[string]interface{}
		json.Unmarshal(resp.Data, &projects)
		fmt.Printf("%-4s %-20s %-32s %-20s\n", "ID", "KEY", "NAME", "UPDATED")
		for _, p := range projects {
			fmt.Printf("%-4.0f %-20s %-32s %-20s\n",
				p["id"], p["key"], p["name"],
				fmt.Sprintf("%v", p["updated_at"])[:20])
		}

	},
}

var projectInfoCmd = &cobra.Command{
	Use: "info <key>", Short: "Show project details", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var resp apiResponse
		must(apiGet("/projects/"+args[0], &resp))
		if useJSON {
			printJSON(resp)
			return
		}

		var p map[string]interface{}
		json.Unmarshal(resp.Data, &p)
		for _, k := range []string{"key", "name", "description", "created_at", "updated_at"} {
			fmt.Printf("%-12s %v\n", k+":", p[k])
		}

	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd, projectInfoCmd)
}
