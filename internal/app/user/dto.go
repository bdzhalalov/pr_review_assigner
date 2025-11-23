package user

type GetPrReviewsRequest struct {
	ReviewerID string `json:"user_id"`
}

type GetPrReviewsResponse struct {
	ReviewerID   string                `json:"user_id"`
	PullRequests []PullRequestResponse `json:"pull_requests"`
}

type PullRequestResponse struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}
