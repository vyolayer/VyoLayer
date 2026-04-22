package apikey

import "time"

/*
==================================================
API Key Response
==================================================
*/

type APIKeyResponse struct {
	ID string `json:"id"`

	OrganizationID string `json:"organization_id"`
	ProjectID      string `json:"project_id"`

	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	Prefix string `json:"prefix"`

	Environment string `json:"environment"`
	Status      string `json:"status"`

	Scopes []string `json:"scopes"`

	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	LastUsedIP string     `json:"last_used_ip,omitempty"`

	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/*
==================================================
Create Response (includes secret once)
==================================================
*/

type CreateAPIKeyResponse struct {
	APIKey APIKeyResponse `json:"api_key"`
	Secret string         `json:"secret"`
}

/*
==================================================
List Response
==================================================
*/

type ListAPIKeysResponse struct {
	Items []APIKeyResponse `json:"items"`

	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

/*
==================================================
Single Get Response
==================================================
*/

type GetAPIKeyResponse struct {
	APIKey APIKeyResponse `json:"api_key"`
}

/*
==================================================
Rotate Response
==================================================
*/

type RotateAPIKeyResponse struct {
	APIKey APIKeyResponse `json:"api_key"`
	Secret string         `json:"secret"`
}

/*
==================================================
Generic Success Response
==================================================
*/

type SuccessResponse struct {
	Success bool `json:"success"`
}
