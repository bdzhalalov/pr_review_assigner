package errors

import "net/http"

type PrExistsError struct {
	Message string
	Code    int
}

func (e *PrExistsError) New() *BaseError {
	return &BaseError{
		Message: "PR already exists",
		Code:    http.StatusConflict,
	}
}
