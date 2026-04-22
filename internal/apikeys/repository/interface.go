package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	apikeymodelv1 "github.com/vyolayer/vyolayer/internal/apikeys/models/v1"
)

type Repository interface {
	Create(
		ctx context.Context,
		item *apikeymodelv1.APIKey,
		scopes []string,
	) error

	List(
		ctx context.Context,
		organizationID uuid.UUID,
		projectID uuid.UUID,
	) ([]apikeymodelv1.APIKey, error)

	Get(
		ctx context.Context,
		id uuid.UUID,
		organizationID uuid.UUID,
		projectID uuid.UUID,
	) (*apikeymodelv1.APIKey, error)

	Update(
		ctx context.Context,
		id uuid.UUID,
		organizationID uuid.UUID,
		projectID uuid.UUID,
		name string,
		description string,
		scopes []string,
	) (*apikeymodelv1.APIKey, error)

	Revoke(
		ctx context.Context,
		id uuid.UUID,
		organizationID uuid.UUID,
		projectID uuid.UUID,
		actorID uuid.UUID,
		revokedAt time.Time,
	) error

	FindByPrefix(
		ctx context.Context,
		prefix string,
	) (*apikeymodelv1.APIKey, []string, error)

	TouchUsage(
		ctx context.Context,
		id uuid.UUID,
	) error

	CreateAuditLog(
		ctx context.Context,
		log *apikeymodelv1.APIKeyAuditLog,
	) error

	UpsertRateLimit(
		ctx context.Context,
		limit *apikeymodelv1.APIKeyRateLimit,
	) error
}
