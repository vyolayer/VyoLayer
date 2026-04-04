package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/usecase"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/pagination"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

type OrganizationHandler struct {
	tenantV1.UnimplementedOrganizationServiceServer
	logger *logger.AppLogger
	orgUC  usecase.OrganizationUseCase
}

func NewOrganizationHandler(
	logger *logger.AppLogger,
	orgUC usecase.OrganizationUseCase,
) *OrganizationHandler {
	return &OrganizationHandler{
		logger: logger,
		orgUC:  orgUC,
	}
}

func (h *OrganizationHandler) CreateOrganization(
	ctx context.Context,
	req *tenantV1.CreateOrganizationRequest,
) (*tenantV1.OrganizationResponse, error) {
	org, member, err := h.orgUC.Create(ctx, req.GetName(), req.GetDescription())

	if err != nil {
		h.logger.ErrorWithErr("Failed to create organization", err)
		return nil, err
	}

	orgDto := mapOrganizationToProto(org)
	memberDto := mapOrganizationMemberToProto(member)
	membersDto := []*tenantV1.OrganizationMember{memberDto}

	h.logger.Debug("Organization created", map[string]any{
		"org":    orgDto,
		"member": memberDto,
	})
	return &tenantV1.OrganizationResponse{
		Organization: orgDto,
		Members:      membersDto,
	}, nil
}

func (h *OrganizationHandler) OnboardOrganization(
	ctx context.Context,
	req *tenantV1.CreateOrganizationRequest,
) (*tenantV1.OrganizationResponse, error) {
	resp, err := h.CreateOrganization(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *OrganizationHandler) ListOrganizations(
	ctx context.Context,
	req *tenantV1.ListOrganizationsRequest,
) (*tenantV1.ListOrganizationsResponse, error) {
	userID, _ := ctxutil.ExtractIAMUserUUID(ctx)
	offset := getOffset(req.GetPageToken())
	limit := getPageSize(req.GetPageSize())

	orgs, nextOffset, err := h.orgUC.List(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	orgsDto := make([]*tenantV1.Organization, len(orgs))
	for i, org := range orgs {
		orgsDto[i] = mapOrganizationToProto(org)
	}

	return &tenantV1.ListOrganizationsResponse{
		Organizations: orgsDto,
		TotalCount:    int32(limit),
		NextPageToken: pagination.EncodePageToken(nextOffset + limit),
	}, nil
}

func (h *OrganizationHandler) GetOrganizationById(
	ctx context.Context,
	req *tenantV1.TenantOrganizationIDRequest,
) (*tenantV1.OrganizationResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	orgWithMembers, err := h.orgUC.GetById(ctx, orgID)
	if err != nil {
		return nil, err
	}

	orgDto := mapOrganizationWithMembersToProto(orgWithMembers)

	membersDto := make([]*tenantV1.OrganizationMember, 0, len(orgWithMembers.Members))
	for _, m := range orgWithMembers.Members {
		memberDto := mapOrganizationMemberToProto(&m)
		membersDto = append(membersDto, memberDto)
	}

	return &tenantV1.OrganizationResponse{
		Organization: orgDto,
		Members:      membersDto,
	}, nil
}

// Update Organization
func (h *OrganizationHandler) UpdateOrganization(
	ctx context.Context,
	req *tenantV1.UpdateOrganizationRequest,
) (*tenantV1.OrganizationResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())

	res, err := h.orgUC.Update(ctx, orgID, req.GetName(), req.GetDescription())
	if err != nil {
		return nil, err
	}

	orgDto := mapOrganizationToProto(res)

	return &tenantV1.OrganizationResponse{
		Organization: orgDto,
	}, nil
}

func (h *OrganizationHandler) DeleteOrganization(
	ctx context.Context,
	req *tenantV1.DeleteOrganizationRequest,
) (*tenantV1.TenantSuccessResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	err := h.orgUC.Delete(ctx, orgID, req.GetConfirmName())
	if err != nil {
		return nil, err
	}

	return &tenantV1.TenantSuccessResponse{
		Message: "Organization deleted successfully",
	}, nil
}

func (h *OrganizationHandler) ArchiveOrganization(
	ctx context.Context,
	req *tenantV1.ArchiveOrganizationRequest,
) (*tenantV1.TenantSuccessResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	err := h.orgUC.Archive(ctx, orgID, req.GetConfirmName())
	if err != nil {
		return nil, err
	}

	return &tenantV1.TenantSuccessResponse{
		Message: "Organization archived successfully",
	}, nil
}

func (h *OrganizationHandler) RestoreOrganization(
	ctx context.Context,
	req *tenantV1.TenantOrganizationIDRequest,
) (*tenantV1.TenantSuccessResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	err := h.orgUC.Restore(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return &tenantV1.TenantSuccessResponse{
		Message: "Organization restored successfully",
	}, nil
}

func (h *OrganizationHandler) TransferOwnership(
	ctx context.Context,
	req *tenantV1.TransferOwnershipRequest,
) (*tenantV1.TenantSuccessResponse, error) {
	currentUserID, _ := ctxutil.ExtractIAMUserUUID(ctx)
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	memberID, _ := uuid.Parse(req.GetNewOwnerMemberId())

	err := h.orgUC.TransferOwnership(ctx, orgID, currentUserID, memberID)
	if err != nil {
		return nil, err
	}

	return &tenantV1.TenantSuccessResponse{
		Message: "Organization ownership transferred successfully",
	}, nil
}

func (h *OrganizationHandler) GetAllPermissions(
	ctx context.Context,
	req *tenantV1.TenantOrganizationIDRequest,
) (*tenantV1.ListOrganizationPermissionsResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	permissions, err := h.orgUC.GetAllPermissions(ctx, orgID)
	if err != nil {
		return nil, err
	}

	permissionsDto := make([]*tenantV1.OrganizationPermission, 0, len(permissions))
	for _, p := range permissions {
		permissionsDto = append(permissionsDto, mapOrganizationPermissionToProto(p))
	}

	return &tenantV1.ListOrganizationPermissionsResponse{
		Permissions: permissionsDto,
	}, nil
}

func (h *OrganizationHandler) GetAllRoles(
	ctx context.Context,
	req *tenantV1.TenantOrganizationIDRequest,
) (*tenantV1.ListOrganizationRolesResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	roles, err := h.orgUC.GetAllRoles(ctx, orgID)
	if err != nil {
		return nil, err
	}

	rolesDto := make([]*tenantV1.OrganizationRole, 0, len(roles))
	for _, r := range roles {
		rolesDto = append(rolesDto, mapOrganizationRoleToProto(r))
	}

	return &tenantV1.ListOrganizationRolesResponse{
		Roles: rolesDto,
	}, nil
}
