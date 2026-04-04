package tenantrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ErrUnimplemented = status.Errorf(codes.Unimplemented, "method not implemented")
)

type organizationMemberInvitationRepo struct {
	gormRepo
}

func NewOrganizationMemberInvitationRepo(
	db *gorm.DB,
	logger *logger.AppLogger,
) OrganizationMemberInvitationRepository {
	return &organizationMemberInvitationRepo{
		gormRepo: gormRepo{
			db:     db,
			logger: logger,
		},
	}
}

// --- Write Implementation ---

func (r *organizationMemberInvitationRepo) Create(ctx context.Context, invitation *domain.OrganizationMemberInvitation) error {
	model := toOrganizationMemberInvitationModel(invitation)
	err := r.db.WithContext(ctx).
		Create(model).
		Error
	if err != nil {
		return ConvertDBError(err, "Failed to create organization member invitation")
	}
	return nil
}

// Accept marks the invitation as accepted and stamps it with the current time.
func (r *organizationMemberInvitationRepo) Accept(ctx context.Context, invitation *domain.OrganizationMemberInvitation) error {
	now := time.Now()
	err := r.db.WithContext(ctx).
		Model(&OrganizationMemberInvitation{}).
		Where("id = ?", invitation.ID).
		Updates(map[string]any{
			"is_accepted": true,
			"accepted_at": now,
			"updated_at":  now,
		}).Error
	if err != nil {
		return ConvertDBError(err, "Failed to accept invitation")
	}
	return nil
}

// Delete soft-deletes an invitation (sets deleted_at via GORM).
func (r *organizationMemberInvitationRepo) Delete(ctx context.Context, invitation *domain.OrganizationMemberInvitation) error {
	err := r.db.WithContext(ctx).
		Where("id = ? AND organization_id = ?", invitation.ID, invitation.OrganizationID).
		Delete(&OrganizationMemberInvitation{}).Error
	if err != nil {
		return ConvertDBError(err, "Failed to delete invitation")
	}
	return nil
}

// --- Read Implementation ---

func (r *organizationMemberInvitationRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.OrganizationMemberInvitation, error) {
	var model OrganizationMemberInvitation
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error
	if err != nil {
		return nil, ConvertDBError(err, "Failed to get invitation by id")
	}
	return toInvitationDomain(&model), nil
}

func (r *organizationMemberInvitationRepo) GetByToken(ctx context.Context, token string) (*domain.OrganizationMemberInvitation, error) {
	var model OrganizationMemberInvitation
	err := r.db.WithContext(ctx).
		Where("token = ?", token).
		First(&model).Error
	if err != nil {
		return nil, ConvertDBError(err, "Failed to get invitation by token")
	}
	return toInvitationDomain(&model), nil
}

// List returns all invitations (including accepted/expired/deleted) for an org.
func (r *organizationMemberInvitationRepo) List(ctx context.Context, organizationID uuid.UUID) ([]*domain.OrganizationMemberInvitation, error) {
	var models []OrganizationMemberInvitation
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Find(&models).Error
	if err != nil {
		return nil, ConvertDBError(err, "Failed to list invitations")
	}
	return toInvitationDomainSlice(models), nil
}

// ListByUserEmail returns all invitations sent to a particular email address.
func (r *organizationMemberInvitationRepo) ListByUserEmail(ctx context.Context, email string) ([]*domain.OrganizationMemberInvitation, error) {
	var models []OrganizationMemberInvitation
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		Find(&models).Error
	if err != nil {
		return nil, ConvertDBError(err, "Failed to list invitations by user email")
	}
	return toInvitationDomainSlice(models), nil
}

// ListPendingByOrg returns only active (not accepted, not expired, not deleted) invitations for an org.
func (r *organizationMemberInvitationRepo) ListPendingByOrg(ctx context.Context, organizationID uuid.UUID) ([]*domain.OrganizationMemberInvitationWithInviter, error) {
	var models []OrganizationMemberInvitation
	err := r.db.WithContext(ctx).
		Preload("Inviter").
		Preload("Inviter.User").
		Where("organization_id = ? AND is_accepted = false AND expired_at > ? AND deleted_at IS NULL", organizationID, time.Now()).
		Find(&models).Error
	if err != nil {
		return nil, ConvertDBError(err, "Failed to list pending invitations")
	}
	return toInvitationDomainWithInviterSlice(models), nil
}

// ListByInvitedBy returns all invitations created by a specific member within an org.
func (r *organizationMemberInvitationRepo) ListByInvitedBy(ctx context.Context, organizationID, invitedByUserID uuid.UUID) ([]*domain.OrganizationMemberInvitation, error) {
	var models []OrganizationMemberInvitation
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND invited_by = ?", organizationID, invitedByUserID).
		Find(&models).Error
	if err != nil {
		return nil, ConvertDBError(err, "Failed to list invitations by inviter")
	}
	return toInvitationDomainSlice(models), nil
}

// --- Helper ---

func toInvitationDomainSlice(models []OrganizationMemberInvitation) []*domain.OrganizationMemberInvitation {
	result := make([]*domain.OrganizationMemberInvitation, len(models))
	for i := range models {
		result[i] = toInvitationDomain(&models[i])
	}
	return result
}

func toInvitationDomainWithInviterSlice(models []OrganizationMemberInvitation) []*domain.OrganizationMemberInvitationWithInviter {
	result := make([]*domain.OrganizationMemberInvitationWithInviter, len(models))
	for i := range models {
		result[i] = toInvitationDomainWithInviter(&models[i])
	}
	return result
}
