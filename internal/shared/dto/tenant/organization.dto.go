package tenant

type Organization struct {
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
