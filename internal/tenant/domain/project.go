package domain

import (
	"time"

	"github.com/google/uuid"
)

// --- Project Domain Model ---
type Project struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Slug           string
	Description    string
	IsActive       bool
	CreatedBy      uuid.UUID
	MaxAPIKeys     uint8
	MaxMembers     uint8
	MemberCount    uint32 // Computed field
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ArchivedAt     *time.Time
}

// --- Constructor ---
func NewProject(orgID, createdBy uuid.UUID, name, slug, description string) *Project {
	now := time.Now()
	return &Project{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		Slug:           slug,
		Description:    description,
		IsActive:       true,
		CreatedBy:      createdBy,
		MaxAPIKeys:     5, // Default
		MaxMembers:     5, // Default
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// --- Smart Getters ---
func (p *Project) GetID() uuid.UUID    { return p.ID }
func (p *Project) GetIDString() string { return p.ID.String() }

func (p *Project) GetOrganizationID() uuid.UUID    { return p.OrganizationID }
func (p *Project) GetOrganizationIDString() string { return p.OrganizationID.String() }

func (p *Project) GetName() string        { return p.Name }
func (p *Project) GetSlug() string        { return p.Slug }
func (p *Project) GetDescription() string { return p.Description }
func (p *Project) GetIsActive() bool      { return p.IsActive }

func (p *Project) GetCreatedBy() uuid.UUID    { return p.CreatedBy }
func (p *Project) GetCreatedByString() string { return p.CreatedBy.String() }

// Safe integer upcasting for Protobuf (uint8 -> uint32)
func (p *Project) GetMaxAPIKeys() uint32  { return uint32(p.MaxAPIKeys) }
func (p *Project) GetMaxMembers() uint32  { return uint32(p.MaxMembers) }
func (p *Project) GetMemberCount() uint32 { return p.MemberCount }

// Safe Time Getters (RFC3339)
func (p *Project) GetCreatedAt() time.Time { return p.CreatedAt }
func (p *Project) GetCreatedAtString() string {
	if p.CreatedAt.IsZero() {
		return ""
	}
	return p.CreatedAt.Format(time.RFC3339)
}

func (p *Project) GetUpdatedAt() time.Time { return p.UpdatedAt }
func (p *Project) GetUpdatedAtString() string {
	if p.UpdatedAt.IsZero() {
		return ""
	}
	return p.UpdatedAt.Format(time.RFC3339)
}

func (p *Project) GetArchivedAt() *time.Time { return p.ArchivedAt }
func (p *Project) GetArchivedAtString() string {
	if p.ArchivedAt == nil || p.ArchivedAt.IsZero() {
		return ""
	}
	return p.ArchivedAt.Format(time.RFC3339)
}

// --- Setters ---
// (Standard setters omitted for brevity, implement as p.Name = name)
