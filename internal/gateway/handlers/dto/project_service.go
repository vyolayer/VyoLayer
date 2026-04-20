package dto

type TProjectService struct {
	ID                 uint64         `json:"id,omitempty"`
	ServiceKey         string         `json:"service_key,omitempty"`
	ServiceName        string         `json:"service_name,omitempty"`
	ServiceDescription string         `json:"service_description,omitempty"`
	ProjectID          string         `json:"project_id,omitempty"`
	ServiceID          uint64         `json:"service_id,omitempty"`
	Status             string         `json:"status,omitempty"`
	Plan               string         `json:"plan,omitempty"`
	Config             map[string]any `json:"config,omitempty"`
	EnabledAt          string         `json:"enabled_at,omitempty"`
	SuspendedAt        string         `json:"suspended_at,omitempty"`
	CreatedAt          string         `json:"created_at,omitempty"`
	UpdatedAt          string         `json:"updated_at,omitempty"`
}

type ProjectServiceResponse struct {
	Service *TProjectService `json:"service,omitempty"`
}

type ListProjectServices struct {
	Services []*TProjectService `json:"services,omitempty"`
}

type EnableProjectServiceRequest struct {
	ServiceKey string `json:"service_key"`
	Plan       string `json:"plan"`
}

type UpdateProjectServiceConfigRequest struct {
	Config map[string]any `json:"config"`
}
