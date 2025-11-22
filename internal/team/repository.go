package team

import (
	"github.com/bdzhalalov/pr-review-assigner/internal/models"
	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (r *TeamRepository) Create(team *models.Team) (*models.Team, error) {
	if err := r.db.Create(&team).Error; err != nil {
		return nil, err
	}

	return team, nil
}

func (r *TeamRepository) GetByName(name string) (*models.Team, error) {
	var team models.Team

	if err := r.db.Where("team_name = ?", name).
		Preload("Members").
		First(&team).Error; err != nil {
		return nil, err
	}

	return &team, nil
}
