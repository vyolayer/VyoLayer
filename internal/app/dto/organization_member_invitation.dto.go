package dto

import (
	"time"
	"worklayer/internal/domain"
)

// CreateInvitationRequestDTO represents the request to create an invitation
type CreateInvitationRequestDTO struct {
	Email   string   `json:"email" validate:"required,email" example:"john@example.com"`
	RoleIDs []string `json:"roleIds" validate:"omitempty,dive,min=1" example:"[\"org_role_123\"]"`
}

// AcceptInvitationRequestDTO represents the request to accept an invitation (optional, can use URL token)
type AcceptInvitationRequestDTO struct {
	Token string `json:"token,omitempty" validate:"omitempty,min=1" example:"abc123token"`
}

// OrganizationMemberInvitationDTO represents an invitation response
type OrganizationMemberInvitationDTO struct {
	ID             string   `json:"id" example:"org_invitation_550e8400-e29b-41d4-a716-446655440000"`
	OrganizationID string   `json:"organizationId" example:"org_550e8400-e29b-41d4-a716-446655440000"`
	InvitedBy      string   `json:"invitedBy" example:"org_member_550e8400-e29b-41d4-a716-446655440000"`
	Email          string   `json:"email" example:"john@example.com"`
	RoleIDs        []string `json:"roleIds" example:"[\"org_role_123\"]"`
	InvitedAt      string   `json:"invitedAt" example:"2023-01-01T00:00:00Z"`
	IsAccepted     bool     `json:"isAccepted" example:"false"`
	AcceptedAt     *string  `json:"acceptedAt,omitempty" example:"2023-01-01T00:00:00Z"`
	ExpiredAt      string   `json:"expiredAt" example:"2023-01-08T00:00:00Z"`
	IsPending      bool     `json:"isPending" example:"true"`
}

// FromDomainOrganizationMemberInvitation converts a domain invitation to DTO
func FromDomainOrganizationMemberInvitation(inv *domain.OrganizationMemberInvitation) OrganizationMemberInvitationDTO {
	if inv == nil {
		return OrganizationMemberInvitationDTO{}
	}

	var acceptedAt *string
	if inv.AcceptedAt != nil {
		str := inv.AcceptedAt.Format(time.RFC3339)
		acceptedAt = &str
	}

	return OrganizationMemberInvitationDTO{
		ID:             inv.ID.String(),
		OrganizationID: inv.OrganizationID.String(),
		InvitedBy:      inv.InvitedBy.String(),
		Email:          inv.Email,
		RoleIDs:        inv.RoleIDs,
		InvitedAt:      inv.InvitedAt.Format(time.RFC3339),
		IsAccepted:     inv.IsAccepted,
		AcceptedAt:     acceptedAt,
		ExpiredAt:      inv.ExpiredAt.Format(time.RFC3339),
		IsPending:      inv.IsPending(),
	}
}
