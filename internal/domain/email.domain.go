package domain

import "strings"

type Email struct {
	value string
}

func NewEmail(value string) (*Email, DomainError) {
	if value == "" {
		return nil, NewDomainError(400, "Invalid email address")
	}
	return &Email{value: value}, nil
}

func (e *Email) String() string {
	return strings.ToLower(strings.TrimSpace(e.value))
}

func (e *Email) IsValid() bool {
	return e.String() != ""
}
