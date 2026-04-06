package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProjectRole string

const (
	ProjectRoleAdmin  ProjectRole = "project_admin"
	ProjectRoleMember ProjectRole = "project_member"
	ProjectRoleViewer ProjectRole = "project_viewer"
)

func (r ProjectRole) String() string {
	return string(r)
}

// --- ProjectMember Domain Model ---
type ProjectMember struct {
	ID        uuid.UUID
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Email     string      // Hydrated from IAM User
	FullName  string      // Hydrated from IAM User
	Role      ProjectRole // e.g., "project_admin", "project_viewer"
	IsActive  bool
	AddedBy   uuid.UUID
	JoinedAt  time.Time
	RemovedAt *time.Time
	RemovedBy *uuid.UUID
}

// --- Constructor ---
func NewProjectMember(projectID, userID, addedBy uuid.UUID, role ProjectRole) *ProjectMember {
	return &ProjectMember{
		ID:        uuid.New(),
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
		IsActive:  true,
		AddedBy:   addedBy,
		JoinedAt:  time.Now(),
	}
}

// --- Smart Getters ---
func (m *ProjectMember) GetID() uuid.UUID    { return m.ID }
func (m *ProjectMember) GetIDString() string { return m.ID.String() }

func (m *ProjectMember) GetProjectID() uuid.UUID    { return m.ProjectID }
func (m *ProjectMember) GetProjectIDString() string { return m.ProjectID.String() }

func (m *ProjectMember) GetUserID() uuid.UUID    { return m.UserID }
func (m *ProjectMember) GetUserIDString() string { return m.UserID.String() }

func (m *ProjectMember) GetEmail() string    { return m.Email }
func (m *ProjectMember) GetFullName() string { return m.FullName }
func (m *ProjectMember) GetRole() string     { return m.Role.String() }
func (m *ProjectMember) GetIsActive() bool   { return m.IsActive }

// Safe Pointer & Time Getters
func (m *ProjectMember) GetJoinedAt() time.Time { return m.JoinedAt }
func (m *ProjectMember) GetJoinedAtString() string {
	if m.JoinedAt.IsZero() {
		return ""
	}
	return m.JoinedAt.Format(time.RFC3339)
}

func (m *ProjectMember) GetRemovedAt() *time.Time { return m.RemovedAt }
func (m *ProjectMember) GetRemovedAtString() string {
	if m.RemovedAt == nil || m.RemovedAt.IsZero() {
		return ""
	}
	return m.RemovedAt.Format(time.RFC3339)
}

func (m *ProjectMember) GetRemovedBy() *uuid.UUID { return m.RemovedBy }
func (m *ProjectMember) GetRemovedByString() string {
	if m.RemovedBy == nil || *m.RemovedBy == uuid.Nil {
		return ""
	}
	return m.RemovedBy.String()
}
