package pullrequest

import (
	prDTO "github.com/bdzhalalov/pr-review-assigner/internal/pullrequest/dto"
	teamDTO "github.com/bdzhalalov/pr-review-assigner/internal/team/dto"
	userDTO "github.com/bdzhalalov/pr-review-assigner/internal/user/dto"
	"github.com/bdzhalalov/pr-review-assigner/pkg/errors"
	"net/http"
)

type CreatePrApp struct {
	teamService TeamServiceInterface
	userService UserServiceInterface
	prService   PrServiceInterface
}

func NewCreatePRApp(tsi TeamServiceInterface, usi UserServiceInterface, psi PrServiceInterface) *CreatePrApp {
	return &CreatePrApp{
		teamService: tsi,
		userService: usi,
		prService:   psi,
	}
}

func (a *CreatePrApp) CreatePullRequest(input AddPrRequest) (AddPrResponse, *errors.BaseError) {
	user, err := a.userService.GetUserByID(userDTO.GetUserByIDDTO{UserID: input.AuthorID})
	if err != nil {
		return AddPrResponse{}, err
	}

	_, err = a.teamService.GetByName(teamDTO.TeamRequestDTO{TeamName: user.TeamName})
	if err != nil {
		return AddPrResponse{}, err
	}

	_, err = a.prService.GetPrByID(prDTO.GetPrByIdDTO{PullRequestID: input.PullRequestID})
	if err == nil {
		return AddPrResponse{}, (&errors.PrExistsError{}).New()
	} else if err.Code != http.StatusNotFound {
		return AddPrResponse{}, err
	}

	pr, err := a.prService.CreatePullRequest(prDTO.PrRequestDTO{
		PullRequestID:   input.PullRequestID,
		PullRequestName: input.PullRequestName,
		AuthorID:        input.AuthorID,
		TeamName:        user.TeamName,
	})
	if err != nil {
		return AddPrResponse{}, err
	}

	output := AddPrResponse{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
	}

	return output, nil
}
