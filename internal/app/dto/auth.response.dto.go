package dto

type TokenResponseDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginUserResponseDTO struct {
	Tokens TokenResponseDTO `json:"tokens"`
	User   UserDTO          `json:"user"`
}
