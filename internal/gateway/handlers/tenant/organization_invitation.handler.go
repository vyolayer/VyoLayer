package tenant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	dto "github.com/vyolayer/vyolayer/internal/shared/dto/tenant"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/response"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrganizationInvitationHandler struct {
	logger *logger.AppLogger
	client tenantV1.OrganizationInvitationServiceClient
	iamJWT jwt.IamJWT
}

func NewOrganizationInvitationHandler(
	logger *logger.AppLogger,
	client tenantV1.OrganizationInvitationServiceClient,
	iamJWT jwt.IamJWT,
) *OrganizationInvitationHandler {
	return &OrganizationInvitationHandler{
		logger: logger.WithContext("Org Invitation Handler"),
		client: client,
		iamJWT: iamJWT,
	}
}

func (h *OrganizationInvitationHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(tenantGRPCTimeout).Handler()

	org := router.Group("/organizations")
	org.Use(grpcCtxMiddleware, middleware.IamJWTVerify(h.iamJWT))

	// Invitation routes that don't require an org context (accept uses token, pending is user-scoped)
	org.Post("/invitations/accept", h.acceptInvitation)
	org.Get("/invitations/pending", h.getPendingByUser)

	// All routes below require a valid organizationID in the path
	orgGroup := org.Group("/:organizationID", middleware.ValidateOrganizationID())

	// Invitations
	orgGroup.
		Post("/invitations", h.createInvitation).
		Get("/invitations", h.listInvitations).
		Get("/invitations/pending", h.getPendingByOrgID).
		Delete("/invitations/:invitationID", h.cancelInvitation)

	h.logger.Info("Organization invitation routes registered", "")
}

func (h *OrganizationInvitationHandler) createInvitation(c *fiber.Ctx) error {

	var req tenantV1.CreateInvitationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.CreateInvitation(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusCreated, resp.GetMessage(), nil)
}

func (h *OrganizationInvitationHandler) listInvitations(c *fiber.Ctx) error {

	req := tenantV1.ListInvitationsRequest{
		OrganizationId: getOrgIDFromLocals(c),
		PageSize:       int32(c.QueryInt("page_size", 0)),
		PageToken:      c.Query("page_token", ""),
	}

	resp, err := h.client.ListInvitations(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	invitationsDto := make([]*dto.OrganizationInvitation, len(resp.GetInvitations()))
	for i, inv := range resp.GetInvitations() {
		invitationsDto[i] = protoInvitationToDTO(inv)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"invitations fetched successfully",
		&dto.ListOrganizationInvitationsResponse{Invitations: invitationsDto},
	)
}

// cancelInvitation cancels a pending invitation by its ID.
func (h *OrganizationInvitationHandler) cancelInvitation(c *fiber.Ctx) error {
	req := tenantV1.CancelInvitationRequest{
		OrganizationId: getOrgIDFromLocals(c),
		InvitationId:   c.Params("invitationID"),
	}

	resp, err := h.client.CancelInvitation(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

// getPendingInvitationsByOrgID returns all pending invitations for the specified orgID.
func (h *OrganizationInvitationHandler) getPendingByOrgID(c *fiber.Ctx) error {
	req := tenantV1.TenantOrganizationIDRequest{
		OrganizationId: getOrgIDFromLocals(c),
	}

	resp, err := h.client.GetPendingByOrg(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	invitationsDto := make([]*dto.OrganizationInvitationForOrg, len(resp.GetInvitations()))
	for i, inv := range resp.GetInvitations() {
		invitationsDto[i] = protoInvitationForOrgToDTO(inv)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"pending invitations fetched successfully",
		&dto.ListOrganizationInvitationsForOrgResponse{Invitations: invitationsDto},
	)
}

// getPendingInvitations returns all pending invitations for the authenticated user.
// This bypasses the org context check – it is scoped to the calling user.
func (h *OrganizationInvitationHandler) getPendingByUser(c *fiber.Ctx) error {
	userEmail, err := ctxutil.ExtractIAMUserEmail(c.UserContext())
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	req := tenantV1.GetPendingInvitationsRequest{
		Email: userEmail,
	}

	resp, err := h.client.GetPendingInvitations(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	invitationsDto := make([]*dto.OrganizationInvitation, len(resp.GetInvitations()))
	for i, inv := range resp.GetInvitations() {
		invitationsDto[i] = protoInvitationToDTO(inv)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"pending invitations fetched successfully",
		&dto.ListOrganizationInvitationsResponse{Invitations: invitationsDto},
	)
}

// acceptInvitation accepts an invitation using the token from the request body.
// It is unauthenticated at the org level (user may not be a member yet).
func (h *OrganizationInvitationHandler) acceptInvitation(c *fiber.Ctx) error {
	var req tenantV1.AcceptInvitationRequest
	token := c.Query("token")
	if token == "" {
		return response.Error(c, status.Error(codes.InvalidArgument, "token is required"))
	}
	req.Token = token

	resp, err := h.client.AcceptInvitation(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}
