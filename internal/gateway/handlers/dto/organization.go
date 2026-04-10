package dto

type TOrganization struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Slug         string `json:"slug,omitempty"`
	Description  string `json:"description,omitempty"`
	IsActive     bool   `json:"is_active,omitempty"`
	OwnerID      string `json:"owner_id,omitempty"`
	MaxMembers   uint32 `json:"max_members,omitempty"`
	MaxProjects  uint32 `json:"max_projects,omitempty"`
	ProjectCount uint32 `json:"project_count,omitempty"`
	MemberCount  uint32 `json:"member_count,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type TOrganizationMember struct {
	ID             string   `json:"id,omitempty"`
	OrganizationID string   `json:"organization_id,omitempty"`
	UserID         string   `json:"user_id,omitempty"`
	FullName       string   `json:"full_name,omitempty"`
	Email          string   `json:"email,omitempty"`
	Roles          []string `json:"roles,omitempty"`
	Status         string   `json:"status,omitempty"`
	JoinedAt       string   `json:"joined_at,omitempty"`
	InvitedAt      string   `json:"invited_at,omitempty"`
	InvitedBy      string   `json:"invited_by,omitempty"`
	DeactivatedBy  string   `json:"deactivated_by,omitempty"`
	DeactivatedAt  string   `json:"deactivated_at,omitempty"`
}

type TOrganizationRole struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	IsSystemRole bool   `json:"is_system_role,omitempty"`
	IsDefault    bool   `json:"is_default,omitempty"`
}

type TOrganizationInvitation struct {
	ID             string   `json:"id,omitempty"`
	OrganizationID string   `json:"organization_id,omitempty"`
	Email          string   `json:"email,omitempty"`
	RoleIDs        []string `json:"role_ids,omitempty"`
	InvitedBy      string   `json:"invited_by,omitempty"`
	InvitedAt      string   `json:"invited_at,omitempty"`
	IsAccepted     bool     `json:"is_accepted,omitempty"`
	AcceptedAt     string   `json:"accepted_at,omitempty"`
	ExpiredAt      string   `json:"expired_at,omitempty"`
	IsPending      bool     `json:"is_pending,omitempty"`
}

type TInvitedBy struct {
	MemberID string `json:"member_id,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty"`
}

type TOrganizationInvitationForOrg struct {
	Invitation *TOrganizationInvitation `json:"invitation,omitempty"`
	InvitedBy  *TInvitedBy              `json:"invited_by,omitempty"`
}

type CreateOrganization struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Organization struct {
	Organization *TOrganization         `json:"organization,omitempty"`
	Members      []*TOrganizationMember `json:"members,omitempty"`
}

type ListOrganizations struct {
	Organizations []*TOrganization `json:"organizations,omitempty"`
	TotalCount    int32            `json:"total_count,omitempty"`
	NextPageToken string           `json:"next_page_token,omitempty"`
}

type ListOrganizationMembers struct {
	Members    []*TOrganizationMember `json:"members,omitempty"`
	TotalCount int32                  `json:"total_count,omitempty"`
}

type ListOrganizationInvitations struct {
	Invitations []*TOrganizationInvitation `json:"invitations,omitempty"`
}

type ListOrganizationInvitationsForOrg struct {
	Invitations []*TOrganizationInvitationForOrg `json:"invitations,omitempty"`
}
