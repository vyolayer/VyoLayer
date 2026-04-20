package console

type ColumnDTO struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Type     string `json:"type"`
	Sortable bool   `json:"sortable"`
	Visible  bool   `json:"visible"`
}

type ActionDTO struct {
	Key     string `json:"key"`
	Label   string `json:"label"`
	Scope   string `json:"scope"`
	Variant string `json:"variant"`
	Danger  bool   `json:"danger"`
}

type FilterDTO struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type ResourceDTO struct {
	Key     string      `json:"key"`
	Label   string      `json:"label"`
	Route   string      `json:"route"`
	Icon    string      `json:"icon"`
	Columns []ColumnDTO `json:"columns"`
	Actions []ActionDTO `json:"actions"`
	Filters []FilterDTO `json:"filters"`
}

type ServiceManifestDTO struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Plan        string `json:"plan"`
	Icon        string `json:"icon"`
}

type ServiceManifestWithResourcesDTO struct {
	ServiceManifestDTO
	Resources []ResourceDTO `json:"resources"`
}
