package apikey

import "time"

/*
==================================================
Create API Key
==================================================
*/

type CreateAPIKeyRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=100"`
	Description string     `json:"description,omitempty"`
	Environment string     `json:"environment" validate:"required,oneof=dev live"`
	Scopes      []string   `json:"scopes,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

/*
==================================================
Update API Key
==================================================
*/

type UpdateAPIKeyRequest struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Scopes      []string   `json:"scopes,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

/*
==================================================
List API Keys (Query Params)
==================================================
*/

type ListAPIKeysQuery struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
	Status string `query:"status"`
}

/*
==================================================
Rotate API Key
==================================================
*/

type RotateAPIKeyRequest struct {
	// no body needed usually, placeholder for extensibility
}

/*
==================================================
Revoke API Key
==================================================
*/

type RevokeAPIKeyRequest struct {
	Reason string `json:"reason,omitempty"`
}
