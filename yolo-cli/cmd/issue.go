package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{Use: "issue", Short: "Manage issues"}

var (
	issueListProject  int
	issueListStatus   string
	issueListExecutor int
	issueListParent   int
)

var issueListCmd = &cobra.Command{
	Use: "list", Short: "List issues",
	Run: func(cmd *cobra.Command, args []string) {
		path := "/issues?"
		if issueListProject > 0 {
			path += fmt.Sprintf("project_id=%d&", issueListProject)
		}

		if issueListStatus != "" {
			path += fmt.Sprintf("status=%s&", issueListStatus)
		}

		if issueListExecutor > 0 {
			path += fmt.Sprintf("executor_id=%d&", issueListExecutor)
		}

		if cmd.Flags().Changed("parent") {
			path += fmt.Sprintf("parent_id=%d&", issueListParent)
		}

		var resp apiResponse
		must(apiGet(strings.TrimRight(path, "&"), &resp))
		if useJSON {
			printJSON(resp)
			return
		}

		var issues []map[string]interface{}
		json.Unmarshal(resp.Data, &issues)
		fmt.Printf("%-18s %-16s %-8s %-20s %-12s %-30s\n", "KEY", "STATUS", "PRIORITY", "TITLE", "DEADLINE", "REPO")
		for _, iss := range issues {
			deadline := "-"
			if d, ok := iss["deadline"]; ok && d != nil {
				deadline = fmt.Sprintf("%v", d)
				if len(deadline) > 10 {
					deadline = deadline[:10]
				}
			}

			repo := "-"
			if rn, ok := iss["repo_name"]; ok && rn != nil && rn != "" {
				repo = fmt.Sprintf("%v", rn)
				if br, ok := iss["branch_name"]; ok && br != nil && br != "" {
					repo += fmt.Sprintf(" (%v)", br)
				}
			}

			fmt.Printf("%-18s %-16s %-8s %-20s %-12s %-30s\n",
				iss["key"], iss["status"], iss["priority"], truncate(fmt.Sprintf("%v", iss["title"]), 20), deadline, truncate(repo, 30))
		}

	},
}

var (
	issueCreateTitle    string
	issueCreateProject  int
	issueCreateDesc     string
	issueCreatePri      string
	issueCreateParent   int
	issueCreateExec     int
	issueCreateDeadline string
)

var issueCreateCmd = &cobra.Command{
	Use: "create", Short: "Create an issue",
	RunE: func(cmd *cobra.Command, args []string) error {
		if issueCreateTitle == "" || issueCreateProject == 0 {
			return fmt.Errorf("title (-t) and project (-p) are required")
		}

		body := map[string]interface{}{
			"project_id": issueCreateProject,
			"title":      issueCreateTitle,
		}

		if issueCreateDesc != "" {
			body["description"] = issueCreateDesc
		}

		if issueCreatePri != "" {
			body["priority"] = issueCreatePri
		}

		if issueCreateParent > 0 {
			body["parent_id"] = issueCreateParent
		}

		if issueCreateExec > 0 {
			body["executor_id"] = issueCreateExec
		}

		if issueCreateDeadline != "" {
			body["deadline"] = issueCreateDeadline
		}

		var resp apiResponse
		must(apiPost("/issues", body, &resp))
		var iss map[string]interface{}
		json.Unmarshal(resp.Data, &iss)
		fmt.Printf("Issue created: %v (%v)\n", iss["key"], iss["title"])
		return nil
	},
}

var issueInfoCmd = &cobra.Command{
	Use: "info <key>", Short: "Show issue details", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var resp apiResponse
		must(apiGet("/issues/"+args[0], &resp))
		if useJSON {
			printJSON(resp)
			return
		}

		var iss map[string]interface{}
		json.Unmarshal(resp.Data, &iss)
		for _, k := range []string{"key", "status", "title", "description", "priority", "deadline"} {
			fmt.Printf("%-12s %v\n", k+":", iss[k])
		}

		repo := "-"
		if rn, ok := iss["repo_name"]; ok && rn != nil && rn != "" {
			repo = fmt.Sprintf("%v", rn)
			if br, ok := iss["branch_name"]; ok && br != nil && br != "" {
				repo += fmt.Sprintf(" (%v)", br)
			}
		}

		fmt.Printf("%-12s %s\n", "Repo:", repo)
	},
}

var (
	issueUpdateStatus   string
	issueUpdateTitle    string
	issueUpdateExec     int64
	issueUpdateDeadline string
	issueUpdateRepoURL  string
	issueUpdateRepoName string
	issueUpdateBranch   string
)

var issueUpdateCmd = &cobra.Command{
	Use: "update <key>", Short: "Update an issue", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		body := map[string]interface{}{}
		if cmd.Flags().Changed("status") {
			body["status"] = issueUpdateStatus
		}

		if cmd.Flags().Changed("title") {
			body["title"] = issueUpdateTitle
		}

		if cmd.Flags().Changed("executor") {
			body["executor_id"] = issueUpdateExec
		}

		if cmd.Flags().Changed("deadline") {
			body["deadline"] = issueUpdateDeadline
		}

		if cmd.Flags().Changed("repo-url") {
			body["repo_url"] = issueUpdateRepoURL
		}

		if cmd.Flags().Changed("repo-name") {
			body["repo_name"] = issueUpdateRepoName
		}

		if cmd.Flags().Changed("branch") {
			body["branch_name"] = issueUpdateBranch
		}

		var resp apiResponse
		must(apiPut("/issues/"+args[0], body, &resp))
		fmt.Printf("Issue updated: %s\n", args[0])
	},
}

var issueDeleteCmd = &cobra.Command{
	Use: "delete <key>", Short: "Delete an issue", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		must(apiDelete("/issues/" + args[0]))
		fmt.Printf("Issue deleted: %s\n", args[0])
	},
}

var (
	todoExecutor int
	todoLimit    int
)

var issueTodoCmd = &cobra.Command{
	Use: "todo", Short: "Show todo list for an executor",
	Run: func(cmd *cobra.Command, args []string) {
		if todoExecutor == 0 {
			fmt.Println("Error: executor (-e) is required")
			return
		}

		limit := todoLimit
		if limit <= 0 {
			limit = 20
		}

		var resp apiResponse
		must(apiGet(fmt.Sprintf("/issues/todo?executor_id=%d&limit=%d", todoExecutor, limit), &resp))
		if useJSON {
			printJSON(resp)
			return
		}

		var issues []map[string]interface{}
		json.Unmarshal(resp.Data, &issues)
		fmt.Printf("%-18s %-16s %-8s %-20s %-12s\n", "KEY", "STATUS", "PRIORITY", "TITLE", "DEADLINE")
		for _, iss := range issues {
			deadline := "-"
			if d, ok := iss["deadline"]; ok && d != nil {
				deadline = fmt.Sprintf("%v", d)
				if len(deadline) > 10 {
					deadline = deadline[:10]
				}
			}

			fmt.Printf("%-18s %-16s %-8s %-20s %-12s\n",
				iss["key"], iss["status"], iss["priority"], truncate(fmt.Sprintf("%v", iss["title"]), 20), deadline)
		}

	},
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n-2]) + ".."
	}

	return s
}

func init() {
	rootCmd.AddCommand(issueCmd)
	issueCmd.AddCommand(issueListCmd, issueCreateCmd, issueInfoCmd, issueUpdateCmd, issueDeleteCmd, issueTodoCmd)

	issueListCmd.Flags().IntVarP(&issueListProject, "project", "p", 0, "filter by project ID")
	issueListCmd.Flags().StringVarP(&issueListStatus, "status", "s", "", "filter by status")
	issueListCmd.Flags().IntVarP(&issueListExecutor, "executor", "e", 0, "filter by executor ID")
	issueListCmd.Flags().IntVar(&issueListParent, "parent", -1, "filter by parent (0 for top-level)")

	issueCreateCmd.Flags().StringVarP(&issueCreateTitle, "title", "t", "", "issue title")
	issueCreateCmd.Flags().IntVarP(&issueCreateProject, "project", "p", 0, "project ID")
	issueCreateCmd.Flags().StringVarP(&issueCreateDesc, "description", "d", "", "description")
	issueCreateCmd.Flags().StringVar(&issueCreatePri, "priority", "", "priority")
	issueCreateCmd.Flags().IntVar(&issueCreateParent, "parent", 0, "parent issue ID")
	issueCreateCmd.Flags().IntVar(&issueCreateExec, "executor", 0, "executor ID")
	issueCreateCmd.Flags().StringVar(&issueCreateDeadline, "deadline", "", "deadline (YYYY-MM-DD)")

	issueUpdateCmd.Flags().StringVarP(&issueUpdateStatus, "status", "s", "", "new status")
	issueUpdateCmd.Flags().StringVarP(&issueUpdateTitle, "title", "t", "", "new title")
	issueUpdateCmd.Flags().Int64Var(&issueUpdateExec, "executor", 0, "executor ID")
	issueUpdateCmd.Flags().StringVar(&issueUpdateDeadline, "deadline", "", "deadline")
	issueUpdateCmd.Flags().StringVar(&issueUpdateRepoURL, "repo-url", "", "repo URL")
	issueUpdateCmd.Flags().StringVar(&issueUpdateRepoName, "repo-name", "", "repo name")
	issueUpdateCmd.Flags().StringVar(&issueUpdateBranch, "branch", "", "branch name")

	issueTodoCmd.Flags().IntVarP(&todoExecutor, "executor", "e", 0, "executor ID (required)")
	issueTodoCmd.Flags().IntVarP(&todoLimit, "limit", "n", 20, "max results")
}
