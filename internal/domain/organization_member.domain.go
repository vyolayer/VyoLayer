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

// Organization member with roles
type OrganizationMemberWithRoles struct {
	OrganizationMember
	Roles []types.OrganizationRoleID
}

func NewOrganizationMemberWithRoles(
	organizationID types.OrganizationID,
	invitedBy *types.OrganizationMemberID,
	invitedAt *time.Time,
	user *User,
	roleIDs []types.OrganizationRoleID,
) *OrganizationMemberWithRoles {
	return &OrganizationMemberWithRoles{
		OrganizationMember: *NewOrganizationMember(
			organizationID,
			invitedBy,
			invitedAt,
			user,
		),
		Roles: roleIDs,
	}
}

func (om *OrganizationMemberWithRoles) AssignRoles(roleIDs []types.OrganizationRoleID) {
	om.Roles = roleIDs
}

func (om *OrganizationMemberWithRoles) RolesString() []string {
	roles := make([]string, len(om.Roles))
	for i, roleID := range om.Roles {
		roles[i] = roleID.InternalID().ID().String()
	}
	return roles
}

type OrganizationPermission struct {
	ID       string
	Resource string
	Action   string
	Group    string
	IsSystem bool
}

func (op OrganizationPermission) Code() string {
	return op.Resource + "." + op.Action
}

type OrganizationRole struct {
	ID           string
	Name         string
	Description  string
	IsSystemRole bool
	IsDefault    bool
}

// Organization member with rbac
type OrganizationMemberWithRBAC struct {
	OrganizationMember
	Roles       map[string]OrganizationRole
	Permissions map[string]OrganizationPermission
}
