package errors

import "net/http"

type TeamExistsError struct {
	Message string
	Code    int
}

func (e *TeamExistsError) New() *BaseError {
	return &BaseError{
		Message: "Team already exists",
		Code:    http.StatusBadRequest,
	}
}
