package errors

import "net/http"

type PrMergedError struct {
	Code    int
	Message string
}

func (e *PrMergedError) New() *BaseError {
	return &BaseError{
		Message: "Cannot reassign on merged PR",
		Code:    http.StatusConflict,
	}
}
