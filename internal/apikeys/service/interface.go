package service

import (
	"context"

	"github.com/google/uuid"
	apikeymodelv1 "github.com/vyolayer/vyolayer/internal/apikeys/models/v1"
)

type Service interface {
	Create(
		ctx context.Context,
		organizationID uuid.UUID,
		projectID uuid.UUID,
		actorID uuid.UUID,
		name string,
		description string,
		environment string,
		scopes []string,
	) (*apikeymodelv1.APIKey, string, error)

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
	) error

	Rotate(
		ctx context.Context,
		id uuid.UUID,
		organizationID uuid.UUID,
		projectID uuid.UUID,
		actorID uuid.UUID,
	) (*apikeymodelv1.APIKey, string, error)

	Validate(
		ctx context.Context,
		secret string,
	) (bool, *apikeymodelv1.APIKey, []string, error)
}
