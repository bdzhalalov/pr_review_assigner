package pullrequest

import (
	"github.com/bdzhalalov/pr-review-assigner/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewPrRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByID(prID string) (*models.PullRequest, error) {
	var pr models.PullRequest

	err := r.db.Preload("AssignedReviewers").
		Where("pull_request_id = ?", prID).
		First(&pr).Error

	if err != nil {
		return nil, err
	}

	return &pr, nil
}

func (r *Repository) GetByReviewerID(reviewerID string) ([]models.PullRequest, error) {
	var prs []models.PullRequest
	err := r.db.Joins("JOIN pull_request_reviewers prr ON prr.pull_request_id = pull_requests.pull_request_id").
		Where("prr.user_id = ?", reviewerID).
		Find(&prs).Error
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func (r *Repository) GetReviewersByTeam(teamName string, authorID string) ([]models.User, error) {
	var reviewers []models.User
	if err := r.db.Where(`
    		team_id = (SELECT id FROM teams WHERE team_name = ? LIMIT 1)
    		AND is_active = TRUE
    		AND user_id != ?
		`, teamName, authorID).
		Order("RAND()").
		Limit(2).
		Find(&reviewers).Error; err != nil {
		return nil, err
	}

	return reviewers, nil
}

func (r *Repository) GetPullRequestReviewers(prID string) ([]models.User, error) {
	var pr models.PullRequest
	if err := r.db.Preload("AssignedReviewers").
		Where("pull_request_id = ?", prID).
		Find(&pr).
		Error; err != nil {
		return nil, err
	}

	return pr.AssignedReviewers, nil
}

func (r *Repository) GetAvailableReviewerForReassign(teamName string, reviewerIds []string) (*models.User, error) {
	var reviewer models.User
	if err := r.db.Where(`
    		team_id = (SELECT id FROM teams WHERE team_name = ? LIMIT 1)
    		AND is_active = TRUE
    		AND user_id NOT IN ?
		`, teamName, reviewerIds).
		Order("RAND()").
		Limit(1).
		First(&reviewer).Error; err != nil {
		return nil, err
	}

	return &reviewer, nil
}

func (r *Repository) ReassignReviewer(prID string, oldReviewerID string, newReviewerID string) error {
	tx := r.db.Begin()

	if err := tx.Exec(`
        DELETE FROM pull_request_reviewers 
        WHERE pull_request_id = ? AND user_id = ?
    `, prID, oldReviewerID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Exec(`
        INSERT INTO pull_request_reviewers (pull_request_id, user_id) 
        VALUES (?, ?)
    `, prID, newReviewerID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *Repository) Create(pr *models.PullRequest) (*models.PullRequest, error) {
	if err := r.db.Create(&pr).Error; err != nil {
		return nil, err
	}

	return pr, nil
}

func (r *Repository) Update(prID string, fields map[string]interface{}) (*models.PullRequest, error) {
	if err := r.db.Model(&models.PullRequest{}).
		Where("pull_request_id = ?", prID).
		Updates(fields).Error; err != nil {
		return nil, err
	}

	//Since MySql does not return an updated record, an additional query is required
	updated, err := r.GetByID(prID)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
