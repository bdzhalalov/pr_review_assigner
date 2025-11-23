package pullrequest

import (
	prDTO "github.com/bdzhalalov/pr-review-assigner/internal/pullrequest/dto"
	teamDTO "github.com/bdzhalalov/pr-review-assigner/internal/team/dto"
	userDTO "github.com/bdzhalalov/pr-review-assigner/internal/user/dto"
	"github.com/bdzhalalov/pr-review-assigner/pkg/errors"
)

type PrServiceInterface interface {
	GetPrByID(input prDTO.GetPrByIdDTO) (prDTO.PrResponseDTO, *errors.BaseError)
	CreatePullRequest(input prDTO.PrRequestDTO) (prDTO.PrResponseDTO, *errors.BaseError)
	ReassignReviewers(input prDTO.ReassignReviewersDTO) *errors.BaseError
}

type TeamServiceInterface interface {
	GetByName(input teamDTO.TeamRequestDTO) (teamDTO.TeamResponseDTO, *errors.BaseError)
}

type UserServiceInterface interface {
	GetUserByID(input userDTO.GetUserByIDDTO) (userDTO.UserResponseDTO, *errors.BaseError)
}

type PRApp struct {
	CreateApp   *CreatePrApp
	ReassignApp *ReassignPrApp
}

func InitPRApps(tsi TeamServiceInterface, usi UserServiceInterface, psi PrServiceInterface) *PRApp {
	createApp := NewCreatePRApp(tsi, usi, psi)

	reassignApp := NewReassignPrApp(usi, psi)
	return &PRApp{
		CreateApp:   createApp,
		ReassignApp: reassignApp,
	}
}
