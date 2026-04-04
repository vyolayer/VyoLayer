package domain

import (
	"time"

	"github.com/google/uuid"
)

// --- Organization Member Struct ---
type OrganizationMember struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	FullName       string
	Email          string
	IsActive       bool
	JoinedAt       *time.Time
	RemovedAt      *time.Time
	RemovedBy      *uuid.UUID
	InvitedBy      *uuid.UUID
	InvitedAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// --- Organization Member With Roles Struct ---
type OrganizationMemberWithRoles struct {
	OrganizationMember
	Roles []OrganizationRole
}

// --- Organization Member With Roles And Permissions Struct ---
type OrganizationMemberWithRolesAndPermissions struct {
	OrganizationMember
	Roles       []OrganizationRole
	Permissions []OrganizationPermission
}

// --- Constructor ---
// NewOrganizationMember creates a new active member with initialized timestamps and a new UUID.
func NewOrganizationMember(orgID, userID uuid.UUID) *OrganizationMember {
	now := time.Now()
	return &OrganizationMember{
		ID:             uuid.New(),
		OrganizationID: orgID,
		UserID:         userID,
		IsActive:       true,
		JoinedAt:       &now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
func (m *OrganizationMember) GetID() uuid.UUID    { return m.ID }
func (m *OrganizationMember) GetIDString() string { return m.ID.String() }

func (m *OrganizationMember) GetOrganizationID() uuid.UUID    { return m.OrganizationID }
func (m *OrganizationMember) GetOrganizationIDString() string { return m.OrganizationID.String() }

func (m *OrganizationMember) GetUserID() uuid.UUID    { return m.UserID }
func (m *OrganizationMember) GetUserIDString() string { return m.UserID.String() }

// Basic Type Getters
func (m *OrganizationMember) GetFullName() string { return m.FullName }
func (m *OrganizationMember) GetEmail() string    { return m.Email }
func (m *OrganizationMember) GetIsActive() bool   { return m.IsActive }

// Status Helper (Computed for Frontend/gRPC convenience)
func (m *OrganizationMember) GetStatus() string {
	if m.RemovedAt != nil {
		return "removed"
	}
	if !m.IsActive && m.InvitedAt != nil {
		return "pending"
	}
	if m.IsActive {
		return "active"
	}
	return "unknown"
}

// Guaranteed Time Getters
func (m *OrganizationMember) GetJoinedAt() time.Time {
	if m.JoinedAt == nil {
		return time.Time{}
	}
	return *m.JoinedAt
}
func (m *OrganizationMember) GetJoinedAtString() string {
	if m.JoinedAt == nil || m.JoinedAt.IsZero() {
		return ""
	}
	return m.JoinedAt.Format(time.RFC3339)
}

func (m *OrganizationMember) GetCreatedAt() time.Time { return m.CreatedAt }
func (m *OrganizationMember) GetCreatedAtString() string {
	if m.CreatedAt.IsZero() {
		return ""
	}
	return m.CreatedAt.Format(time.RFC3339)
}

func (m *OrganizationMember) GetUpdatedAt() time.Time { return m.UpdatedAt }
func (m *OrganizationMember) GetUpdatedAtString() string {
	if m.UpdatedAt.IsZero() {
		return ""
	}
	return m.UpdatedAt.Format(time.RFC3339)
}

// Nullable / Optional Pointer Getters
func (m *OrganizationMember) GetRemovedBy() *uuid.UUID { return m.RemovedBy }
func (m *OrganizationMember) GetRemovedByString() string {
	if m.RemovedBy == nil || *m.RemovedBy == uuid.Nil {
		return ""
	}
	return m.RemovedBy.String()
}

func (m *OrganizationMember) GetRemovedAt() *time.Time { return m.RemovedAt }
func (m *OrganizationMember) GetRemovedAtString() string {
	if m.RemovedAt == nil || m.RemovedAt.IsZero() {
		return ""
	}
	return m.RemovedAt.Format(time.RFC3339)
}

func (m *OrganizationMember) GetInvitedBy() *uuid.UUID { return m.InvitedBy }
func (m *OrganizationMember) GetInvitedByString() string {
	if m.InvitedBy == nil || *m.InvitedBy == uuid.Nil {
		return ""
	}
	return m.InvitedBy.String()
}

func (m *OrganizationMember) GetInvitedAt() *time.Time { return m.InvitedAt }
func (m *OrganizationMember) GetInvitedAtString() string {
	if m.InvitedAt == nil || m.InvitedAt.IsZero() {
		return ""
	}
	return m.InvitedAt.Format(time.RFC3339)
}

// --- Setters ---

func (m *OrganizationMember) SetID(id uuid.UUID)                { m.ID = id }
func (m *OrganizationMember) SetOrganizationID(orgID uuid.UUID) { m.OrganizationID = orgID }
func (m *OrganizationMember) SetUserID(userID uuid.UUID)        { m.UserID = userID }
func (m *OrganizationMember) SetFullName(name string)           { m.FullName = name }
func (m *OrganizationMember) SetEmail(email string)             { m.Email = email }
func (m *OrganizationMember) SetIsActive(isActive bool)         { m.IsActive = isActive }
func (m *OrganizationMember) SetJoinedAt(joinedAt *time.Time)   { m.JoinedAt = joinedAt }
func (m *OrganizationMember) SetRemovedAt(removedAt *time.Time) { m.RemovedAt = removedAt }
func (m *OrganizationMember) SetRemovedBy(removedBy *uuid.UUID) { m.RemovedBy = removedBy }
func (m *OrganizationMember) SetInvitedBy(invitedBy *uuid.UUID) { m.InvitedBy = invitedBy }
func (m *OrganizationMember) SetInvitedAt(invitedAt *time.Time) { m.InvitedAt = invitedAt }
func (m *OrganizationMember) SetCreatedAt(createdAt time.Time)  { m.CreatedAt = createdAt }
func (m *OrganizationMember) SetUpdatedAt(updatedAt time.Time)  { m.UpdatedAt = updatedAt }
