package errors

import "net/http"

// ErrorCode represents a unique error code for the application
type ErrorCode string

// Error Categories and Codes
const (
	// Authentication & Authorization Errors (AUTH)
	ErrAuthInvalidCredentials ErrorCode = "ERR_AUTH_INVALID_CREDENTIALS"
	ErrAuthUnauthorized       ErrorCode = "ERR_AUTH_UNAUTHORIZED"
	ErrAuthTokenExpired       ErrorCode = "ERR_AUTH_TOKEN_EXPIRED"
	ErrAuthTokenInvalid       ErrorCode = "ERR_AUTH_TOKEN_INVALID"
	ErrAuthSessionNotFound    ErrorCode = "ERR_AUTH_SESSION_NOT_FOUND"
	ErrAuthSessionExpired     ErrorCode = "ERR_AUTH_SESSION_EXPIRED"
	ErrAuthForbidden          ErrorCode = "ERR_AUTH_FORBIDDEN"
	ErrAuthPasswordMismatch   ErrorCode = "ERR_AUTH_PASSWORD_MISMATCH"
	ErrAuthAccountLocked      ErrorCode = "ERR_AUTH_ACCOUNT_LOCKED"
	ErrAuthAccountNotVerified ErrorCode = "ERR_AUTH_ACCOUNT_NOT_VERIFIED"

	// Validation Errors (VALIDATION)
	ErrValidationFailed          ErrorCode = "ERR_VALIDATION_FAILED"
	ErrValidationRequiredField   ErrorCode = "ERR_VALIDATION_REQUIRED_FIELD"
	ErrValidationInvalidFormat   ErrorCode = "ERR_VALIDATION_INVALID_FORMAT"
	ErrValidationInvalidEmail    ErrorCode = "ERR_VALIDATION_INVALID_EMAIL"
	ErrValidationInvalidPassword ErrorCode = "ERR_VALIDATION_INVALID_PASSWORD"
	ErrValidationInvalidInput    ErrorCode = "ERR_VALIDATION_INVALID_INPUT"
	ErrValidationTooShort        ErrorCode = "ERR_VALIDATION_TOO_SHORT"
	ErrValidationTooLong         ErrorCode = "ERR_VALIDATION_TOO_LONG"
	ErrValidationOutOfRange      ErrorCode = "ERR_VALIDATION_OUT_OF_RANGE"

	// Database & Repository Errors (DB)
	ErrDBRecordNotFound      ErrorCode = "ERR_DB_RECORD_NOT_FOUND"
	ErrDBDuplicateKey        ErrorCode = "ERR_DB_DUPLICATE_KEY"
	ErrDBConstraintViolation ErrorCode = "ERR_DB_CONSTRAINT_VIOLATION"
	ErrDBQueryFailed         ErrorCode = "ERR_DB_QUERY_FAILED"
	ErrDBConnectionFailed    ErrorCode = "ERR_DB_CONNECTION_FAILED"
	ErrDBTransactionFailed   ErrorCode = "ERR_DB_TRANSACTION_FAILED"
	ErrDBMigrationFailed     ErrorCode = "ERR_DB_MIGRATION_FAILED"

	// Resource Errors (RESOURCE)
	ErrResourceNotFound      ErrorCode = "ERR_RESOURCE_NOT_FOUND"
	ErrResourceAlreadyExists ErrorCode = "ERR_RESOURCE_ALREADY_EXISTS"
	ErrResourceConflict      ErrorCode = "ERR_RESOURCE_CONFLICT"
	ErrResourceLocked        ErrorCode = "ERR_RESOURCE_LOCKED"
	ErrResourceDeleted       ErrorCode = "ERR_RESOURCE_DELETED"

	// User Errors (USER)
	ErrUserNotFound      ErrorCode = "ERR_USER_NOT_FOUND"
	ErrUserAlreadyExists ErrorCode = "ERR_USER_ALREADY_EXISTS"
	ErrUserInactive      ErrorCode = "ERR_USER_INACTIVE"
	ErrUserDeleted       ErrorCode = "ERR_USER_DELETED"

	// Organization Errors (ORGANIZATION)
	ErrOrganizationFull                ErrorCode = "ERR_ORGANIZATION_FULL"
	ErrOrganizationNotFound            ErrorCode = "ERR_ORGANIZATION_NOT_FOUND"
	ErrOrganizationDeleted             ErrorCode = "ERR_ORGANIZATION_DELETED"
	ErrOrganizationNotActive           ErrorCode = "ERR_ORGANIZATION_NOT_ACTIVE"
	ErrOrganizationNotVerified         ErrorCode = "ERR_ORGANIZATION_NOT_VERIFIED"
	ErrOrganizationNotOwner            ErrorCode = "ERR_ORGANIZATION_NOT_OWNER"
	ErrOrganizationInfoNotLoadedFromDB ErrorCode = "ERR_ORGANIZATION_INFO_NOT_LOADED_FROM_DB"

	ErrOrganizationMemberAlreadyExists ErrorCode = "ERR_ORGANIZATION_MEMBER_ALREADY_EXISTS"
	ErrOrganizationMemberNotFound      ErrorCode = "ERR_ORGANIZATION_MEMBER_NOT_FOUND"
	ErrOrganizationMemberDeleted       ErrorCode = "ERR_ORGANIZATION_MEMBER_DELETED"
	ErrOrganizationMemberNotActive     ErrorCode = "ERR_ORGANIZATION_MEMBER_NOT_ACTIVE"
	ErrOrganizationMemberNotVerified   ErrorCode = "ERR_ORGANIZATION_MEMBER_NOT_VERIFIED"
	ErrOrganizationMemberNotOwner      ErrorCode = "ERR_ORGANIZATION_MEMBER_NOT_OWNER"

	ErrInvitationNotFound        ErrorCode = "ERR_INVITATION_NOT_FOUND"
	ErrInvitationExpired         ErrorCode = "ERR_INVITATION_EXPIRED"
	ErrInvitationAlreadyAccepted ErrorCode = "ERR_INVITATION_ALREADY_ACCEPTED"
	ErrInvitationAlreadyExists   ErrorCode = "ERR_INVITATION_ALREADY_EXISTS"
	ErrInvitationInvalid         ErrorCode = "ERR_INVITATION_INVALID"

	// Project Errors (PROJECT)
	ErrProjectNotFound            ErrorCode = "ERR_PROJECT_NOT_FOUND"
	ErrProjectNotActive           ErrorCode = "ERR_PROJECT_NOT_ACTIVE"
	ErrProjectFull                ErrorCode = "ERR_PROJECT_FULL"
	ErrProjectLimitReached        ErrorCode = "ERR_PROJECT_LIMIT_REACHED"
	ErrProjectMemberAlreadyExists ErrorCode = "ERR_PROJECT_MEMBER_ALREADY_EXISTS"
	ErrProjectMemberNotFound      ErrorCode = "ERR_PROJECT_MEMBER_NOT_FOUND"
	ErrProjectMemberNotActive     ErrorCode = "ERR_PROJECT_MEMBER_NOT_ACTIVE"
	ErrProjectInfoNotLoaded       ErrorCode = "ERR_PROJECT_INFO_NOT_LOADED"
	ErrProjectInvitationNotFound  ErrorCode = "ERR_PROJECT_INVITATION_NOT_FOUND"
	ErrProjectInvitationExpired   ErrorCode = "ERR_PROJECT_INVITATION_EXPIRED"
	ErrProjectInvitationAccepted  ErrorCode = "ERR_PROJECT_INVITATION_ALREADY_ACCEPTED"
	ErrProjectInvitationExists    ErrorCode = "ERR_PROJECT_INVITATION_ALREADY_EXISTS"

	// API Key Errors (API_KEY)
	ErrApiKeyNotFound     ErrorCode = "ERR_API_KEY_NOT_FOUND"
	ErrApiKeyRevoked      ErrorCode = "ERR_API_KEY_REVOKED"
	ErrApiKeyExpired      ErrorCode = "ERR_API_KEY_EXPIRED"
	ErrApiKeyInvalid      ErrorCode = "ERR_API_KEY_INVALID"
	ErrApiKeyLimitReached ErrorCode = "ERR_API_KEY_LIMIT_REACHED"
	ErrApiKeyRateLimited  ErrorCode = "ERR_API_KEY_RATE_LIMITED"

	// Request Errors (REQUEST)
	ErrRequestInvalidBody    ErrorCode = "ERR_REQUEST_INVALID_BODY"
	ErrRequestInvalidParams  ErrorCode = "ERR_REQUEST_INVALID_PARAMS"
	ErrRequestInvalidHeaders ErrorCode = "ERR_REQUEST_INVALID_HEADERS"
	ErrRequestTooLarge       ErrorCode = "ERR_REQUEST_TOO_LARGE"
	ErrRequestTimeout        ErrorCode = "ERR_REQUEST_TIMEOUT"
	ErrRequestRateLimited    ErrorCode = "ERR_REQUEST_RATE_LIMITED"

	// External Service Errors (EXTERNAL)
	ErrExternalServiceUnavailable ErrorCode = "ERR_EXTERNAL_SERVICE_UNAVAILABLE"
	ErrExternalServiceTimeout     ErrorCode = "ERR_EXTERNAL_SERVICE_TIMEOUT"
	ErrExternalServiceError       ErrorCode = "ERR_EXTERNAL_SERVICE_ERROR"
	ErrExternalAPIError           ErrorCode = "ERR_EXTERNAL_API_ERROR"

	// Internal Errors (INTERNAL)
	ErrInternalUnexpected     ErrorCode = "ERR_INTERNAL_UNEXPECTED"
	ErrInternalNotImplemented ErrorCode = "ERR_INTERNAL_NOT_IMPLEMENTED"
	ErrInternalConfiguration  ErrorCode = "ERR_INTERNAL_CONFIGURATION"
	ErrInternalEncryption     ErrorCode = "ERR_INTERNAL_ENCRYPTION"
	ErrInternalDecryption     ErrorCode = "ERR_INTERNAL_DECRYPTION"
	ErrInternalHashing        ErrorCode = "ERR_INTERNAL_HASHING"
	ErrInternalSerialization  ErrorCode = "ERR_INTERNAL_SERIALIZATION"

	// Business Logic Errors (BUSINESS)
	ErrBusinessRuleViolation       ErrorCode = "ERR_BUSINESS_RULE_VIOLATION"
	ErrBusinessInvalidState        ErrorCode = "ERR_BUSINESS_INVALID_STATE"
	ErrBusinessOperationNotAllowed ErrorCode = "ERR_BUSINESS_OPERATION_NOT_ALLOWED"
)

// ErrorCodeMetadata contains metadata about an error code
type ErrorCodeMetadata struct {
	Code           ErrorCode
	DefaultMessage string
	HTTPStatus     int
	Severity       Severity
}

// errorCodeRegistry maps error codes to their metadata
var errorCodeRegistry = map[ErrorCode]ErrorCodeMetadata{
	// Authentication & Authorization
	ErrAuthInvalidCredentials: {
		Code:           ErrAuthInvalidCredentials,
		DefaultMessage: "Invalid credentials provided",
		HTTPStatus:     http.StatusUnauthorized,
		Severity:       SeverityWarning,
	},
	ErrAuthUnauthorized: {
		Code:           ErrAuthUnauthorized,
		DefaultMessage: "Unauthorized access",
		HTTPStatus:     http.StatusUnauthorized,
		Severity:       SeverityWarning,
	},
	ErrAuthTokenExpired: {
		Code:           ErrAuthTokenExpired,
		DefaultMessage: "Authentication token has expired",
		HTTPStatus:     http.StatusUnauthorized,
		Severity:       SeverityInfo,
	},
	ErrAuthTokenInvalid: {
		Code:           ErrAuthTokenInvalid,
		DefaultMessage: "Invalid authentication token",
		HTTPStatus:     http.StatusUnauthorized,
		Severity:       SeverityWarning,
	},
	ErrAuthSessionNotFound: {
		Code:           ErrAuthSessionNotFound,
		DefaultMessage: "Session not found",
		HTTPStatus:     http.StatusUnauthorized,
		Severity:       SeverityInfo,
	},
	ErrAuthSessionExpired: {
		Code:           ErrAuthSessionExpired,
		DefaultMessage: "Session has expired",
		HTTPStatus:     http.StatusUnauthorized,
		Severity:       SeverityInfo,
	},
	ErrAuthForbidden: {
		Code:           ErrAuthForbidden,
		DefaultMessage: "Access to this resource is forbidden",
		HTTPStatus:     http.StatusForbidden,
		Severity:       SeverityWarning,
	},
	ErrAuthPasswordMismatch: {
		Code:           ErrAuthPasswordMismatch,
		DefaultMessage: "Password does not match",
		HTTPStatus:     http.StatusUnauthorized,
		Severity:       SeverityWarning,
	},
	ErrAuthAccountLocked: {
		Code:           ErrAuthAccountLocked,
		DefaultMessage: "Account is locked",
		HTTPStatus:     http.StatusForbidden,
		Severity:       SeverityWarning,
	},
	ErrAuthAccountNotVerified: {
		Code:           ErrAuthAccountNotVerified,
		DefaultMessage: "Account is not verified",
		HTTPStatus:     http.StatusForbidden,
		Severity:       SeverityInfo,
	},

	// Validation
	ErrValidationFailed: {
		Code:           ErrValidationFailed,
		DefaultMessage: "Validation failed",
		HTTPStatus:     http.StatusUnprocessableEntity,
		Severity:       SeverityInfo,
	},
	ErrValidationRequiredField: {
		Code:           ErrValidationRequiredField,
		DefaultMessage: "Required field is missing",
		HTTPStatus:     http.StatusUnprocessableEntity,
		Severity:       SeverityInfo,
	},
	ErrValidationInvalidFormat: {
		Code:           ErrValidationInvalidFormat,
		DefaultMessage: "Invalid format",
		HTTPStatus:     http.StatusUnprocessableEntity,
		Severity:       SeverityInfo,
	},
	ErrValidationInvalidEmail: {
		Code:           ErrValidationInvalidEmail,
		DefaultMessage: "Invalid email address",
		HTTPStatus:     http.StatusUnprocessableEntity,
		Severity:       SeverityInfo,
	},
	ErrValidationInvalidPassword: {
		Code:           ErrValidationInvalidPassword,
		DefaultMessage: "Invalid password format",
		HTTPStatus:     http.StatusUnprocessableEntity,
		Severity:       SeverityInfo,
	},
	ErrValidationInvalidInput: {
		Code:           ErrValidationInvalidInput,
		DefaultMessage: "Invalid input provided",
		HTTPStatus:     http.StatusBadRequest,
		Severity:       SeverityInfo,
	},
	ErrValidationTooShort: {
		Code:           ErrValidationTooShort,
		DefaultMessage: "Value is too short",
		HTTPStatus:     http.StatusUnprocessableEntity,
		Severity:       SeverityInfo,
	},
	ErrValidationTooLong: {
		Code:           ErrValidationTooLong,
		DefaultMessage: "Value is too long",
		HTTPStatus:     http.StatusUnprocessableEntity,
		Severity:       SeverityInfo,
	},
	ErrValidationOutOfRange: {
		Code:           ErrValidationOutOfRange,
		DefaultMessage: "Value is out of allowed range",
		HTTPStatus:     http.StatusUnprocessableEntity,
		Severity:       SeverityInfo,
	},

	// Database
	ErrDBRecordNotFound: {
		Code:           ErrDBRecordNotFound,
		DefaultMessage: "Record not found in database",
		HTTPStatus:     http.StatusNotFound,
		Severity:       SeverityInfo,
	},
	ErrDBDuplicateKey: {
		Code:           ErrDBDuplicateKey,
		DefaultMessage: "Duplicate key violation",
		HTTPStatus:     http.StatusConflict,
		Severity:       SeverityWarning,
	},
	ErrDBConstraintViolation: {
		Code:           ErrDBConstraintViolation,
		DefaultMessage: "Database constraint violation",
		HTTPStatus:     http.StatusConflict,
		Severity:       SeverityWarning,
	},
	ErrDBQueryFailed: {
		Code:           ErrDBQueryFailed,
		DefaultMessage: "Database query failed",
		HTTPStatus:     http.StatusInternalServerError,
		Severity:       SeverityError,
	},
	ErrDBConnectionFailed: {
		Code:           ErrDBConnectionFailed,
		DefaultMessage: "Failed to connect to database",
		HTTPStatus:     http.StatusServiceUnavailable,
		Severity:       SeverityCritical,
	},
	ErrDBTransactionFailed: {
		Code:           ErrDBTransactionFailed,
		DefaultMessage: "Database transaction failed",
		HTTPStatus:     http.StatusInternalServerError,
		Severity:       SeverityError,
	},
	ErrDBMigrationFailed: {
		Code:           ErrDBMigrationFailed,
		DefaultMessage: "Database migration failed",
		HTTPStatus:     http.StatusInternalServerError,
		Severity:       SeverityCritical,
	},

	// Resource
	ErrResourceNotFound: {
		Code:           ErrResourceNotFound,
		DefaultMessage: "Resource not found",
		HTTPStatus:     http.StatusNotFound,
		Severity:       SeverityInfo,
	},
	ErrResourceAlreadyExists: {
		Code:           ErrResourceAlreadyExists,
		DefaultMessage: "Resource already exists",
		HTTPStatus:     http.StatusConflict,
		Severity:       SeverityWarning,
	},
	ErrResourceConflict: {
		Code:           ErrResourceConflict,
		DefaultMessage: "Resource conflict",
		HTTPStatus:     http.StatusConflict,
		Severity:       SeverityWarning,
	},
	ErrResourceLocked: {
		Code:           ErrResourceLocked,
		DefaultMessage: "Resource is locked",
		HTTPStatus:     http.StatusLocked,
		Severity:       SeverityWarning,
	},
	ErrResourceDeleted: {
		Code:           ErrResourceDeleted,
		DefaultMessage: "Resource has been deleted",
		HTTPStatus:     http.StatusGone,
		Severity:       SeverityInfo,
	},

	// User
	ErrUserNotFound: {
		Code:           ErrUserNotFound,
		DefaultMessage: "User not found",
		HTTPStatus:     http.StatusNotFound,
		Severity:       SeverityInfo,
	},
	ErrUserAlreadyExists: {
		Code:           ErrUserAlreadyExists,
		DefaultMessage: "User already exists",
		HTTPStatus:     http.StatusConflict,
		Severity:       SeverityWarning,
	},
	ErrUserInactive: {
		Code:           ErrUserInactive,
		DefaultMessage: "User account is inactive",
		HTTPStatus:     http.StatusForbidden,
		Severity:       SeverityWarning,
	},
	ErrUserDeleted: {
		Code:           ErrUserDeleted,
		DefaultMessage: "User account has been deleted",
		HTTPStatus:     http.StatusGone,
		Severity:       SeverityInfo,
	},

	// Request
	ErrRequestInvalidBody: {
		Code:           ErrRequestInvalidBody,
		DefaultMessage: "Invalid request body",
		HTTPStatus:     http.StatusBadRequest,
		Severity:       SeverityInfo,
	},
	ErrRequestInvalidParams: {
		Code:           ErrRequestInvalidParams,
		DefaultMessage: "Invalid request parameters",
		HTTPStatus:     http.StatusBadRequest,
		Severity:       SeverityInfo,
	},
	ErrRequestInvalidHeaders: {
		Code:           ErrRequestInvalidHeaders,
		DefaultMessage: "Invalid request headers",
		HTTPStatus:     http.StatusBadRequest,
		Severity:       SeverityInfo,
	},
	ErrRequestTooLarge: {
		Code:           ErrRequestTooLarge,
		DefaultMessage: "Request payload too large",
		HTTPStatus:     http.StatusRequestEntityTooLarge,
		Severity:       SeverityWarning,
	},
	ErrRequestTimeout: {
		Code:           ErrRequestTimeout,
		DefaultMessage: "Request timeout",
		HTTPStatus:     http.StatusRequestTimeout,
		Severity:       SeverityWarning,
	},
	ErrRequestRateLimited: {
		Code:           ErrRequestRateLimited,
		DefaultMessage: "Rate limit exceeded",
		HTTPStatus:     http.StatusTooManyRequests,
		Severity:       SeverityWarning,
	},

	// External Service
	ErrExternalServiceUnavailable: {
		Code:           ErrExternalServiceUnavailable,
		DefaultMessage: "External service is unavailable",
		HTTPStatus:     http.StatusServiceUnavailable,
		Severity:       SeverityError,
	},
	ErrExternalServiceTimeout: {
		Code:           ErrExternalServiceTimeout,
		DefaultMessage: "External service timeout",
		HTTPStatus:     http.StatusGatewayTimeout,
		Severity:       SeverityError,
	},
	ErrExternalServiceError: {
		Code:           ErrExternalServiceError,
		DefaultMessage: "External service error",
		HTTPStatus:     http.StatusBadGateway,
		Severity:       SeverityError,
	},
	ErrExternalAPIError: {
		Code:           ErrExternalAPIError,
		DefaultMessage: "External API error",
		HTTPStatus:     http.StatusBadGateway,
		Severity:       SeverityError,
	},

	// Internal
	ErrInternalUnexpected: {
		Code:           ErrInternalUnexpected,
		DefaultMessage: "An unexpected error occurred",
		HTTPStatus:     http.StatusInternalServerError,
		Severity:       SeverityError,
	},
	ErrInternalNotImplemented: {
		Code:           ErrInternalNotImplemented,
		DefaultMessage: "Feature not implemented",
		HTTPStatus:     http.StatusNotImplemented,
		Severity:       SeverityWarning,
	},
	ErrInternalConfiguration: {
		Code:           ErrInternalConfiguration,
		DefaultMessage: "Configuration error",
		HTTPStatus:     http.StatusInternalServerError,
		Severity:       SeverityCritical,
	},
	ErrInternalEncryption: {
		Code:           ErrInternalEncryption,
		DefaultMessage: "Encryption failed",
		HTTPStatus:     http.StatusInternalServerError,
		Severity:       SeverityError,
	},
	ErrInternalDecryption: {
		Code:           ErrInternalDecryption,
		DefaultMessage: "Decryption failed",
		HTTPStatus:     http.StatusInternalServerError,
		Severity:       SeverityError,
	},
	ErrInternalHashing: {
		Code:           ErrInternalHashing,
		DefaultMessage: "Hashing operation failed",
		HTTPStatus:     http.StatusInternalServerError,
		Severity:       SeverityError,
	},
	ErrInternalSerialization: {
		Code:           ErrInternalSerialization,
		DefaultMessage: "Serialization failed",
		HTTPStatus:     http.StatusInternalServerError,
		Severity:       SeverityError,
	},

	// Business Logic
	ErrBusinessRuleViolation: {
		Code:           ErrBusinessRuleViolation,
		DefaultMessage: "Business rule violation",
		HTTPStatus:     http.StatusUnprocessableEntity,
		Severity:       SeverityWarning,
	},
	ErrBusinessInvalidState: {
		Code:           ErrBusinessInvalidState,
		DefaultMessage: "Invalid business state",
		HTTPStatus:     http.StatusConflict,
		Severity:       SeverityWarning,
	},
	ErrBusinessOperationNotAllowed: {
		Code:           ErrBusinessOperationNotAllowed,
		DefaultMessage: "Operation not allowed",
		HTTPStatus:     http.StatusForbidden,
		Severity:       SeverityWarning,
	},

	// Project
	ErrProjectNotFound: {
		Code:           ErrProjectNotFound,
		DefaultMessage: "Project not found",
		HTTPStatus:     http.StatusNotFound,
		Severity:       SeverityInfo,
	},
	ErrProjectNotActive: {
		Code:           ErrProjectNotActive,
		DefaultMessage: "Project is not active",
		HTTPStatus:     http.StatusForbidden,
		Severity:       SeverityWarning,
	},
	ErrProjectFull: {
		Code:           ErrProjectFull,
		DefaultMessage: "Project has reached maximum member capacity",
		HTTPStatus:     http.StatusConflict,
		Severity:       SeverityWarning,
	},
	ErrProjectLimitReached: {
		Code:           ErrProjectLimitReached,
		DefaultMessage: "Organization has reached maximum project limit",
		HTTPStatus:     http.StatusConflict,
		Severity:       SeverityWarning,
	},
	ErrProjectMemberAlreadyExists: {
		Code:           ErrProjectMemberAlreadyExists,
		DefaultMessage: "User is already a member of this project",
		HTTPStatus:     http.StatusConflict,
		Severity:       SeverityWarning,
	},
	ErrProjectMemberNotFound: {
		Code:           ErrProjectMemberNotFound,
		DefaultMessage: "Project member not found",
		HTTPStatus:     http.StatusNotFound,
		Severity:       SeverityInfo,
	},
	ErrProjectMemberNotActive: {
		Code:           ErrProjectMemberNotActive,
		DefaultMessage: "Project member is not active",
		HTTPStatus:     http.StatusForbidden,
		Severity:       SeverityWarning,
	},
	ErrProjectInfoNotLoaded: {
		Code:           ErrProjectInfoNotLoaded,
		DefaultMessage: "Project member information not loaded from database",
		HTTPStatus:     http.StatusInternalServerError,
		Severity:       SeverityError,
	},
	ErrProjectInvitationNotFound: {
		Code:           ErrProjectInvitationNotFound,
		DefaultMessage: "Project invitation not found",
		HTTPStatus:     http.StatusNotFound,
		Severity:       SeverityInfo,
	},
	ErrProjectInvitationExpired: {
		Code:           ErrProjectInvitationExpired,
		DefaultMessage: "Project invitation has expired",
		HTTPStatus:     http.StatusGone,
		Severity:       SeverityInfo,
	},
	ErrProjectInvitationAccepted: {
		Code:           ErrProjectInvitationAccepted,
		DefaultMessage: "Project invitation has already been accepted",
		HTTPStatus:     http.StatusConflict,
		Severity:       SeverityWarning,
	},
	ErrProjectInvitationExists: {
		Code:           ErrProjectInvitationExists,
		DefaultMessage: "A project invitation for this email already exists",
		HTTPStatus:     http.StatusConflict,
		Severity:       SeverityWarning,
	},

	// API Key
	ErrApiKeyNotFound: {
		Code:           ErrApiKeyNotFound,
		DefaultMessage: "API key not found",
		HTTPStatus:     http.StatusNotFound,
		Severity:       SeverityInfo,
	},
	ErrApiKeyRevoked: {
		Code:           ErrApiKeyRevoked,
		DefaultMessage: "API key has been revoked",
		HTTPStatus:     http.StatusUnauthorized,
		Severity:       SeverityWarning,
	},
	ErrApiKeyExpired: {
		Code:           ErrApiKeyExpired,
		DefaultMessage: "API key has expired",
		HTTPStatus:     http.StatusUnauthorized,
		Severity:       SeverityWarning,
	},
	ErrApiKeyInvalid: {
		Code:           ErrApiKeyInvalid,
		DefaultMessage: "Invalid API key",
		HTTPStatus:     http.StatusUnauthorized,
		Severity:       SeverityWarning,
	},
	ErrApiKeyLimitReached: {
		Code:           ErrApiKeyLimitReached,
		DefaultMessage: "Project has reached maximum API key limit",
		HTTPStatus:     http.StatusConflict,
		Severity:       SeverityWarning,
	},
	ErrApiKeyRateLimited: {
		Code:           ErrApiKeyRateLimited,
		DefaultMessage: "API key rate limit exceeded",
		HTTPStatus:     http.StatusTooManyRequests,
		Severity:       SeverityWarning,
	},
}

// GetMetadata returns the metadata for a given error code
func GetMetadata(code ErrorCode) ErrorCodeMetadata {
	if metadata, ok := errorCodeRegistry[code]; ok {
		return metadata
	}
	// Return default metadata for unknown codes
	return ErrorCodeMetadata{
		Code:           code,
		DefaultMessage: "An error occurred",
		HTTPStatus:     http.StatusInternalServerError,
		Severity:       SeverityError,
	}
}

// IsRegistered checks if an error code is registered
func IsRegistered(code ErrorCode) bool {
	_, ok := errorCodeRegistry[code]
	return ok
}
