package team

import (
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

func (r *TeamRepository) Create(team *Team) (*Team, error) {
	if err := r.db.Create(team).Error; err != nil {
		return nil, err
	}

	return team, nil
}

func (r *TeamRepository) GetByName(name string) (*Team, error) {
	var team Team

	if err := r.db.Where("team_name = ?", name).First(&team).Error; err != nil {
		return nil, err
	}

	return &team, nil
}
