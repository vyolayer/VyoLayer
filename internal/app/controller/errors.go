package controller

import "worklayer/internal/domain"

func NewControllerError(code int, message string) error {
	return domain.NewControllerError(code, message)
}
