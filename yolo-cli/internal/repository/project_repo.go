package repository

import (
	"context"

	"yolo-team/yolo-cli/internal/db"
	"yolo-team/yolo-cli/internal/model"
	"yolo-team/yolo-cli/internal/query"
	"yolo-team/yolo-cli/internal/utils/keyutils"
)

type ProjectRepo struct{}

func NewProjectRepo() *ProjectRepo {
	return &ProjectRepo{}
}

func (r *ProjectRepo) List() ([]model.Project, error) {
	projects, err := db.G[model.Project]().
		Order(query.Project.CreatedAt.Desc()).
		Find(context.Background())
	return projects, err
}

func (r *ProjectRepo) Create(p *model.Project) (*model.Project, error) {
	if err := db.G[model.Project]().Create(context.Background(), p); err != nil {
		return nil, err
	}

	p.Key = keyutils.Project(p.ID)
	if _, err := db.G[model.Project]().
		Where(query.Project.ID.Eq(p.ID)).
		Update(context.Background(), "key", p.Key); err != nil {
		return nil, err
	}

	return p, nil
}

func (r *ProjectRepo) GetByKey(key string) (*model.Project, error) {
	p, err := db.G[model.Project]().
		Where(query.Project.Key.Eq(key)).
		First(context.Background())
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProjectRepo) GetByID(id int64) (*model.Project, error) {
	p, err := db.G[model.Project]().
		Where(query.Project.ID.Eq(id)).
		First(context.Background())
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProjectRepo) Save(p *model.Project) error {
	_, err := db.G[model.Project]().
		Where(query.Project.ID.Eq(p.ID)).
		Updates(context.Background(), *p)
	return err
}

func (r *ProjectRepo) Delete(key string) error {
	_, err := db.G[model.Project]().
		Where(query.Project.Key.Eq(key)).
		Delete(context.Background())
	return err
}
