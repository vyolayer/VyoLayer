package tenant

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	dto "github.com/vyolayer/vyolayer/internal/shared/dto/tenant"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/response"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

const (
	tenantGRPCTimeout = 10 * time.Second
)

var (
	ErrInvalidBody  = errors.BadRequest("invalid request body")
	ErrInvalidOrgID = errors.BadRequest("invalid organization id")
	ErrInvalidSlug  = errors.BadRequest("invalid slug")
)

type OrganizationHandler struct {
	logger *logger.AppLogger
	client tenantV1.OrganizationServiceClient
	iamJWT jwt.IamJWT
}

func NewOrganizationHandler(
	logger *logger.AppLogger,
	client tenantV1.OrganizationServiceClient,
	iamJWT jwt.IamJWT,
) *OrganizationHandler {
	return &OrganizationHandler{
		logger: logger.WithContext("Org Handler"),
		client: client,
		iamJWT: iamJWT,
	}
}

func (h *OrganizationHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(tenantGRPCTimeout).Handler()

	org := router.Group("/organizations")
	org.Use(grpcCtxMiddleware)
	org.Use(middleware.IamJWTVerify(h.iamJWT))

	org.
		Post("/onboarding", h.onboarding).
		Post("/", h.create).
		Get("/", h.list)

	org.Get("slug/:slug", h.getBySlug)

	// All routes below require a valid organizationID in the path
	orgGroup := org.Group("/:organizationID", middleware.ValidateOrganizationID())

	// Organization lifecycle
	orgGroup.
		Get("/", h.getById).
		Patch("/", h.update).
		Delete("/", h.delete).
		Delete("/archive", h.archive).
		Post("/restore", h.restore).
		Post("/transfer-ownership", h.transferOwnership)

	// Roles and permissions
	orgGroup.Get("/roles", h.listRoles).
		Get("/permissions", h.listPermissions)

	h.logger.Info("Organization routes registered", "")
}

// ─── Organization ────────────────────────────────────────────────────────────

func (h *OrganizationHandler) create(c *fiber.Ctx) error {
	var req tenantV1.CreateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}

	resp, err := h.client.CreateOrganization(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	h.logger.Debug("Organization created", resp)

	return response.SuccessWithMessage(
		c,
		fiber.StatusCreated,
		"organization created successfully",
		&dto.CreateOrganizationResponse{
			Name:        req.GetName(),
			Description: req.GetDescription(),
		},
	)
}

// onboarding is called only when the user is not yet a member of any org.
func (h *OrganizationHandler) onboarding(c *fiber.Ctx) error {
	var req tenantV1.CreateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}

	resp, err := h.client.OnboardOrganization(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	h.logger.Debug("Organization onboarded", resp)

	dtoResp := protoOrgResponseToDTO(resp)

	return response.SuccessWithMessage(
		c,
		fiber.StatusCreated,
		"organization onboarded successfully",
		&dto.OnboardOrganizationResponse{
			Organization: dtoResp.Organization,
			Members:      dtoResp.Members,
		},
	)
}

func (h *OrganizationHandler) getById(c *fiber.Ctx) error {
	req := tenantV1.TenantOrganizationIDRequest{
		OrganizationId: getOrgIDFromLocals(c),
	}

	resp, err := h.client.GetOrganizationById(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	h.logger.Debug("Organization fetched by id", resp)
	orgDto := protoOrgResponseToDTO(resp)

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization fetched successfully",
		orgDto,
	)
}

func (h *OrganizationHandler) getBySlug(c *fiber.Ctx) error {
	var (
		slug string
		in   tenantV1.OrganizationSlugRequest
	)

	slug = c.Params("slug")
	if slug == "" {
		return response.Error(c, ErrInvalidSlug)
	}
	in.Slug = slug

	resp, err := h.client.GetBySlug(c.UserContext(), &in)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	h.logger.Debug("Organization fetched by slug", resp)
	orgDto := protoOrgResponseToDTO(resp)

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization fetched successfully",
		orgDto,
	)
}

func (h *OrganizationHandler) list(c *fiber.Ctx) error {
	req := tenantV1.ListOrganizationsRequest{
		PageSize:  int32(c.QueryInt("page_size", 0)),
		PageToken: c.Query("page_token", ""),
	}

	resp, err := h.client.ListOrganizations(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	orgsDto := make([]*dto.Organization, len(resp.GetOrganizations()))
	for i, org := range resp.GetOrganizations() {
		orgsDto[i] = protoOrgToDTO(org)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organizations fetched successfully",
		&dto.ListOrganizationsResponse{
			Organizations: orgsDto,
			TotalCount:    resp.GetTotalCount(),
			NextPageToken: resp.GetNextPageToken(),
		},
	)
}

func (h *OrganizationHandler) update(c *fiber.Ctx) error {
	var req tenantV1.UpdateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.UpdateOrganization(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	orgDto := protoOrgResponseToDTO(resp)

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization updated successfully",
		orgDto,
	)
}

func (h *OrganizationHandler) archive(c *fiber.Ctx) error {
	var req tenantV1.ArchiveOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.ArchiveOrganization(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

func (h *OrganizationHandler) restore(c *fiber.Ctx) error {
	req := tenantV1.TenantOrganizationIDRequest{
		OrganizationId: getOrgIDFromLocals(c),
	}

	resp, err := h.client.RestoreOrganization(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

func (h *OrganizationHandler) delete(c *fiber.Ctx) error {
	var req tenantV1.DeleteOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.DeleteOrganization(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

func (h *OrganizationHandler) transferOwnership(c *fiber.Ctx) error {
	var req tenantV1.TransferOwnershipRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.TransferOwnership(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

func (h *OrganizationHandler) listRoles(c *fiber.Ctx) error {
	var req tenantV1.TenantOrganizationIDRequest
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.GetAllRoles(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	rolesDto := make([]*dto.OrganizationRole, len(resp.GetRoles()))
	for i, r := range resp.GetRoles() {
		rolesDto[i] = protoOrgRoleToDTO(r)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization roles fetched successfully",
		rolesDto,
	)
}

// Get all permissions
func (h *OrganizationHandler) listPermissions(c *fiber.Ctx) error {
	var req tenantV1.TenantOrganizationIDRequest
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.GetAllPermissions(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	permsDto := make([]*dto.OrganizationPerm, len(resp.GetPermissions()))
	for i, p := range resp.GetPermissions() {
		permsDto[i] = protoPermToDTO(p)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization permissions fetched successfully",
		permsDto,
	)
}
