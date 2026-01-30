package dto

// RegisterUserDTO is a struct that holds the data for registering a new user
type RegisterUserDTO struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=20,containsany=!@#$%^&*"`
	FullName string `json:"fullName" validate:"required,min=3,max=100"`
}

type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required"`
}
