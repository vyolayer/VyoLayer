package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationMemberInvitation struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	InvitedBy      uuid.UUID
	Email          string
	Token          string
	RoleIDs        []uuid.UUID
	InvitedAt      time.Time
	IsAccepted     bool
	AcceptedAt     *time.Time
	ExpiredAt      time.Time
	DeletedBy      *uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Inviter struct {
	MemberID uuid.UUID
	FullName string
	Email    string
}

type OrganizationMemberInvitationWithInviter struct {
	OrganizationMemberInvitation
	Inviter Inviter
}

// --- Constructor ---

// NewOrganizationMemberInvitation creates a new pending invitation.
// It automatically sets the expiration time (e.g., 7 days from now).
func NewOrganizationMemberInvitation(orgID, invitedBy uuid.UUID, email, token string, roleIDs []uuid.UUID, expirationDuration time.Duration) *OrganizationMemberInvitation {
	now := time.Now()

	// Ensure roleIDs is never nil to avoid panics later
	if roleIDs == nil {
		roleIDs = make([]uuid.UUID, 0)
	}

	return &OrganizationMemberInvitation{
		ID:             uuid.New(),
		OrganizationID: orgID,
		InvitedBy:      invitedBy,
		Email:          email,
		Token:          token,
		RoleIDs:        roleIDs,
		InvitedAt:      now,
		IsAccepted:     false,
		ExpiredAt:      now.Add(expirationDuration),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// --- Business Logic Helpers ---

func (i *OrganizationMemberInvitation) IsExpired() bool {
	return time.Now().After(i.ExpiredAt)
}

func (i *OrganizationMemberInvitation) IsPending() bool {
	return !i.IsAccepted && !i.IsExpired() && i.DeletedBy == nil
}

// --- Smart Getters ---

// UUID Getters
func (i *OrganizationMemberInvitation) GetID() uuid.UUID    { return i.ID }
func (i *OrganizationMemberInvitation) GetIDString() string { return i.ID.String() }

func (i *OrganizationMemberInvitation) GetOrganizationID() uuid.UUID { return i.OrganizationID }
func (i *OrganizationMemberInvitation) GetOrganizationIDString() string {
	return i.OrganizationID.String()
}

func (i *OrganizationMemberInvitation) GetInvitedBy() uuid.UUID    { return i.InvitedBy }
func (i *OrganizationMemberInvitation) GetInvitedByString() string { return i.InvitedBy.String() }

// Basic Type Getters
func (i *OrganizationMemberInvitation) GetEmail() string        { return i.Email }
func (i *OrganizationMemberInvitation) GetToken() string        { return i.Token }
func (i *OrganizationMemberInvitation) GetIsAccepted() bool     { return i.IsAccepted }
func (i *OrganizationMemberInvitation) GetRoleIDs() []uuid.UUID { return i.RoleIDs }

// Helper for Protobuf `repeated string role_ids`
func (i *OrganizationMemberInvitation) GetRoleIDStrings() []string {
	if len(i.RoleIDs) == 0 {
		return []string{}
	}
	strs := make([]string, len(i.RoleIDs))
	for idx, roleID := range i.RoleIDs {
		strs[idx] = roleID.String()
	}
	return strs
}

// Guaranteed Time Getters
func (i *OrganizationMemberInvitation) GetInvitedAt() time.Time { return i.InvitedAt }
func (i *OrganizationMemberInvitation) GetInvitedAtString() string {
	if i.InvitedAt.IsZero() {
		return ""
	}
	return i.InvitedAt.Format(time.RFC3339)
}

func (i *OrganizationMemberInvitation) GetExpiredAt() time.Time { return i.ExpiredAt }
func (i *OrganizationMemberInvitation) GetExpiredAtString() string {
	if i.ExpiredAt.IsZero() {
		return ""
	}
	return i.ExpiredAt.Format(time.RFC3339)
}

func (i *OrganizationMemberInvitation) GetCreatedAt() time.Time { return i.CreatedAt }
func (i *OrganizationMemberInvitation) GetCreatedAtString() string {
	if i.CreatedAt.IsZero() {
		return ""
	}
	return i.CreatedAt.Format(time.RFC3339)
}

func (i *OrganizationMemberInvitation) GetUpdatedAt() time.Time { return i.UpdatedAt }
func (i *OrganizationMemberInvitation) GetUpdatedAtString() string {
	if i.UpdatedAt.IsZero() {
		return ""
	}
	return i.UpdatedAt.Format(time.RFC3339)
}

// Nullable / Optional Pointer Getters
func (i *OrganizationMemberInvitation) GetAcceptedAt() *time.Time { return i.AcceptedAt }
func (i *OrganizationMemberInvitation) GetAcceptedAtString() string {
	if i.AcceptedAt == nil || i.AcceptedAt.IsZero() {
		return ""
	}
	return i.AcceptedAt.Format(time.RFC3339)
}

func (i *OrganizationMemberInvitation) GetDeletedBy() *uuid.UUID { return i.DeletedBy }
func (i *OrganizationMemberInvitation) GetDeletedByString() string {
	if i.DeletedBy == nil || *i.DeletedBy == uuid.Nil {
		return ""
	}
	return i.DeletedBy.String()
}

// --- Setters ---

func (i *OrganizationMemberInvitation) SetID(id uuid.UUID)                { i.ID = id }
func (i *OrganizationMemberInvitation) SetOrganizationID(orgID uuid.UUID) { i.OrganizationID = orgID }
func (i *OrganizationMemberInvitation) SetInvitedBy(invitedBy uuid.UUID)  { i.InvitedBy = invitedBy }
func (i *OrganizationMemberInvitation) SetEmail(email string)             { i.Email = email }
func (i *OrganizationMemberInvitation) SetToken(token string)             { i.Token = token }
func (i *OrganizationMemberInvitation) SetRoleIDs(roleIDs []uuid.UUID)    { i.RoleIDs = roleIDs }
func (i *OrganizationMemberInvitation) SetInvitedAt(invitedAt time.Time)  { i.InvitedAt = invitedAt }
func (i *OrganizationMemberInvitation) SetIsAccepted(isAccepted bool)     { i.IsAccepted = isAccepted }
func (i *OrganizationMemberInvitation) SetAcceptedAt(acceptedAt *time.Time) {
	i.AcceptedAt = acceptedAt
}
func (i *OrganizationMemberInvitation) SetExpiredAt(expiredAt time.Time)  { i.ExpiredAt = expiredAt }
func (i *OrganizationMemberInvitation) SetDeletedBy(deletedBy *uuid.UUID) { i.DeletedBy = deletedBy }
func (i *OrganizationMemberInvitation) SetCreatedAt(createdAt time.Time)  { i.CreatedAt = createdAt }
func (i *OrganizationMemberInvitation) SetUpdatedAt(updatedAt time.Time)  { i.UpdatedAt = updatedAt }

func NewOrganizationMemberInvitationWithInviter(invitation OrganizationMemberInvitation, inviter Inviter) OrganizationMemberInvitationWithInviter {
	return OrganizationMemberInvitationWithInviter{
		OrganizationMemberInvitation: invitation,
		Inviter:                      inviter,
	}
}

func (i *OrganizationMemberInvitationWithInviter) GetInviter() Inviter {
	return i.Inviter
}
