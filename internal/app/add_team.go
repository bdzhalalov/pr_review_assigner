package app

import (
	"github.com/bdzhalalov/pr-review-assigner/internal/team"
	"github.com/bdzhalalov/pr-review-assigner/internal/user"
	"github.com/bdzhalalov/pr-review-assigner/pkg/errors"
)

type TeamApp struct {
	teamService team.Service
	userService user.Service
}

func NewTeamApp(ts team.Service, us user.Service) *TeamApp {
	return &TeamApp{
		teamService: ts,
		userService: us,
	}
}

func (a *TeamApp) AddTeam(input team.AddTeamDTO) (*team.Team, *errors.BaseError) {
	if existing, _ := a.teamService.GetByName(input.TeamName); existing != nil {
		var teamExistsError errors.BaseAbstractError = &errors.TeamExistsError{}
		return nil, teamExistsError.New()
	}

	users := make([]user.User, 0, len(input.Members))
	for _, m := range input.Members {
		users = append(users, user.User{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	if err := a.userService.EnsureUsers(users); err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(input.Members))
	for _, m := range input.Members {
		userIDs = append(userIDs, m.UserID)
	}

	dbUsers, err := a.userService.GetUsersByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	teamModel := &team.Team{
		TeamName: input.TeamName,
		Members:  dbUsers,
	}

	createdTeam, err := a.teamService.Create(teamModel)
	if err != nil {
		return nil, err
	}

	return createdTeam, nil
}
