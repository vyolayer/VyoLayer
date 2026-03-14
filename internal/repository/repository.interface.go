package repository

import (
	"vyolayer/internal/domain"
	"vyolayer/internal/platform/database/models"
	"vyolayer/internal/platform/database/types"
	"vyolayer/pkg/errors"
)

// IDs type
type (
	// User IDs
	UserID = types.UserID

	// Organization IDs
	OrgID = types.OrganizationID

	// Organization Role IDs
	OrgRoleID       = types.OrganizationRoleID
	OrgPermissionID = types.OrganizationPermissionID

	// Organization Member IDs
	OrgMemberID     = types.OrganizationMemberID
	OrgMemberRoleID = types.MemberOrganizationRoleID

	// Invitation IDs
	InvitationID = types.OrganizationMemberInvitationID

	// Project IDs
	ProjectID           = types.ProjectID
	ProjectMemberID     = types.ProjectMemberID
	ProjectInvitationID = types.ProjectInvitationID

	// API Key IDs
	ApiKeyID = types.ApiKeyID
)

type (
	// Base Models
	TBaseModel  = models.BaseModel
	TTimeStamps = models.TimeStamps

	// User Management
	TUser    = models.User
	TSession = models.UserSession

	// Organization Management
	TOrganization                 = models.Organization
	TOrganizationMember           = models.OrganizationMember
	TOrganizationMemberInvitation = models.OrganizationMemberInvitation

	// Organization Role Management
	TOrganizationRole           = models.OrganizationRole
	TMemberOrganizationRole     = models.MemberOrganizationRole
	TOrganizationPermission     = models.OrganizationPermission
	TOrganizationRolePermission = models.OrganizationRolePermission

	// Project Management
	TProject           = models.Project
	TProjectMember     = models.ProjectMember
	TProjectInvitation = models.ProjectInvitation

	// API Key Management
	TApiKey         = models.ApiKey
	TApiKeyUsageLog = models.ApiKeyUsageLog
)

type UserRepository interface {
	CreateUser(user domain.User) (*domain.User, *errors.AppError)
	FindByEmail(email string) (*domain.User, *errors.AppError)
	FindById(id types.UserID) (*domain.User, *errors.AppError)
}

type SessionRepository interface {
	Save(session *models.UserSession) *errors.AppError
	FindByUserId(userId types.UserID) (*models.UserSession, *errors.AppError)
	FindByTokenHash(hashedToken string) (*models.UserSession, *errors.AppError)
	RotateByTokenHash(oldHashedToken string, newSession *models.UserSession) *errors.AppError
	DeleteByTokenHash(tokenHash string) *errors.AppError
}
