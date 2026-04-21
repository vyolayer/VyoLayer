package tenant

type OrganizationRole struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsSystemRole bool   `json:"is_system_role"`
	IsDefault    bool   `json:"is_default"`
}

type OrganizationPerm struct {
	ID           string `json:"id"`
	Resource     string `json:"resource"`
	Action       string `json:"action"`
	Code         string `json:"code"`
	Group        string `json:"group"`
	IsSystemPerm bool   `json:"is_system_perm"`
}
