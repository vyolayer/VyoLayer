package dto

import (
	"time"

	"github.com/vyolayer/vyolayer/internal/domain"
)

// ── Request DTOs ─────────────────────────────────────────────────────────────

type CreateApiKeyRequestDTO struct {
	Name string `json:"name" validate:"required,min=1,max=100" example:"My API Key"`
	Mode string `json:"mode" validate:"required,oneof=dev live" example:"dev"`
}

// ── Response DTOs ────────────────────────────────────────────────────────────

// ApiKeyDTO shows the key in a masked form — only the prefix is visible.
type ApiKeyDTO struct {
	ID           string  `json:"id" example:"api_key_550e8400-e29b-41d4-a716-446655440000"`
	ProjectID    string  `json:"projectId" example:"project_550e8400-e29b-41d4-a716-446655440000"`
	Name         string  `json:"name" example:"My API Key"`
	KeyPrefix    string  `json:"keyPrefix" example:"wl_live_ab3f1234"`
	Mode         string  `json:"mode" example:"live"`
	CreatedBy    string  `json:"createdBy" example:"user_550e8400-e29b-41d4-a716-446655440000"`
	ExpiresAt    *string `json:"expiresAt,omitempty" example:"2024-01-01T00:00:00Z"`
	LastUsedAt   *string `json:"lastUsedAt,omitempty"`
	IsRevoked    bool    `json:"isRevoked" example:"false"`
	RevokedAt    *string `json:"revokedAt,omitempty"`
	RequestLimit int     `json:"requestLimit" example:"1000"`
	RateLimit    int     `json:"rateLimit" example:"60"`
	CreatedAt    string  `json:"createdAt" example:"2023-01-01T00:00:00Z"`
}

// ApiKeyCreatedDTO is returned only once at creation — includes the raw key.
type ApiKeyCreatedDTO struct {
	ApiKeyDTO
	RawKey string `json:"rawKey" example:"wl_live_a1b2c3d4e5f6..."`
}

// ── Domain-to-DTO converters ─────────────────────────────────────────────────

func FromDomainApiKey(ak *domain.ApiKey) *ApiKeyDTO {
	if ak == nil {
		return nil
	}

	var expiresAt, lastUsedAt, revokedAt *string
	if ak.ExpiresAt != nil {
		s := ak.ExpiresAt.Format(time.RFC3339)
		expiresAt = &s
	}
	if ak.LastUsedAt != nil {
		s := ak.LastUsedAt.Format(time.RFC3339)
		lastUsedAt = &s
	}
	if ak.RevokedAt != nil {
		s := ak.RevokedAt.Format(time.RFC3339)
		revokedAt = &s
	}

	return &ApiKeyDTO{
		ID:           ak.ID.String(),
		ProjectID:    ak.ProjectID.String(),
		Name:         ak.Name,
		KeyPrefix:    ak.KeyPrefix,
		Mode:         ak.Mode,
		CreatedBy:    ak.CreatedBy.String(),
		ExpiresAt:    expiresAt,
		LastUsedAt:   lastUsedAt,
		IsRevoked:    ak.IsRevoked(),
		RevokedAt:    revokedAt,
		RequestLimit: ak.RequestLimit,
		RateLimit:    ak.RateLimit,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}
}

func FromDomainApiKeyCreated(ak *domain.ApiKey, rawKey string) *ApiKeyCreatedDTO {
	dto := FromDomainApiKey(ak)
	if dto == nil {
		return nil
	}
	return &ApiKeyCreatedDTO{
		ApiKeyDTO: *dto,
		RawKey:    rawKey,
	}
}
