package user

import (
	"errors"
	"github.com/bdzhalalov/pr-review-assigner/internal/models"
	"github.com/bdzhalalov/pr-review-assigner/internal/user/dto"
	customErrors "github.com/bdzhalalov/pr-review-assigner/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Create(user *models.User) (*models.User, error)
	Update(userId string, fields map[string]interface{}) (*models.User, error)
	GetByUserID(userId string) (*models.User, error)
	GetByIDs(ids []string) ([]models.User, error)
	UpsertUsers(users []models.User) error
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

func (s *Service) GetUserByID(input dto.GetUserByIDDTO) (dto.UserResponseDTO, *customErrors.BaseError) {
	user, err := s.repo.GetByUserID(input.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponseDTO{}, (&customErrors.NotFoundError{}).New()
		}

		s.logger.Errorf("Error while getting user by ID: %s", err)

		return dto.UserResponseDTO{}, (&customErrors.InternalServerError{}).New()
	}

	output := dto.UserResponseDTO{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.Team.TeamName,
		IsActive: user.IsActive,
	}

	return output, nil
}

func (s *Service) GetUsersByIDs(input dto.GetUsersByIdsDTO) ([]dto.UserResponseDTO, *customErrors.BaseError) {
	users, err := s.repo.GetByIDs(input.IDs)
	if err != nil {
		s.logger.Errorf("error getting users by ids: %s", err.Error())

		return nil, (&customErrors.InternalServerError{}).New()
	}

	output := s.getDTOFromStruct(users)

	return output, nil
}

func (s *Service) UpdateUserActivity(input dto.UpdateUserActivityDTO) (dto.UserResponseDTO, *customErrors.BaseError) {
	_, err := s.repo.GetByUserID(input.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponseDTO{}, (&customErrors.NotFoundError{}).New()
		}

		s.logger.Errorf("Error while getting user by ID: %s", err)

		return dto.UserResponseDTO{}, (&customErrors.InternalServerError{}).New()
	}

	updatedUser, err := s.repo.Update(input.UserID, map[string]interface{}{
		"is_active": input.IsActive,
	})
	if err != nil {
		s.logger.Errorf("Error while updating user activity: %s", err)

		return dto.UserResponseDTO{}, (&customErrors.InternalServerError{}).New()
	}

	userDTO := dto.UserResponseDTO{
		UserID:   updatedUser.UserID,
		Username: updatedUser.Username,
		TeamName: updatedUser.Team.TeamName,
		IsActive: updatedUser.IsActive,
	}

	return userDTO, nil
}

func (s *Service) CreateOrUpdate(input []dto.UserRequestDTO) *customErrors.BaseError {
	users := s.getStructFromDTO(input)

	err := s.repo.UpsertUsers(users)
	if err != nil {
		return (&customErrors.InternalServerError{}).New()
	}

	return nil
}

func (s *Service) getDTOFromStruct(users []models.User) []dto.UserResponseDTO {
	DTOs := make([]dto.UserResponseDTO, 0, len(users))

	for _, u := range users {
		DTOs = append(DTOs, dto.UserResponseDTO{
			UserID:   u.UserID,
			Username: u.Username,
			TeamName: u.Team.TeamName,
			IsActive: u.IsActive,
		})
	}

	return DTOs
}

func (s *Service) getStructFromDTO(DTOs []dto.UserRequestDTO) []models.User {
	users := make([]models.User, 0, len(DTOs))
	for _, userDTO := range DTOs {
		users = append(users, models.User{
			UserID:   userDTO.UserID,
			Username: userDTO.Username,
			TeamID:   userDTO.TeamID,
			IsActive: userDTO.IsActive,
		})
	}

	return users
}
