package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/vyolayer/vyolayer/internal/console/model"
	"github.com/vyolayer/vyolayer/internal/console/repository"
	"github.com/vyolayer/vyolayer/internal/shared/dto/console"
)

type ManifestService interface {
	GetProjectServiceManifest(ctx context.Context, projectID uuid.UUID, serviceKey string) (*console.ServiceManifestWithResourcesDTO, error)
	ListProjectServices(ctx context.Context, projectID uuid.UUID) ([]console.ServiceManifestDTO, error)
}

type manifestService struct {
	projectServiceRepo repository.ProjectServiceRepository
	resourceRepo       repository.ResourceRepository
	overrideRepo       repository.OverrideRepository
}

func NewManifestService(
	projectServiceRepo repository.ProjectServiceRepository,
	resourceRepo repository.ResourceRepository,
	overrideRepo repository.OverrideRepository,
) ManifestService {
	return &manifestService{
		projectServiceRepo: projectServiceRepo,
		resourceRepo:       resourceRepo,
		overrideRepo:       overrideRepo,
	}
}

func (s *manifestService) GetProjectServiceManifest(ctx context.Context, projectID uuid.UUID, serviceKey string) (*console.ServiceManifestWithResourcesDTO, error) {
	// Load active project service
	ps, err := s.projectServiceRepo.GetActiveByProjectAndKey(ctx, projectID, serviceKey)
	if err != nil {
		return nil, err
	}

	// Fetch resources
	resources, err := s.resourceRepo.ListByServiceID(ctx, ps.ServiceID)
	if err != nil {
		return nil, err
	}

	var resourceIDs []uint64
	for _, r := range resources {
		resourceIDs = append(resourceIDs, r.ID)
	}

	// Batch fetch relationships
	columns, err := s.resourceRepo.ListColumns(ctx, resourceIDs)
	if err != nil {
		return nil, err
	}
	actions, err := s.resourceRepo.ListActions(ctx, resourceIDs)
	if err != nil {
		return nil, err
	}
	filters, err := s.resourceRepo.ListFilters(ctx, resourceIDs)
	if err != nil {
		return nil, err
	}

	// // Fetch overrides
	// overrides, err := s.overrideRepo.ListOverrides(ctx, projectID, resourceIDs)
	// if err != nil {
	// 	return nil, err
	// }

	// overrideMap := make(map[uint64]model.ProjectResourceOverride)
	// for _, o := range overrides {
	// 	overrideMap[o.ResourceID] = o
	// }

	// Group relations by resource ID
	colMap := make(map[uint64][]model.ServiceResourceColumn)
	for _, c := range columns {
		colMap[c.ResourceID] = append(colMap[c.ResourceID], c)
	}
	actionMap := make(map[uint64][]model.ServiceResourceAction)
	for _, a := range actions {
		actionMap[a.ResourceID] = append(actionMap[a.ResourceID], a)
	}
	filterMap := make(map[uint64][]model.ServiceResourceFilter)
	for _, f := range filters {
		filterMap[f.ResourceID] = append(filterMap[f.ResourceID], f)
	}

	// Build DTO
	manifestDTO := &console.ServiceManifestWithResourcesDTO{
		ServiceManifestDTO: console.ServiceManifestDTO{
			Key:         ps.Service.Key,
			Name:        ps.Service.Name,
			Description: ps.Service.Description,
			Status:      ps.Status,
			Plan:        ps.Plan,
			Icon:        ps.Service.Icon,
		},
	}

	for _, res := range resources {
		// skip globally hidden
		if !res.IsVisible {
			continue
		}

		// resOverride, hasOverride := overrideMap[res.ID]

		// skip if hidden by override
		// if hasOverride && resOverride.IsVisible != nil && !*resOverride.IsVisible {
		// 	continue
		// }

		resDTO := console.ResourceDTO{
			Key:     res.Key,
			Label:   res.Label,
			Route:   res.Route,
			Icon:    res.Icon,
			Columns: []console.ColumnDTO{},
			Actions: []console.ActionDTO{},
			Filters: []console.FilterDTO{},
		}

		// if hasOverride && resOverride.CustomLabel != "" {
		// 	resDTO.Label = resOverride.CustomLabel
		// }

		// build columns
		colOverrides := make(map[string]map[string]interface{})
		// if hasOverride && len(resOverride.ColumnOverrides) > 0 {
		// 	_ = json.Unmarshal(resOverride.ColumnOverrides, &colOverrides)
		// }
		for _, c := range colMap[res.ID] {
			visible := c.Visible
			label := c.Label

			if cOverride, ok := colOverrides[c.Key]; ok {
				if v, ok := cOverride["visible"].(bool); ok {
					visible = v
				}
				if l, ok := cOverride["label"].(string); ok {
					label = l
				}
			}

			if !visible {
				continue
			}

			resDTO.Columns = append(resDTO.Columns, console.ColumnDTO{
				Key:      c.Key,
				Label:    label,
				Type:     c.Type,
				Sortable: c.Sortable,
				Visible:  visible,
			})
		}

		// build actions
		actionOverrides := make(map[string]map[string]interface{})
		// if hasOverride && len(resOverride.ActionOverrides) > 0 {
		// 	_ = json.Unmarshal(resOverride.ActionOverrides, &actionOverrides)
		// }
		for _, a := range actionMap[res.ID] {
			visible := true
			label := a.Label

			if aOverride, ok := actionOverrides[a.Key]; ok {
				if v, ok := aOverride["disabled"].(bool); ok && v {
					visible = false
				}
				if l, ok := aOverride["label"].(string); ok {
					label = l
				}
			}

			if !visible {
				continue
			}

			resDTO.Actions = append(resDTO.Actions, console.ActionDTO{
				Key:     a.Key,
				Label:   label,
				Scope:   a.Scope,
				Variant: a.Variant,
				Danger:  a.IsDanger,
			})
		}

		// build filters
		for _, f := range filterMap[res.ID] {
			resDTO.Filters = append(resDTO.Filters, console.FilterDTO{
				Key:   f.Key,
				Label: f.Label,
				Type:  f.Type,
			})
		}

		manifestDTO.Resources = append(manifestDTO.Resources, resDTO)
	}

	return manifestDTO, nil
}

func (s *manifestService) ListProjectServices(ctx context.Context, projectID uuid.UUID) ([]console.ServiceManifestDTO, error) {
	pss, err := s.projectServiceRepo.ListActiveByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var results []console.ServiceManifestDTO
	for _, ps := range pss {
		results = append(results, console.ServiceManifestDTO{
			Key:         ps.Service.Key,
			Name:        ps.Service.Name,
			Description: ps.Service.Description,
			Status:      ps.Status,
			Plan:        ps.Plan,
			Icon:        ps.Service.Icon,
		})
	}
	return results, nil
}
