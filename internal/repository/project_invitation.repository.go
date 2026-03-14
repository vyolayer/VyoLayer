package repository

import (
	"context"
	"vyolayer/internal/platform/database/types"
	"vyolayer/pkg/errors"

	"gorm.io/gorm"
)

type ProjectInvitationRepository interface {
	Create(ctx context.Context, invitation *TProjectInvitation) *errors.AppError
	FindByID(ctx context.Context, invitationID types.ProjectInvitationID) (*TProjectInvitation, *errors.AppError)
	FindByToken(ctx context.Context, token string) (*TProjectInvitation, *errors.AppError)
	FindPendingByProjectID(ctx context.Context, projectID types.ProjectID) ([]TProjectInvitation, *errors.AppError)
	FindPendingByEmail(ctx context.Context, email string) ([]TProjectInvitation, *errors.AppError)
	Accept(ctx context.Context, invitationID types.ProjectInvitationID) *errors.AppError
	Delete(ctx context.Context, invitationID types.ProjectInvitationID, deletedBy types.ProjectMemberID) *errors.AppError
}

type projectInvitationRepository struct {
	db *gorm.DB
}

func NewProjectInvitationRepository(db *gorm.DB) ProjectInvitationRepository {
	return &projectInvitationRepository{db: db}
}

func (r *projectInvitationRepository) Create(
	ctx context.Context,
	invitation *TProjectInvitation,
) *errors.AppError {
	err := r.db.WithContext(ctx).Create(invitation).Error
	if err != nil {
		return ConvertDBError(err, "creating project invitation")
	}
	return nil
}

func (r *projectInvitationRepository) FindByID(
	ctx context.Context,
	invitationID types.ProjectInvitationID,
) (*TProjectInvitation, *errors.AppError) {
	var invitation TProjectInvitation
	err := r.db.
		Where("id = ?", invitationID.InternalID().ID()).
		First(&invitation).Error
	if err != nil {
		return nil, ConvertDBError(err, "finding project invitation by ID")
	}
	return &invitation, nil
}

func (r *projectInvitationRepository) FindByToken(
	ctx context.Context,
	token string,
) (*TProjectInvitation, *errors.AppError) {
	var invitation TProjectInvitation
	err := r.db.
		Where("token = ? AND is_accepted = ? AND deleted_at IS NULL", token, false).
		First(&invitation).Error
	if err != nil {
		return nil, ConvertDBError(err, "finding project invitation by token")
	}
	return &invitation, nil
}

func (r *projectInvitationRepository) FindPendingByProjectID(
	ctx context.Context,
	projectID types.ProjectID,
) ([]TProjectInvitation, *errors.AppError) {
	var invitations []TProjectInvitation
	err := r.db.
		Where("project_id = ? AND is_accepted = ? AND deleted_at IS NULL",
			projectID.InternalID().ID(), false).
		Find(&invitations).Error
	if err != nil {
		return nil, ConvertDBError(err, "listing pending project invitations")
	}
	return invitations, nil
}

func (r *projectInvitationRepository) FindPendingByEmail(
	ctx context.Context,
	email string,
) ([]TProjectInvitation, *errors.AppError) {
	var invitations []TProjectInvitation
	err := r.db.
		Where("email = ? AND is_accepted = ? AND deleted_at IS NULL", email, false).
		Find(&invitations).Error
	if err != nil {
		return nil, ConvertDBError(err, "listing pending project invitations by email")
	}
	return invitations, nil
}

func (r *projectInvitationRepository) Accept(
	ctx context.Context,
	invitationID types.ProjectInvitationID,
) *errors.AppError {
	err := r.db.WithContext(ctx).
		Model(&TProjectInvitation{}).
		Where("id = ?", invitationID.InternalID().ID()).
		Updates(map[string]interface{}{
			"is_accepted": true,
			"accepted_at": gorm.Expr("NOW()"),
		}).Error
	if err != nil {
		return ConvertDBError(err, "accepting project invitation")
	}
	return nil
}

func (r *projectInvitationRepository) Delete(
	ctx context.Context,
	invitationID types.ProjectInvitationID,
	deletedBy types.ProjectMemberID,
) *errors.AppError {
	deletedByID := deletedBy.InternalID().ID()
	err := r.db.WithContext(ctx).
		Model(&TProjectInvitation{}).
		Where("id = ?", invitationID.InternalID().ID()).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("NOW()"),
			"deleted_by": deletedByID,
		}).Error
	if err != nil {
		return ConvertDBError(err, "deleting project invitation")
	}
	return nil
}
