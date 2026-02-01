package repository

import "worklayer/internal/domain"

// RepositoryError is an interface that represents a repository error.
type RepositoryError *domain.AppError

// NewRepositoryError creates a new repository error.
func NewRepositoryError(code int, message string) RepositoryError {
	return domain.NewRepositoryError(code, message)
}
