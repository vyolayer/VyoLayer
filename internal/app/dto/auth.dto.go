package dto

// # Auth DTOs
//
// RegisterUserSchema is a struct that holds the data for registering a new user
type RegisterUserSchema struct {
	Email    string `json:"email" validate:"required,email,max=255" example:"subhajit@worklayer.com"`
	Password string `json:"password" validate:"required,min=8,max=20,containsany=!@#$%^&*" example:"Password!123"`
	FullName string `json:"fullName" validate:"required,min=3,max=100" example:"Subhajit Pramanik"`
}

// LoginUserSchema is a struct that holds the data for logging in a user
type LoginUserSchema struct {
	Email    string `json:"email" validate:"required,email,max=255" example:"subhajit@worklayer.com"`
	Password string `json:"password" validate:"required" example:"Password!123"`
}

// TokenResponseDTO is a struct that holds the data for tokens
type TokenResponseDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LoginUserResponseDTO is a struct that holds the data for logging in a user
type LoginUserResponseDTO struct {
	TokenResponseDTO
	User UserDTO `json:"user"`
}

// RefreshSessionResponseDTO is a struct that holds the data for refreshing a session
type RefreshSessionResponseDTO struct {
	TokenResponseDTO
}
