package model

import (
	"fmt"
	"strings"
	"time"
)

type Document struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Key       string    `gorm:"column:key;type:varchar(64);uniqueIndex;not null" json:"key"`
	ProjectID int64     `gorm:"column:project_id;not null;index" json:"project_id"`
	Title     string    `gorm:"column:title;type:varchar(128);not null" json:"title"`
	Creator   string    `gorm:"column:creator;type:varchar(64);default:'人类'" json:"creator"`
	SortOrder int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Project *Project `gorm:"foreignKey:ProjectID" json:"-"`
}

func (Document) TableName() string { return "project_documents" }

func NewDocument(projectID int64, title string) *Document {
	return &Document{
		ProjectID: projectID,
		Title:     strings.TrimSpace(title),
	}
}

func (d *Document) Validate() error {
	if d.Title == "" {
		return fmt.Errorf("title is required")
	}

	if len(d.Title) > 128 {
		return fmt.Errorf("title must not exceed 128 characters")
	}

	return nil
}
