package grpc

import (
	"time"

	model "github.com/vyolayer/vyolayer/internal/apikeys/models/v1"
	apikeyv1 "github.com/vyolayer/vyolayer/proto/apikey/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoAPIKey(m *model.APIKey, scopes []string) *apikeyv1.APIKey {
	if m == nil {
		return nil
	}

	out := &apikeyv1.APIKey{
		Id:             m.ID.String(),
		OrganizationId: m.OrganizationID.String(),
		ProjectId:      m.ProjectID.String(),
		Name:           m.Name,
		Description:    m.Description,
		Prefix:         m.Prefix,
		Environment:    m.Environment,
		Status:         m.Status,
		Scopes:         scopes,
		CreatedBy:      m.CreatedBy.String(),
		LastUsedIp:     m.LastUsedIP,
		LastUsedUa:     m.LastUsedUA,
		CreatedAt:      toProtoTS(m.CreatedAt),
		UpdatedAt:      toProtoTS(m.UpdatedAt),
	}

	if m.LastUsedAt != nil {
		out.LastUsedAt = toProtoTS(*m.LastUsedAt)
	}
	if m.ExpiresAt != nil {
		out.ExpiresAt = toProtoTS(*m.ExpiresAt)
	}
	if m.RevokedBy != nil {
		out.RevokedBy = m.RevokedBy.String()
	}
	if m.RevokedAt != nil {
		out.RevokedAt = toProtoTS(*m.RevokedAt)
	}

	return out
}

func toProtoTS(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

// func toProtoEnv(v string) apikeyv1.APIKeyEnvironment {
// 	switch v {
// 	case model.APIKeyModeLive:
// 		return apikeyv1.APIKeyEnvironment_API_KEY_ENVIRONMENT_LIVE
// 	default:
// 		return apikeyv1.APIKeyEnvironment_API_KEY_ENVIRONMENT_DEV
// 	}
// }

// func toProtoStatus(v string) apikeyv1.APIKeyStatus {
// 	switch v {
// 	case model.APIKeyStatusRevoked:
// 		return apikeyv1.APIKeyStatus_API_KEY_STATUS_REVOKED
// 	case model.APIKeyStatusExpired:
// 		return apikeyv1.APIKeyStatus_API_KEY_STATUS_EXPIRED
// 	case model.APIKeyStatusDisabled:
// 		return apikeyv1.APIKeyStatus_API_KEY_STATUS_DISABLED
// 	default:
// 		return apikeyv1.APIKeyStatus_API_KEY_STATUS_ACTIVE
// 	}
// }
