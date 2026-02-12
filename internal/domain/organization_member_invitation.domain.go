package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"
	"worklayer/internal/platform/database/types"
	"worklayer/pkg/errors"
)

type OrganizationMemberInvitation struct {
	// Info
	ID             types.OrganizationMemberInvitationID
	OrganizationID types.OrganizationID
	InvitedBy      types.OrganizationMemberID
	Email          string
	Token          string
	RoleIDs        []string

	// Status
	InvitedAt  time.Time
	IsAccepted bool
	AcceptedAt *time.Time
	ExpiredAt  time.Time
	DeletedBy  *types.OrganizationMemberID
}

// NewOrganizationMemberInvitation creates a new organization member invitation
func NewOrganizationMemberInvitation(
	organizationID types.OrganizationID,
	invitedBy types.OrganizationMemberID,
	email string,
	roleIDs []string,
	expirationDays int,
) (*OrganizationMemberInvitation, *errors.AppError) {
	if email == "" {
		return nil, ValidationError("Email is required for invitation")
	}

	token, err := generateInvitationToken()
	if err != nil {
		return nil, errors.InternalWrap(err, "Failed to generate invitation token")
	}

	now := time.Now()
	expiredAt := now.AddDate(0, 0, expirationDays)

	id := types.NewOrganizationMemberInvitationID()

	return &OrganizationMemberInvitation{
		ID:             id,
		OrganizationID: organizationID,
		InvitedBy:      invitedBy,
		Email:          email,
		Token:          token,
		RoleIDs:        roleIDs,
		InvitedAt:      now,
		IsAccepted:     false,
		AcceptedAt:     nil,
		ExpiredAt:      expiredAt,
		DeletedBy:      nil,
	}, nil
}

// ReconstructOrganizationMemberInvitation reconstructs an invitation from database data
func ReconstructOrganizationMemberInvitation(
	id types.OrganizationMemberInvitationID,
	organizationID types.OrganizationID,
	invitedBy types.OrganizationMemberID,
	email, token string,
	roleIDs []string,
	invitedAt time.Time,
	isAccepted bool,
	acceptedAt *time.Time,
	expiredAt time.Time,
	deletedBy *types.OrganizationMemberID,
) *OrganizationMemberInvitation {
	return &OrganizationMemberInvitation{
		ID:             id,
		OrganizationID: organizationID,
		InvitedBy:      invitedBy,
		Email:          email,
		Token:          token,
		RoleIDs:        roleIDs,
		InvitedAt:      invitedAt,
		IsAccepted:     isAccepted,
		AcceptedAt:     acceptedAt,
		ExpiredAt:      expiredAt,
		DeletedBy:      deletedBy,
	}
}

// Accept marks the invitation as accepted
func (omi *OrganizationMemberInvitation) Accept() *errors.AppError {
	if omi.IsAccepted {
		return InvitationAlreadyAcceptedError(omi.ID.String())
	}

	if omi.IsExpired() {
		return InvitationExpiredError()
	}

	now := time.Now()
	omi.IsAccepted = true
	omi.AcceptedAt = &now

	return nil
}

// IsExpired checks if the invitation has expired
func (omi *OrganizationMemberInvitation) IsExpired() bool {
	return time.Now().After(omi.ExpiredAt)
}

// IsPending checks if the invitation is still pending
func (omi *OrganizationMemberInvitation) IsPending() bool {
	return !omi.IsAccepted && !omi.IsExpired()
}

// Validate validates the invitation
func (omi *OrganizationMemberInvitation) Validate() *errors.AppError {
	if omi.Email == "" {
		return ValidationError("Invitation email is required")
	}

	if omi.Token == "" {
		return ValidationError("Invitation token is required")
	}

	if omi.ExpiredAt.Before(omi.InvitedAt) {
		return ValidationError("Expiration date cannot be before invitation date")
	}

	return nil
}

// generateInvitationToken generates a secure random token for the invitation
func generateInvitationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
