package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var executorCmd = &cobra.Command{Use: "executor", Short: "Manage executors"}

var executorListCmd = &cobra.Command{
	Use: "list", Short: "List all executors",
	Run: func(cmd *cobra.Command, args []string) {
		var resp apiResponse
		must(apiGet("/executors", &resp))
		if useJSON {
			printJSON(resp)
			return
		}

		var executors []map[string]interface{}
		json.Unmarshal(resp.Data, &executors)
		fmt.Printf("%-4s %-20s %-12s %-20s\n", "ID", "NAME", "ROLE", "CREATED")
		for _, e := range executors {
			fmt.Printf("%-4.0f %-20s %-12s %-20s\n",
				e["id"], e["name"], e["role"],
				fmt.Sprintf("%v", e["created_at"])[:20])
		}

	},
}

var docCmd = &cobra.Command{Use: "doc", Short: "Manage documents"}

var (
	docListProject string
	docCreateProj  string
	docCreateTitle string
)

var docListCmd = &cobra.Command{
	Use: "list", Short: "List documents in a project",
	Run: func(cmd *cobra.Command, args []string) {
		var resp apiResponse
		must(apiGet("/projects/"+docListProject+"/documents", &resp))
		if useJSON {
			printJSON(resp)
			return
		}

		var docs []map[string]interface{}
		json.Unmarshal(resp.Data, &docs)
		fmt.Printf("%-18s %-20s %-20s\n", "KEY", "TITLE", "UPDATED")
		for _, d := range docs {
			fmt.Printf("%-18s %-20s %-20s\n", d["key"], d["title"],
				fmt.Sprintf("%v", d["updated_at"])[:20])
		}

	},
}

var docCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a document",
	RunE: func(cmd *cobra.Command, args []string) error {
		body := map[string]interface{}{
			"title":   docCreateTitle,
			"content": "",
		}

		var resp apiResponse
		must(apiPost("/projects/"+docCreateProj+"/documents", body, &resp))
		var doc map[string]interface{}
		json.Unmarshal(resp.Data, &doc)
		fmt.Printf("Document created: %v\n", doc["key"])
		fmt.Printf("Path: %v\n", doc["file_path"])
		return nil
	},
}

var docInfoCmd = &cobra.Command{
	Use: "info <key>", Short: "Show document details", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var resp apiResponse
		must(apiGet("/documents/"+args[0], &resp))
		if useJSON {
			printJSON(resp)
			return
		}

		var doc map[string]interface{}
		json.Unmarshal(resp.Data, &doc)
		fmt.Printf("File: %v\n\n", doc["file_path"])
		fmt.Printf("%v\n", doc["content"])
	},
}

func init() {
	rootCmd.AddCommand(executorCmd)
	executorCmd.AddCommand(executorListCmd)

	rootCmd.AddCommand(docCmd)
	docCmd.AddCommand(docListCmd, docCreateCmd, docInfoCmd)

	docListCmd.Flags().StringVarP(&docListProject, "project", "p", "", "project key (required)")
	docCreateCmd.Flags().StringVarP(&docCreateProj, "project", "p", "", "project key (required)")
	docCreateCmd.Flags().StringVarP(&docCreateTitle, "title", "t", "", "document title (required)")
}
