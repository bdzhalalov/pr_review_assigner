package add_team

import (
	teamDTO "github.com/bdzhalalov/pr-review-assigner/internal/team/dto"
	userDTO "github.com/bdzhalalov/pr-review-assigner/internal/user/dto"
	"github.com/bdzhalalov/pr-review-assigner/pkg/errors"
	"net/http"
)

type TeamServiceInterface interface {
	Create(input teamDTO.TeamDTO) (teamDTO.TeamDTO, *errors.BaseError)
	GetByName(input teamDTO.GetTeamByNameDTO) (teamDTO.TeamDTO, *errors.BaseError)
}

type UserServiceInterface interface {
	GetUsersByIDs(input userDTO.GetUsersByIdsDTO) ([]userDTO.UserDTO, *errors.BaseError)
	EnsureUsers(input []userDTO.UserDTO) *errors.BaseError
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
	teamNameDTO := teamDTO.GetTeamByNameDTO{
		TeamName: request.TeamName,
	}
	_, err := a.teamService.GetByName(teamNameDTO)
	if err == nil {
		return AddTeamResponse{}, (&errors.TeamExistsError{}).New()
	} else if err.Code != http.StatusNotFound {
		return AddTeamResponse{}, (&errors.InternalServerError{}).New()
	}

	users := make([]userDTO.UserDTO, 0, len(request.Members))
	for _, m := range request.Members {
		users = append(users, userDTO.UserDTO{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
			TeamName: request.TeamName,
		})
	}

	if err := a.userService.EnsureUsers(users); err != nil {
		return AddTeamResponse{}, err
	}

	dto := a.getDTOFromRequest(request)

	createdTeamDTO, err := a.teamService.Create(dto)
	if err != nil {
		return AddTeamResponse{}, err
	}

	response := a.getResponseFromDTO(createdTeamDTO)
	return response, nil
}

func (a *TeamApp) getDTOFromRequest(request AddTeamRequest) teamDTO.TeamDTO {
	dto := teamDTO.TeamDTO{
		TeamName: request.TeamName,
		Members:  make([]teamDTO.MemberDTO, 0, len(request.Members)),
	}
	for _, member := range request.Members {
		dto.Members = append(dto.Members, teamDTO.MemberDTO{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return dto
}

func (a *TeamApp) getResponseFromDTO(dto teamDTO.TeamDTO) AddTeamResponse {
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
