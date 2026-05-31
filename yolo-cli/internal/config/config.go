package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Settings struct {
	Workspace string `json:"workspace"`
}

type DBConfig struct {
	Dialect string `json:"dialect"`
	DSN     string `json:"dsn"`
}

type Config struct {
	ServerPort int
	Settings   *Settings
	DB         *DBConfig
}

func settingsPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".yolo-team", "settings.json")
}

func LoadSettings() (*Settings, error) {
	data, err := os.ReadFile(settingsPath())
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
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath(), data, 0644)
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
