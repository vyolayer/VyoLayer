package dto

import (
	"time"

	"github.com/vyolayer/vyolayer/internal/domain"
)

// ── Request DTOs ─────────────────────────────────────────────────────────────

type CreateProjectRequestDTO struct {
	Name        string `json:"name" validate:"required,min=3,max=100" example:"My Project"`
	Description string `json:"description" validate:"max=500" example:"A cool project"`
}

type UpdateProjectRequestDTO struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=100" example:"My Project"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500" example:"Updated description"`
}

type DeleteProjectRequestDTO struct {
	ConfirmName string `json:"confirmName" validate:"required" example:"My Project"`
}

type AddProjectMemberRequestDTO struct {
	UserID string `json:"userId" validate:"required" example:"user_550e8400-e29b-41d4-a716-446655440000"`
	Role   string `json:"role" validate:"required,oneof=admin member viewer" example:"member"`
}

type ChangeProjectMemberRoleRequestDTO struct {
	Role string `json:"role" validate:"required,oneof=admin member viewer" example:"admin"`
}

// ── Response DTOs ────────────────────────────────────────────────────────────

type ProjectDTO struct {
	ID             string `json:"id" example:"project_550e8400-e29b-41d4-a716-446655440000"`
	OrganizationID string `json:"organizationId" example:"org_550e8400-e29b-41d4-a716-446655440000"`
	Name           string `json:"name" example:"My Project"`
	Slug           string `json:"slug" example:"my-project"`
	Description    string `json:"description" example:"A cool project"`
	IsActive       bool   `json:"isActive" example:"true"`
	CreatedBy      string `json:"createdBy" example:"user_550e8400-e29b-41d4-a716-446655440000"`
	MaxApiKeys     int    `json:"maxApiKeys" example:"5"`
	MaxMembers     int    `json:"maxMembers" example:"10"`
	MemberCount    int    `json:"memberCount" example:"3"`
	CreatedAt      string `json:"createdAt" example:"2023-01-01T00:00:00Z"`
}

type ProjectResponseDTO struct {
	ProjectDTO
	Members []ProjectMemberDTO `json:"members"`
}

type ProjectMemberDTO struct {
	ID        string  `json:"id" example:"project_member_550e8400-e29b-41d4-a716-446655440000"`
	UserID    string  `json:"userId" example:"user_550e8400-e29b-41d4-a716-446655440000"`
	Email     string  `json:"email" example:"john@example.com"`
	FullName  string  `json:"fullName" example:"John Doe"`
	Role      string  `json:"role" example:"admin"`
	IsActive  bool    `json:"isActive" example:"true"`
	JoinedAt  string  `json:"joinedAt" example:"2023-01-01T00:00:00Z"`
	RemovedAt *string `json:"removedAt,omitempty"`
}

// ── Domain-to-DTO converters ─────────────────────────────────────────────────

func FromDomainProject(p *domain.Project) *ProjectDTO {
	if p == nil {
		return nil
	}

	memberCount := 0
	maxMembers := 0
	if members, err := p.GetMembers(); err == nil {
		memberCount = len(members)
	}
	maxMembers = p.MemberInfo.MaxNoOfMembers

	return &ProjectDTO{
		ID:             p.ID.String(),
		OrganizationID: p.OrganizationID.String(),
		Name:           p.Name,
		Slug:           p.Slug,
		Description:    p.Description,
		IsActive:       p.IsActive,
		CreatedBy:      p.CreatedBy.String(),
		MaxApiKeys:     p.MaxApiKeys,
		MaxMembers:     maxMembers,
		MemberCount:    memberCount,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}
}

func FromDomainProjectWithMembers(p *domain.Project) *ProjectResponseDTO {
	dto := FromDomainProject(p)
	if dto == nil {
		return nil
	}

	members, _ := p.GetMembers()
	memberDTOs := make([]ProjectMemberDTO, 0, len(members))
	for _, m := range members {
		memberDTOs = append(memberDTOs, FromDomainProjectMember(&m))
	}

	return &ProjectResponseDTO{
		ProjectDTO: *dto,
		Members:    memberDTOs,
	}
}

func FromDomainProjectMember(m *domain.ProjectMember) ProjectMemberDTO {
	if m == nil {
		return ProjectMemberDTO{}
	}

	var removedAt *string
	if m.RemovedAt != nil {
		str := m.RemovedAt.Format(time.RFC3339)
		removedAt = &str
	}

	return ProjectMemberDTO{
		ID:        m.ID.String(),
		UserID:    m.UserID.String(),
		Email:     m.Email,
		FullName:  m.FullName,
		Role:      m.Role,
		IsActive:  m.IsActive,
		JoinedAt:  m.JoinedAt.Format(time.RFC3339),
		RemovedAt: removedAt,
	}
}
