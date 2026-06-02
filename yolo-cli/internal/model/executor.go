package model

import (
	"fmt"
	"strings"
	"time"
)

type Executor struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(64);uniqueIndex;not null" json:"name"`
	Role      string    `gorm:"column:role;type:varchar(16);not null;default:developer" json:"role"`
	Soul      string    `gorm:"column:soul;type:text;default:''" json:"soul"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Executor) TableName() string { return "executors" }

const (
	RoleLeader    = "leader"
	RoleArchitect = "architect"
	RoleDeveloper = "developer"
	RoleQA        = "qa"
)

var AllowedRoles = []string{RoleLeader, RoleArchitect, RoleDeveloper, RoleQA}

func NewExecutor(name, soul string) *Executor {
	return &Executor{
		Name: strings.TrimSpace(name),
		Role: RoleDeveloper,
		Soul: soul,
	}
}

func (e *Executor) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("name is required")
	}

	if len(e.Name) > 64 {
		return fmt.Errorf("name must not exceed 64 characters")
	}

	if strings.Contains(e.Name, " ") {
		return fmt.Errorf("name must not contain spaces")
	}

	return nil
}

func (e *Executor) SetRole(role string) error {
	for _, r := range AllowedRoles {
		if r == role {
			e.Role = role
			return nil
		}
	}

	return fmt.Errorf("invalid role: %s", role)
}
