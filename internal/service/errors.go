package service

import "worklayer/internal/domain"

type ServiceError *domain.AppError

func NewServiceError(code int, message string) ServiceError {
	return domain.NewServiceError(code, message)
}
