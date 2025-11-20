package team

import (
	"github.com/bdzhalalov/pr-review-assigner/pkg/database"
)

type TeamRepository struct {
	db *database.Database
}

func NewTeamRepository(db *database.Database) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (r *TeamRepository) Create(team *Team) (*Team, error) {
	if err := r.db.Connection.Create(team).Error; err != nil {
		return nil, err
	}

	if err := r.db.Connection.Model(team).Association("Members").Replace(team.Members).Error; err != nil {
		return nil, err
	}

	return team, nil
}

func (r *TeamRepository) GetByName(name string) (*Team, error) {
	var team Team

	if err := r.db.Connection.Preload("Members").Where("name = ?", name).First(&team).Error; err != nil {
		return nil, err
	}

	return &team, nil
}
