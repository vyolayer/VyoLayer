package domain

import (
	"time"
	"worklayer/internal/platform/database/types"
	"worklayer/pkg/errors"
)

type OrganizationMember struct {
	// Info
	ID             types.OrganizationMemberID
	OrganizationID types.OrganizationID
	UserID         types.UserID

	// User Info
	Email    string
	FullName string
	IsActive bool

	// Status
	JoinedAt      time.Time
	InvitedBy     *types.OrganizationMemberID
	InvitedAt     *time.Time
	DeactivatedBy *types.OrganizationMemberID
	DeactivatedAt *time.Time
}

func NewOrganizationMember(
	organizationID types.OrganizationID,
	invitedBy *types.OrganizationMemberID,
	invitedAt *time.Time,
	user *User,
) *OrganizationMember {
	id := types.NewOrganizationMemberID()
	return &OrganizationMember{
		ID:             id,
		OrganizationID: organizationID,
		UserID:         user.ID,
		Email:          user.Email,
		FullName:       user.FullName,
		IsActive:       true,
		JoinedAt:       time.Now(),
		InvitedBy:      invitedBy,
		InvitedAt:      invitedAt,
	}
}

// ReconstructOrganizationMember reconstructs an organization member from database data
func ReconstructOrganizationMember(
	id types.OrganizationMemberID,
	organizationID types.OrganizationID,
	userID types.UserID,
	email, fullName string,
	isActive bool,
	joinedAt time.Time,
	invitedBy *types.OrganizationMemberID,
	invitedAt *time.Time,
	deactivatedBy *types.OrganizationMemberID,
	deactivatedAt *time.Time,
) *OrganizationMember {
	return &OrganizationMember{
		ID:             id,
		OrganizationID: organizationID,
		UserID:         userID,
		Email:          email,
		FullName:       fullName,
		IsActive:       isActive,
		JoinedAt:       joinedAt,
		InvitedBy:      invitedBy,
		InvitedAt:      invitedAt,
		DeactivatedBy:  deactivatedBy,
		DeactivatedAt:  deactivatedAt,
	}
}

// Deactivate deactivates the organization member
func (om *OrganizationMember) Deactivate(deactivatedBy types.OrganizationMemberID) *errors.AppError {
	if !om.IsActive {
		return OrganizationMemberNotActiveError()
	}

	now := time.Now()
	om.IsActive = false
	om.DeactivatedBy = &deactivatedBy
	om.DeactivatedAt = &now

	return nil
}

// Reactivate reactivates the organization member
func (om *OrganizationMember) Reactivate() *errors.AppError {
	if om.IsActive {
		return nil // Already active
	}

	om.IsActive = true
	om.DeactivatedBy = nil
	om.DeactivatedAt = nil

	return nil
}

// UpdateUserInfo updates the cached user information
func (om *OrganizationMember) UpdateUserInfo(email, fullName string) {
	om.Email = email
	om.FullName = fullName
}

// Validate validates the organization member
func (om *OrganizationMember) Validate() *errors.AppError {
	if om.Email == "" {
		return ValidationError("Member email is required")
	}

	if om.FullName == "" {
		return ValidationError("Member full name is required")
	}

	return nil
}
