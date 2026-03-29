package main

import (
	accountmodelv1 "github.com/vyolayer/vyolayer/internal/account/models/v1"
	iammodelv1 "github.com/vyolayer/vyolayer/internal/iam/models/v1"

	"github.com/vyolayer/vyolayer/internal/config"
	"github.com/vyolayer/vyolayer/internal/platform/database"
	"github.com/vyolayer/vyolayer/internal/platform/database/models"
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
	db.Exec("CREATE SCHEMA IF NOT EXISTS iam;")

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

		// IAM models
		iammodelv1.Avatar{},
		iammodelv1.User{},
		iammodelv1.Session{},
		iammodelv1.VerificationToken{},
		iammodelv1.PasswordResetToken{},

		// Account service models
		accountmodelv1.ServiceUser{},
		accountmodelv1.ServiceUserAvatar{},
		accountmodelv1.ServiceUserSession{},
		accountmodelv1.ServiceUserVerificationToken{},
		accountmodelv1.ServiceUserLoginAttempt{},
		accountmodelv1.ServiceUserAccountLock{},
	)
	if err != nil {
		panic(err)
	}
}
