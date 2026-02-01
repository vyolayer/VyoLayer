package dto

import (
	"time"
	"worklayer/internal/domain"
)

type UserDTO struct {
	ID              string `json:"id"`
	FullName        string `json:"fullName"`
	Email           string `json:"email"`
	Status          string `json:"status"`
	IsEmailVerified bool   `json:"isEmailVerified"`
	JoinedAt        string `json:"joinedAt"`
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
