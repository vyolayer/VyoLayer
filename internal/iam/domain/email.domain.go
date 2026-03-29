package domain

import "strings"

type Email struct {
	Value string
}

func NewEmail(s string) *Email {
	v := strings.ToLower(strings.TrimSpace(s))

	if v == "" {
		return nil
	}

	return &Email{Value: v}
}

func (e *Email) IsValid() bool {
	if e == nil || e.Value == "" {
		return false
	}

	if !strings.Contains(e.Value, "@") {
		return false
	}

	return true
}

func (e *Email) String() string {
	return e.Value
}
