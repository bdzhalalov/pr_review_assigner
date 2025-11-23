package errors

import "net/http"

type PrNoReviewCandidateError struct {
	Code    int
	Message string
}

func (e *PrNoReviewCandidateError) New() *BaseError {
	return &BaseError{
		Message: "No active replacement candidate in team",
		Code:    http.StatusConflict,
	}
}
