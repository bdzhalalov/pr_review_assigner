package team

import (
	"errors"
	customErrors "github.com/bdzhalalov/pr-review-assigner/pkg/errors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
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

func (s *Service) Create(team *Team) (*Team, *customErrors.BaseError) {
	res, err := s.repo.Create(team)

	if err != nil {
		s.logger.Errorf("Error while creating new team: %s", err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			var notFoundError customErrors.BaseAbstractError = &customErrors.NotFoundError{}
			return nil, notFoundError.New()
		}
		var internalError customErrors.BaseAbstractError = &customErrors.InternalServerError{}
		return nil, internalError.New()
	}

	return res, nil
}

func (s *Service) GetByName(name string) (*Team, *customErrors.BaseError) {
	team, err := s.repo.GetByName(name)
	if err != nil {
		s.logger.Errorf("Error while getting team by name: %s", err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			var notFoundError customErrors.BaseAbstractError = &customErrors.NotFoundError{}
			return nil, notFoundError.New()
		}
		var internalError customErrors.BaseAbstractError = &customErrors.InternalServerError{}
		return nil, internalError.New()
	}

	return team, nil
}
