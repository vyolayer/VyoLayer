package repository

import "gorm.io/gorm"

type Registry struct {
	User    UserRepository
	Session SessionRepository
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		User:    NewUserRepository(db),
		Session: NewSessionRepository(db),
	}
}
