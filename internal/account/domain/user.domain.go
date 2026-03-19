package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/shared/auth"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

func (s UserStatus) String() string {
	return string(s)
}

type User struct {
	ID              uuid.UUID
	ProjectID       uuid.UUID
	Email           string
	Username        string
	FirstName       string
	LastName        string
	HashedPassword  string
	Status          UserStatus
	IsEmailVerified bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	LastLoginAt     *time.Time
	Avatar          *Avatar
}

type UserUpdate struct {
	FirstName       string
	LastName        string
	IsEmailVerified bool
	LastLoginAt     *time.Time
	Status          UserStatus
}

func NewUser(
	projectID uuid.UUID,
	email, username, password, firstName, lastName string,
) *User {
	id := uuid.New()

	passHash, _ := auth.GenerateHash(password)

	return &User{
		ID:              id,
		ProjectID:       projectID,
		Email:           strings.ToLower(strings.TrimSpace(email)),
		Username:        strings.ToLower(strings.TrimSpace(username)),
		FirstName:       strings.TrimSpace(firstName),
		LastName:        strings.TrimSpace(lastName),
		HashedPassword:  passHash,
		IsEmailVerified: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		LastLoginAt:     nil,
		Avatar:          nil,
	}
}

func (u *User) IsActive() bool   { return u.Status.String() == UserStatusActive.String() }
func (u *User) IsVerified() bool { return u.IsEmailVerified }

// Get full name
func (u *User) FullName() string { return u.FirstName + " " + u.LastName }

// Init Avatar
func (u *User) InitAvatar() {
	u.Avatar = NewAvatar(u.FirstName)
}

// VerifyPassword checks if the provided password matches the user's hashed password
func (u *User) VerifyPassword(password string) bool {
	return auth.CheckHash(password, u.HashedPassword)
}

// User verification
func (u *User) VerifyEmail() {
	u.IsEmailVerified = true
	u.UpdatedAt = time.Now()
	u.Status = UserStatusActive
}

// Is same password
func (u *User) IsSamePassword(password string) bool {
	return auth.CheckHash(password, u.HashedPassword)
}

// Set new password
func (u *User) ChangePassword(password string) error {
	newHash, e := auth.GenerateHash(password)
	if e != nil {
		return e
	}
	u.HashedPassword = newHash
	u.UpdatedAt = time.Now()
	return nil
}
