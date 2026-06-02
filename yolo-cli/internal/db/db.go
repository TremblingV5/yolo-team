package db

import (
	"fmt"

	"yolo-team/yolo-cli/internal/config"
	"yolo-team/yolo-cli/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func G[T any]() gorm.Interface[T] {
	return gorm.G[T](DB)
}

func Init(workspace string) error {
	dbPath := config.DBPath(workspace)
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	if err := DB.AutoMigrate(
		&model.Project{},
		&model.Executor{},
		&model.Issue{},
		&model.Document{},
	); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	// Convert old status values to new 4-state system
	oldToNew := map[string]string{
		"design":         model.StatusInProgress,
		"review":         model.StatusInProgress,
		"implementation": model.StatusInProgress,
		"qa":             model.StatusInProgress,
		"pending_review": model.StatusInProgress,
	}
	for old, new := range oldToNew {
		DB.Model(&model.Issue{}).Where("status = ?", old).Update("status", new)
	}

	return nil
}
