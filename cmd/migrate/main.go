package main

import (
	"worklayer/internal/config"
	"worklayer/internal/platform/database"
	"worklayer/internal/platform/database/models"
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
	)
	if err != nil {
		panic(err)
	}
}
