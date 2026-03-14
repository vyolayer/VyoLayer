package domain

import (
	"strings"
	"vyolayer/pkg/errors"
)

type Email struct {
	value string
}

func NewEmail(value string) (*Email, *errors.AppError) {
	if value == "" {
		return nil, InvalidEmailError(value)
	}
	return &Email{value: value}, nil
}

func (e *Email) String() string {
	return strings.ToLower(strings.TrimSpace(e.value))
}

func (e *Email) IsValid() bool {
	return e.String() != ""
}
