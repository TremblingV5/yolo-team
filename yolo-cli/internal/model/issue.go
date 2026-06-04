package model

import (
	"fmt"
	"time"
)

type Issue struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Key         string     `gorm:"column:key;type:varchar(64);uniqueIndex;not null" json:"key"`
	ProjectID   int64      `gorm:"column:project_id;not null;index" json:"project_id"`
	Status      string     `gorm:"column:status;type:varchar(16);not null;default:created" json:"status"`
	Title       string     `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Description string     `gorm:"column:description;type:varchar(512);default:''" json:"description"`
	Deadline    *time.Time `gorm:"column:deadline;default:null" json:"deadline"`
	Priority    string     `gorm:"column:priority;type:varchar(16);default:medium" json:"priority"`
	ExecutorID  *int64     `gorm:"column:executor_id;default:null" json:"executor_id"`
	RepoURL     string     `gorm:"column:repo_url;type:varchar(500);default:''" json:"repo_url"`
	RepoName    string     `gorm:"column:repo_name;type:varchar(255);default:''" json:"repo_name"`
	BranchName  string     `gorm:"column:branch_name;type:varchar(255);default:''" json:"branch_name"`
	SortOrder   int        `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Project   *Project   `gorm:"foreignKey:ProjectID" json:"-"`
	Executor  *Executor  `gorm:"foreignKey:ExecutorID" json:"executor,omitempty"`
	Documents []Document `gorm:"many2many:issue_documents;" json:"documents,omitempty"`
}

func (Issue) TableName() string { return "issues" }

const (
	StatusCreated    = "created"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
	StatusArchived   = "archived"
)

const (
	PriorityLow      = "low"
	PriorityMedium   = "medium"
	PriorityHigh     = "high"
	PriorityCritical = "critical"
)

var AllStatuses = []string{
	StatusCreated, StatusInProgress, StatusDone, StatusArchived,
}

func NewIssue(projectID int64, title string) *Issue {
	return &Issue{
		ProjectID: projectID,
		Title:     title,
		Status:    StatusCreated,
		Priority:  PriorityMedium,
	}
}

func (i *Issue) Validate() error {
	if i.Title == "" {
		return fmt.Errorf("title is required")
	}

	if len(i.Title) > 200 {
		return fmt.Errorf("title must not exceed 200 characters")
	}

	if len(i.Description) > 512 {
		return fmt.Errorf("description must not exceed 512 characters")
	}

	return nil
}

func (i *Issue) ApplyUpdate(req *UpdateIssueRequest) {
	if req.Title != nil {
		i.Title = *req.Title
	}

	if req.Description != nil {
		i.Description = *req.Description
	}

	if req.Status != nil {
		i.Status = *req.Status
	}

	if req.Priority != nil {
		i.Priority = *req.Priority
	}

	if req.ExecutorID != nil {
		i.ExecutorID = req.ExecutorID
	}

	if req.Deadline != nil {
		i.Deadline = req.Deadline
	}

	if req.RepoURL != nil {
		i.RepoURL = *req.RepoURL
	}

	if req.RepoName != nil {
		i.RepoName = *req.RepoName
	}

	if req.BranchName != nil {
		i.BranchName = *req.BranchName
	}
}

type UpdateIssueRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	Priority    *string    `json:"priority"`
	ExecutorID  *int64     `json:"executor_id"`
	Deadline    *time.Time `json:"deadline"`
	RepoURL     *string    `json:"repo_url"`
	RepoName    *string    `json:"repo_name"`
	BranchName  *string    `json:"branch_name"`
}
