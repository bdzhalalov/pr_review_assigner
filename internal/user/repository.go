package user

import "github.com/bdzhalalov/pr-review-assigner/pkg/database"

type UserRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *User) (*User, error) {
	if err := r.db.Connection.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByUserID(userId string) (*User, error) {
	var user User
	if err := r.db.Connection.Where("user_id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByIDs(ids []string) ([]User, error) {
	var users []User
	if err := r.db.Connection.Where("user_id in (?)", ids).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) CreateBatch(users []User) error {
	if err := r.db.Connection.Create(&users).Error; err != nil {
		return err
	}

	return nil
}
