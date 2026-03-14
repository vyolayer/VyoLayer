package domain

import (
	"time"

	"github.com/vyolayer/vyolayer/internal/platform/database/types"
	"github.com/vyolayer/vyolayer/pkg/errors"
)

// API Key mode constants
const (
	ApiKeyModeLive = "live"
	ApiKeyModeDev  = "dev"
)

// Default limits per mode
const (
	DevRequestLimitPerDay  = 1000
	DevRateLimitPerMinute  = 60
	LiveRequestLimitPerDay = 0 // 0 = unlimited
	LiveRateLimitPerMinute = 600
)

// ValidApiKeyModes is the set of allowed mode values
var ValidApiKeyModes = map[string]bool{
	ApiKeyModeLive: true,
	ApiKeyModeDev:  true,
}

// IsValidApiKeyMode checks whether a mode string is valid
func IsValidApiKeyMode(mode string) bool {
	return ValidApiKeyModes[mode]
}

// DefaultLimitsForMode returns the default request and rate limits for a given mode
func DefaultLimitsForMode(mode string) (requestLimit int, rateLimit int) {
	switch mode {
	case ApiKeyModeLive:
		return LiveRequestLimitPerDay, LiveRateLimitPerMinute
	default:
		return DevRequestLimitPerDay, DevRateLimitPerMinute
	}
}

// ---------------------------------------------------------------------------
// ApiKey
// ---------------------------------------------------------------------------

type ApiKey struct {
	ID             types.ApiKeyID
	ProjectID      types.ProjectID
	OrganizationID types.OrganizationID

	Name      string
	KeyPrefix string // visible prefix for identification (e.g. "wl_live_ab3f")
	KeyHash   string // SHA-256 of the full key, stored in DB

	Mode      string // "live" or "dev"
	CreatedBy types.UserID

	ExpiresAt  *time.Time
	LastUsedAt *time.Time
	RevokedAt  *time.Time
	RevokedBy  *types.UserID

	RequestLimit int // max requests per day (0 = unlimited)
	RateLimit    int // max requests per minute
}

func NewApiKey(
	projectID types.ProjectID,
	organizationID types.OrganizationID,
	name string,
	keyPrefix string,
	keyHash string,
	mode string,
	createdBy types.UserID,
	expiresAt *time.Time,
) *ApiKey {
	requestLimit, rateLimit := DefaultLimitsForMode(mode)

	return &ApiKey{
		ID:             types.NewApiKeyID(),
		ProjectID:      projectID,
		OrganizationID: organizationID,
		Name:           name,
		KeyPrefix:      keyPrefix,
		KeyHash:        keyHash,
		Mode:           mode,
		CreatedBy:      createdBy,
		ExpiresAt:      expiresAt,
		RequestLimit:   requestLimit,
		RateLimit:      rateLimit,
	}
}

func ReconstructApiKey(
	id types.ApiKeyID,
	projectID types.ProjectID,
	organizationID types.OrganizationID,
	name, keyPrefix, keyHash, mode string,
	createdBy types.UserID,
	expiresAt, lastUsedAt, revokedAt *time.Time,
	revokedBy *types.UserID,
	requestLimit, rateLimit int,
) *ApiKey {
	return &ApiKey{
		ID:             id,
		ProjectID:      projectID,
		OrganizationID: organizationID,
		Name:           name,
		KeyPrefix:      keyPrefix,
		KeyHash:        keyHash,
		Mode:           mode,
		CreatedBy:      createdBy,
		ExpiresAt:      expiresAt,
		LastUsedAt:     lastUsedAt,
		RevokedAt:      revokedAt,
		RevokedBy:      revokedBy,
		RequestLimit:   requestLimit,
		RateLimit:      rateLimit,
	}
}

func (ak *ApiKey) IsRevoked() bool {
	return ak.RevokedAt != nil
}

func (ak *ApiKey) IsExpired() bool {
	if ak.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ak.ExpiresAt)
}

func (ak *ApiKey) IsUsable() *errors.AppError {
	if ak.IsRevoked() {
		return ApiKeyRevokedError(ak.ID.String())
	}
	if ak.IsExpired() {
		return ApiKeyExpiredError(ak.ID.String())
	}
	return nil
}

func (ak *ApiKey) Revoke(revokedBy types.UserID) *errors.AppError {
	if ak.IsRevoked() {
		return ApiKeyRevokedError(ak.ID.String())
	}
	now := time.Now()
	ak.RevokedAt = &now
	ak.RevokedBy = &revokedBy
	return nil
}

func (ak *ApiKey) Validate() *errors.AppError {
	if ak.Name == "" {
		return ValidationError("API key name is required")
	}
	if !IsValidApiKeyMode(ak.Mode) {
		return ValidationError("API key mode must be 'dev' or 'live'")
	}
	return nil
}

// ---------------------------------------------------------------------------
// ApiKeyUsageLog
// ---------------------------------------------------------------------------

type ApiKeyUsageLog struct {
	ApiKeyID   types.ApiKeyID
	ProjectID  types.ProjectID
	Endpoint   string
	Method     string
	StatusCode int
	IPAddress  string
	UserAgent  string
}
