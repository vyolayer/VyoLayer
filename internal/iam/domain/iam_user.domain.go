package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusPending   UserStatus = "pending"
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

func (s UserStatus) String() string {
	return string(s)
}

type IAMUser struct {
	ID              uuid.UUID
	Email           *Email
	Password        *Password
	FullName        string
	Avatar          *IAMUserAvatar
	IsEmailVerified bool
	Status          UserStatus
	Timestamps      *Timestamps
}

func NewIAMUser(email, password, fullName string) *IAMUser {
	return &IAMUser{
		ID:              uuid.New(),
		Email:           NewEmail(email),
		Password:        NewPassword(password),
		FullName:        fullName,
		IsEmailVerified: false,
		Timestamps:      NewTimestamps(),
		Status:          UserStatusPending,
	}
}

// ReconstructIAMUser rebuilds a domain user from persisted storage values.
func ReconstructIAMUser(
	id uuid.UUID,
	email, password, fullName string,
	isEmailVerified bool, status string,
	createdAt, updatedAt time.Time,
) *IAMUser {
	return &IAMUser{
		ID:              id,
		Email:           NewEmail(email),
		Password:        ReconstructPassword(password),
		FullName:        fullName,
		IsEmailVerified: isEmailVerified,
		Timestamps: &Timestamps{
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		Status: UserStatus(status),
	}
}

// Getter
func (u *IAMUser) GetID() uuid.UUID    { return u.ID }
func (u *IAMUser) GetEmail() string    { return u.Email.String() }
func (u *IAMUser) GetPassword() string { return u.Password.String() }
func (u *IAMUser) GetFullName() string { return u.FullName }
func (u *IAMUser) GetStatus() string   { return u.Status.String() }

func (u *IAMUser) VerifyEmail() {
	u.IsEmailVerified = true
	u.Timestamps.Update()
}

func (u *IAMUser) SetFullName(fullName string) {
	u.FullName = fullName
	u.Timestamps.Update()
}

func (u *IAMUser) SetStatus(status UserStatus) {
	u.Status = status
	u.Timestamps.Update()
}

// avatar
func (u *IAMUser) InitAvatar(avatar *IAMUserAvatar) {
	u.Avatar = avatar
}
