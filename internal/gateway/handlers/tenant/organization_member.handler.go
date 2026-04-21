package tenant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	dto "github.com/vyolayer/vyolayer/internal/shared/dto/tenant"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/response"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

type OrganizationMemberHandler struct {
	logger *logger.AppLogger
	client tenantV1.OrganizationMemberServiceClient
	iamJWT jwt.IamJWT
}

func NewOrganizationMemberHandler(
	logger *logger.AppLogger,
	client tenantV1.OrganizationMemberServiceClient,
	iamJWT jwt.IamJWT,
) *OrganizationMemberHandler {
	return &OrganizationMemberHandler{
		logger: logger.WithContext("Org Member Handler"),
		client: client,
		iamJWT: iamJWT,
	}
}

func (h *OrganizationMemberHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(tenantGRPCTimeout).Handler()

	orgMemberGroup := router.Group("/organizations/:organizationID/members")
	orgMemberGroup.Use(
		grpcCtxMiddleware,
		middleware.IamJWTVerify(h.iamJWT),
		middleware.ValidateOrganizationID(),
	)

	// Members
	orgMemberGroup.
		Get("/", h.listMembers).
		Get("/me", h.getCurrentMember).
		Get("/:memberID", h.getMemberByID).
		Delete("/:memberID", h.removeMember)

	h.logger.Info("Organization member routes registered", "")
}

func (h *OrganizationMemberHandler) getCurrentMember(c *fiber.Ctx) error {
	resp, err := h.client.GetCurrentMember(
		c.UserContext(),
		&tenantV1.TenantOrganizationIDRequest{OrganizationId: getOrgIDFromLocals(c)},
	)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	member := resp.GetMember()
	if member == nil {
		return response.Error(c, errors.NotFound("member not found"))
	}

	permsDto := make([]string, len(resp.GetPermissions()))
	for i, p := range resp.GetPermissions() {
		permsDto[i] = p.GetCode()
	}

	rolesDto := make([]string, len(resp.GetRoles()))
	for i, r := range resp.GetRoles() {
		rolesDto[i] = r.GetName()
	}

	memberDto := &dto.OrganizationMemberWithRBACResponse{
		OrganizationMember: *protoMemberToDTO(member),
		Roles:              rolesDto,
		Perms:              permsDto,
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"member fetched successfully",
		memberDto,
	)
}

func (h *OrganizationMemberHandler) listMembers(c *fiber.Ctx) error {
	req := tenantV1.TenantOrganizationIDRequest{
		OrganizationId: getOrgIDFromLocals(c),
	}

	resp, err := h.client.GetAllMembersByOrg(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	membersDto := make([]*dto.OrganizationMember, len(resp.GetMembers()))
	for i, m := range resp.GetMembers() {
		membersDto[i] = protoMemberToDTO(m)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"members fetched successfully",
		&dto.ListOrganizationMembersResponse{
			Members:    membersDto,
			TotalCount: resp.GetTotalCount(),
		},
	)
}

func (h *OrganizationMemberHandler) getMemberByID(c *fiber.Ctx) error {
	req := tenantV1.GetOrganizationMemberByIdRequest{
		OrganizationId: getOrgIDFromLocals(c),
		MemberId:       c.Params("memberID"),
	}

	resp, err := h.client.GetMemberById(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	member := resp.GetMember()
	if member == nil {
		return response.Error(c, errors.NotFound("member not found"))
	}

	memberDto := protoMemberToDTO(member)

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"member fetched successfully",
		memberDto,
	)
}

func (h *OrganizationMemberHandler) removeMember(c *fiber.Ctx) error {
	req := tenantV1.RemoveOrganizationMemberRequest{
		OrganizationId: getOrgIDFromLocals(c),
		MemberId:       c.Params("memberID"),
	}

	resp, err := h.client.RemoveMember(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}
