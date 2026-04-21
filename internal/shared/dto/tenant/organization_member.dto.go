package tenant

type OrganizationMember struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id"`
	UserID         string   `json:"user_id"`
	FullName       string   `json:"full_name"`
	Email          string   `json:"email"`
	Roles          []string `json:"roles"`
	Status         string   `json:"status"`
	JoinedAt       string   `json:"joined_at"`
	InvitedAt      string   `json:"invited_at"`
	InvitedBy      string   `json:"invited_by"`
	DeactivatedBy  string   `json:"deactivated_by"`
	DeactivatedAt  string   `json:"deactivated_at"`
}

type OrganizationMemberWithRBAC struct {
	OrganizationMember
	Roles []string `json:"roles"`
	Perms []string `json:"perms"`
}
