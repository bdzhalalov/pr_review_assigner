package errors

import (
	"fmt"
	"net/http"
)

type PrReviewerNotAssignedError struct {
	Code    int
	Message string
}

func (e *PrReviewerNotAssignedError) New(detail string) *BaseError {
	return &BaseError{
		Message: fmt.Sprintf("Reviewer is not assigned to this PR: %s", detail),
		Code:    http.StatusConflict,
	}
}
