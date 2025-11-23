package pullrequest

import (
	prDTO "github.com/bdzhalalov/pr-review-assigner/internal/pullrequest/dto"
	userDTO "github.com/bdzhalalov/pr-review-assigner/internal/user/dto"
	"github.com/bdzhalalov/pr-review-assigner/pkg/errors"
)

type ReassignPrApp struct {
	userService UserServiceInterface
	prService   PrServiceInterface
}

func NewReassignPrApp(usi UserServiceInterface, psi PrServiceInterface) *ReassignPrApp {
	return &ReassignPrApp{
		userService: usi,
		prService:   psi,
	}
}

func (a *ReassignPrApp) ReassignPullRequest(input ReassignPrRequest) (ReassignPrResponse, *errors.BaseError) {
	user, err := a.userService.GetUserByID(userDTO.GetUserByIDDTO{UserID: input.OldReviewerID})
	if err != nil {
		return ReassignPrResponse{}, err
	}
	dto := prDTO.GetPrByIdDTO{PullRequestID: input.PullRequestID}

	pr, err := a.prService.GetPrByID(dto)
	if err != nil {
		return ReassignPrResponse{}, err
	}

	if pr.Status == "MERGED" {
		return ReassignPrResponse{}, (&errors.PrMergedError{}).New()
	}

	if user.IsActive == false {
		return ReassignPrResponse{}, (&errors.PrReviewerNotAssignedError{}).New("User is not active")
	}

	err = a.prService.ReassignReviewers(prDTO.ReassignReviewersDTO{
		PullRequestID: pr.PullRequestID,
		TeamName:      user.TeamName,
		ReviewerID:    input.OldReviewerID,
		AuthorID:      pr.AuthorID,
	})
	if err != nil {
		return ReassignPrResponse{}, err
	}

	updatedPr, err := a.prService.GetPrByID(dto)

	return a.getResponseFromDTO(updatedPr, input.OldReviewerID), err
}

func (a *ReassignPrApp) getResponseFromDTO(pr prDTO.PrResponseDTO, oldReviewerID string) ReassignPrResponse {
	prResponse := AddPrResponse{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
	}

	return ReassignPrResponse{
		PullRequest: prResponse,
		ReplacedBy:  oldReviewerID,
	}
}
