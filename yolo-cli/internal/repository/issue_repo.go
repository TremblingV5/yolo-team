package repository

import (
	"context"

	"gorm.io/gorm"

	"yolo-team/yolo-cli/internal/db"
	"yolo-team/yolo-cli/internal/model"
	"yolo-team/yolo-cli/internal/query"
	"yolo-team/yolo-cli/internal/utils/keyutils"
)

type IssueRepo struct{}

func NewIssueRepo() *IssueRepo {
	return &IssueRepo{}
}

type IssueFilter struct {
	ProjectID  *int64
	Status     *string
	ExecutorID *int64
}

func (r *IssueRepo) List(filter IssueFilter) ([]model.Issue, error) {
	q := db.DB.Model(&model.Issue{}).Preload("Executor").Preload("Documents")
	if filter.ProjectID != nil {
		q = q.Where("project_id = ?", *filter.ProjectID)
	}

	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}

	if filter.ExecutorID != nil {
		q = q.Where("executor_id = ?", *filter.ExecutorID)
	}

	var issues []model.Issue
	err := q.Order("sort_order ASC, created_at DESC").Find(&issues).Error
	return issues, err
}

func (r *IssueRepo) Create(issue *model.Issue) error {
	if err := db.G[model.Issue]().Create(context.Background(), issue); err != nil {
		return err
	}

	issue.Key = keyutils.Issue(issue.ID)
	if _, err := db.G[model.Issue]().
		Where(query.Issue.ID.Eq(issue.ID)).
		Update(context.Background(), "key", issue.Key); err != nil {
		return err
	}

	return nil
}

func (r *IssueRepo) GetByKey(key string) (*model.Issue, error) {
	var i model.Issue
	err := db.DB.
		Preload("Executor").
		Preload("Documents").
		Where("key = ?", key).First(&i).Error
	if err != nil {
		return nil, err
	}

	// Load children separately to avoid circular preload issues
	var children []model.Issue
	db.DB.Where("parent_id = ?", i.ID).Find(&children)
	i.Children = children

	return &i, nil
}

func (r *IssueRepo) GetByID(id int64) (*model.Issue, error) {
	i, err := db.G[model.Issue]().
		Where(query.Issue.ID.Eq(id)).
		First(context.Background())
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (r *IssueRepo) Save(issue *model.Issue) error {
	// Use Select to include pointer fields that may be set to nil (like ParentID)
	_, err := db.G[model.Issue]().
		Where(query.Issue.ID.Eq(issue.ID)).
		Select("*").
		Updates(context.Background(), *issue)
	return err
}

func (r *IssueRepo) Delete(key string) error {
	issue, err := r.GetByKey(key)
	if err != nil {
		return err
	}

	_, err = db.G[model.Issue]().
		Where(query.Issue.ID.Eq(issue.ID)).
		Delete(context.Background())
	return err
}

func (r *IssueRepo) GetTodo(executorID int64, limit int) ([]model.Issue, error) {
	issues, err := db.G[model.Issue]().
		Where(query.Issue.ExecutorID.Eq(executorID)).
		Where(query.Issue.Status.NotIn(model.StatusDone, model.StatusArchived)).
		Order(query.Issue.Deadline.Asc()).
		Order(gorm.Expr("CASE priority WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 END ASC")).
		Limit(limit).
		Find(context.Background())
	return issues, err
}

func (r *IssueRepo) TopIssuesByProject(projectID int64) ([]model.Issue, error) {
	issues, err := db.G[model.Issue]().
		Where(query.Issue.ProjectID.Eq(projectID)).
		Order(query.Issue.SortOrder.Asc()).
		Order(query.Issue.CreatedAt.Desc()).
		Find(context.Background())
	return issues, err
}

func (r *IssueRepo) LinkDocument(issueID int64, documentID int64) error {
	return db.DB.Model(&model.Issue{ID: issueID}).
		Association("Documents").
		Append(&model.Document{ID: documentID})
}

func (r *IssueRepo) UnlinkDocument(issueID int64, documentID int64) error {
	return db.DB.Model(&model.Issue{ID: issueID}).
		Association("Documents").
		Delete(&model.Document{ID: documentID})
}

func (r *IssueRepo) ListByParent(parentID int64) ([]model.Issue, error) {
	var issues []model.Issue
	err := db.DB.Where("parent_id = ?", parentID).Find(&issues).Error
	return issues, err
}
