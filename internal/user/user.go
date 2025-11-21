package user

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    string `gorm:"uniqueIndex, size:64" json:"user_id"`
	Username  string `gorm:"size:255"`
	IsActive  bool
	TeamName  string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
