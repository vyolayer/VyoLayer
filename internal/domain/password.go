package domain

import (
	"worklayer/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

type Password interface {
	HashedPassword() string
	CheckPassword(oldHashedPassword string) bool
}

type password struct {
	hashedPassword string
}

func NewPassword(value string) (Password, *errors.AppError) {
	if len(value) < 8 {
		return nil, InvalidPasswordError("Password must be at least 8 characters")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return nil, PasswordHashFailedError(err)
	}
	return &password{hashedPassword: string(hashedPassword)}, nil
}

func (p *password) HashedPassword() string {
	return p.hashedPassword
}

func (p *password) CheckPassword(oldHashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashedPassword), []byte(oldHashedPassword))
	return err == nil
}
