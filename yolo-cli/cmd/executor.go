package cmd

import (
	"fmt"

	"yolo-team/yolo-cli/internal/handler"

	"github.com/spf13/cobra"
)

var executorCmd = &cobra.Command{Use: "executor", Short: "Manage executors"}

var (
	execCreateName string
	execCreateSoul string
)

var executorCreateCmd = &cobra.Command{
	Use:    "create",
	Short:  "Create an executor",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if execCreateName == "" {
			return fmt.Errorf("name (-n) is required")
		}
		h := handler.NewExecutorHandler()
		e, err := h.Create(ctx, handler.CreateExecutorReq{Name: execCreateName, Soul: execCreateSoul})
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": e})
			return nil
		}
		fmt.Printf("Executor created: %s (ID: %d)\n", e.Name, e.ID)
		return nil
	},
}

var executorListCmd = &cobra.Command{
	Use: "list", Short: "List all executors",
	RunE: func(cmd *cobra.Command, args []string) error {
		h := handler.NewExecutorHandler()
		executors, err := h.List(ctx, struct{}{})
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": executors})
			return nil
		}

		fmt.Printf("%-4s %-20s %-12s %-20s\n", "ID", "NAME", "ROLE", "CREATED")
		for _, e := range executors {
			fmt.Printf("%-4d %-20s %-12s %-20s\n",
				e.ID, e.Name, e.Role, fmtTime(e.CreatedAt))
		}
		return nil
	},
}

var docCmd = &cobra.Command{Use: "doc", Short: "Manage documents"}

var (
	docListProject string
	docCreateProj  string
	docCreateTitle string
	docCreateContent string
)

var docListCmd = &cobra.Command{
	Use: "list", Short: "List documents in a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		h := handler.NewDocumentHandler(workspace())
		docs, err := h.List(ctx, handler.ListDocReq{ProjectKey: docListProject})
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": docs})
			return nil
		}

		fmt.Printf("%-18s %-20s %-10s %-20s\n", "KEY", "TITLE", "CREATOR", "UPDATED")
		for _, d := range docs {
			fmt.Printf("%-18s %-20s %-10s %-20s\n", d.Key, d.Title, d.Creator, fmtTime(d.UpdatedAt))
		}
		return nil
	},
}

var docCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a document",
	RunE: func(cmd *cobra.Command, args []string) error {
		h := handler.NewDocumentHandler(workspace())
		doc, err := h.Create(ctx, handler.CreateDocReq{
			ProjectKey: docCreateProj,
			Title:      docCreateTitle,
			Content:    docCreateContent,
			Creator:    "CLI",
		})
		if err != nil {
			return err
		}
		fmt.Printf("Document created: %s\n", doc.Key)
		fmt.Printf("Path: %s\n", doc.FilePath)
		return nil
	},
}

var docInfoCmd = &cobra.Command{
	Use: "info <key>", Short: "Show document details", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h := handler.NewDocumentHandler(workspace())
		doc, err := h.Get(ctx, handler.GetDocReq{Key: args[0]})
		if err != nil {
			return err
		}
		if useJSON {
			printJSON(map[string]interface{}{"data": doc})
			return nil
		}

		fmt.Printf("File: %s\n\n", doc.FilePath)
		fmt.Printf("%s\n", doc.Content)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(executorCmd)
	executorCmd.AddCommand(executorCreateCmd, executorListCmd)

	executorCreateCmd.Flags().StringVarP(&execCreateName, "name", "n", "", "executor name (required)")
	executorCreateCmd.Flags().StringVarP(&execCreateSoul, "soul", "s", "", "executor soul/prompt")

	rootCmd.AddCommand(docCmd)
	docCmd.AddCommand(docListCmd, docCreateCmd, docInfoCmd)

	docListCmd.Flags().StringVarP(&docListProject, "project", "p", "", "project key (required)")
	docCreateCmd.Flags().StringVarP(&docCreateProj, "project", "p", "", "project key (required)")
	docCreateCmd.Flags().StringVarP(&docCreateTitle, "title", "t", "", "document title (required)")
	docCreateCmd.Flags().StringVarP(&docCreateContent, "content", "c", "", "document content")
}
