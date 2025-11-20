package user

import (
	"github.com/bdzhalalov/pr-review-assigner/internal/team"
	"time"
)

type User struct {
	ID        uint        `json:"-"`
	UserID    string      `gorm:"uniqueIndex" json:"user_id"`
	Username  string      `json:"username"`
	IsActive  bool        `json:"is_active"`
	Teams     []team.Team `gorm:"many2many:team_members;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"teams"`
	CreatedAt time.Time   `json:"-"`
	UpdatedAt time.Time   `json:"-"`
}
