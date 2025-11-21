package team

import (
	"time"
)

type Team struct {
	ID        uint         `gorm:"primaryKey"`
	TeamName  string       `gorm:"uniqueIndex, size:255"`
	Members   []TeamMember `gorm:"json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TeamMember struct {
	UserID   string `gorm:"size:64"`
	Username string `gorm:"size:255"`
	IsActive bool
}
