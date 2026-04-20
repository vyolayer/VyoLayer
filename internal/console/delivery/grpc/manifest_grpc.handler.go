package grpc

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/vyolayer/vyolayer/internal/console/service"
	consolev1 "github.com/vyolayer/vyolayer/proto/console/v1"
)

type ManifestServer struct {
	consolev1.UnimplementedProjectServiceManifestServer
	manifestService service.ManifestService
}

func NewManifestServer(manifestService service.ManifestService) *ManifestServer {
	return &ManifestServer{
		manifestService: manifestService,
	}
}

func (s *ManifestServer) GetProjectServiceManifest(
	ctx context.Context,
	req *consolev1.GetProjectServiceManifestRequest,
) (*consolev1.GetProjectServiceManifestResponse, error) {
	projectID, err := uuid.Parse(req.ProjectId)
	if err != nil {
		log.Printf("invalid project_id: %v", err)
		return nil, err
	}

	dto, err := s.manifestService.GetProjectServiceManifest(ctx, projectID, req.ServiceKey)
	if err != nil {
		log.Printf("error getting manifest: %v", err)
		return nil, err
	}

	// map DTO to proto
	data := &consolev1.ManifestData{
		Key:         dto.Key,
		Name:        dto.Name,
		Description: dto.Description,
		Status:      dto.Status,
		Plan:        dto.Plan,
		Icon:        dto.Icon,
	}

	for _, res := range dto.Resources {
		protoRes := &consolev1.ResourceData{
			Key:   res.Key,
			Label: res.Label,
			Route: res.Route,
			Icon:  res.Icon,
		}

		for _, col := range res.Columns {
			protoRes.Columns = append(protoRes.Columns, &consolev1.ColumnData{
				Key:      col.Key,
				Label:    col.Label,
				Type:     col.Type,
				Sortable: col.Sortable,
				Visible:  col.Visible,
			})
		}
		for _, act := range res.Actions {
			protoRes.Actions = append(protoRes.Actions, &consolev1.ActionData{
				Key:     act.Key,
				Label:   act.Label,
				Scope:   act.Scope,
				Variant: act.Variant,
				Danger:  act.Danger,
			})
		}
		for _, fil := range res.Filters {
			protoRes.Filters = append(protoRes.Filters, &consolev1.FilterData{
				Key:   fil.Key,
				Label: fil.Label,
				Type:  fil.Type,
			})
		}

		data.Resources = append(data.Resources, protoRes)
	}

	return &consolev1.GetProjectServiceManifestResponse{
		Success: true,
		Data:    data,
	}, nil
}

func (s *ManifestServer) ListProjectServices(
	ctx context.Context,
	req *consolev1.ListProjectServicesRequest,
) (*consolev1.ListProjectServicesResponse, error) {
	projectID, err := uuid.Parse(req.ProjectId)
	if err != nil {
		log.Printf("invalid project_id: %v", err)
		return nil, err
	}

	dtos, err := s.manifestService.ListProjectServices(ctx, projectID)
	if err != nil {
		log.Printf("error listing services: %v", err)
		return nil, err
	}

	res := &consolev1.ListProjectServicesResponse{
		Success: true,
	}

	for _, dto := range dtos {
		res.Data = append(res.Data, &consolev1.ProjectServiceData{
			Key:         dto.Key,
			Name:        dto.Name,
			Status:      dto.Status,
			Plan:        dto.Plan,
			Icon:        dto.Icon,
			Description: dto.Description,
		})
	}

	return res, nil
}
