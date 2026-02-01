package domain

import (
	"log"
	"time"
	"worklayer/internal/platform/database/types"
)

var (
	ErrInvalidEmail    = NewDomainError(400, "Invalid email address")
	ErrPasswordWeak    = NewDomainError(400, "Password must be at least 8 characters")
	ErrInvalidPassword = NewDomainError(400, "Invalid password")
)

type User struct {
	ID              types.UserID
	FullName        string
	Email           string
	HashedPassword  string
	IsActive        bool
	IsEmailVerified bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewUser(email, rawPassword, fullName string) (*User, DomainError) {
	hashedPassword, err := NewPassword(rawPassword)

	if err != nil {
		return nil, err
	}

	theEmail, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	id := types.NewUserID()
	return &User{
		ID:              id,
		FullName:        fullName,
		Email:           theEmail.String(),
		HashedPassword:  hashedPassword.HashedPassword(),
		IsActive:        true,
		IsEmailVerified: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

func ReconstructUser(
	id types.UserID,
	email, hashedPassword, fullName string,
	isActive, isEmailVerified bool,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		ID:              id,
		FullName:        fullName,
		Email:           email,
		HashedPassword:  hashedPassword,
		IsActive:        isActive,
		IsEmailVerified: isEmailVerified,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

// VerifyPassword verifies the password.
func (u *User) VerifyPassword(password string) bool {
	p, err := NewPassword(password)
	if err != nil {
		log.Println("DOMAIN.USER :: VerifyPassword() err: ", err)
		return false
	}

	return p.CheckPassword(u.HashedPassword)
}

// ChangePassword changes the user's password.
func (u *User) ChangePassword(oldPassword, newPassword string) DomainError {
	// Verify old password
	if !u.VerifyPassword(oldPassword) {
		return ErrInvalidPassword
	}

	// Check if new password is weak
	if len(newPassword) < 8 {
		return ErrPasswordWeak
	}

	// Hash new password
	pass, err := NewPassword(newPassword)
	if err != nil {
		return err
	}
	u.HashedPassword = pass.HashedPassword()
	u.UpdatedAt = time.Now()

	return nil
}
