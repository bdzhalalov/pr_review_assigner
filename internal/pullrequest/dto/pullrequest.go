package dto

type GetPrByIdDTO struct {
	PullRequestID string
}

type PrRequestDTO struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
	TeamName        string
}

type PrResponseDTO struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

type PrUpdateStatusDTO struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignReviewersDTO struct {
	PullRequestID string
	TeamName      string
	ReviewerID    string
	AuthorID      string
}
