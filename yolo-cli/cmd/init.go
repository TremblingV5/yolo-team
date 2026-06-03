package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"yolo-team/yolo-cli/internal/config"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var initWorkspace string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize yolo-team workspace",
	Long:  `Interactive setup for your yolo-team workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if initWorkspace != "" {
			absPath, _ := filepath.Abs(initWorkspace)
			os.MkdirAll(absPath, 0755)
			os.MkdirAll(filepath.Join(absPath, "documents"), 0755)
			if err := config.SaveSettings(&config.Settings{Workspace: absPath}); err != nil {
				return err
			}
			fmt.Printf("Initialized: %s\n", absPath)
			return nil
		}

		fmt.Println()
		fmt.Println("  ╔══════════════════════════════════╗")
		fmt.Println("  ║     Yolo-Team Setup Wizard       ║")
		fmt.Println("  ╚══════════════════════════════════╝")
		fmt.Println()

		defaultPath := filepath.Join(config.DefaultDir, "workspace")

		sel := promptui.Select{
			Label: "Choose workspace location",
			Items: []string{
				fmt.Sprintf("Default: %s", defaultPath),
				"Custom path (enter manually)",
			},
			Templates: &promptui.SelectTemplates{
				Active:   "▸ {{ . | cyan }}",
				Inactive: "  {{ . }}",
			},
		}

		idx, _, err := sel.Run()
		if err != nil {
			return err
		}

		wsPath := defaultPath
		if idx == 1 {
			validate := func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("path cannot be empty")
				}
				return nil
			}

			prompt := promptui.Prompt{
				Label:    "Workspace path",
				Default:  defaultPath,
				Validate: validate,
				Templates: &promptui.PromptTemplates{
					Prompt:  "{{ . | bold }} ",
					Valid:   "{{ . | green }} ",
					Invalid: "{{ . | red }} ",
				},
			}

			wsPath, err = prompt.Run()
			if err != nil {
				return err
			}
		}

		absPath, err := filepath.Abs(wsPath)
		if err != nil {
			return fmt.Errorf("invalid path: %w", err)
		}

		fmt.Println()
		fmt.Printf("  %s  %s\n", promptui.IconGood, absPath)
		fmt.Printf("  %s  Confirm? (Y/n): ", promptui.IconWarn)

		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "" && strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
			fmt.Println("  Canceled")
			return nil
		}

		if err := os.MkdirAll(absPath, 0755); err != nil {
			return fmt.Errorf("failed to create workspace: %w", err)
		}
		os.MkdirAll(filepath.Join(absPath, "documents"), 0755)

		s := &config.Settings{Workspace: absPath}
		if err := config.SaveSettings(s); err != nil {
			return fmt.Errorf("failed to save settings: %w", err)
		}

		fmt.Println()
		fmt.Printf("  %s Setup complete!\n", promptui.IconGood)
		fmt.Printf("    Config:   %s\n", config.SettingsPath())
		fmt.Printf("    Workspace: %s\n", absPath)
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&initWorkspace, "workspace", "w", "", "workspace path (skip interactive)")
}
