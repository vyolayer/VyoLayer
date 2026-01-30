package service

import "worklayer/internal/utils/response"

// ServiceError interface for custom error type
type ServiceError interface {
	Error() *response.ErrorResponse
	Code() int
}

type serviceError struct {
	err *response.ErrorResponse
}

func (se *serviceError) Error() *response.ErrorResponse {
	return se.err
}

func (se *serviceError) Code() int {
	return se.err.StatusCode
}

func NewServiceError(err *response.ErrorResponse) ServiceError {
	return &serviceError{err: err}
}
