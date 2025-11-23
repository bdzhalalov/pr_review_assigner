package team

import (
	"errors"
	"github.com/bdzhalalov/pr-review-assigner/internal/models"
	"github.com/bdzhalalov/pr-review-assigner/internal/team/dto"
	customErrors "github.com/bdzhalalov/pr-review-assigner/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Create(team *models.Team) (*models.Team, error)
	GetByName(id string) (*models.Team, error)
}

type Service struct {
	repo   RepositoryInterface
	logger *logrus.Logger
}

func NewTeamService(repo RepositoryInterface, logger *logrus.Logger) *Service {

	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(input dto.TeamRequestDTO) (dto.TeamResponseDTO, *customErrors.BaseError) {
	team := &models.Team{
		TeamName: input.TeamName,
	}

	res, err := s.repo.Create(team)

	if err != nil {
		s.logger.Errorf("Error while creating new team: %s", err)

		return dto.TeamResponseDTO{}, (&customErrors.InternalServerError{}).New()
	}

	return s.getDTOFromStruct(res), nil
}

func (s *Service) GetByName(input dto.TeamRequestDTO) (dto.TeamResponseDTO, *customErrors.BaseError) {
	team, err := s.repo.GetByName(input.TeamName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.TeamResponseDTO{}, (&customErrors.NotFoundError{}).New("Team not found")
		}

		s.logger.Errorf("Error while getting team by name: %s", err)

		return dto.TeamResponseDTO{}, (&customErrors.InternalServerError{}).New()
	}

	return s.getDTOFromStruct(team), nil
}

func (s *Service) getDTOFromStruct(team *models.Team) dto.TeamResponseDTO {
	teamDTO := dto.TeamResponseDTO{
		ID:       team.ID,
		TeamName: team.TeamName,
		Members:  make([]dto.MemberDTO, 0, len(team.Members)),
	}

	for _, m := range team.Members {
		teamDTO.Members = append(teamDTO.Members, dto.MemberDTO{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	return teamDTO
}
