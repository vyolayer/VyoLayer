package seed

import (
	"log"

	"gorm.io/gorm"

	"github.com/vyolayer/vyolayer/internal/console/model"
)

func SeedManifests(db *gorm.DB) error {
	log.Println("Seeding console service manifests...")

	// 1. Ensure 'auth' service exists
	authService := model.Service{
		Key:         "auth",
		Name:        "Authentication",
		Description: "User authentication system",
		Category:    "Core",
		Icon:        "shield",
		Version:     "1.0.0",
		Status:      "ready",
		IsPublic:    true,
		IsInternal:  true,
		SortOrder:   10,
	}

	if err := db.Where(model.Service{Key: "auth"}).FirstOrCreate(&authService).Error; err != nil {
		return err
	}

	// 2. Resources
	usersRes := model.ServiceResource{
		ServiceID:   authService.ID,
		Key:         "users",
		Label:       "Users",
		Description: "Manage authenticated users",
		Icon:        "users",
		Route:       "/projects/:id/users",
		Method:      "GET",
		SortOrder:   10,
		IsVisible:   true,
	}
	db.Where(model.ServiceResource{ServiceID: authService.ID, Key: "users"}).FirstOrCreate(&usersRes)

	sessionsRes := model.ServiceResource{
		ServiceID:   authService.ID,
		Key:         "sessions",
		Label:       "Sessions",
		Description: "Active user sessions",
		Icon:        "key",
		Route:       "/projects/:id/auth/sessions",
		Method:      "GET",
		SortOrder:   20,
		IsVisible:   true,
	}
	db.Where(model.ServiceResource{ServiceID: authService.ID, Key: "sessions"}).FirstOrCreate(&sessionsRes)

	auditLogsRes := model.ServiceResource{
		ServiceID:   authService.ID,
		Key:         "audit_logs",
		Label:       "Audit Logs",
		Description: "Authentication events",
		Icon:        "activity",
		Route:       "/projects/:id/auth/audit_logs",
		Method:      "GET",
		SortOrder:   30,
		IsVisible:   true,
	}
	db.Where(model.ServiceResource{ServiceID: authService.ID, Key: "audit_logs"}).FirstOrCreate(&auditLogsRes)

	// 3. User Columns
	userCols := []model.ServiceResourceColumn{
		{ResourceID: usersRes.ID, Key: "avatar", Label: "Avatar", Type: "image", Visible: true, SortOrder: 10},
		{ResourceID: usersRes.ID, Key: "name", Label: "Name", Type: "text", Sortable: true, Visible: true, SortOrder: 20},
		{ResourceID: usersRes.ID, Key: "email", Label: "Email", Type: "text", Sortable: true, Visible: true, SortOrder: 30},
		{ResourceID: usersRes.ID, Key: "status", Label: "Status", Type: "badge", Sortable: true, Visible: true, SortOrder: 40},
		{ResourceID: usersRes.ID, Key: "last_login", Label: "Last Login", Type: "datetime", Sortable: true, Visible: true, SortOrder: 50},
	}
	for _, col := range userCols {
		db.Where(model.ServiceResourceColumn{ResourceID: usersRes.ID, Key: col.Key}).FirstOrCreate(&col)
	}

	// 4. Sessions Columns
	sessionCols := []model.ServiceResourceColumn{
		{ResourceID: sessionsRes.ID, Key: "user", Label: "User", Type: "text", Sortable: true, Visible: true, SortOrder: 10},
		{ResourceID: sessionsRes.ID, Key: "device", Label: "Device", Type: "text", Sortable: false, Visible: true, SortOrder: 20},
		{ResourceID: sessionsRes.ID, Key: "ip", Label: "IP Address", Type: "text", Sortable: false, Visible: true, SortOrder: 30},
		{ResourceID: sessionsRes.ID, Key: "created_at", Label: "Created At", Type: "datetime", Sortable: true, Visible: true, SortOrder: 40},
		{ResourceID: sessionsRes.ID, Key: "expires_at", Label: "Expires At", Type: "datetime", Sortable: true, Visible: true, SortOrder: 50},
		{ResourceID: sessionsRes.ID, Key: "status", Label: "Status", Type: "badge", Sortable: true, Visible: true, SortOrder: 60},
	}
	for _, col := range sessionCols {
		db.Where(model.ServiceResourceColumn{ResourceID: sessionsRes.ID, Key: col.Key}).FirstOrCreate(&col)
	}

	log.Println("Seeding console service manifests completed.")
	return nil
}
