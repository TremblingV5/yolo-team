package ai

import (
	"os"

	"yolo-team/yolo-cli/internal/config"
)

// Config holds AI provider configuration.
type Config struct {
	BaseURL      string
	APIKey       string
	Model        string
	SystemPrompt string
	Workspace    string
}

// LoadConfig reads AI config from settings and overrides with environment variables.
func LoadConfig(settings *config.Settings) *Config {
	cfg := &Config{
		BaseURL:      "https://api.openai.com/v1",
		Model:        "gpt-4o",
		SystemPrompt: "You are Yolo-Team, an AI assistant for project management.",
	}

	if settings != nil {
		if settings.AI.BaseURL != "" {
			cfg.BaseURL = settings.AI.BaseURL
		}
		if settings.AI.APIKey != "" {
			cfg.APIKey = settings.AI.APIKey
		}
		if settings.AI.Model != "" {
			cfg.Model = settings.AI.Model
		}
		if settings.AI.SystemPrompt != "" {
			cfg.SystemPrompt = settings.AI.SystemPrompt
		}
	}

	// Environment variables override settings file.
	if v := os.Getenv("OPENAI_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		cfg.APIKey = v
	}
	if v := os.Getenv("OPENAI_MODEL"); v != "" {
		cfg.Model = v
	}

	return cfg
}
