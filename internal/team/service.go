package team

import (
	"errors"
	"github.com/bdzhalalov/pr-review-assigner/internal/team/dto"
	customErrors "github.com/bdzhalalov/pr-review-assigner/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Create(team *Team) (*Team, error)
	GetByName(id string) (*Team, error)
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

func (s *Service) Create(input dto.TeamDTO) (dto.TeamDTO, *customErrors.BaseError) {
	team := s.getStructFromDto(input)

	res, err := s.repo.Create(team)

	if err != nil {
		s.logger.Errorf("Error while creating new team: %s", err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.TeamDTO{}, (&customErrors.NotFoundError{}).New()
		}
		return dto.TeamDTO{}, (&customErrors.InternalServerError{}).New()
	}

	output := s.getDTOFromStruct(res)

	return output, nil
}

func (s *Service) GetByName(input dto.GetTeamByNameDTO) (dto.TeamDTO, *customErrors.BaseError) {
	team, err := s.repo.GetByName(input.TeamName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.TeamDTO{}, (&customErrors.NotFoundError{}).New()
		}

		s.logger.Errorf("Error while getting team by name: %s", err)

		return dto.TeamDTO{}, (&customErrors.InternalServerError{}).New()
	}

	output := s.getDTOFromStruct(team)

	return output, nil
}

func (s *Service) getStructFromDto(dto dto.TeamDTO) *Team {
	team := &Team{
		TeamName: dto.TeamName,
		Members:  make([]TeamMember, 0, len(dto.Members)),
	}
	for _, member := range dto.Members {
		team.Members = append(team.Members, TeamMember{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return team
}

func (s *Service) getDTOFromStruct(team *Team) dto.TeamDTO {
	teamDTO := dto.TeamDTO{
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
