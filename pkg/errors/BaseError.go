package errors

type BaseError struct {
	Message string
	Code    int
}

type BaseAbstractError interface {
	New() *BaseError
}
