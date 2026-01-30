package email

import "strings"

type Email struct {
	value string
}

func NewEmail(value string) *Email {
	return &Email{value: value}
}

func (e *Email) Value() string {
	return strings.ToLower(strings.TrimSpace(e.value))
}
