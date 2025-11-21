package user

import (
	"github.com/bdzhalalov/pr-review-assigner/internal/user/dto"
	"github.com/bdzhalalov/pr-review-assigner/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RepositoryInterface interface {
	Create(user *User) (*User, error)
	GetByUserID(userId string) (*User, error)
	GetByIDs(ids []string) ([]User, error)
	CreateBatch(users []*User) error
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

func (s *Service) GetUsersByIDs(input dto.GetUsersByIdsDTO) ([]dto.UserDTO, *errors.BaseError) {
	users, err := s.repo.GetByIDs(input.IDs)
	if err != nil {
		s.logger.Errorf("error getting users by ids: %s", err.Error())

		return nil, (&errors.InternalServerError{}).New()
	}

	output := s.getDTOFromStruct(users)

	return output, nil
}

func (s *Service) EnsureUsers(input []dto.UserDTO) *errors.BaseError {
	users := s.getStructFromDTO(input)

	if len(users) == 0 {
		return nil
	}

	ids := make([]string, len(users))
	for i, u := range users {
		ids[i] = u.UserID
	}

	existing, err := s.repo.GetByIDs(ids)
	if err != nil {
		s.logger.Errorf("error getting existing users: %s", err)
		return (&errors.InternalServerError{}).New()
	}

	existingMap := make(map[string]struct{}, len(existing))
	for _, u := range existing {
		existingMap[u.UserID] = struct{}{}
	}

	toCreate := make([]*User, 0, len(users))
	for _, u := range users {
		if _, found := existingMap[u.UserID]; !found {
			toCreate = append(toCreate, &u)
			existingMap[u.UserID] = struct{}{}
		}
	}

	if len(toCreate) > 0 {
		if err := s.repo.CreateBatch(toCreate); err != nil {
			s.logger.Errorf("error creating users: %s", err)
			return (&errors.InternalServerError{}).New()
		}
	}

	return nil
}

func (s *Service) getDTOFromStruct(users []User) []dto.UserDTO {
	DTOs := make([]dto.UserDTO, 0, len(users))

	for _, u := range users {
		DTOs = append(DTOs, dto.UserDTO{
			UserID:   u.UserID,
			Username: u.Username,
			TeamName: u.TeamName,
			IsActive: u.IsActive,
		})
	}

	return DTOs
}

func (s *Service) getStructFromDTO(DTOs []dto.UserDTO) []User {
	users := make([]User, 0, len(DTOs))
	for _, userDTO := range DTOs {
		users = append(users, User{
			UserID:   userDTO.UserID,
			Username: userDTO.Username,
			TeamName: userDTO.TeamName,
			IsActive: userDTO.IsActive,
		})
	}

	return users
}
