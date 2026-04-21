package tenant

type CreateOrganizationResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OnboardOrganizationResponse struct {
	Organization *Organization         `json:"organization"`
	Members      []*OrganizationMember `json:"members,omitempty"`
}

type OrganizationDetailResponse struct {
	Organization *Organization         `json:"organization"`
	Members      []*OrganizationMember `json:"members,omitempty"`
}

type ListOrganizationsResponse struct {
	Organizations []*Organization `json:"organizations"`
	TotalCount    int32           `json:"total_count"`
	NextPageToken string          `json:"next_page_token"`
}

type OrganizationRolesResponse struct {
	Roles []*OrganizationRole
}

type OrganizationPermissionsResponse struct {
	Permissions []*OrganizationPerm
}

type ListOrganizationInvitationsResponse struct {
	Invitations []*OrganizationInvitation `json:"invitations"`
}

type ListOrganizationInvitationsForOrgResponse struct {
	Invitations []*OrganizationInvitationForOrg `json:"invitations"`
}

type OrganizationMemberWithRBACResponse struct {
	OrganizationMember
	Roles []string `json:"roles"`
	Perms []string `json:"perms"`
}

type ListOrganizationMembersResponse struct {
	Members    []*OrganizationMember `json:"members"`
	TotalCount int32                 `json:"total_count"`
}
