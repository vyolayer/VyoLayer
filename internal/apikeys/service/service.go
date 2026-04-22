package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	model "github.com/vyolayer/vyolayer/internal/apikeys/models/v1"
	"github.com/vyolayer/vyolayer/internal/apikeys/repository"
	"github.com/vyolayer/vyolayer/internal/apikeys/util"
)

type service struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(
	ctx context.Context,
	organizationID uuid.UUID,
	projectID uuid.UUID,
	actorID uuid.UUID,
	name string,
	description string,
	environment string,
	scopes []string,
) (*model.APIKey, string, error) {

	gen, err := util.GenerateAPIKey(environment)
	if err != nil {
		return nil, "", err
	}

	now := time.Now()

	item := &model.APIKey{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		ProjectID:      projectID,
		Name:           name,
		Description:    description,
		Prefix:         gen.Prefix,
		SecretHash:     util.HashSecret(gen.Secret),
		Environment:    environment,
		Status:         model.APIKeyStatusActive,
		CreatedBy:      actorID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.Create(ctx, item, scopes); err != nil {
		return nil, "", err
	}

	return item, gen.Secret, nil
}

func (s *service) List(
	ctx context.Context,
	organizationID uuid.UUID,
	projectID uuid.UUID,
) ([]model.APIKey, error) {
	return s.repo.List(ctx, organizationID, projectID)
}

func (s *service) Get(
	ctx context.Context,
	id uuid.UUID,
	organizationID uuid.UUID,
	projectID uuid.UUID,
) (*model.APIKey, error) {
	return s.repo.Get(ctx, id, organizationID, projectID)
}

func (s *service) Update(
	ctx context.Context,
	id uuid.UUID,
	organizationID uuid.UUID,
	projectID uuid.UUID,
	name string,
	description string,
	scopes []string,
) (*model.APIKey, error) {
	return s.repo.Update(
		ctx,
		id,
		organizationID,
		projectID,
		name,
		description,
		scopes,
	)
}

func (s *service) Revoke(
	ctx context.Context,
	id uuid.UUID,
	organizationID uuid.UUID,
	projectID uuid.UUID,
	actorID uuid.UUID,
) error {
	return s.repo.Revoke(
		ctx,
		id,
		organizationID,
		projectID,
		actorID,
		time.Now(),
	)
}

func (s *service) Rotate(
	ctx context.Context,
	id uuid.UUID,
	organizationID uuid.UUID,
	projectID uuid.UUID,
	actorID uuid.UUID,
) (*model.APIKey, string, error) {

	oldKey, err := s.repo.Get(ctx, id, organizationID, projectID)
	if err != nil {
		return nil, "", err
	}

	if err := s.Revoke(ctx, id, organizationID, projectID, actorID); err != nil {
		return nil, "", err
	}

	return s.Create(
		ctx,
		organizationID,
		projectID,
		actorID,
		oldKey.Name,
		oldKey.Description,
		oldKey.Environment,
		[]string{},
	)
}

func (s *service) Validate(
	ctx context.Context,
	secret string,
) (bool, *model.APIKey, []string, error) {

	if err := util.ValidateAPIKeyFormat(secret); err != nil {
		return false, nil, nil, nil
	}

	prefix := util.ExtractPrefix(secret)

	item, scopes, err := s.repo.FindByPrefix(ctx, prefix)
	if err != nil {
		return false, nil, nil, err
	}

	if item == nil {
		return false, nil, nil, nil
	}

	if !item.IsUsable() {
		return false, nil, nil, nil
	}

	if !util.VerifySecret(secret, item.SecretHash) {
		return false, nil, nil, nil
	}

	_ = s.repo.TouchUsage(ctx, item.ID)

	return true, item, scopes, nil
}
