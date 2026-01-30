package dto

import (
	"worklayer/internal/domain"
)

type UserID = domain.UserID

type UserDTO struct {
	ID              UserID `json:"id"`
	Email           string `json:"email"`
	IsActive        bool   `json:"isActive"`
	IsEmailVerified bool   `json:"isEmailVerified"`
	FullName        string `json:"fullName"`
	CreatedAt       string `json:"createdAt"`
}
