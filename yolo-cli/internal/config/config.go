package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var DefaultDir string

func init() {
	home, _ := os.UserHomeDir()
	DefaultDir = filepath.Join(home, ".yolo-team")
}

type AIConfig struct {
	BaseURL      string `json:"base_url"`
	APIKey       string `json:"api_key"`
	Model        string `json:"model"`
	SystemPrompt string `json:"system_prompt"`
}

type Settings struct {
	Workspace string  `json:"workspace"`
	AI        AIConfig `json:"ai"`
}

func SettingsPath() string {
	return filepath.Join(DefaultDir, "settings.json")
}

func LoadSettings() (*Settings, error) {
	data, err := os.ReadFile(SettingsPath())
	if err != nil {
		return nil, err
	}

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func SaveSettings(s *Settings) error {
	os.MkdirAll(DefaultDir, 0755)
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(SettingsPath(), data, 0644)
}

func DBPath(workspace string) string {
	return filepath.Join(workspace, "yolo.db")
}

func DocumentsDir(workspace string) string {
	return filepath.Join(workspace, "documents")
}

func (s *Settings) ProjectDocumentsDir(projectID int64) string {
	return filepath.Join(s.Workspace, "documents", formatInt(projectID))
}

func (s *Settings) DocumentPath(projectID int64, title string) string {
	return filepath.Join(s.ProjectDocumentsDir(projectID), title+".md")
}

func formatInt(n int64) string {
	return fmt.Sprintf("%d", n)
}
