package cmd

import (
	"fmt"
	"time"

	"yolo-team/yolo-cli/internal/handler"
	"yolo-team/yolo-cli/internal/model"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		h := handler.NewIssueHandler()
		req := handler.ListIssueReq{}
		if issueListProject > 0 {
			v := int64(issueListProject)
			req.ProjectID = &v
		}
		if issueListStatus != "" {
			req.Status = &issueListStatus
		}
		if issueListExecutor > 0 {
			v := int64(issueListExecutor)
			req.ExecutorID = &v
		}
		if cmd.Flags().Changed("parent") {
			v := int64(issueListParent)
			req.ParentID = &v
		}

		issues, err := h.List(ctx, req)
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": issues})
			return nil
		}

		fmt.Printf("%-18s %-14s %-8s %-24s %-12s %s\n", "KEY", "STATUS", "PRIORITY", "TITLE", "DEADLINE", "REPO")
		for _, iss := range issues {
			deadline := "-"
			if iss.Deadline != nil {
				deadline = iss.Deadline.Format("2006-01-02")
			}
			repo := "-"
			if iss.RepoName != "" {
				repo = iss.RepoName
				if iss.BranchName != "" {
					repo += fmt.Sprintf(" (%s)", iss.BranchName)
				}
			}
			fmt.Printf("%-18s %-14s %-8s %-24s %-12s %s\n",
				iss.Key, iss.Status, iss.Priority, trunc(iss.Title, 24), deadline, trunc(repo, 30))
		}
		return nil
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
		h := handler.NewIssueHandler()
		req := handler.CreateIssueReq{
			ProjectID:   int64(issueCreateProject),
			Title:       issueCreateTitle,
			Description: issueCreateDesc,
			Priority:    issueCreatePri,
			Deadline:    issueCreateDeadline,
		}
		if issueCreateParent > 0 {
			v := int64(issueCreateParent)
			req.ParentID = &v
		}
		if issueCreateExec > 0 {
			v := int64(issueCreateExec)
			req.ExecutorID = &v
		}

		iss, err := h.Create(ctx, req)
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": iss})
			return nil
		}
		fmt.Printf("Issue created: %s (%s)\n", iss.Key, iss.Title)
		return nil
	},
}

var issueInfoCmd = &cobra.Command{
	Use: "info <key>", Short: "Show issue details", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h := handler.NewIssueHandler()
		iss, err := h.Get(ctx, handler.GetIssueReq{Key: args[0]})
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": iss})
			return nil
		}

		fmt.Printf("%-12s %s\n", "Key:", iss.Key)
		fmt.Printf("%-12s %s\n", "Status:", iss.Status)
		fmt.Printf("%-12s %s\n", "Title:", iss.Title)
		fmt.Printf("%-12s %s\n", "Priority:", iss.Priority)
		fmt.Printf("%-12s %s\n", "Description:", iss.Description)
		if iss.Deadline != nil {
			fmt.Printf("%-12s %s\n", "Deadline:", iss.Deadline.Format("2006-01-02"))
		}
		if iss.ExecutorID != nil && iss.Executor != nil {
			fmt.Printf("%-12s %s (%s)\n", "Executor:", iss.Executor.Name, iss.Executor.Role)
		}
		repo := "-"
		if iss.RepoName != "" {
			repo = iss.RepoName
			if iss.BranchName != "" {
				repo += fmt.Sprintf(" (%s)", iss.BranchName)
			}
		}
		fmt.Printf("%-12s %s\n", "Repo:", repo)
		return nil
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
	RunE: func(cmd *cobra.Command, args []string) error {
		h := handler.NewIssueHandler()
		body := model.UpdateIssueRequest{}
		if cmd.Flags().Changed("status") {
			body.Status = &issueUpdateStatus
		}
		if cmd.Flags().Changed("title") {
			body.Title = &issueUpdateTitle
		}
		if cmd.Flags().Changed("executor") {
			body.ExecutorID = &issueUpdateExec
		}
		if cmd.Flags().Changed("deadline") {
			t, err := time.Parse("2006-01-02", issueUpdateDeadline)
			if err == nil {
				body.Deadline = &t
			}
		}
		if cmd.Flags().Changed("repo-url") {
			body.RepoURL = &issueUpdateRepoURL
		}
		if cmd.Flags().Changed("repo-name") {
			body.RepoName = &issueUpdateRepoName
		}
		if cmd.Flags().Changed("branch") {
			body.BranchName = &issueUpdateBranch
		}

		_, err := h.Update(ctx, handler.UpdateIssueReq{Key: args[0], UpdateIssueRequest: body})
		if err != nil {
			return err
		}
		fmt.Printf("Issue updated: %s\n", args[0])
		return nil
	},
}

var issueDeleteCmd = &cobra.Command{
	Use: "delete <key>", Short: "Delete an issue", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h := handler.NewIssueHandler()
		_, err := h.Delete(ctx, handler.DeleteIssueReq{Key: args[0]})
		if err != nil {
			return err
		}
		fmt.Printf("Issue deleted: %s\n", args[0])
		return nil
	},
}

var (
	todoExecutor int
	todoLimit    int
)

var issueTodoCmd = &cobra.Command{
	Use: "todo", Short: "Show todo list for an executor",
	RunE: func(cmd *cobra.Command, args []string) error {
		if todoExecutor == 0 {
			return fmt.Errorf("executor (-e) is required")
		}
		limit := todoLimit
		if limit <= 0 {
			limit = 20
		}
		h := handler.NewIssueHandler()
		issues, err := h.Todo(ctx, handler.TodoIssueReq{ExecutorID: int64(todoExecutor), Limit: limit})
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": issues})
			return nil
		}

		fmt.Printf("%-18s %-14s %-8s %-24s %-12s\n", "KEY", "STATUS", "PRIORITY", "TITLE", "DEADLINE")
		for _, iss := range issues {
			deadline := "-"
			if iss.Deadline != nil {
				deadline = iss.Deadline.Format("2006-01-02")
			}
			fmt.Printf("%-18s %-14s %-8s %-24s %-12s\n",
				iss.Key, iss.Status, iss.Priority, trunc(iss.Title, 24), deadline)
		}
		return nil
	},
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
