package mapper

import (
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/models"
	"worklayer/internal/platform/database/types"
)

// ToDomainApiKey converts an ApiKey model to a domain ApiKey.
func ToDomainApiKey(m *models.ApiKey) *domain.ApiKey {
	if m == nil {
		return nil
	}

	akID, _ := types.ReconstructApiKeyID(m.ID.String())
	projectID, _ := types.ReconstructProjectID(m.ProjectID.String())
	orgID, _ := types.ReconstructOrganizationID(m.OrganizationID.String())
	createdBy, _ := types.ReconstructUserID(m.CreatedBy.String())

	var revokedBy *types.UserID
	if m.IsRevoked() {
		id, _ := types.ReconstructUserID(m.RevokedBy.String())
		revokedBy = &id
	}

	return domain.ReconstructApiKey(
		akID,
		projectID,
		orgID,
		m.Name,
		m.KeyPrefix,
		m.KeyHash,
		m.Mode,
		createdBy,
		m.ExpiresAt,
		m.LastUsedAt,
		m.RevokedAt,
		revokedBy,
		m.RequestLimit,
		m.RateLimit,
	)
}
