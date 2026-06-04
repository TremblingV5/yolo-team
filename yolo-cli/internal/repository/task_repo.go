package repository

import (
	"context"

	"yolo-team/yolo-cli/internal/db"
	"yolo-team/yolo-cli/internal/model"
	"yolo-team/yolo-cli/internal/query"
	"yolo-team/yolo-cli/internal/utils/keyutils"
)

type TaskRepo struct{}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{}
}

func (r *TaskRepo) ListByIssue(issueID int64) ([]model.Task, error) {
	issues, err := db.G[model.Task]().
		Where(query.Task.IssueID.Eq(issueID)).
		Order(query.Task.SortOrder.Asc()).
		Order(query.Task.CreatedAt.Asc()).
		Find(context.Background())
	return issues, err
}

func (r *TaskRepo) Create(task *model.Task) error {
	if err := db.G[model.Task]().Create(context.Background(), task); err != nil {
		return err
	}

	task.Key = keyutils.Task(task.ID)
	if _, err := db.G[model.Task]().
		Where(query.Task.ID.Eq(task.ID)).
		Update(context.Background(), "key", task.Key); err != nil {
		return err
	}

	return nil
}

func (r *TaskRepo) GetByKey(key string) (*model.Task, error) {
	t, err := db.G[model.Task]().
		Where(query.Task.Key.Eq(key)).
		First(context.Background())
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepo) GetByID(id int64) (*model.Task, error) {
	t, err := db.G[model.Task]().
		Where(query.Task.ID.Eq(id)).
		First(context.Background())
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepo) Update(task *model.Task) error {
	_, err := db.G[model.Task]().
		Where(query.Task.ID.Eq(task.ID)).
		Updates(context.Background(), *task)
	return err
}

func (r *TaskRepo) DeleteByKey(key string) error {
	_, err := db.G[model.Task]().
		Where(query.Task.Key.Eq(key)).
		Delete(context.Background())
	return err
}
