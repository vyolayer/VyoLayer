package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	apikeymodelv1 "github.com/vyolayer/vyolayer/internal/apikeys/models/v1"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(
	ctx context.Context,
	item *apikeymodelv1.APIKey,
	scopes []string,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(item).Error; err != nil {
			return err
		}

		for _, scope := range scopes {
			row := apikeymodelv1.APIKeyScope{
				ApiKeyID: item.ID,
				Scope:    scope,
			}

			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *repository) List(
	ctx context.Context,
	organizationID uuid.UUID,
	projectID uuid.UUID,
) ([]apikeymodelv1.APIKey, error) {

	var items []apikeymodelv1.APIKey

	err := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Where("project_id = ?", projectID).
		Order("created_at desc").
		Find(&items).Error

	return items, err
}

func (r *repository) Get(
	ctx context.Context,
	id uuid.UUID,
	organizationID uuid.UUID,
	projectID uuid.UUID,
) (*apikeymodelv1.APIKey, error) {

	var item apikeymodelv1.APIKey

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Where("organization_id = ?", organizationID).
		Where("project_id = ?", projectID).
		First(&item).Error

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *repository) Update(
	ctx context.Context,
	id uuid.UUID,
	organizationID uuid.UUID,
	projectID uuid.UUID,
	name string,
	description string,
	scopes []string,
) (*apikeymodelv1.APIKey, error) {

	var item apikeymodelv1.APIKey

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("id = ?", id).
			Where("organization_id = ?", organizationID).
			Where("project_id = ?", projectID).
			First(&item).Error; err != nil {
			return err
		}

		item.Name = name
		item.Description = description
		item.UpdatedAt = time.Now()

		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		if err := tx.
			Where("api_key_id = ?", item.ID).
			Delete(&apikeymodelv1.APIKeyScope{}).Error; err != nil {
			return err
		}

		for _, scope := range scopes {
			row := apikeymodelv1.APIKeyScope{
				ApiKeyID: item.ID,
				Scope:    scope,
			}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *repository) Revoke(
	ctx context.Context,
	id uuid.UUID,
	organizationID uuid.UUID,
	projectID uuid.UUID,
	actorID uuid.UUID,
	revokedAt time.Time,
) error {

	return r.db.WithContext(ctx).
		Model(&apikeymodelv1.APIKey{}).
		Where("id = ?", id).
		Where("organization_id = ?", organizationID).
		Where("project_id = ?", projectID).
		Updates(map[string]any{
			"status":     apikeymodelv1.APIKeyStatusRevoked,
			"revoked_by": actorID,
			"revoked_at": revokedAt,
			"updated_at": time.Now(),
		}).Error
}

func (r *repository) FindByPrefix(
	ctx context.Context,
	prefix string,
) (*apikeymodelv1.APIKey, []string, error) {

	var item apikeymodelv1.APIKey

	err := r.db.WithContext(ctx).
		Where("prefix = ?", prefix).
		First(&item).Error

	if err != nil {
		return nil, nil, err
	}

	var rows []apikeymodelv1.APIKeyScope

	err = r.db.WithContext(ctx).
		Where("api_key_id = ?", item.ID).
		Find(&rows).Error

	if err != nil {
		return nil, nil, err
	}

	scopes := make([]string, 0, len(rows))
	for _, row := range rows {
		scopes = append(scopes, row.Scope)
	}

	return &item, scopes, nil
}

func (r *repository) TouchUsage(
	ctx context.Context,
	id uuid.UUID,
) error {
	return r.db.WithContext(ctx).
		Model(&apikeymodelv1.APIKey{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"last_used_at": time.Now(),
			"updated_at":   time.Now(),
		}).Error
}
