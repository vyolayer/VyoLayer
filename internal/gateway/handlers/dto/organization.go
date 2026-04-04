package dto

type TOrganization struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	IsActive     bool   `json:"is_active"`
	OwnerID      string `json:"owner_id"`
	MaxMembers   uint32 `json:"max_members"`
	MaxProjects  uint32 `json:"max_projects"`
	ProjectCount uint32 `json:"project_count"`
	MemberCount  uint32 `json:"member_count"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type TOrganizationMember struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id,omitempty"`
	UserID         string   `json:"user_id"`
	FullName       string   `json:"full_name"`
	Email          string   `json:"email"`
	Roles          []string `json:"roles"`
	Status         string   `json:"status"`
	JoinedAt       string   `json:"joined_at"`
}

type TOrganizationRole struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsSystemRole bool   `json:"is_system_role"`
	IsDefault    bool   `json:"is_default"`
}

type TOrganizationInvitation struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id"`
	Email          string   `json:"email"`
	RoleIDs        []string `json:"role_ids"`
	InvitedBy      string   `json:"invited_by"`
	InvitedAt      string   `json:"invited_at"`
	IsAccepted     bool     `json:"is_accepted"`
	AcceptedAt     string   `json:"accepted_at,omitempty"`
	ExpiredAt      string   `json:"expired_at"`
	IsPending      bool     `json:"is_pending"`
}

type TInvitedBy struct {
	MemberID string `json:"member_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type TOrganizationInvitationForOrg struct {
	Invitation *TOrganizationInvitation `json:"invitation"`
	InvitedBy  *TInvitedBy              `json:"invited_by"`
}

type CreateOrganization struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Organization struct {
	Organization *TOrganization         `json:"organization"`
	Members      []*TOrganizationMember `json:"members"`
}

type ListOrganizations struct {
	Organizations []*TOrganization `json:"organizations"`
	TotalCount    int32            `json:"total_count"`
	NextPageToken string           `json:"next_page_token"`
}

type ListOrganizationMembers struct {
	Members    []*TOrganizationMember `json:"members"`
	TotalCount int32                  `json:"total_count"`
}

type ListOrganizationInvitations struct {
	Invitations []*TOrganizationInvitation `json:"invitations"`
}

type ListOrganizationInvitationsForOrg struct {
	Invitations []*TOrganizationInvitationForOrg `json:"invitations"`
}
