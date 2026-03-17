package main

import (
	"github.com/vyolayer/vyolayer/internal/config"
	"github.com/vyolayer/vyolayer/internal/platform/database"
	"github.com/vyolayer/vyolayer/internal/platform/database/models"
	servicemodelv1 "github.com/vyolayer/vyolayer/pkg/postgres/models/service/account/v1"
)

func main() {
	cfg, err := config.Load("config/config.dev.yaml")
	if err != nil {
		panic(err)
	}

	db, err := database.Init(&cfg.Database)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	db.Exec("CREATE SCHEMA IF NOT EXISTS account_service;")

	err = db.AutoMigrate(
		// IAM models
		models.User{},
		models.UserSession{},

		// Organization models
		models.Organization{},
		models.OrganizationMember{},
		models.OrganizationMemberInvitation{},

		// Organization RBAC models
		models.OrganizationRole{},
		models.OrganizationPermission{},
		models.OrganizationRolePermission{},
		models.MemberOrganizationRole{},

		// Project models
		models.Project{},
		models.ProjectMember{},
		models.ProjectInvitation{},

		// API Key models
		models.ApiKey{},
		models.ApiKeyUsageLog{},

		// Audit log models
		models.AuditLog{},

		// Account service models
		servicemodelv1.ServiceUser{},
		servicemodelv1.ServiceUserAvatar{},
		servicemodelv1.ServiceUserSession{},
		servicemodelv1.ServiceUserVerificationToken{},
		servicemodelv1.ServiceUserLoginAttempt{},
		servicemodelv1.ServiceUserAccountLock{},
	)
	if err != nil {
		panic(err)
	}
}
