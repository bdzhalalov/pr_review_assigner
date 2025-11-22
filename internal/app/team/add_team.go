package team

import (
	teamDTO "github.com/bdzhalalov/pr-review-assigner/internal/team/dto"
	userDTO "github.com/bdzhalalov/pr-review-assigner/internal/user/dto"
	"github.com/bdzhalalov/pr-review-assigner/pkg/errors"
	"net/http"
)

type TeamServiceInterface interface {
	Create(input teamDTO.TeamRequestDTO) (teamDTO.TeamResponseDTO, *errors.BaseError)
	GetByName(input teamDTO.TeamRequestDTO) (teamDTO.TeamResponseDTO, *errors.BaseError)
}

type UserServiceInterface interface {
	CreateOrUpdate(input []userDTO.UserRequestDTO) *errors.BaseError
}

type TeamApp struct {
	teamService TeamServiceInterface
	userService UserServiceInterface
}

func NewTeamApp(tsi TeamServiceInterface, usi UserServiceInterface) *TeamApp {
	return &TeamApp{
		teamService: tsi,
		userService: usi,
	}
}

func (a *TeamApp) AddTeam(request AddTeamRequest) (AddTeamResponse, *errors.BaseError) {
	dto := teamDTO.TeamRequestDTO{
		TeamName: request.TeamName,
	}
	_, err := a.teamService.GetByName(dto)
	if err == nil {
		return AddTeamResponse{}, (&errors.TeamExistsError{}).New()
	} else if err.Code != http.StatusNotFound {
		return AddTeamResponse{}, err
	}

	createdTeam, err := a.teamService.Create(dto)
	if err != nil {
		return AddTeamResponse{}, err
	}

	users := make([]userDTO.UserRequestDTO, 0, len(request.Members))
	for _, m := range request.Members {
		users = append(users, userDTO.UserRequestDTO{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
			TeamID:   createdTeam.ID,
		})
	}

	if err := a.userService.CreateOrUpdate(users); err != nil {
		return AddTeamResponse{}, err
	}

	teamWithMembers, err := a.teamService.GetByName(dto)
	if err != nil {
		return AddTeamResponse{}, err
	}

	output := a.getResponseFromDTO(teamWithMembers)

	return output, nil
}

func (a *TeamApp) getResponseFromDTO(dto teamDTO.TeamResponseDTO) AddTeamResponse {
	response := AddTeamResponse{
		TeamName: dto.TeamName,
		Members:  make([]TeamMember, 0, len(dto.Members)),
	}
	for _, member := range dto.Members {
		response.Members = append(response.Members, TeamMember{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return response
}
