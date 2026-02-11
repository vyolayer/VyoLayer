package dto

import (
	"time"
	"worklayer/internal/domain"
)

// CreateOrganizationRequestDTO represents the request to create an organization
type CreateOrganizationRequestDTO struct {
	Name        string `json:"name" validate:"required,min=3,max=100" example:"Acme Corp"`
	Description string `json:"description" validate:"max=500" example:"Building amazing products"`
}

// OrganizationMemberDTO represents a member of an organization
type OrganizationMemberDTO struct {
	ID            string  `json:"id" example:"orgmem_550e8400-e29b-41d4-a716-446655440000"`
	UserID        string  `json:"userId" example:"user_550e8400-e29b-41d4-a716-446655440000"`
	Email         string  `json:"email" example:"john@example.com"`
	FullName      string  `json:"fullName" example:"John Doe"`
	IsActive      bool    `json:"isActive" example:"true"`
	JoinedAt      string  `json:"joinedAt" example:"2023-01-01T00:00:00Z"`
	InvitedBy     *string `json:"invitedBy,omitempty" example:"orgmem_550e8400-e29b-41d4-a716-446655440001"`
	InvitedAt     *string `json:"invitedAt,omitempty" example:"2023-01-01T00:00:00Z"`
	DeactivatedBy *string `json:"deactivatedBy,omitempty"`
	DeactivatedAt *string `json:"deactivatedAt,omitempty"`
}

// OrganizationDTO represents an organization without members list
type OrganizationDTO struct {
	ID            string  `json:"id" example:"org_550e8400-e29b-41d4-a716-446655440000"`
	Name          string  `json:"name" example:"Acme Corp"`
	Slug          string  `json:"slug" example:"acme-corp"`
	Description   string  `json:"description" example:"Building amazing products"`
	OwnerID       string  `json:"ownerId" example:"user_550e8400-e29b-41d4-a716-446655440000"`
	IsActive      bool    `json:"isActive" example:"true"`
	MaxProjects   int     `json:"maxProjects" example:"5"`
	MaxMembers    int     `json:"maxMembers" example:"10"`
	MemberCount   int     `json:"memberCount" example:"3"`
	DeactivatedBy *string `json:"deactivatedBy,omitempty"`
	DeactivatedAt *string `json:"deactivatedAt,omitempty"`
	CreatedAt     string  `json:"createdAt" example:"2023-01-01T00:00:00Z"`
}

// OrganizationResponseDTO represents the full organization response with members
type OrganizationResponseDTO struct {
	OrganizationDTO
	Members []OrganizationMemberDTO `json:"members"`
}

// FromDomainOrganization converts a domain organization to DTO
func FromDomainOrganization(org *domain.Organization) *OrganizationDTO {
	if org == nil {
		return nil
	}

	var deactivatedBy *string
	if org.DeactivatedBy != nil {
		str := (*org.DeactivatedBy).String()
		deactivatedBy = &str
	}

	var deactivatedAt *string
	if org.DeactivatedAt != nil {
		str := org.DeactivatedAt.Format(time.RFC3339)
		deactivatedAt = &str
	}

	return &OrganizationDTO{
		ID:            org.ID.String(),
		Name:          org.Name,
		Slug:          org.Slug,
		Description:   org.Description,
		OwnerID:       org.OwnerID.String(),
		IsActive:      org.IsActive,
		MaxProjects:   org.MaxProjects,
		MaxMembers:    org.MemberInfo.MaxNoOfMembers,
		MemberCount:   org.MemberInfo.NoOfMembers,
		DeactivatedBy: deactivatedBy,
		DeactivatedAt: deactivatedAt,
		CreatedAt:     time.Now().Format(time.RFC3339), // Will be set properly when retrieved from DB
	}
}

// FromDomainOrganizationWithMembers converts a domain organization with members to response DTO
func FromDomainOrganizationWithMembers(org *domain.Organization) OrganizationResponseDTO {
	orgDTO := FromDomainOrganization(org)

	members, _ := org.GetMembers()
	memberDTOs := make([]OrganizationMemberDTO, 0, len(members))

	for _, member := range members {
		memberDTOs = append(memberDTOs, FromDomainOrganizationMember(&member))
	}

	return OrganizationResponseDTO{
		OrganizationDTO: *orgDTO,
		Members:         memberDTOs,
	}
}

// FromDomainOrganizationMember converts a domain organization member to DTO
func FromDomainOrganizationMember(member *domain.OrganizationMember) OrganizationMemberDTO {
	if member == nil {
		return OrganizationMemberDTO{}
	}

	var invitedBy *string
	if member.InvitedBy != nil {
		str := (*member.InvitedBy).String()
		invitedBy = &str
	}

	var invitedAt *string
	if member.InvitedAt != nil {
		str := member.InvitedAt.Format(time.RFC3339)
		invitedAt = &str
	}

	var deactivatedBy *string
	if member.DeactivatedBy != nil {
		str := (*member.DeactivatedBy).String()
		deactivatedBy = &str
	}

	var deactivatedAt *string
	if member.DeactivatedAt != nil {
		str := member.DeactivatedAt.Format(time.RFC3339)
		deactivatedAt = &str
	}

	return OrganizationMemberDTO{
		ID:            member.ID.String(),
		UserID:        member.UserID.String(),
		Email:         member.Email,
		FullName:      member.FullName,
		IsActive:      member.IsActive,
		JoinedAt:      member.JoinedAt.Format(time.RFC3339),
		InvitedBy:     invitedBy,
		InvitedAt:     invitedAt,
		DeactivatedBy: deactivatedBy,
		DeactivatedAt: deactivatedAt,
	}
}
