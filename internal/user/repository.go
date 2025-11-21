package user

import (
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *User) (*User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByUserID(userId string) (*User, error) {
	var user User
	if err := r.db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByIDs(ids []string) ([]User, error) {
	var users []User
	if err := r.db.Where("user_id in (?)", ids).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) CreateBatch(users []*User) error {
	if err := r.db.Create(users).Error; err != nil {
		return err
	}

	return nil
}
