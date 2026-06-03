package cmd

import (
	"fmt"

	"yolo-team/yolo-cli/internal/handler"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{Use: "project", Short: "Manage projects"}

var (
	projCreateName string
	projCreateDesc string
)

var projectCreateCmd = &cobra.Command{
	Use:    "create",
	Short:  "Create a project",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if projCreateName == "" {
			return fmt.Errorf("name (-n) is required")
		}
		h := handler.NewProjectHandler()
		p, err := h.Create(ctx, handler.CreateProjectReq{Name: projCreateName, Description: projCreateDesc})
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": p})
			return nil
		}
		fmt.Printf("Project created: %s (%s)\n", p.Key, p.Name)
		return nil
	},
}

var projectListCmd = &cobra.Command{
	Use: "list", Short: "List all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		h := handler.NewProjectHandler()
		projects, err := h.List(ctx, struct{}{})
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": projects})
			return nil
		}

		fmt.Printf("%-4s %-20s %-32s %-20s\n", "ID", "KEY", "NAME", "UPDATED")
		for _, p := range projects {
			fmt.Printf("%-4d %-20s %-32s %-20s\n",
				p.ID, p.Key, p.Name, fmtTime(p.UpdatedAt))
		}
		return nil
	},
}

var projectInfoCmd = &cobra.Command{
	Use: "info <key>", Short: "Show project details", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h := handler.NewProjectHandler()
		p, err := h.Get(ctx, handler.GetProjectReq{Key: args[0]})
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": p})
			return nil
		}

		fmt.Printf("%-12s %s\n", "Key:", p.Key)
		fmt.Printf("%-12s %s\n", "Name:", p.Name)
		fmt.Printf("%-12s %s\n", "Description:", p.Description)
		fmt.Printf("%-12s %s\n", "Created:", fmtTime(p.CreatedAt))
		fmt.Printf("%-12s %s\n", "Updated:", fmtTime(p.UpdatedAt))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectCreateCmd, projectListCmd, projectInfoCmd)

	projectCreateCmd.Flags().StringVarP(&projCreateName, "name", "n", "", "project name (required)")
	projectCreateCmd.Flags().StringVarP(&projCreateDesc, "description", "d", "", "project description")
}
