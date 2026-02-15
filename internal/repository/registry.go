package repository

import "gorm.io/gorm"

type Registry struct {
	User                         UserRepository
	Session                      SessionRepository
	Organization                 OrganizationRepository
	OrganizationMember           OrganizationMemberRepository
	OrganizationMemberInvitation OrganizationMemberInvitationRepository
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		User:                         NewUserRepository(db),
		Session:                      NewSessionRepository(db),
		Organization:                 NewOrganizationRepository(db),
		OrganizationMember:           NewOrganizationMemberRepository(db),
		OrganizationMemberInvitation: NewOrganizationMemberInvitationRepository(db),
	}
}
