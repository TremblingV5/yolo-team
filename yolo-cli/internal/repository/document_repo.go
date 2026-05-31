package repository

import (
	"context"

	"yolo-team/yolo-cli/internal/db"
	"yolo-team/yolo-cli/internal/model"
	"yolo-team/yolo-cli/internal/query"
	"yolo-team/yolo-cli/internal/utils/keyutils"
)

type DocumentRepo struct{}

func NewDocumentRepo() *DocumentRepo {
	return &DocumentRepo{}
}

func (r *DocumentRepo) ListByProject(projectID int64) ([]model.Document, error) {
	docs, err := db.G[model.Document]().
		Where(query.Document.ProjectID.Eq(projectID)).
		Order(query.Document.SortOrder.Asc()).
		Order(query.Document.CreatedAt.Desc()).
		Find(context.Background())
	return docs, err
}

func (r *DocumentRepo) Create(doc *model.Document) (*model.Document, error) {
	if err := db.G[model.Document]().Create(context.Background(), doc); err != nil {
		return nil, err
	}

	doc.Key = keyutils.Document(doc.ID)
	if _, err := db.G[model.Document]().
		Where(query.Document.ID.Eq(doc.ID)).
		Update(context.Background(), "key", doc.Key); err != nil {
		return nil, err
	}

	return doc, nil
}

func (r *DocumentRepo) GetByKey(key string) (*model.Document, error) {
	d, err := db.G[model.Document]().
		Where(query.Document.Key.Eq(key)).
		First(context.Background())
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *DocumentRepo) GetByID(id int64) (*model.Document, error) {
	d, err := db.G[model.Document]().
		Where(query.Document.ID.Eq(id)).
		First(context.Background())
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *DocumentRepo) Save(doc *model.Document) error {
	_, err := db.G[model.Document]().
		Where(query.Document.ID.Eq(doc.ID)).
		Updates(context.Background(), *doc)
	return err
}

func (r *DocumentRepo) Delete(key string) error {
	_, err := db.G[model.Document]().
		Where(query.Document.Key.Eq(key)).
		Delete(context.Background())
	return err
}
