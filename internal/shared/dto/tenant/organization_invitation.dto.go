package tenant

type OrganizationInvitation struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id"`
	Email          string   `json:"email"`
	RoleIDs        []string `json:"role_ids"`
	InvitedBy      string   `json:"invited_by"`
	InvitedAt      string   `json:"invited_at"`
	IsAccepted     bool     `json:"is_accepted"`
	AcceptedAt     string   `json:"accepted_at"`
	ExpiredAt      string   `json:"expired_at"`
	IsPending      bool     `json:"is_pending"`
}

type InvitedBy struct {
	MemberID string `json:"member_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type OrganizationInvitationForOrg struct {
	Invitation *OrganizationInvitation `json:"invitation"`
	InvitedBy  *InvitedBy              `json:"invited_by"`
}
