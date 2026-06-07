package cmd

import (
	"encoding/json"
	"fmt"

	"yolo-team/yolo-cli/internal/common"
	"yolo-team/yolo-cli/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "查看和修改 yolo-team 配置",
}

var configAICmd = &cobra.Command{
	Use:   "ai",
	Short: "View or update AI configuration",
	Long: `查看或修改 AI 配置（base_url、api_key、model）。

配置存储在 ~/.yolo-team/settings.json 中。
环境变量（OPENAI_BASE_URL, OPENAI_API_KEY, OPENAI_MODEL）优先级更高。

示例:
  yolo config ai --show                  查看当前 AI 配置
  yolo config ai --api-key sk-xxx        设置 API Key
  yolo config ai --base-url https://...  设置 Base URL
  yolo config ai --model gpt-4o          设置模型
  yolo config ai --api-key sk-xxx --model gpt-4o  同时设置多个`,
	RunE: func(cmd *cobra.Command, args []string) error {
		show, _ := cmd.Flags().GetBool("show")

		s, err := config.LoadSettings()
		if err != nil {
			return fmt.Errorf("load settings: %w", err)
		}

		if show {
			return showAIConfig(s)
		}

		changed := false

		if v, _ := cmd.Flags().GetString("api-key"); v != "" {
			s.AI.APIKey = v
			changed = true
		}
		if v, _ := cmd.Flags().GetString("base-url"); v != "" {
			s.AI.BaseURL = v
			changed = true
		}
		if v, _ := cmd.Flags().GetString("model"); v != "" {
			s.AI.Model = v
			changed = true
		}

		if !changed {
			return showAIConfig(s)
		}

		if err := config.SaveSettings(s); err != nil {
			return fmt.Errorf("save settings: %w", err)
		}

		fmt.Println("AI configuration updated successfully.")
		return nil
	},
}

func showAIConfig(s *config.Settings) error {
	data, err := json.MarshalIndent(s.AI, "", "  ")
	if err != nil {
		return common.NewAppError(50001, err.Error())
	}
	fmt.Println(string(data))
	return nil
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configAICmd)

	configAICmd.Flags().Bool("show", false, "显示当前 AI 配置")
	configAICmd.Flags().String("api-key", "", "设置 API Key")
	configAICmd.Flags().String("base-url", "", "设置 API Base URL")
	configAICmd.Flags().String("model", "", "设置模型名称")
}
