package user

import (
	"github.com/bdzhalalov/pr-review-assigner/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Update(userId string, fields map[string]interface{}) (*models.User, error) {
	if err := r.db.Model(&models.User{}).Where("user_id = ?", userId).Updates(fields).Error; err != nil {
		return nil, err
	}

	//Since MySql does not return an updated record, an additional query is required
	updated, err := r.GetByUserID(userId)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (r *UserRepository) GetByUserID(userId string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Team").Where("user_id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpsertUsers(users []models.User) error {
	updateColumns := []string{"username", "is_active", "team_id"}

	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(&users).Error; err != nil {
		return err
	}

	return nil
}
