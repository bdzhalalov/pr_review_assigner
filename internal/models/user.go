package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    string `gorm:"uniqueIndex;size:64"`
	Username  string `gorm:"size:256"`
	IsActive  bool
	TeamID    uint `gorm:"size:256"`
	Team      Team `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
