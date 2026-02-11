package domain

type permission struct {
	Id       string
	resource string
	action   string
	Group    string
	IsSystem bool
}

func NewPermission(id string, resource string, action string, group string, isSystem bool) permission {
	return permission{
		Id:       id,
		resource: resource,
		action:   action,
		Group:    group,
		IsSystem: isSystem,
	}
}

func (p permission) Code() string {
	return p.resource + "." + p.action
}

type PermissionManager interface {
	HasPermission(permission string) bool
	HasAnyPermission(permissions []string) bool
	HasAllPermissions(permissions []string) bool
}

type permissionManager struct {
	permissions map[string]permission
}

func NewPermissionManager(permissions []permission) PermissionManager {
	permissionsMap := make(map[string]permission)
	for _, permission := range permissions {
		permissionsMap[permission.Code()] = permission
	}
	return &permissionManager{permissions: permissionsMap}
}

func (pm *permissionManager) HasPermission(permission string) bool {
	_, ok := pm.permissions[permission]
	return ok
}

func (pm *permissionManager) HasAnyPermission(permissions []string) bool {
	for _, permission := range permissions {
		if pm.HasPermission(permission) {
			return true
		}
	}
	return false
}

func (pm *permissionManager) HasAllPermissions(permissions []string) bool {
	for _, permission := range permissions {
		if !pm.HasPermission(permission) {
			return false
		}
	}
	return true
}
