package tenant

type Project struct {
	ID             string `json:"id,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
	Name           string `json:"name,omitempty"`
	Slug           string `json:"slug,omitempty"`
	Description    string `json:"description,omitempty"`
	IsActive       bool   `json:"is_active,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
	MaxAPIKeys     uint32 `json:"max_api_keys,omitempty"`
	MaxMembers     uint32 `json:"max_members,omitempty"`
	MemberCount    uint32 `json:"member_count,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
}

type ProjectMember struct {
	ID        string  `json:"id,omitempty"`
	UserID    string  `json:"user_id,omitempty"`
	Email     string  `json:"email,omitempty"`
	FullName  string  `json:"full_name,omitempty"`
	Role      string  `json:"role,omitempty"`
	IsActive  bool    `json:"is_active,omitempty"`
	JoinedAt  string  `json:"joined_at,omitempty"`
	RemovedAt *string `json:"removed_at,omitempty"`
}

type ProjectResponse struct {
	Project *Project         `json:"project,omitempty"`
	Members []*ProjectMember `json:"members,omitempty"`
}

type ListProjectsResponse struct {
	Projects      []*Project `json:"projects,omitempty"`
	TotalCount    int32      `json:"total_count,omitempty"`
	NextPageToken string     `json:"next_page_token,omitempty"`
}

type ListProjectMembersResponse struct {
	Members       []*ProjectMember `json:"members,omitempty"`
	TotalCount    int32            `json:"total_count,omitempty"`
	NextPageToken string           `json:"next_page_token,omitempty"`
}
