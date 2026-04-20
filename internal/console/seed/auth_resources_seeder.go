package seed

import (
	"errors"

	"github.com/vyolayer/vyolayer/internal/console/model"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func SeedAuthResources(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var auth model.Service

		if err := tx.Where("key = ?", "auth").First(&auth).Error; err != nil {
			return errors.New("auth service not found, seed services first")
		}

		/*
			------------------------------------
			Resource: Users
			------------------------------------
		*/
		users := model.ServiceResource{
			ServiceID:   auth.ID,
			Key:         "users",
			Label:       "Users",
			Description: "Manage project users and identities.",
			Icon:        "users",
			Route:       "/projects/:id/users",
			Method:      "GET",
			SortOrder:   1,
			IsVisible:   true,
			Supports: datatypes.JSON([]byte(`{
				"create": true,
				"edit": true,
				"delete": true,
				"search": true,
				"filter": true,
				"pagination": true,
				"export": true
			}`)),
		}

		if err := tx.Where(
			"service_id = ? AND key = ?",
			auth.ID,
			"users",
		).FirstOrCreate(&users).Error; err != nil {
			return err
		}

		userColumns := []model.ServiceResourceColumn{
			{ResourceID: users.ID, Key: "avatar", Label: "Avatar", Type: "image", Visible: true, SortOrder: 1},
			{ResourceID: users.ID, Key: "name", Label: "Name", Type: "text", Sortable: true, Visible: true, SortOrder: 2},
			{ResourceID: users.ID, Key: "email", Label: "Email", Type: "text", Sortable: true, Visible: true, SortOrder: 3},
			{ResourceID: users.ID, Key: "status", Label: "Status", Type: "badge", Sortable: true, Visible: true, SortOrder: 4},
			{ResourceID: users.ID, Key: "last_login", Label: "Last Login", Type: "datetime", Sortable: true, Visible: true, SortOrder: 5},
			{ResourceID: users.ID, Key: "created_at", Label: "Created", Type: "datetime", Sortable: true, Visible: false, SortOrder: 6},
		}

		for _, item := range userColumns {
			if err := tx.Where(
				"resource_id = ? AND key = ?",
				item.ResourceID,
				item.Key,
			).FirstOrCreate(&item).Error; err != nil {
				return err
			}
		}

		userActions := []model.ServiceResourceAction{
			{ResourceID: users.ID, Key: "create", Label: "Add User", Scope: "page", Variant: "primary", Route: "/projects/:id/users", Method: "POST", SortOrder: 1},
			{ResourceID: users.ID, Key: "export", Label: "Export CSV", Scope: "page", Variant: "secondary", Route: "/projects/:id/users/export", Method: "GET", SortOrder: 2},
			{ResourceID: users.ID, Key: "view", Label: "View", Scope: "row", Variant: "secondary", Route: "/projects/:id/users/:userId", Method: "GET", SortOrder: 3},
			{ResourceID: users.ID, Key: "edit", Label: "Edit", Scope: "row", Variant: "secondary", Route: "/projects/:id/users/:userId", Method: "PATCH", SortOrder: 4},
			{ResourceID: users.ID, Key: "delete", Label: "Delete", Scope: "row", Variant: "danger", IsDanger: true, Route: "/projects/:id/users/:userId", Method: "DELETE", SortOrder: 5},
		}

		for _, item := range userActions {
			if err := tx.Where(
				"resource_id = ? AND key = ?",
				item.ResourceID,
				item.Key,
			).FirstOrCreate(&item).Error; err != nil {
				return err
			}
		}

		userFilters := []model.ServiceResourceFilter{
			{
				ResourceID: users.ID,
				Key:        "status",
				Label:      "Status",
				Type:       "select",
				Options:    datatypes.JSON([]byte(`["active","pending","suspended"]`)),
				SortOrder:  1,
			},
			{
				ResourceID: users.ID,
				Key:        "email_verified",
				Label:      "Email Verified",
				Type:       "boolean",
				SortOrder:  2,
			},
		}

		for _, item := range userFilters {
			if err := tx.Where(
				"resource_id = ? AND key = ?",
				item.ResourceID,
				item.Key,
			).FirstOrCreate(&item).Error; err != nil {
				return err
			}
		}

		/*
			------------------------------------
			Resource: Sessions
			------------------------------------
		*/
		sessions := model.ServiceResource{
			ServiceID:   auth.ID,
			Key:         "sessions",
			Label:       "Sessions",
			Description: "View and revoke active sessions.",
			Icon:        "monitor-smartphone",
			Route:       "/projects/:id/sessions",
			Method:      "GET",
			SortOrder:   2,
			IsVisible:   true,
			Supports: datatypes.JSON([]byte(`{
				"create": false,
				"edit": false,
				"delete": true,
				"search": true,
				"filter": true,
				"pagination": true
			}`)),
		}

		if err := tx.Where(
			"service_id = ? AND key = ?",
			auth.ID,
			"sessions",
		).FirstOrCreate(&sessions).Error; err != nil {
			return err
		}

		sessionColumns := []model.ServiceResourceColumn{
			{ResourceID: sessions.ID, Key: "user", Label: "User", Type: "text", Visible: true, SortOrder: 1},
			{ResourceID: sessions.ID, Key: "device", Label: "Device", Type: "text", Visible: true, SortOrder: 2},
			{ResourceID: sessions.ID, Key: "ip", Label: "IP Address", Type: "text", Visible: true, SortOrder: 3},
			{ResourceID: sessions.ID, Key: "created_at", Label: "Created", Type: "datetime", Sortable: true, Visible: true, SortOrder: 4},
			{ResourceID: sessions.ID, Key: "expires_at", Label: "Expires", Type: "datetime", Sortable: true, Visible: true, SortOrder: 5},
			{ResourceID: sessions.ID, Key: "status", Label: "Status", Type: "badge", Visible: true, SortOrder: 6},
		}

		for _, item := range sessionColumns {
			if err := tx.Where(
				"resource_id = ? AND key = ?",
				item.ResourceID,
				item.Key,
			).FirstOrCreate(&item).Error; err != nil {
				return err
			}
		}

		sessionActions := []model.ServiceResourceAction{
			{
				ResourceID: sessions.ID,
				Key:        "revoke_all",
				Label:      "Revoke All Sessions",
				Scope:      "page",
				Variant:    "danger",
				IsDanger:   true,
				Route:      "/projects/:id/sessions/revoke-all",
				Method:     "POST",
				SortOrder:  1,
			},
			{
				ResourceID: sessions.ID,
				Key:        "revoke",
				Label:      "Revoke",
				Scope:      "row",
				Variant:    "danger",
				IsDanger:   true,
				Route:      "/projects/:id/sessions/:sessionId",
				Method:     "DELETE",
				SortOrder:  2,
			},
		}

		for _, item := range sessionActions {
			if err := tx.Where(
				"resource_id = ? AND key = ?",
				item.ResourceID,
				item.Key,
			).FirstOrCreate(&item).Error; err != nil {
				return err
			}
		}

		sessionFilters := []model.ServiceResourceFilter{
			{
				ResourceID: sessions.ID,
				Key:        "status",
				Label:      "Status",
				Type:       "select",
				Options:    datatypes.JSON([]byte(`["active","expired","revoked"]`)),
				SortOrder:  1,
			},
		}

		for _, item := range sessionFilters {
			if err := tx.Where(
				"resource_id = ? AND key = ?",
				item.ResourceID,
				item.Key,
			).FirstOrCreate(&item).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
