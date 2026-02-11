package repository

import (
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/models"
	"worklayer/internal/platform/database/types"
	"worklayer/pkg/errors"
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
)

type (
	// Base Models
	TBaseModel  = models.BaseModel
	TTimeStamps = models.TimeStamps

	// User Management
	TUser    = models.User
	TSession = models.UserSession

	// Organization Management
	TOrganization       = models.Organization
	TOrganizationMember = models.OrganizationMember

	// Organization Role Management
	TOrganizationRole           = models.OrganizationRole
	TMemberOrganizationRole     = models.MemberOrganizationRole
	TOrganizationPermission     = models.OrganizationPermission
	TOrganizationRolePermission = models.OrganizationRolePermission
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
