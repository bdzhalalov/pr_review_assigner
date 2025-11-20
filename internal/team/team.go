package team

import (
	"github.com/bdzhalalov/pr-review-assigner/internal/user"
	"time"
)

type Team struct {
	ID        uint        `json:"-"`
	TeamName  string      `gorm:"uniqueIndex" json:"team_name"`
	Members   []user.User `gorm:"many2many:team_members;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"members"`
	CreatedAt time.Time   `json:"-"`
	UpdatedAt time.Time   `json:"-"`
}
