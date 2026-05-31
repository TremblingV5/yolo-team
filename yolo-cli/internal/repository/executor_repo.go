package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"yolo-team/yolo-cli/internal/db"
	"yolo-team/yolo-cli/internal/model"
	"yolo-team/yolo-cli/internal/query"
)

type ExecutorRepo struct{}

func NewExecutorRepo() *ExecutorRepo {
	return &ExecutorRepo{}
}

func (r *ExecutorRepo) List() ([]model.Executor, error) {
	executors, err := db.G[model.Executor]().
		Order(query.Executor.CreatedAt.Desc()).
		Find(context.Background())
	return executors, err
}

func (r *ExecutorRepo) Create(e *model.Executor) (*model.Executor, error) {
	if err := db.G[model.Executor]().Create(context.Background(), e); err != nil {
		return nil, err
	}

	return e, nil
}

func (r *ExecutorRepo) GetByID(id int64) (*model.Executor, error) {
	e, err := db.G[model.Executor]().
		Where(query.Executor.ID.Eq(id)).
		First(context.Background())
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *ExecutorRepo) GetByName(name string) (*model.Executor, error) {
	e, err := db.G[model.Executor]().
		Where(query.Executor.Name.Eq(name)).
		First(context.Background())
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *ExecutorRepo) DeleteByName(name string) error {
	_, err := db.G[model.Executor]().
		Where(query.Executor.Name.Eq(name)).
		Delete(context.Background())
	return err
}

func (r *ExecutorRepo) UpdateRole(id int64, role string) error {
	_, err := db.G[model.Executor]().
		Where(query.Executor.ID.Eq(id)).
		Update(context.Background(), "role", role)
	if err != nil {
		return fmt.Errorf("update role: %w", err)
	}

	return nil
}

func (r *ExecutorRepo) FindByNamePrefix(prefix string, limit int) ([]model.Executor, error) {
	executors, err := db.G[model.Executor]().
		Where(query.Executor.Name.Like(prefix + "%")).
		Order(query.Executor.Name.Asc()).
		Limit(limit).
		Find(context.Background())
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if executors == nil {
		return []model.Executor{}, nil
	}

	return executors, nil
}
