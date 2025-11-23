package pullrequest

import (
	"errors"
	"github.com/bdzhalalov/pr-review-assigner/internal/models"
	"github.com/bdzhalalov/pr-review-assigner/internal/pullrequest/dto"
	customErrors "github.com/bdzhalalov/pr-review-assigner/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type RepositoryInterface interface {
	GetByID(prID string) (*models.PullRequest, error)
	GetByReviewerID(userID string) ([]models.PullRequest, error)
	GetReviewersByTeam(teamName string, authorId string) ([]models.User, error)
	GetPullRequestReviewers(prID string) ([]models.User, error)
	GetAvailableReviewerForReassign(teamName string, reviewerIds []string) (*models.User, error)
	ReassignReviewer(prID string, oldReviewerID string, newReviewerID string) error
	Create(pr *models.PullRequest) (*models.PullRequest, error)
	Update(prID string, fields map[string]interface{}) (*models.PullRequest, error)
}

type Service struct {
	repo   RepositoryInterface
	logger *logrus.Logger
}

func NewPrService(repo RepositoryInterface, logger *logrus.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) GetPrByID(input dto.GetPrByIdDTO) (dto.PrResponseDTO, *customErrors.BaseError) {
	pr, err := s.repo.GetByID(input.PullRequestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.PrResponseDTO{}, (&customErrors.NotFoundError{}).New("Pull Request not found")
		}

		s.logger.Errorf("Error while getting pull request by ID: %s", err)

		return dto.PrResponseDTO{}, (&customErrors.InternalServerError{}).New()
	}

	return s.getDTOFromStruct(pr), nil
}

func (s *Service) GetPrByReviewerID(input dto.GetPrByReviewerIdDTO) ([]dto.PrResponseDTO, *customErrors.BaseError) {
	pullRequests, err := s.repo.GetByReviewerID(input.ReviewerID)
	if err != nil {
		return []dto.PrResponseDTO{}, (&customErrors.InternalServerError{}).New()
	}

	DTOs := make([]dto.PrResponseDTO, 0, len(pullRequests))

	for _, pr := range pullRequests {
		DTOs = append(DTOs, dto.PrResponseDTO{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		})
	}

	return DTOs, nil
}

func (s *Service) CreatePullRequest(input dto.PrRequestDTO) (dto.PrResponseDTO, *customErrors.BaseError) {
	reviewers, err := s.repo.GetReviewersByTeam(input.TeamName, input.AuthorID)
	if err != nil {
		return dto.PrResponseDTO{}, (&customErrors.InternalServerError{}).New()
	}

	pr := &models.PullRequest{
		PullRequestID:     input.PullRequestID,
		PullRequestName:   input.PullRequestName,
		AuthorID:          input.AuthorID,
		AssignedReviewers: reviewers,
	}

	result, err := s.repo.Create(pr)
	if err != nil {
		return dto.PrResponseDTO{}, (&customErrors.InternalServerError{}).New()
	}

	return s.getDTOFromStruct(result), nil
}

func (s *Service) UpdatePullRequestStatus(input dto.PrUpdateStatusDTO) (dto.PrResponseDTO, *customErrors.BaseError) {
	pr, err := s.repo.GetByID(input.PullRequestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.PrResponseDTO{}, (&customErrors.NotFoundError{}).New("Pull Request not found")
		}

		s.logger.Errorf("Error while getting pull request by ID: %s", err)

		return dto.PrResponseDTO{}, (&customErrors.InternalServerError{}).New()
	}

	if pr.Status == "MERGED" {
		return s.getDTOFromStruct(pr), nil
	}

	updatedPr, err := s.repo.Update(input.PullRequestID, map[string]interface{}{
		"status":    "MERGED",
		"merged_at": time.Now().UTC(),
	})

	return s.getDTOFromStruct(updatedPr), nil
}

func (s *Service) ReassignReviewers(input dto.ReassignReviewersDTO) *customErrors.BaseError {
	reviewers, err := s.repo.GetPullRequestReviewers(input.PullRequestID)
	if err != nil {
		return (&customErrors.InternalServerError{}).New()
	}
	if s.containsUserID(reviewers, input.ReviewerID) == false {
		return (&customErrors.PrReviewerNotAssignedError{}).New("User is not reviewer")
	}

	excludedIds := make([]string, 0, len(reviewers)+1)
	for _, reviewer := range reviewers {
		excludedIds = append(excludedIds, reviewer.UserID)
	}
	excludedIds = append(excludedIds, input.AuthorID)

	availableReviewer, err := s.repo.GetAvailableReviewerForReassign(input.TeamName, excludedIds)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return (&customErrors.PrNoReviewCandidateError{}).New()
		}

		return (&customErrors.InternalServerError{}).New()
	}

	err = s.repo.ReassignReviewer(input.PullRequestID, input.ReviewerID, availableReviewer.UserID)
	if err != nil {
		return (&customErrors.InternalServerError{}).New()
	}

	return nil
}

func (s *Service) getDTOFromStruct(pr *models.PullRequest) dto.PrResponseDTO {
	prDTO := dto.PrResponseDTO{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: make([]string, 0, len(pr.AssignedReviewers)),
	}

	for _, reviewer := range pr.AssignedReviewers {
		prDTO.AssignedReviewers = append(prDTO.AssignedReviewers, reviewer.UserID)
	}

	return prDTO
}

func (s *Service) containsUserID(users []models.User, id string) bool {
	for _, u := range users {
		if u.UserID == id {
			return true
		}
	}
	return false
}
