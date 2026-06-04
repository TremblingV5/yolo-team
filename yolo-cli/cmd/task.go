package cmd

import (
	"fmt"

	"yolo-team/yolo-cli/internal/handler"

	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{Use: "task", Short: "Manage tasks under an issue"}

var (
	taskListIssue string
)

var taskListCmd = &cobra.Command{
	Use: "list", Short: "List tasks for an issue",
	RunE: func(cmd *cobra.Command, args []string) error {
		if taskListIssue == "" {
			return fmt.Errorf("issue key (-i) is required")
		}
		h := handler.NewIssueHandler()
		tasks, err := h.ListTasks(ctx, handler.ListTaskReq{IssueKey: taskListIssue})
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": tasks})
			return nil
		}

		fmt.Printf("%-18s %-14s %-24s %s\n", "KEY", "STATUS", "TITLE", "DESCRIPTION")
		for _, t := range tasks {
			fmt.Printf("%-18s %-14s %-24s %s\n", t.Key, t.Status, trunc(t.Title, 24), trunc(t.Description, 30))
		}
		return nil
	},
}

var (
	taskCreateIssue string
	taskCreateTitle string
	taskCreateDesc  string
)

var taskCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a task under an issue",
	RunE: func(cmd *cobra.Command, args []string) error {
		if taskCreateIssue == "" || taskCreateTitle == "" {
			return fmt.Errorf("issue key (-i) and title (-t) are required")
		}
		h := handler.NewIssueHandler()
		task, err := h.CreateTask(ctx, handler.CreateTaskReq{
			IssueKey:    taskCreateIssue,
			Title:       taskCreateTitle,
			Description: taskCreateDesc,
		})
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": task})
			return nil
		}
		fmt.Printf("Task created: %s (%s)\n", task.Key, task.Title)
		return nil
	},
}

var (
	taskUpdateKey    string
	taskUpdateTitle  string
	taskUpdateDesc   string
	taskUpdateStatus string
)

var taskUpdateCmd = &cobra.Command{
	Use: "update <issue-key>", Short: "Update a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		if taskUpdateKey == "" {
			return fmt.Errorf("task key (-k) is required")
		}
		h := handler.NewIssueHandler()
		req := handler.UpdateTaskReq{TaskKey: taskUpdateKey}
		if cmd.Flags().Changed("title") {
			req.Title = &taskUpdateTitle
		}
		if cmd.Flags().Changed("description") {
			req.Description = &taskUpdateDesc
		}
		if cmd.Flags().Changed("status") {
			req.Status = &taskUpdateStatus
		}
		task, err := h.UpdateTask(ctx, req)
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": task})
			return nil
		}
		fmt.Printf("Task updated: %s (%s) [%s]\n", task.Key, task.Title, task.Status)
		return nil
	},
}

var taskDeleteKey string

var taskDeleteCmd = &cobra.Command{
	Use: "delete <issue-key>", Short: "Delete a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		if taskDeleteKey == "" {
			return fmt.Errorf("task key (-k) is required")
		}
		h := handler.NewIssueHandler()
		_, err := h.DeleteTask(ctx, handler.DeleteTaskReq{TaskKey: taskDeleteKey})
		if err != nil {
			return err
		}
		fmt.Printf("Task deleted: %s\n", taskDeleteKey)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskListCmd, taskCreateCmd, taskUpdateCmd, taskDeleteCmd)

	taskListCmd.Flags().StringVarP(&taskListIssue, "issue", "i", "", "issue key (required)")

	taskCreateCmd.Flags().StringVarP(&taskCreateIssue, "issue", "i", "", "issue key (required)")
	taskCreateCmd.Flags().StringVarP(&taskCreateTitle, "title", "t", "", "task title (required)")
	taskCreateCmd.Flags().StringVarP(&taskCreateDesc, "description", "d", "", "task description")

	taskUpdateCmd.Flags().StringVarP(&taskUpdateKey, "key", "k", "", "task key (required)")
	taskUpdateCmd.Flags().StringVarP(&taskUpdateTitle, "title", "t", "", "new title")
	taskUpdateCmd.Flags().StringVarP(&taskUpdateDesc, "description", "d", "", "new description")
	taskUpdateCmd.Flags().StringVarP(&taskUpdateStatus, "status", "s", "", "new status (created/in_progress/done)")

	taskDeleteCmd.Flags().StringVarP(&taskDeleteKey, "key", "k", "", "task key (required)")
}
