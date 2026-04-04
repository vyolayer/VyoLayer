package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/usecase"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/logger"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

type OrganizationInvitationHandler struct {
	tenantV1.UnimplementedOrganizationInvitationServiceServer
	logger                *logger.AppLogger
	orgMemberInvitationUC usecase.OrganizationMemberInvitationUseCase
}

func NewOrganizationInvitationHandler(
	logger *logger.AppLogger,
	orgMemberInvitationUC usecase.OrganizationMemberInvitationUseCase,
) *OrganizationInvitationHandler {
	return &OrganizationInvitationHandler{
		logger:                logger,
		orgMemberInvitationUC: orgMemberInvitationUC,
	}
}

func (h *OrganizationInvitationHandler) CreateInvitation(
	ctx context.Context,
	req *tenantV1.CreateInvitationRequest,
) (*tenantV1.TenantSuccessResponse, error) {
	userID, _ := ctxutil.ExtractIAMUserUUID(ctx)
	orgID, _ := uuid.Parse(req.GetOrganizationId())

	err := h.orgMemberInvitationUC.Create(ctx, orgID, req.GetEmail(), req.GetRoleIds(), userID)
	if err != nil {
		return nil, err
	}

	return &tenantV1.TenantSuccessResponse{
		Message: "Invitation sent successfully",
	}, nil
}

func (h *OrganizationInvitationHandler) ListInvitations(
	ctx context.Context,
	req *tenantV1.ListInvitationsRequest,
) (*tenantV1.ListOrganizationInvitationsResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	invitations, err := h.orgMemberInvitationUC.ListByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}

	invitationsDto := make([]*tenantV1.OrganizationMemberInvitation, len(invitations))
	for i, invitation := range invitations {
		invitationsDto[i] = mapOrganizationMemberInvitationToProto(&invitation)
	}

	return &tenantV1.ListOrganizationInvitationsResponse{
		Invitations: invitationsDto,
	}, nil
}

func (h *OrganizationInvitationHandler) GetPendingByOrg(
	ctx context.Context,
	req *tenantV1.TenantOrganizationIDRequest,
) (*tenantV1.ListOrgInvitationsResponseForOrg, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	invitations, err := h.orgMemberInvitationUC.ListPendingByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}

	invitationsDto := make([]*tenantV1.OrganizationMemberInvitationForOrg, len(invitations))
	for i := range invitations {
		invitationsDto[i] = mapInvitationWithInviterToProto(&invitations[i])
	}

	return &tenantV1.ListOrgInvitationsResponseForOrg{
		Invitations: invitationsDto,
	}, nil
}

func (h *OrganizationInvitationHandler) GetPendingInvitations(
	ctx context.Context,
	req *tenantV1.GetPendingInvitationsRequest,
) (*tenantV1.ListOrganizationInvitationsResponse, error) {
	invitations, err := h.orgMemberInvitationUC.ListByUserEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	invitationsDto := make([]*tenantV1.OrganizationMemberInvitation, len(invitations))
	for i, invitation := range invitations {
		invitationsDto[i] = mapOrganizationMemberInvitationToProto(&invitation)
	}

	return &tenantV1.ListOrganizationInvitationsResponse{
		Invitations: invitationsDto,
	}, nil
}

func (h *OrganizationInvitationHandler) AcceptInvitation(
	ctx context.Context,
	req *tenantV1.AcceptInvitationRequest,
) (*tenantV1.TenantSuccessResponse, error) {
	userID, _ := ctxutil.ExtractIAMUserUUID(ctx)
	err := h.orgMemberInvitationUC.Accept(ctx, userID, req.GetToken())
	if err != nil {
		return nil, err
	}

	return &tenantV1.TenantSuccessResponse{
		Message: "Invitation accepted successfully",
	}, nil
}

func (h *OrganizationInvitationHandler) CancelInvitation(
	ctx context.Context,
	req *tenantV1.CancelInvitationRequest,
) (*tenantV1.TenantSuccessResponse, error) {
	userID, _ := ctxutil.ExtractIAMUserUUID(ctx)
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	err := h.orgMemberInvitationUC.CancelByOrgMember(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}

	return &tenantV1.TenantSuccessResponse{
		Message: "Invitation cancelled successfully",
	}, nil
}
