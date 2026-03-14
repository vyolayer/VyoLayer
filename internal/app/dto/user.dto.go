package dto

import (
	"time"

	"github.com/vyolayer/vyolayer/internal/domain"
)

type UserDTO struct {
	ID              string `json:"id" example:"user_550e8400-e29b-41d4-a716-446655440000"`
	FullName        string `json:"fullName" example:"Subhajit Pramanik"`
	Email           string `json:"email" example:"subhajit@vyolayer.com"`
	Status          string `json:"status" example:"active" enums:"active,inactive"`
	IsEmailVerified bool   `json:"isEmailVerified" example:"true"`
	JoinedAt        string `json:"joinedAt" example:"2023-01-01T00:00:00Z"`
}

func FromDomainUser(user *domain.User) UserDTO {
	if user == nil {
		return UserDTO{}
	}

	status := "inactive"
	if user.IsActive {
		status = "active"
	}

	return UserDTO{
		ID:              user.ID.String(),
		Email:           user.Email,
		Status:          status,
		IsEmailVerified: user.IsEmailVerified,
		FullName:        user.FullName,
		JoinedAt:        user.CreatedAt.Format(time.RFC3339),
	}
}
