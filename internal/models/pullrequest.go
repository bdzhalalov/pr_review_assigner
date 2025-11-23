package models

import (
	"time"
)

type PullRequest struct {
	PullRequestID     string `gorm:"primaryKey;size:256"`
	PullRequestName   string `gorm:"size:256"`
	AuthorID          string `gorm:";index;size:256"`
	Status            string `gorm:"type:enum('OPEN','MERGED');default:'OPEN'"`
	AssignedReviewers []User `gorm:"many2many:pull_request_reviewers;joinForeignKey:PullRequestID;joinReferences:UserID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	MergedAt          *time.Time
}
