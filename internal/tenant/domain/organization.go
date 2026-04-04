package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/pkg/utils"
)

const (
	defaultMaxMembers  = 10
	defaultMaxProjects = 1
)

// --- Organization Struct ---
type Organization struct {
	ID            uuid.UUID
	Name          string
	Slug          string
	Description   string
	IsActive      bool
	OwnerID       uuid.UUID
	MaxMembers    uint32
	MaxProjects   uint32
	ProjectCount  uint32
	MemberCount   uint32
	DeactivatedBy *uuid.UUID
	DeactivatedAt *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ArchivedAt    *time.Time
}

type OrganizationWithMember struct {
	Organization
	Members []OrganizationMember
}

// --- Constructor ---
func NewOrganization(ownerID uuid.UUID, name, description string) *Organization {
	id := uuid.New()
	now := time.Now()
	slugify := utils.ToSlug(name).Slugify()

	return &Organization{
		ID:            id,
		Name:          name,
		Slug:          slugify.String(),
		Description:   description,
		IsActive:      true,
		OwnerID:       ownerID,
		MaxMembers:    defaultMaxMembers,
		MaxProjects:   defaultMaxProjects,
		ProjectCount:  0,
		MemberCount:   1,
		DeactivatedBy: nil,
		DeactivatedAt: nil,
		CreatedAt:     now,
		UpdatedAt:     now,
		ArchivedAt:    nil,
	}
}

// --- Getters ---
func (o *Organization) GetID() uuid.UUID    { return o.ID }
func (o *Organization) GetIDString() string { return o.ID.String() }

func (o *Organization) GetName() string        { return o.Name }
func (o *Organization) GetSlug() string        { return o.Slug }
func (o *Organization) GetDescription() string { return o.Description }
func (o *Organization) GetIsActive() bool      { return o.IsActive }

func (o *Organization) GetOwnerID() uuid.UUID    { return o.OwnerID }
func (o *Organization) GetOwnerIDString() string { return o.OwnerID.String() }

func (o *Organization) GetMaxMembers() uint32   { return o.MaxMembers }
func (o *Organization) GetMaxProjects() uint32  { return o.MaxProjects }
func (o *Organization) GetProjectCount() uint32 { return o.ProjectCount }
func (o *Organization) GetMemberCount() uint32  { return o.MemberCount }

// Safe Pointer Getters
func (o *Organization) GetDeactivatedBy() *uuid.UUID { return o.DeactivatedBy }
func (o *Organization) GetDeactivatedByString() string {
	if o.DeactivatedBy == nil || *o.DeactivatedBy == uuid.Nil {
		return ""
	}
	return o.DeactivatedBy.String()
}

// Safe Time Getters (Formats to RFC3339 for Protobuf)
func (o *Organization) GetCreatedAt() time.Time { return o.CreatedAt }
func (o *Organization) GetCreatedAtString() string {
	if o.CreatedAt.IsZero() {
		return ""
	}
	return o.CreatedAt.Format(time.RFC3339)
}

func (o *Organization) GetUpdatedAt() time.Time { return o.UpdatedAt }
func (o *Organization) GetUpdatedAtString() string {
	if o.UpdatedAt.IsZero() {
		return ""
	}
	return o.UpdatedAt.Format(time.RFC3339)
}

// Safe Nullable Time Getters
func (o *Organization) GetDeactivatedAt() *time.Time { return o.DeactivatedAt }
func (o *Organization) GetDeactivatedAtString() string {
	if o.DeactivatedAt == nil || o.DeactivatedAt.IsZero() {
		return ""
	}
	return o.DeactivatedAt.Format(time.RFC3339)
}

func (o *Organization) GetArchivedAt() *time.Time { return o.ArchivedAt }
func (o *Organization) GetArchivedAtString() string {
	if o.ArchivedAt == nil || o.ArchivedAt.IsZero() {
		return ""
	}
	return o.ArchivedAt.Format(time.RFC3339)
}

// --- Setters ---

func (o *Organization) SetID(id uuid.UUID)                        { o.ID = id }
func (o *Organization) SetName(name string)                       { o.Name = name }
func (o *Organization) SetSlug(slug string)                       { o.Slug = slug }
func (o *Organization) SetDescription(desc string)                { o.Description = desc }
func (o *Organization) SetIsActive(isActive bool)                 { o.IsActive = isActive }
func (o *Organization) SetOwnerID(ownerID uuid.UUID)              { o.OwnerID = ownerID }
func (o *Organization) SetMaxMembers(max uint32)                  { o.MaxMembers = max }
func (o *Organization) SetMaxProjects(max uint32)                 { o.MaxProjects = max }
func (o *Organization) SetProjectCount(count uint32)              { o.ProjectCount = count }
func (o *Organization) SetMemberCount(count uint32)               { o.MemberCount = count }
func (o *Organization) SetDeactivatedBy(deactivatedBy *uuid.UUID) { o.DeactivatedBy = deactivatedBy }
func (o *Organization) SetDeactivatedAt(deactivatedAt *time.Time) { o.DeactivatedAt = deactivatedAt }
func (o *Organization) SetCreatedAt(createdAt time.Time)          { o.CreatedAt = createdAt }
func (o *Organization) SetUpdatedAt(updatedAt time.Time)          { o.UpdatedAt = updatedAt }
func (o *Organization) SetArchivedAt(archivedAt *time.Time)       { o.ArchivedAt = archivedAt }
