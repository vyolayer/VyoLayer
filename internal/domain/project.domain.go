package domain

import (
	"time"

	"github.com/vyolayer/vyolayer/internal/platform/database/types"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/utils"
)

// Project RBAC role constants (minimal)
const (
	ProjectRoleAdmin  = "admin"
	ProjectRoleMember = "member"
	ProjectRoleViewer = "viewer"
)

// Default project limits
var (
	DefaultProjectMaxApiKeys = 5
	DefaultProjectMaxMembers = 10
)

// ValidProjectRoles is the set of allowed project role values
var ValidProjectRoles = map[string]bool{
	ProjectRoleAdmin:  true,
	ProjectRoleMember: true,
	ProjectRoleViewer: true,
}

// IsValidProjectRole checks whether a role string is a valid project role
func IsValidProjectRole(role string) bool {
	return ValidProjectRoles[role]
}

// ---------------------------------------------------------------------------
// Project
// ---------------------------------------------------------------------------

type projectMemberInfo struct {
	isInfoLoaded   bool
	MaxNoOfMembers int
	NoOfMembers    int
	Members        []ProjectMember
}

type Project struct {
	ID             types.ProjectID
	OrganizationID types.OrganizationID
	Name           string
	Slug           string
	Description    string
	IsActive       bool
	CreatedBy      types.UserID
	MaxApiKeys     int

	MemberInfo projectMemberInfo
}

func NewProject(
	organizationID types.OrganizationID,
	creator *User,
	name string,
	description string,
	maxApiKeys *int,
	maxMembers *int,
) *Project {
	id := types.NewProjectID()
	slug := utils.ToSlug(name).Slugify().String()

	if maxApiKeys == nil {
		maxApiKeys = &DefaultProjectMaxApiKeys
	}
	if maxMembers == nil {
		maxMembers = &DefaultProjectMaxMembers
	}

	member := NewProjectMember(id, creator.ID, creator.ID, ProjectRoleAdmin)

	return &Project{
		ID:             id,
		OrganizationID: organizationID,
		Name:           name,
		Slug:           slug,
		Description:    description,
		IsActive:       true,
		CreatedBy:      creator.ID,
		MaxApiKeys:     *maxApiKeys,
		MemberInfo: projectMemberInfo{
			isInfoLoaded:   true,
			MaxNoOfMembers: *maxMembers,
			NoOfMembers:    1,
			Members:        []ProjectMember{*member},
		},
	}
}

// ReconstructProject rebuilds a Project from stored data
func ReconstructProject(
	id types.ProjectID,
	organizationID types.OrganizationID,
	name, slug, description string,
	isActive bool,
	createdBy types.UserID,
	maxApiKeys, maxMembers int,
) *Project {
	return &Project{
		ID:             id,
		OrganizationID: organizationID,
		Name:           name,
		Slug:           slug,
		Description:    description,
		IsActive:       isActive,
		CreatedBy:      createdBy,
		MaxApiKeys:     maxApiKeys,
		MemberInfo: projectMemberInfo{
			isInfoLoaded:   false,
			MaxNoOfMembers: maxMembers,
			NoOfMembers:    0,
			Members:        []ProjectMember{},
		},
	}
}

func (p *Project) Validate() *errors.AppError {
	if p.Name == "" {
		return ValidationError("Project name is required")
	}
	if p.Slug == "" {
		return ValidationError("Project slug is required")
	}
	if p.MaxApiKeys < 0 {
		return ValidationError("Max API keys cannot be negative")
	}
	if p.MemberInfo.MaxNoOfMembers < 1 {
		return ValidationError("Max members must be at least 1")
	}
	return nil
}

func (p *Project) UpdateName(name string) {
	p.Name = name
	p.Slug = utils.ToSlug(name).Slugify().String()
}

func (p *Project) UpdateDescription(description string) {
	p.Description = description
}

func (p *Project) Deactivate() *errors.AppError {
	if !p.IsActive {
		return ProjectNotActiveError(p.ID.String())
	}
	p.IsActive = false
	return nil
}

func (p *Project) Reactivate() *errors.AppError {
	if p.IsActive {
		return nil
	}
	p.IsActive = true
	return nil
}

func (p *Project) CanAddMember() bool {
	return p.MemberInfo.NoOfMembers < p.MemberInfo.MaxNoOfMembers
}

func (p *Project) LoadMembers(members []ProjectMember) {
	p.MemberInfo.Members = members
	p.MemberInfo.NoOfMembers = len(members)
	p.MemberInfo.isInfoLoaded = true
}

func (p *Project) GetMembers() ([]ProjectMember, *errors.AppError) {
	if !p.MemberInfo.isInfoLoaded {
		return nil, ProjectMembersNotLoadedError()
	}
	return p.MemberInfo.Members, nil
}

func (p *Project) IsMember(userID types.UserID) bool {
	if !p.MemberInfo.isInfoLoaded {
		return false
	}
	for _, m := range p.MemberInfo.Members {
		if m.UserID.String() == userID.String() {
			return true
		}
	}
	return false
}

func (p *Project) GetMemberByUserID(userID types.UserID) (*ProjectMember, *errors.AppError) {
	if !p.MemberInfo.isInfoLoaded {
		return nil, ProjectMembersNotLoadedError()
	}
	for _, m := range p.MemberInfo.Members {
		if m.UserID.String() == userID.String() {
			return &m, nil
		}
	}
	return nil, ProjectMemberNotFoundError(userID.String())
}

// ---------------------------------------------------------------------------
// Project Member
// ---------------------------------------------------------------------------

type ProjectMember struct {
	ID        types.ProjectMemberID
	ProjectID types.ProjectID
	UserID    types.UserID
	Role      string
	IsActive  bool

	// User info (populated when loaded via joins)
	Email    string
	FullName string

	// Tracking
	AddedBy   types.UserID
	JoinedAt  time.Time
	RemovedAt *time.Time
	RemovedBy *types.UserID
}

func NewProjectMember(
	projectID types.ProjectID,
	userID types.UserID,
	addedBy types.UserID,
	role string,
) *ProjectMember {
	return &ProjectMember{
		ID:        types.NewProjectMemberID(),
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
		IsActive:  true,
		AddedBy:   addedBy,
		JoinedAt:  time.Now(),
	}
}

func ReconstructProjectMember(
	id types.ProjectMemberID,
	projectID types.ProjectID,
	userID types.UserID,
	role string,
	email, fullName string,
	isActive bool,
	addedBy types.UserID,
	joinedAt time.Time,
	removedAt *time.Time,
	removedBy *types.UserID,
) *ProjectMember {
	return &ProjectMember{
		ID:        id,
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
		Email:     email,
		FullName:  fullName,
		IsActive:  isActive,
		AddedBy:   addedBy,
		JoinedAt:  joinedAt,
		RemovedAt: removedAt,
		RemovedBy: removedBy,
	}
}

func (pm *ProjectMember) Deactivate(removedBy types.UserID) *errors.AppError {
	if !pm.IsActive {
		return ProjectMemberNotActiveError()
	}
	now := time.Now()
	pm.IsActive = false
	pm.RemovedAt = &now
	pm.RemovedBy = &removedBy
	return nil
}

func (pm *ProjectMember) IsAdmin() bool {
	return pm.Role == ProjectRoleAdmin
}

func (pm *ProjectMember) IsMember() bool {
	return pm.Role == ProjectRoleMember || pm.IsAdmin()
}

func (pm *ProjectMember) IsViewer() bool {
	return pm.Role == ProjectRoleViewer || pm.IsMember()
}

// ---------------------------------------------------------------------------
// Project Invitation
// ---------------------------------------------------------------------------

type ProjectInvitation struct {
	ID        types.ProjectInvitationID
	ProjectID types.ProjectID
	InvitedBy types.ProjectMemberID
	Email     string
	Role      string
	Token     string

	InvitedAt  time.Time
	IsAccepted bool
	AcceptedAt *time.Time
	ExpiredAt  time.Time
}

func NewProjectInvitation(
	projectID types.ProjectID,
	invitedBy types.ProjectMemberID,
	email string,
	role string,
	token string,
	expiresAt time.Time,
) *ProjectInvitation {
	return &ProjectInvitation{
		ID:         types.NewProjectInvitationID(),
		ProjectID:  projectID,
		InvitedBy:  invitedBy,
		Email:      email,
		Role:       role,
		Token:      token,
		InvitedAt:  time.Now(),
		IsAccepted: false,
		ExpiredAt:  expiresAt,
	}
}

func (pi *ProjectInvitation) IsExpired() bool {
	return time.Now().After(pi.ExpiredAt)
}

func (pi *ProjectInvitation) IsPending() bool {
	return !pi.IsAccepted && !pi.IsExpired()
}

func (pi *ProjectInvitation) Accept() *errors.AppError {
	if pi.IsAccepted {
		return ProjectInvitationAlreadyAcceptedError(pi.ID.String())
	}
	if pi.IsExpired() {
		return ProjectInvitationExpiredError()
	}
	now := time.Now()
	pi.IsAccepted = true
	pi.AcceptedAt = &now
	return nil
}
