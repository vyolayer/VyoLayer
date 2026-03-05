package repository

import "gorm.io/gorm"

type Registry struct {
	User                         UserRepository
	Session                      SessionRepository
	Organization                 OrganizationRepository
	OrganizationMember           OrganizationMemberRepository
	OrganizationMemberInvitation OrganizationMemberInvitationRepository
	OrganizationRBAC             OrganizationRBACRepository
	AuditLog                     AuditLogRepository
	Project                      ProjectRepository
	ProjectMember                ProjectMemberRepository
	ProjectInvitation            ProjectInvitationRepository
	ApiKey                       ApiKeyRepository
	ApiKeyUsageLog               ApiKeyUsageLogRepository
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		User:                         NewUserRepository(db),
		Session:                      NewSessionRepository(db),
		Organization:                 NewOrganizationRepository(db),
		OrganizationMember:           NewOrganizationMemberRepository(db),
		OrganizationMemberInvitation: NewOrganizationMemberInvitationRepository(db),
		OrganizationRBAC:             NewOrganizationRBACRepository(db),
		AuditLog:                     NewAuditLogRepository(db),
		Project:                      NewProjectRepository(db),
		ProjectMember:                NewProjectMemberRepository(db),
		ProjectInvitation:            NewProjectInvitationRepository(db),
		ApiKey:                       NewApiKeyRepository(db),
		ApiKeyUsageLog:               NewApiKeyUsageLogRepository(db),
	}
}
