package user

import (
	"github.com/bdzhalalov/pr-review-assigner/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RepositoryInterface interface {
	Create(user *User) (*User, error)
	GetByUserID(userId string) (*User, error)
	GetByIDs(ids []string) ([]User, error)
	CreateBatch(users []User) error
}

type Service struct {
	repo   RepositoryInterface
	logger *logrus.Logger
}

func NewUserService(repo RepositoryInterface, logger *logrus.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) GetUsersByIDs(ids []string) ([]User, *errors.BaseError) {
	users, err := s.repo.GetByIDs(ids)
	if err != nil {
		s.logger.Errorf("error getting users by ids: %s", err.Error())

		var internalError errors.BaseAbstractError = &errors.InternalServerError{}
		return nil, internalError.New()
	}

	return users, nil
}

func (s *Service) EnsureUsers(users []User) *errors.BaseError {
	ids := make([]string, len(users))
	for i, u := range users {
		ids[i] = u.UserID
	}

	existing, err := s.repo.GetByIDs(ids)
	if err != nil {
		s.logger.Errorf("error getting existing users: %s", err.Error())

		var internalError errors.BaseAbstractError = &errors.InternalServerError{}
		return internalError.New()
	}

	existingMap := make(map[string]struct{}, len(existing))
	for _, u := range existing {
		existingMap[u.UserID] = struct{}{}
	}

	toCreate := []User{}
	for _, u := range users {
		if _, ok := existingMap[u.UserID]; !ok {
			toCreate = append(toCreate, u)
		}
	}

	if len(toCreate) > 0 {
		createErrors := s.repo.CreateBatch(toCreate)
		if createErrors != nil {
			s.logger.Errorf("error creating users: %s", createErrors.Error())

			var internalError errors.BaseAbstractError = &errors.InternalServerError{}
			return internalError.New()
		}
	}

	return nil
}
