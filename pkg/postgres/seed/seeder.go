package seed

import (
	"fmt"

	"github.com/google/uuid"
	tenantmodelv1 "github.com/vyolayer/vyolayer/internal/tenant/models/v1"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Run(db *gorm.DB) error {
	// 1. Seed Permissions
	fmt.Println("🌱 Seeding Permissions...")
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&OrgPermissions).Error; err != nil {
		return err
	}

	// 2. Seed Roles (System Global)
	fmt.Println("🌱 Seeding Global Roles...")
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&OrgRoles).Error; err != nil {
		return err
	}

	// 3. Map Permissions to Roles
	fmt.Println("🔗 Linking Permissions to Roles...")
	return mapPermissions(db)
}

func mapPermissions(db *gorm.DB) error {
	// A. Fetch freshly created records to get their UUIDs
	var roles []tenantmodelv1.OrganizationRole
	var perms []tenantmodelv1.OrganizationPermission

	// Get only system roles/perms to be safe
	db.Where("is_system = ?", true).Find(&roles)
	db.Where("is_system = ?", true).Find(&perms)

	// B. Create Helper Maps
	roleMap := make(map[string]uuid.UUID)
	for _, r := range roles {
		roleMap[r.Name] = r.ID // Assuming ID field handles value retrieval
	}

	// C. Define Logic
	var links []tenantmodelv1.OrganizationRolePermission

	for _, p := range perms {
		// --- 1. OWNER & ADMIN (Get Everything) ---
		links = append(links,
			tenantmodelv1.OrganizationRolePermission{RoleID: roleMap[RoleOwner], PermissionID: p.ID},
			tenantmodelv1.OrganizationRolePermission{RoleID: roleMap[RoleAdmin], PermissionID: p.ID},
		)

		// Audit & Roles are not accessible to any role
		if p.Resource == ResourceAudit || p.Resource == ResourceRole {
			continue
		}

		// --- 2. VIEWER (Read Only) ---
		if isReadAction(p.Action) {
			links = append(links,
				tenantmodelv1.OrganizationRolePermission{RoleID: roleMap[RoleViewer], PermissionID: p.ID},
			)
		}

		// --- 3. MEMBER (Read + Work) ---
		// Rule A: Full Access to Projects
		if p.Resource == ResourceProject {
			links = append(links,
				tenantmodelv1.OrganizationRolePermission{RoleID: roleMap[RoleMember], PermissionID: p.ID},
			)
			continue // Done for this permission
		}

		// Rule B: Read-Only Access to Organization & Team
		// (Allows: List Members, View Org Details)
		// (Blocks: Invite Member, Remove Member, Update Org)
		if (p.Resource == ResourceOrganization || p.Resource == ResourceMember) && isReadAction(p.Action) {
			links = append(links,
				tenantmodelv1.OrganizationRolePermission{RoleID: roleMap[RoleMember], PermissionID: p.ID},
			)
		}
	}

	// D. Batch Insert
	if len(links) > 0 {
		return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&links).Error
	}

	return nil
}

func isReadAction(action string) bool {
	return action == ActionRead || action == ActionView || action == ActionList
}
