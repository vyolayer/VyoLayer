package main

import (
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
	)
	if err != nil {
		panic(err)
	}
}
