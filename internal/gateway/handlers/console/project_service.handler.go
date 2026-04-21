package console

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/response"
	consolev1 "github.com/vyolayer/vyolayer/proto/console/v1"

	dto "github.com/vyolayer/vyolayer/internal/shared/dto/console"
)

const (
	grpcTimeout = 10 * time.Second
)

type ProjectServiceHandler struct {
	logger *logger.AppLogger
	client consolev1.ProjectServiceManifestClient
	iamJWT jwt.IamJWT
}

func NewProjectServiceHandler(
	logger *logger.AppLogger,
	client consolev1.ProjectServiceManifestClient,
	iamJWT jwt.IamJWT,
) *ProjectServiceHandler {
	return &ProjectServiceHandler{
		logger: logger.WithContext("ConsoleProjectServiceHandler"),
		client: client,
		iamJWT: iamJWT,
	}
}

func (h *ProjectServiceHandler) RegisterRoutes(router fiber.Router) {
	// grpc ctx timeout
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(grpcTimeout)

	// /console/projects/:projectID/services
	services := router.Group("/console/projects/:projectID/services")
	services.Use(grpcCtxMiddleware.Handler())
	services.Use(middleware.IamJWTVerify(h.iamJWT))
	services.Use(middleware.ValidateProjectID())

	services.Get("/", h.list)
	services.Get("/:serviceKey", h.get)

	h.logger.Info("ProjectService routes registered", "")
}

func (h *ProjectServiceHandler) list(c *fiber.Ctx) error {
	req := &consolev1.ListProjectServicesRequest{
		ProjectId: getProjectIDFromLocals(c),
	}

	resp, err := h.client.ListProjectServices(c.UserContext(), req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	servicesDTO := make([]dto.ServiceManifestDTO, len(resp.GetData()))
	for i, service := range resp.GetData() {
		servicesDTO[i] = dto.ServiceManifestDTO{
			Key:         service.GetKey(),
			Name:        service.GetName(),
			Status:      service.GetStatus(),
			Plan:        service.GetPlan(),
			Icon:        service.GetIcon(),
			Description: service.GetDescription(),
		}
	}

	return response.Success(c, servicesDTO)
}

func (h *ProjectServiceHandler) get(c *fiber.Ctx) error {
	req := &consolev1.GetProjectServiceManifestRequest{
		ProjectId:  getProjectIDFromLocals(c),
		ServiceKey: c.Params("serviceKey"),
	}

	grpcResp, err := h.client.GetProjectServiceManifest(c.UserContext(), req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	grpcData := grpcResp.GetData()
	if grpcData == nil {
		return response.Error(c, errors.NotFound("service not found"))
	}

	var resourcesDTO []dto.ResourceDTO
	for _, res := range grpcData.GetResources() {
		var columnsDTO []dto.ColumnDTO
		for _, col := range res.GetColumns() {
			columnsDTO = append(columnsDTO, dto.ColumnDTO{
				Key:      col.GetKey(),
				Label:    col.GetLabel(),
				Type:     col.GetType(),
				Sortable: col.GetSortable(),
				Visible:  col.GetVisible(),
			})
		}

		var actionsDTO []dto.ActionDTO
		for _, act := range res.GetActions() {
			actionsDTO = append(actionsDTO, dto.ActionDTO{
				Key:     act.GetKey(),
				Label:   act.GetLabel(),
				Scope:   act.GetScope(),
				Variant: act.GetVariant(),
				Danger:  act.GetDanger(),
			})
		}

		var filtersDTO []dto.FilterDTO
		for _, fil := range res.GetFilters() {
			filtersDTO = append(filtersDTO, dto.FilterDTO{
				Key:   fil.GetKey(),
				Label: fil.GetLabel(),
				Type:  fil.GetType(),
			})
		}

		resourcesDTO = append(resourcesDTO, dto.ResourceDTO{
			Key:     res.GetKey(),
			Label:   res.GetLabel(),
			Route:   res.GetRoute(),
			Icon:    res.GetIcon(),
			Columns: columnsDTO,
			Actions: actionsDTO,
			Filters: filtersDTO,
		})
	}

	serviceManifestDTO := &dto.ServiceManifestWithResourcesDTO{
		ServiceManifestDTO: dto.ServiceManifestDTO{
			Key:         grpcData.GetKey(),
			Name:        grpcData.GetName(),
			Description: grpcData.GetDescription(),
			Status:      grpcData.GetStatus(),
			Plan:        grpcData.GetPlan(),
			Icon:        grpcData.GetIcon(),
		},
		Resources: resourcesDTO,
	}

	return response.Success(c, serviceManifestDTO)
}
