package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/usecase"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/logger"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

type OrganizationMemberHandler struct {
	tenantV1.UnimplementedOrganizationMemberServiceServer
	logger      *logger.AppLogger
	orgMemberUC usecase.OrganizationMemberUseCase
}

func NewOrganizationMemberHandler(
	logger *logger.AppLogger,
	orgMemberUC usecase.OrganizationMemberUseCase,
) *OrganizationMemberHandler {
	return &OrganizationMemberHandler{
		logger:      logger,
		orgMemberUC: orgMemberUC,
	}
}

// Get current member
func (h *OrganizationMemberHandler) GetCurrentMember(
	ctx context.Context,
	req *tenantV1.TenantOrganizationIDRequest,
) (*tenantV1.OrganizationMemberWithRBACResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	currentUserID, _ := ctxutil.ExtractIAMUserUUID(ctx)
	member, err := h.orgMemberUC.GetByUserID(ctx, orgID, currentUserID)
	if err != nil {
		return nil, err
	}

	memberDto := mapOrganizationMemberWithRolesAndPermissionsToProto(member)

	return &tenantV1.OrganizationMemberWithRBACResponse{
		Member: memberDto,
	}, nil
}

// Get all members by org
func (h *OrganizationMemberHandler) GetAllMembersByOrg(
	ctx context.Context,
	req *tenantV1.TenantOrganizationIDRequest,
) (*tenantV1.ListOrganizationMembersResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	members, err := h.orgMemberUC.GetAllMembersByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}

	membersDto := make([]*tenantV1.OrganizationMember, len(members))
	for i, member := range members {
		membersDto[i] = mapOrganizationMemberWithRolesToProto(member)
	}

	return &tenantV1.ListOrganizationMembersResponse{
		Members: membersDto,
	}, nil
}

func (h *OrganizationMemberHandler) GetMemberById(
	ctx context.Context,
	req *tenantV1.GetOrganizationMemberByIdRequest,
) (*tenantV1.OrganizationMemberResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	memberID, _ := uuid.Parse(req.GetMemberId())
	member, err := h.orgMemberUC.GetById(ctx, orgID, memberID)
	if err != nil {
		return nil, err
	}

	memberDto := mapOrganizationMemberWithRolesAndPermissionsToProto(member)

	return &tenantV1.OrganizationMemberResponse{
		Member: memberDto,
	}, nil
}

func (h *OrganizationMemberHandler) RemoveMember(
	ctx context.Context,
	req *tenantV1.RemoveOrganizationMemberRequest,
) (*tenantV1.TenantSuccessResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	memberID, _ := uuid.Parse(req.GetMemberId())
	currentUserID, _ := ctxutil.ExtractIAMUserUUID(ctx)

	err := h.orgMemberUC.RemoveMember(ctx, orgID, memberID, currentUserID)
	if err != nil {
		return nil, err
	}

	return &tenantV1.TenantSuccessResponse{
		Message: "Member removed successfully",
	}, nil
}

// func (h *OrganizationMemberHandler) ChangeMemberRole(
// 	ctx context.Context,
// 	req *tenantV1.ChangeOrganizationMemberRoleRequest,
// ) (*tenantV1.TenantSuccessResponse, error) {

// }
