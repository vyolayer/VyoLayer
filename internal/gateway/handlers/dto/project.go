package dto

// ─── Project ──────────────────────────────────────────────────────────────────

type TProject struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Description    string `json:"description"`
	IsActive       bool   `json:"is_active"`
	CreatedBy      string `json:"created_by"`
	MaxAPIKeys     uint32 `json:"max_api_keys"`
	MaxMembers     uint32 `json:"max_members"`
	MemberCount    uint32 `json:"member_count"`
	CreatedAt      string `json:"created_at"`
}

type TProjectMember struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Email     string  `json:"email"`
	FullName  string  `json:"full_name"`
	Role      string  `json:"role"`
	IsActive  bool    `json:"is_active"`
	JoinedAt  string  `json:"joined_at"`
	RemovedAt *string `json:"removed_at,omitempty"`
}

type ProjectResponse struct {
	Project *TProject         `json:"project"`
	Members []*TProjectMember `json:"members,omitempty"`
}

type ListProjects struct {
	Projects      []*TProject `json:"projects"`
	TotalCount    int32       `json:"total_count"`
	NextPageToken string      `json:"next_page_token,omitempty"`
}

type ListProjectMembers struct {
	Members       []*TProjectMember `json:"members"`
	TotalCount    int32             `json:"total_count"`
	NextPageToken string            `json:"next_page_token,omitempty"`
}
