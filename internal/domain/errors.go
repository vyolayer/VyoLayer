package domain

import "vyolayer/pkg/errors"

// Domain-specific error helpers
// These provide convenient constructors for domain layer errors

// User errors
var (
	ErrUserNotFound       = errors.ErrUserNotFound
	ErrUserAlreadyExists  = errors.ErrUserAlreadyExists
	ErrUserInactive       = errors.ErrUserInactive
	ErrInvalidCredentials = errors.ErrAuthInvalidCredentials
	ErrUserNotVerified    = errors.ErrAuthAccountNotVerified
	ErrSessionNotFound    = errors.ErrAuthSessionNotFound
	ErrSessionExpired     = errors.ErrAuthSessionExpired
	ErrTokenExpired       = errors.ErrAuthTokenExpired
	ErrTokenInvalid       = errors.ErrAuthTokenInvalid
	ErrPasswordHashFailed = errors.ErrInternalHashing
)

// UserNotFoundError creates a user not found error
func UserNotFoundError(userID string) *errors.AppError {
	return errors.UserNotFound(userID)
}

// UserAlreadyExistsError creates a user already exists error
func UserAlreadyExistsError(email string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrUserAlreadyExists, "User with email '%s' already exists", email).
		WithMetadata("email", email)
}

// InvalidCredentialsError creates an invalid credentials error
func InvalidCredentialsError() *errors.AppError {
	return errors.InvalidCredentials()
}

// UserNotVerifiedError creates a user not verified error
func UserNotVerifiedError() *errors.AppError {
	return errors.AccountNotVerified()
}

// SessionNotFoundError creates a session not found error
func SessionNotFoundError() *errors.AppError {
	return errors.SessionNotFound()
}

// SessionExpiredError creates a session expired error
func SessionExpiredError() *errors.AppError {
	return errors.SessionExpired()
}

// TokenExpiredError creates a token expired error
func TokenExpiredError() *errors.AppError {
	return errors.TokenExpired()
}

// TokenInvalidError creates an invalid token error
func TokenInvalidError(reason string) *errors.AppError {
	return errors.TokenInvalid(reason)
}

// PasswordHashFailedError creates a password hash failed error
func PasswordHashFailedError(err error) *errors.AppError {
	return errors.Wrap(err, errors.ErrInternalHashing, "Failed to hash password")
}

// ValidationError creates a validation error
func ValidationError(message string) *errors.AppError {
	return errors.ValidationFailed(message)
}

// InvalidEmailError creates an invalid email error
func InvalidEmailError(email string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrValidationInvalidEmail, "Invalid email: %s", email).
		WithMetadata("email", email)
}

// InvalidPasswordError creates an invalid password error
func InvalidPasswordError(reason string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrValidationInvalidPassword, "Invalid password: %s", reason).
		WithMetadata("reason", reason)
}

// Organization errors
var (
	ErrOrganizationNotFound            = errors.ErrOrganizationNotFound
	ErrOrganizationNotActive           = errors.ErrOrganizationNotActive
	ErrOrganizationNotOwner            = errors.ErrOrganizationNotOwner
	ErrOrganizationFull                = errors.ErrOrganizationFull
	ErrOrganizationInfoNotLoadedFromDB = errors.ErrOrganizationInfoNotLoadedFromDB
	ErrOrganizationMemberAlreadyExists = errors.ErrOrganizationMemberAlreadyExists
	ErrOrganizationMemberNotFound      = errors.ErrOrganizationMemberNotFound
	ErrOrganizationMemberNotActive     = errors.ErrOrganizationMemberNotActive
)

// OrganizationNotFoundError creates an organization not found error
func OrganizationNotFoundError(orgID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrOrganizationNotFound, "Organization '%s' not found", orgID).
		WithMetadata("organization_id", orgID)
}

// OrganizationNotActiveError creates an organization not active error
func OrganizationNotActiveError() *errors.AppError {
	return errors.NewWithMessage(errors.ErrOrganizationNotActive, "Organization is not active")
}

// OrganizationNotOwnerError creates an organization not owner error
func OrganizationNotOwnerError(userID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrOrganizationNotOwner, "User '%s' is not the owner of this organization", userID).
		WithMetadata("user_id", userID)
}

// OrganizationFullError creates an organization full error
func OrganizationFullError() *errors.AppError {
	return errors.NewWithMessage(errors.ErrOrganizationFull, "Organization has reached maximum member capacity")
}

// OrganizationMembersNotLoadedError creates an error for when members are not loaded
func OrganizationMembersNotLoadedError() *errors.AppError {
	return errors.NewWithMessage(errors.ErrOrganizationInfoNotLoadedFromDB, "Organization member information not loaded from database")
}

// OrganizationMemberAlreadyExistsError creates an error for when a member already exists
func OrganizationMemberAlreadyExistsError(userID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrOrganizationMemberAlreadyExists, "User '%s' is already a member of this organization", userID).
		WithMetadata("user_id", userID)
}

// OrganizationMemberNotFoundError creates a member not found error
func OrganizationMemberNotFoundError(memberID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrOrganizationMemberNotFound, "Organization member '%s' not found", memberID).
		WithMetadata("member_id", memberID)
}

// OrganizationMemberNotActiveError creates a member not active error
func OrganizationMemberNotActiveError() *errors.AppError {
	return errors.NewWithMessage(errors.ErrOrganizationMemberNotActive, "Organization member is not active")
}

// OrganizationCannotRemoveOwnerError creates an error for when trying to remove the owner
func OrganizationCannotRemoveOwnerError() *errors.AppError {
	return errors.NewWithMessage(errors.ErrOrganizationNotOwner, "Cannot remove the organization owner")
}

// OrganizationLastOwnerError creates an error for when an action would leave the org without an owner
func OrganizationLastOwnerError() *errors.AppError {
	return errors.BadRequest("Cannot perform this action: organization must always have at least one owner")
}

// OrganizationSlugConflictError creates an error for slug uniqueness violation
func OrganizationSlugConflictError(slug string) *errors.AppError {
	return errors.Conflict("Slug '%s' is already taken", slug)
}

// OrganizationDeleteConfirmationError creates an error for failed delete confirmation
func OrganizationDeleteConfirmationError() *errors.AppError {
	return errors.BadRequest("Organization name does not match. Please type the exact organization name to confirm deletion.")
}

// Invitation errors
var (
	ErrInvitationNotFound        = errors.ErrInvitationNotFound
	ErrInvitationExpired         = errors.ErrInvitationExpired
	ErrInvitationAlreadyAccepted = errors.ErrInvitationAlreadyAccepted
	ErrInvitationAlreadyExists   = errors.ErrInvitationAlreadyExists
	ErrInvitationInvalid         = errors.ErrInvitationInvalid
)

// InvitationNotFoundError creates an invitation not found error
func InvitationNotFoundError(invitationID string) *errors.AppError {
	return errors.InvitationNotFound(invitationID)
}

// InvitationExpiredError creates an invitation expired error
func InvitationExpiredError() *errors.AppError {
	return errors.InvitationExpired()
}

// InvitationAlreadyAcceptedError creates an invitation already accepted error
func InvitationAlreadyAcceptedError(invitationID string) *errors.AppError {
	return errors.InvitationAlreadyAccepted(invitationID)
}

// InvitationAlreadyExistsError creates an invitation already exists error
func InvitationAlreadyExistsError(email, orgID string) *errors.AppError {
	return errors.InvitationAlreadyExists(email, orgID)
}

// InvitationInvalidError creates an invalid invitation error
func InvitationInvalidError(reason string) *errors.AppError {
	return errors.InvitationInvalid(reason)
}

// Project errors
var (
	ErrProjectNotFound            = errors.ErrProjectNotFound
	ErrProjectNotActive           = errors.ErrProjectNotActive
	ErrProjectFull                = errors.ErrProjectFull
	ErrProjectLimitReached        = errors.ErrProjectLimitReached
	ErrProjectMemberAlreadyExists = errors.ErrProjectMemberAlreadyExists
	ErrProjectMemberNotFound      = errors.ErrProjectMemberNotFound
	ErrProjectMemberNotActive     = errors.ErrProjectMemberNotActive
	ErrProjectInfoNotLoaded       = errors.ErrProjectInfoNotLoaded
)

// ProjectNotFoundError creates a project not found error
func ProjectNotFoundError(projectID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrProjectNotFound, "Project '%s' not found", projectID).
		WithMetadata("project_id", projectID)
}

// ProjectNotActiveError creates a project not active error
func ProjectNotActiveError(projectID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrProjectNotActive, "Project '%s' is not active", projectID).
		WithMetadata("project_id", projectID)
}

// ProjectFullError creates a project full error
func ProjectFullError() *errors.AppError {
	return errors.NewWithMessage(errors.ErrProjectFull, "Project has reached maximum member capacity")
}

// ProjectLimitReachedError creates a project limit reached error
func ProjectLimitReachedError() *errors.AppError {
	return errors.NewWithMessage(errors.ErrProjectLimitReached, "Organization has reached maximum project limit")
}

// ProjectMemberAlreadyExistsError creates a project member already exists error
func ProjectMemberAlreadyExistsError(userID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrProjectMemberAlreadyExists, "User '%s' is already a member of this project", userID).
		WithMetadata("user_id", userID)
}

// ProjectMemberNotFoundError creates a project member not found error
func ProjectMemberNotFoundError(memberID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrProjectMemberNotFound, "Project member '%s' not found", memberID).
		WithMetadata("member_id", memberID)
}

// ProjectMemberNotActiveError creates a project member not active error
func ProjectMemberNotActiveError() *errors.AppError {
	return errors.NewWithMessage(errors.ErrProjectMemberNotActive, "Project member is not active")
}

// ProjectMembersNotLoadedError creates an error for when project members are not loaded
func ProjectMembersNotLoadedError() *errors.AppError {
	return errors.NewWithMessage(errors.ErrProjectInfoNotLoaded, "Project member information not loaded from database")
}

// ProjectSlugConflictError creates a project slug conflict error
func ProjectSlugConflictError(slug string) *errors.AppError {
	return errors.Conflict("Slug '%s' is already taken in this organization", slug)
}

// ProjectDeleteConfirmationError creates a project delete confirmation error
func ProjectDeleteConfirmationError() *errors.AppError {
	return errors.BadRequest("Project name does not match. Please type the exact project name to confirm deletion.")
}

// Project invitation errors
func ProjectInvitationNotFoundError(invitationID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrProjectInvitationNotFound, "Project invitation '%s' not found", invitationID).
		WithMetadata("invitation_id", invitationID)
}

func ProjectInvitationExpiredError() *errors.AppError {
	return errors.New(errors.ErrProjectInvitationExpired)
}

func ProjectInvitationAlreadyAcceptedError(invitationID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrProjectInvitationAccepted, "Project invitation '%s' has already been accepted", invitationID).
		WithMetadata("invitation_id", invitationID)
}

func ProjectInvitationAlreadyExistsError(email, projectID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrProjectInvitationExists, "An invitation for '%s' to project '%s' already exists", email, projectID).
		WithMetadata("email", email).
		WithMetadata("project_id", projectID)
}

// API Key errors
var (
	ErrApiKeyNotFound     = errors.ErrApiKeyNotFound
	ErrApiKeyRevoked      = errors.ErrApiKeyRevoked
	ErrApiKeyExpired      = errors.ErrApiKeyExpired
	ErrApiKeyInvalid      = errors.ErrApiKeyInvalid
	ErrApiKeyLimitReached = errors.ErrApiKeyLimitReached
	ErrApiKeyRateLimited  = errors.ErrApiKeyRateLimited
)

// ApiKeyNotFoundError creates an API key not found error
func ApiKeyNotFoundError(apiKeyID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrApiKeyNotFound, "API key '%s' not found", apiKeyID).
		WithMetadata("api_key_id", apiKeyID)
}

// ApiKeyRevokedError creates an API key revoked error
func ApiKeyRevokedError(apiKeyID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrApiKeyRevoked, "API key '%s' has been revoked", apiKeyID).
		WithMetadata("api_key_id", apiKeyID)
}

// ApiKeyExpiredError creates an API key expired error
func ApiKeyExpiredError(apiKeyID string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrApiKeyExpired, "API key '%s' has expired", apiKeyID).
		WithMetadata("api_key_id", apiKeyID)
}

// ApiKeyInvalidError creates an invalid API key error
func ApiKeyInvalidError() *errors.AppError {
	return errors.New(errors.ErrApiKeyInvalid)
}

// ApiKeyLimitReachedError creates an API key limit reached error
func ApiKeyLimitReachedError() *errors.AppError {
	return errors.NewWithMessage(errors.ErrApiKeyLimitReached, "Project has reached maximum API key limit")
}

// ApiKeyRateLimitedError creates an API key rate limited error
func ApiKeyRateLimitedError() *errors.AppError {
	return errors.NewWithMessage(errors.ErrApiKeyRateLimited, "API key rate limit exceeded")
}
