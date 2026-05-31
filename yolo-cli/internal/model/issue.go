package model

import (
	"fmt"
	"time"
)

type Issue struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Key         string     `gorm:"column:key;type:varchar(64);uniqueIndex;not null" json:"key"`
	ProjectID   int64      `gorm:"column:project_id;not null;index" json:"project_id"`
	ParentID    *int64     `gorm:"column:parent_id;default:null" json:"parent_id"`
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

	Children []Issue   `gorm:"-" json:"children,omitempty"`
	Project  *Project  `gorm:"foreignKey:ProjectID" json:"-"`
	Parent   *Issue    `gorm:"foreignKey:ParentID" json:"-"`
	Executor *Executor `gorm:"foreignKey:ExecutorID" json:"-"`
}

func (Issue) TableName() string { return "issues" }

const (
	StatusCreated        = "created"
	StatusDesign         = "design"
	StatusReview         = "review"
	StatusImplementation = "implementation"
	StatusQA             = "qa"
	StatusPendingReview  = "pending_review"
	StatusDone           = "done"
	StatusArchived       = "archived"
)

const (
	PriorityLow      = "low"
	PriorityMedium   = "medium"
	PriorityHigh     = "high"
	PriorityCritical = "critical"
)

var AllStatuses = []string{
	StatusCreated, StatusDesign, StatusReview, StatusImplementation,
	StatusQA, StatusPendingReview, StatusDone, StatusArchived,
}

var StatusesCLIForbidden = []string{StatusReview, StatusPendingReview, StatusArchived}

type transition struct {
	from  string
	to    string
	roles []string
	cliOK bool
}

var transitions = []transition{
	{StatusCreated, StatusDesign, []string{RoleLeader}, true},
	{StatusDesign, StatusReview, []string{RoleLeader, RoleArchitect}, true},
	{StatusReview, StatusImplementation, []string{RoleLeader}, false},
	{StatusReview, StatusDesign, []string{RoleLeader, RoleQA}, false},
	{StatusImplementation, StatusQA, []string{RoleDeveloper}, true},
	{StatusQA, StatusPendingReview, []string{RoleQA}, true},
	{StatusQA, StatusImplementation, []string{RoleQA}, true},
	{StatusPendingReview, StatusDone, []string{RoleLeader}, false},
	{StatusPendingReview, StatusDesign, []string{RoleLeader}, false},
	{StatusPendingReview, StatusImplementation, []string{RoleLeader}, false},
	{StatusPendingReview, StatusQA, []string{RoleLeader}, false},
	{StatusDone, StatusArchived, nil, false},
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

func (i *Issue) CanTransitionTo(targetStatus, executorRole string, fromCLI bool) error {
	if i.Status == targetStatus {
		return nil
	}

	for _, t := range transitions {
		if t.from == i.Status && t.to == targetStatus {
			if fromCLI && !t.cliOK {
				return fmt.Errorf("status transition '%s -> %s' is not allowed from CLI", i.Status, targetStatus)
			}

			if t.roles == nil {
				return nil
			}

			if executorRole == "" {
				return fmt.Errorf("no executor assigned")
			}

			for _, r := range t.roles {
				if r == executorRole {
					return nil
				}
			}

			return fmt.Errorf("role '%s' is not allowed to move from '%s' to '%s'", executorRole, i.Status, targetStatus)
		}
	}

	return fmt.Errorf("transition from '%s' to '%s' is not allowed", i.Status, targetStatus)
}

func (i *Issue) CanBeParent() bool {
	return i.ParentID == nil
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
