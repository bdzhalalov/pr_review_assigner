package user

import (
	prDTO "github.com/bdzhalalov/pr-review-assigner/internal/pullrequest/dto"
	userDTO "github.com/bdzhalalov/pr-review-assigner/internal/user/dto"
	"github.com/bdzhalalov/pr-review-assigner/pkg/errors"
)

type PrServiceInterface interface {
	GetPrByReviewerID(input prDTO.GetPrByReviewerIdDTO) ([]prDTO.PrResponseDTO, *errors.BaseError)
}

type UserServiceInterface interface {
	GetUserByID(input userDTO.GetUserByIDDTO) (userDTO.UserResponseDTO, *errors.BaseError)
}

type UserApp struct {
	prService   PrServiceInterface
	userService UserServiceInterface
}

func NewUserApp(psi PrServiceInterface, usi UserServiceInterface) *UserApp {
	return &UserApp{
		prService:   psi,
		userService: usi,
	}
}

func (a *UserApp) GetUserPullRequestReviews(input GetPrReviewsRequest) (GetPrReviewsResponse, *errors.BaseError) {
	_, err := a.userService.GetUserByID(userDTO.GetUserByIDDTO{UserID: input.ReviewerID})
	if err != nil {
		return GetPrReviewsResponse{}, err
	}

	pullRequests, err := a.prService.GetPrByReviewerID(prDTO.GetPrByReviewerIdDTO{ReviewerID: input.ReviewerID})
	if err != nil {
		return GetPrReviewsResponse{}, err
	}

	return a.makeResponse(input.ReviewerID, pullRequests), nil
}

func (a *UserApp) makeResponse(reviewerID string, prs []prDTO.PrResponseDTO) GetPrReviewsResponse {
	response := GetPrReviewsResponse{
		ReviewerID:   reviewerID,
		PullRequests: make([]PullRequestResponse, 0, len(prs)),
	}

	for _, pr := range prs {
		response.PullRequests = append(response.PullRequests, PullRequestResponse{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		})
	}

	return response
}
