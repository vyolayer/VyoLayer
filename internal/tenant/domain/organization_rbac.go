package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// --- Roles ---
type OrganizationRole struct {
	ID             uuid.UUID
	Name           string
	Description    string
	IsSystem       bool
	IsDefault      bool
	HierarchyLevel uint32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewOrganizationRole(name, description string, isSystem, isDefault bool, hierarchyLevel uint32) *OrganizationRole {
	now := time.Now()
	return &OrganizationRole{
		ID:             uuid.New(),
		Name:           name,
		Description:    description,
		IsSystem:       isSystem,
		IsDefault:      isDefault,
		HierarchyLevel: hierarchyLevel,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// --- Smart Getters ---
func (r *OrganizationRole) GetID() uuid.UUID    { return r.ID }
func (r *OrganizationRole) GetIDString() string { return r.ID.String() }

func (r *OrganizationRole) GetName() string           { return r.Name }
func (r *OrganizationRole) GetDescription() string    { return r.Description }
func (r *OrganizationRole) GetIsSystem() bool         { return r.IsSystem }
func (r *OrganizationRole) GetIsDefault() bool        { return r.IsDefault }
func (r *OrganizationRole) GetHierarchyLevel() uint32 { return r.HierarchyLevel }

func (r *OrganizationRole) GetCreatedAt() time.Time { return r.CreatedAt }
func (r *OrganizationRole) GetCreatedAtString() string {
	if r.CreatedAt.IsZero() {
		return ""
	}
	return r.CreatedAt.Format(time.RFC3339)
}

func (r *OrganizationRole) GetUpdatedAt() time.Time { return r.UpdatedAt }
func (r *OrganizationRole) GetUpdatedAtString() string {
	if r.UpdatedAt.IsZero() {
		return ""
	}
	return r.UpdatedAt.Format(time.RFC3339)
}

// --- Setters ---
func (r *OrganizationRole) SetID(id uuid.UUID)             { r.ID = id }
func (r *OrganizationRole) SetName(name string)            { r.Name = name }
func (r *OrganizationRole) SetDescription(desc string)     { r.Description = desc }
func (r *OrganizationRole) SetIsSystem(isSystem bool)      { r.IsSystem = isSystem }
func (r *OrganizationRole) SetIsDefault(isDefault bool)    { r.IsDefault = isDefault }
func (r *OrganizationRole) SetHierarchyLevel(level uint32) { r.HierarchyLevel = level }
func (r *OrganizationRole) SetCreatedAt(t time.Time)       { r.CreatedAt = t }
func (r *OrganizationRole) SetUpdatedAt(t time.Time)       { r.UpdatedAt = t }

// --- Permissions ---
type OrganizationPermission struct {
	ID          uuid.UUID
	Resource    string
	Action      string
	Code        string // e.g., "organization:update" (Matches Protobuf)
	Group       string
	Description string
	IsSystem    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// --- Constructor ---
func NewOrganizationPermission(resource, action, group, description string, isSystem bool) *OrganizationPermission {
	now := time.Now()
	return &OrganizationPermission{
		ID:          uuid.New(),
		Resource:    resource,
		Action:      action,
		Code:        fmt.Sprintf("%s:%s", resource, action), // Auto-generate the gRPC code
		Group:       group,
		Description: description,
		IsSystem:    isSystem,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// --- Smart Getters ---
func (p *OrganizationPermission) GetID() uuid.UUID    { return p.ID }
func (p *OrganizationPermission) GetIDString() string { return p.ID.String() }

func (p *OrganizationPermission) GetResource() string { return p.Resource }
func (p *OrganizationPermission) GetAction() string   { return p.Action }
func (p *OrganizationPermission) GetCode() string {
	if p.Code == "" {
		p.Code = fmt.Sprintf("%s.%s", p.Resource, p.Action)
	}
	return p.Code
}
func (p *OrganizationPermission) GetGroup() string       { return p.Group }
func (p *OrganizationPermission) GetDescription() string { return p.Description }
func (p *OrganizationPermission) GetIsSystem() bool      { return p.IsSystem }

func (p *OrganizationPermission) GetCreatedAt() time.Time { return p.CreatedAt }
func (p *OrganizationPermission) GetCreatedAtString() string {
	if p.CreatedAt.IsZero() {
		return ""
	}
	return p.CreatedAt.Format(time.RFC3339)
}

func (p *OrganizationPermission) GetUpdatedAt() time.Time { return p.UpdatedAt }
func (p *OrganizationPermission) GetUpdatedAtString() string {
	if p.UpdatedAt.IsZero() {
		return ""
	}
	return p.UpdatedAt.Format(time.RFC3339)
}

// --- Setters ---
func (p *OrganizationPermission) SetID(id uuid.UUID)          { p.ID = id }
func (p *OrganizationPermission) SetResource(resource string) { p.Resource = resource }
func (p *OrganizationPermission) SetAction(action string)     { p.Action = action }
func (p *OrganizationPermission) SetCode(code string)         { p.Code = code }
func (p *OrganizationPermission) SetGroup(group string)       { p.Group = group }
func (p *OrganizationPermission) SetDescription(desc string)  { p.Description = desc }
func (p *OrganizationPermission) SetIsSystem(isSystem bool)   { p.IsSystem = isSystem }
func (p *OrganizationPermission) SetCreatedAt(t time.Time)    { p.CreatedAt = t }
func (p *OrganizationPermission) SetUpdatedAt(t time.Time)    { p.UpdatedAt = t }

// --- Role to Permission Mapping ---
type OrganizationRolePermission struct {
	ID           uuid.UUID
	RoleID       uuid.UUID
	PermissionID uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// --- Constructor ---
func NewOrganizationRolePermission(roleID, permissionID uuid.UUID) *OrganizationRolePermission {
	now := time.Now()
	return &OrganizationRolePermission{
		ID:           uuid.New(),
		RoleID:       roleID,
		PermissionID: permissionID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// --- Smart Getters ---
func (rp *OrganizationRolePermission) GetID() uuid.UUID    { return rp.ID }
func (rp *OrganizationRolePermission) GetIDString() string { return rp.ID.String() }

func (rp *OrganizationRolePermission) GetRoleID() uuid.UUID    { return rp.RoleID }
func (rp *OrganizationRolePermission) GetRoleIDString() string { return rp.RoleID.String() }

func (rp *OrganizationRolePermission) GetPermissionID() uuid.UUID    { return rp.PermissionID }
func (rp *OrganizationRolePermission) GetPermissionIDString() string { return rp.PermissionID.String() }

func (rp *OrganizationRolePermission) GetCreatedAt() time.Time { return rp.CreatedAt }
func (rp *OrganizationRolePermission) GetCreatedAtString() string {
	if rp.CreatedAt.IsZero() {
		return ""
	}
	return rp.CreatedAt.Format(time.RFC3339)
}

func (rp *OrganizationRolePermission) GetUpdatedAt() time.Time { return rp.UpdatedAt }
func (rp *OrganizationRolePermission) GetUpdatedAtString() string {
	if rp.UpdatedAt.IsZero() {
		return ""
	}
	return rp.UpdatedAt.Format(time.RFC3339)
}

// --- Setters ---
func (rp *OrganizationRolePermission) SetID(id uuid.UUID)               { rp.ID = id }
func (rp *OrganizationRolePermission) SetRoleID(roleID uuid.UUID)       { rp.RoleID = roleID }
func (rp *OrganizationRolePermission) SetPermissionID(permID uuid.UUID) { rp.PermissionID = permID }
func (rp *OrganizationRolePermission) SetCreatedAt(t time.Time)         { rp.CreatedAt = t }
func (rp *OrganizationRolePermission) SetUpdatedAt(t time.Time)         { rp.UpdatedAt = t }

// --- Member to Role Mapping ---
type MemberOrganizationRole struct {
	ID             uuid.UUID
	MemberID       uuid.UUID
	OrganizationID uuid.UUID
	RoleID         uuid.UUID
	GrantedBy      *uuid.UUID
	GrantedAt      time.Time
	RevokedBy      *uuid.UUID
	RevokedAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// --- Constructor ---
func NewMemberOrganizationRole(orgID, memberID, roleID, grantedBy uuid.UUID) *MemberOrganizationRole {
	now := time.Now()
	return &MemberOrganizationRole{
		ID:             uuid.New(),
		OrganizationID: orgID,
		MemberID:       memberID,
		RoleID:         roleID,
		GrantedBy:      &grantedBy,
		GrantedAt:      now, // Initialize the pointer safely
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// --- Smart Getters ---
func (mr *MemberOrganizationRole) GetID() uuid.UUID    { return mr.ID }
func (mr *MemberOrganizationRole) GetIDString() string { return mr.ID.String() }

func (mr *MemberOrganizationRole) GetMemberID() uuid.UUID    { return mr.MemberID }
func (mr *MemberOrganizationRole) GetMemberIDString() string { return mr.MemberID.String() }

func (mr *MemberOrganizationRole) GetOrganizationID() uuid.UUID    { return mr.OrganizationID }
func (mr *MemberOrganizationRole) GetOrganizationIDString() string { return mr.OrganizationID.String() }

func (mr *MemberOrganizationRole) GetRoleID() uuid.UUID    { return mr.RoleID }
func (mr *MemberOrganizationRole) GetRoleIDString() string { return mr.RoleID.String() }

// Safe Null/Empty checks for UUIDs
func (mr *MemberOrganizationRole) GetGrantedBy() uuid.UUID { return *mr.GrantedBy }
func (mr *MemberOrganizationRole) GetGrantedByString() string {
	if mr.GrantedBy == nil {
		return ""
	}
	return mr.GrantedBy.String()
}

func (mr *MemberOrganizationRole) GetRevokedBy() uuid.UUID { return *mr.RevokedBy }
func (mr *MemberOrganizationRole) GetRevokedByString() string {
	if mr.RevokedBy == nil {
		return ""
	}
	return mr.RevokedBy.String()
}

// Safe Nullable Pointer Time Getters
func (mr *MemberOrganizationRole) GetGrantedAt() *time.Time { return &mr.GrantedAt }
func (mr *MemberOrganizationRole) GetGrantedAtString() string {
	if mr.GrantedAt.IsZero() {
		return ""
	}
	return mr.GrantedAt.Format(time.RFC3339)
}

func (mr *MemberOrganizationRole) GetRevokedAt() *time.Time { return mr.RevokedAt }
func (mr *MemberOrganizationRole) GetRevokedAtString() string {
	if mr.RevokedAt == nil || mr.RevokedAt.IsZero() {
		return ""
	}
	return mr.RevokedAt.Format(time.RFC3339)
}

func (mr *MemberOrganizationRole) GetCreatedAt() time.Time { return mr.CreatedAt }
func (mr *MemberOrganizationRole) GetCreatedAtString() string {
	if mr.CreatedAt.IsZero() {
		return ""
	}
	return mr.CreatedAt.Format(time.RFC3339)
}

func (mr *MemberOrganizationRole) GetUpdatedAt() time.Time { return mr.UpdatedAt }
func (mr *MemberOrganizationRole) GetUpdatedAtString() string {
	if mr.UpdatedAt.IsZero() {
		return ""
	}
	return mr.UpdatedAt.Format(time.RFC3339)
}

// --- Setters ---
func (mr *MemberOrganizationRole) SetID(id uuid.UUID)                { mr.ID = id }
func (mr *MemberOrganizationRole) SetMemberID(memberID uuid.UUID)    { mr.MemberID = memberID }
func (mr *MemberOrganizationRole) SetOrganizationID(orgID uuid.UUID) { mr.OrganizationID = orgID }
func (mr *MemberOrganizationRole) SetRoleID(roleID uuid.UUID)        { mr.RoleID = roleID }
func (mr *MemberOrganizationRole) SetGrantedBy(grantedBy uuid.UUID)  { mr.GrantedBy = &grantedBy }
func (mr *MemberOrganizationRole) SetGrantedAt(t *time.Time)         { mr.GrantedAt = *t }
func (mr *MemberOrganizationRole) SetRevokedBy(revokedBy uuid.UUID)  { mr.RevokedBy = &revokedBy }
func (mr *MemberOrganizationRole) SetRevokedAt(t *time.Time)         { mr.RevokedAt = t }
func (mr *MemberOrganizationRole) SetCreatedAt(t time.Time)          { mr.CreatedAt = t }
func (mr *MemberOrganizationRole) SetUpdatedAt(t time.Time)          { mr.UpdatedAt = t }

// IsActive is a convenient business logic helper
func (mr *MemberOrganizationRole) IsActive() bool {
	return mr.RevokedAt == nil
}
