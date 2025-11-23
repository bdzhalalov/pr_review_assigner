package models

import (
	"time"
)

type Team struct {
	ID        uint   `gorm:"primaryKey"`
	TeamName  string `gorm:"uniqueIndex;size:256"`
	Members   []User `gorm:"foreignKey:TeamID;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
