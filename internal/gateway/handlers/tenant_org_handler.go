package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/handlers/dto"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
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
	router.Use(grpcCtxMiddleware(tenantGRPCTimeout))
	router.Use(middleware.IamJWTVerify(h.iamJWT))

	org := router.Group("/organizations")
	org.
		Post("/onboarding", h.onboarding).
		Post("/", h.create).
		Get("/", h.list)

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
		&dto.CreateOrganization{
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

	return response.SuccessWithMessage(
		c,
		fiber.StatusCreated,
		"organization onboarded successfully",
		resp,
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

	org := resp.GetOrganization()
	orgDto := protoOrgToDTO(org)

	members := resp.GetMembers()
	memberListDto := make([]*dto.TOrganizationMember, len(members))
	for i, m := range members {
		memberListDto[i] = protoMemberToDTO(m)
	}

	h.logger.Debug("Organization fetched by id", resp)

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization fetched successfully",
		&dto.Organization{
			Organization: orgDto,
			Members:      memberListDto,
		},
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

	orgsDto := make([]*dto.TOrganization, len(resp.GetOrganizations()))
	for i, org := range resp.GetOrganizations() {
		orgsDto[i] = protoOrgToDTO(org)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organizations fetched successfully",
		&dto.ListOrganizations{
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

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization updated successfully",
		protoOrgToDTO(resp.GetOrganization()),
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

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization roles fetched successfully",
		resp.GetRoles(),
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

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization permissions fetched successfully",
		resp.GetPermissions(),
	)
}
