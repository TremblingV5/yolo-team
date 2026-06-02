package model

import (
	"fmt"
	"strings"
	"time"
)

type Project struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Key         string    `gorm:"column:key;type:varchar(64);uniqueIndex;not null" json:"key"`
	Name        string    `gorm:"column:name;type:varchar(32);not null" json:"name"`
	Description string    `gorm:"column:description;type:varchar(128);default:''" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Project) TableName() string { return "projects" }

func NewProject(name, description string) *Project {
	return &Project{
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
	}
}

func (p *Project) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}

	if len(p.Name) > 32 {
		return fmt.Errorf("name must not exceed 32 characters")
	}

	if len(p.Description) > 128 {
		return fmt.Errorf("description must not exceed 128 characters")
	}

	return nil
}

func (p *Project) ApplyUpdate(name, description *string) {
	if name != nil {
		p.Name = strings.TrimSpace(*name)
	}

	if description != nil {
		p.Description = strings.TrimSpace(*description)
	}
}
