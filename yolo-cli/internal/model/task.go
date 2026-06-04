package model

import (
	"fmt"
	"time"
)

type Task struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Key         string    `gorm:"column:key;type:varchar(64);uniqueIndex;not null" json:"key"`
	IssueID     int64     `gorm:"column:issue_id;not null;index" json:"issue_id"`
	Title       string    `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Description string    `gorm:"column:description;type:varchar(512);default:''" json:"description"`
	Status      string    `gorm:"column:status;type:varchar(16);not null;default:created" json:"status"`
	SortOrder   int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Task) TableName() string { return "tasks" }

const (
	TaskStatusCreated    = "created"
	TaskStatusInProgress = "in_progress"
	TaskStatusDone       = "done"
)

var AllTaskStatuses = []string{
	TaskStatusCreated, TaskStatusInProgress, TaskStatusDone,
}

func NewTask(issueID int64, title string) *Task {
	return &Task{
		IssueID: issueID,
		Title:   title,
		Status:  TaskStatusCreated,
	}
}

func (t *Task) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("title is required")
	}

	if len(t.Title) > 200 {
		return fmt.Errorf("title must not exceed 200 characters")
	}

	if len(t.Description) > 512 {
		return fmt.Errorf("description must not exceed 512 characters")
	}

	return nil
}
