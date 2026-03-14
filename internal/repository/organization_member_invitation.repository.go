package repository

import (
	"context"
	"time"

	"github.com/vyolayer/vyolayer/internal/domain"
	"github.com/vyolayer/vyolayer/internal/platform/database/mapper"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"gorm.io/gorm"
)

type OrganizationMemberInvitationRepository interface {
	Create(
		ctx context.Context,
		invitation *domain.OrganizationMemberInvitation,
	) *errors.AppError

	GetByID(
		ctx context.Context,
		invitationID InvitationID,
	) (*domain.OrganizationMemberInvitation, *errors.AppError)

	GetByToken(
		ctx context.Context,
		token string,
	) (*domain.OrganizationMemberInvitation, *errors.AppError)

	GetByOrgID(
		ctx context.Context,
		orgID OrgID,
	) ([]domain.OrganizationMemberInvitation, *errors.AppError)

	GetPendingByEmail(
		ctx context.Context,
		email string,
	) ([]domain.OrganizationMemberInvitation, *errors.AppError)

	Update(
		ctx context.Context,
		invitation *domain.OrganizationMemberInvitation,
	) *errors.AppError

	Delete(
		ctx context.Context,
		invitationID InvitationID,
		deletedBy UserID,
	) *errors.AppError

	ExistsByEmailAndOrg(
		ctx context.Context,
		email string,
		orgID OrgID,
	) (bool, *errors.AppError)
}

type organizationMemberInvitationRepository struct {
	db *gorm.DB
}

func NewOrganizationMemberInvitationRepository(db *gorm.DB) OrganizationMemberInvitationRepository {
	return &organizationMemberInvitationRepository{db: db}
}

// Create creates a new invitation
func (r *organizationMemberInvitationRepository) Create(
	ctx context.Context,
	invitation *domain.OrganizationMemberInvitation,
) *errors.AppError {
	model := mapper.ToModelOrganizationMemberInvitation(invitation)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return ConvertDBError(err, "creating invitation")
	}

	return nil
}

// GetByID gets an invitation by ID
func (r *organizationMemberInvitationRepository) GetByID(
	ctx context.Context,
	invitationID InvitationID,
) (*domain.OrganizationMemberInvitation, *errors.AppError) {
	var invitation TOrganizationMemberInvitation
	err := r.db.WithContext(ctx).
		Where("id = ?", invitationID.InternalID().String()).
		First(&invitation).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.InvitationNotFound(invitationID.String())
		}
		return nil, ConvertDBError(err, "getting invitation by ID")
	}

	return mapper.ToDomainOrganizationMemberInvitation(&invitation), nil
}

// GetByToken gets an invitation by token
func (r *organizationMemberInvitationRepository) GetByToken(
	ctx context.Context,
	token string,
) (*domain.OrganizationMemberInvitation, *errors.AppError) {
	var invitation TOrganizationMemberInvitation
	err := r.db.WithContext(ctx).
		Where("token = ?", token).
		First(&invitation).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.InvitationNotFound("token")
		}
		return nil, ConvertDBError(err, "getting invitation by token")
	}

	return mapper.ToDomainOrganizationMemberInvitation(&invitation), nil
}

// GetByOrgID gets all invitations for an organization
func (r *organizationMemberInvitationRepository) GetByOrgID(
	ctx context.Context,
	orgID OrgID,
) ([]domain.OrganizationMemberInvitation, *errors.AppError) {
	var invitations []TOrganizationMemberInvitation
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID.InternalID().String()).
		Order("invited_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, ConvertDBError(err, "getting invitations by organization ID")
	}

	domainInvitations := make([]domain.OrganizationMemberInvitation, 0, len(invitations))
	for _, inv := range invitations {
		domainInv := mapper.ToDomainOrganizationMemberInvitation(&inv)
		if domainInv != nil {
			domainInvitations = append(domainInvitations, *domainInv)
		}
	}

	return domainInvitations, nil
}

// GetPendingByEmail gets all pending invitations for a user's email
func (r *organizationMemberInvitationRepository) GetPendingByEmail(
	ctx context.Context,
	email string,
) ([]domain.OrganizationMemberInvitation, *errors.AppError) {
	var invitations []TOrganizationMemberInvitation
	err := r.db.WithContext(ctx).
		Where("email = ? AND is_accepted = ? AND expired_at > ? AND deleted_at IS NULL AND deleted_by IS NULL", email, false, time.Now()).
		Order("invited_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, ConvertDBError(err, "getting pending invitations by email")
	}

	domainInvitations := make([]domain.OrganizationMemberInvitation, len(invitations))
	for i, inv := range invitations {
		domainInv := mapper.ToDomainOrganizationMemberInvitation(&inv)
		if domainInv != nil {
			domainInvitations[i] = *domainInv
		}
	}

	return domainInvitations, nil
}

// Update updates an invitation
func (r *organizationMemberInvitationRepository) Update(
	ctx context.Context,
	invitation *domain.OrganizationMemberInvitation,
) *errors.AppError {
	model := mapper.ToModelOrganizationMemberInvitation(invitation)

	err := r.db.WithContext(ctx).
		Model(&TOrganizationMemberInvitation{}).
		Where("id = ?", model.ID).
		Updates(model).Error

	if err != nil {
		return ConvertDBError(err, "updating invitation")
	}

	return nil
}

// Delete soft deletes an invitation
func (r *organizationMemberInvitationRepository) Delete(
	ctx context.Context,
	invitationID InvitationID,
	deletedBy UserID,
) *errors.AppError {
	deletedByUUID := deletedBy.InternalID().ID()

	err := r.db.WithContext(ctx).
		Model(&TOrganizationMemberInvitation{}).
		Where("id = ?", invitationID.InternalID().String()).
		Updates(map[string]interface{}{
			"deleted_by": deletedByUUID,
			"deleted_at": time.Now(),
		}).Error

	if err != nil {
		return ConvertDBError(err, "deleting invitation")
	}

	return nil
}

// ExistsByEmailAndOrg checks if a pending invitation already exists for the email and organization
func (r *organizationMemberInvitationRepository) ExistsByEmailAndOrg(
	ctx context.Context,
	email string,
	orgID OrgID,
) (bool, *errors.AppError) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&TOrganizationMemberInvitation{}).
		Where("email = ? AND organization_id = ? AND is_accepted = ? AND expired_at > ? AND deleted_at IS NULL",
			email, orgID.InternalID().String(), false, time.Now()).
		Count(&count).Error

	if err != nil {
		return false, ConvertDBError(err, "checking invitation existence")
	}

	return count > 0, nil
}
