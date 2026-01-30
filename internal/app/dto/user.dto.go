package dto

type UserDTO struct {
	ID              uint   `json:"id"`
	Email           string `json:"email"`
	IsActive        bool   `json:"isActive"`
	IsEmailVerified bool   `json:"isEmailVerified"`
	FullName        string `json:"fullName"`
	CreatedAt       string `json:"createdAt"`
}
