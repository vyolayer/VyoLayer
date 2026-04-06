package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	"github.com/vyolayer/vyolayer/internal/tenant/usecase"
	"github.com/vyolayer/vyolayer/pkg/logger"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

type ProjectHandler struct {
	tenantV1.UnimplementedProjectServiceServer
	logger          *logger.AppLogger
	orgUC           usecase.OrganizationUseCase
	orgMemberUC     usecase.OrganizationMemberUseCase
	projectUC       domain.ProjectUseCase
	projectMemberUC domain.ProjectMemberUseCase
}

// NewProjectHandler injects BOTH use cases required by the Protobuf service
func NewProjectHandler(
	logger *logger.AppLogger,
	orgUC usecase.OrganizationUseCase,
	orgMemberUC usecase.OrganizationMemberUseCase,
	projectUC domain.ProjectUseCase,
	projectMemberUC domain.ProjectMemberUseCase,
) *ProjectHandler {
	return &ProjectHandler{
		logger:          logger,
		orgUC:           orgUC,
		orgMemberUC:     orgMemberUC,
		projectUC:       projectUC,
		projectMemberUC: projectMemberUC,
	}
}

// --- Helper to extract UserID from JWT context (Set by your Auth Interceptor) ---
func (h *ProjectHandler) getUserID(ctx context.Context) (uuid.UUID, error) {
	// Assume extractUserIDFromContext pulls the "sub" or "user_id" claim from the context metadata
	userID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return uuid.Nil, status.Error(codes.Unauthenticated, "invalid or missing user identity")
	}
	return userID, nil
}

// ==========================================
// PROJECT CRUD OPERATIONS
// ==========================================

func (h *ProjectHandler) CreateProject(ctx context.Context, req *tenantV1.CreateProjectRequest) (*tenantV1.ProjectResponse, error) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(req.GetOrganizationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid organization id format")
	}

	org, err := h.orgUC.GetById(ctx, orgID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "organization not found")
	}

	orgMember, err := h.orgMemberUC.GetByUserID(ctx, org.ID, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "organization member not found")
	}

	project, err := h.projectUC.Create(ctx, org.ID, orgMember.ID, req.GetName(), req.GetDescription())
	if err != nil {
		h.logger.Error("handler: failed to create project", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// add member
	member, err := h.projectMemberUC.AddMember(ctx, orgID, project.ID, userID, orgMember.ID, "project_admin")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tenantV1.ProjectResponse{
		Project: mapProjectToProto(project),
		Members: []*tenantV1.ProjectMember{mapProjectMemberToProto(member)},
	}, nil
}

func (h *ProjectHandler) GetProject(ctx context.Context, req *tenantV1.GetProjectRequest) (*tenantV1.ProjectResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	projectID, err := uuid.Parse(req.GetProjectId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid project id format")
	}

	project, err := h.projectUC.Get(ctx, orgID, projectID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &tenantV1.ProjectResponse{Project: mapProjectToProto(project)}, nil
}

func (h *ProjectHandler) ListProjects(ctx context.Context, req *tenantV1.ListProjectsRequest) (*tenantV1.ListProjectsResponse, error) {
	orgID, err := uuid.Parse(req.GetOrganizationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid organization id")
	}

	// Simple pagination conversion (you can upgrade page_token to keyset pagination later)
	limit := req.GetPageSize()
	offset := int32(0) // Default offset, expand based on your page_token logic

	projects, totalCount, err := h.projectUC.List(ctx, orgID, limit, offset)
	if err != nil {
		h.logger.Error("handler: failed to list projects", err)
		return nil, status.Error(codes.Internal, "failed to fetch projects")
	}

	var protoProjects []*tenantV1.Project
	for _, p := range projects {
		protoProjects = append(protoProjects, mapProjectToProto(p))
	}

	return &tenantV1.ListProjectsResponse{
		Projects:   protoProjects,
		TotalCount: totalCount,
		// NextPageToken: "calculate_next_token_here",
	}, nil
}

func (h *ProjectHandler) UpdateProject(ctx context.Context, req *tenantV1.UpdateProjectRequest) (*tenantV1.ProjectResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	projectID, err := uuid.Parse(req.GetProjectId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid project id")
	}

	// req.Name and req.Description are already *string because they are "optional" in proto3!
	project, err := h.projectUC.Update(ctx, orgID, projectID, req.Name, req.Description)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tenantV1.ProjectResponse{Project: mapProjectToProto(project)}, nil
}

// func (h *ProjectHandler) ArchiveProject(ctx context.Context, req *tenantV1.ArchiveProjectRequest) (*tenantV1.TenantSuccessResponse, error) {
// 	orgID, _ := uuid.Parse(req.GetOrganizationId())
// 	projectID, _ := uuid.Parse(req.GetProjectId())

// 	if err := h.projectUC.Archive(ctx, orgID, projectID); err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &tenantV1.TenantSuccessResponse{Success: true, Message: "Project archived successfully"}, nil
// }

// func (h *ProjectHandler) RestoreProject(ctx context.Context, req *tenantV1.RestoreProjectRequest) (*tenantV1.TenantSuccessResponse, error) {
// 	orgID, _ := uuid.Parse(req.GetOrganizationId())
// 	projectID, _ := uuid.Parse(req.GetProjectId())

// 	if err := h.projectUC.Restore(ctx, orgID, projectID); err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &tenantV1.TenantSuccessResponse{Success: true, Message: "Project restored successfully"}, nil
// }

func (h *ProjectHandler) DeleteProject(ctx context.Context, req *tenantV1.DeleteProjectRequest) (*tenantV1.TenantSuccessResponse, error) {
	orgID, _ := uuid.Parse(req.GetOrganizationId())
	projectID, _ := uuid.Parse(req.GetProjectId())

	if err := h.projectUC.Delete(ctx, orgID, projectID, req.GetConfirmName()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &tenantV1.TenantSuccessResponse{Message: "Project permanently deleted"}, nil
}

// ==========================================
// PROJECT MEMBER OPERATIONS
// ==========================================

func (h *ProjectHandler) AddMember(ctx context.Context, req *tenantV1.AddProjectMemberRequest) (*tenantV1.ProjectMemberResponse, error) {
	callerID, err := h.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	orgID, _ := uuid.Parse(req.GetOrganizationId())
	projectID, _ := uuid.Parse(req.GetProjectId())
	targetUserID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid target user id")
	}

	member, err := h.projectMemberUC.AddMember(ctx, orgID, projectID, targetUserID, callerID, req.GetRole())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &tenantV1.ProjectMemberResponse{Member: mapProjectMemberToProto(member)}, nil
}

func (h *ProjectHandler) ListMembers(ctx context.Context, req *tenantV1.ListProjectMembersRequest) (*tenantV1.ListProjectMembersResponse, error) {
	projectID, _ := uuid.Parse(req.GetProjectId())

	members, totalCount, err := h.projectMemberUC.ListMembers(ctx, projectID, req.GetPageSize(), 0)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list members")
	}

	var protoMembers []*tenantV1.ProjectMember
	for _, m := range members {
		protoMembers = append(protoMembers, mapProjectMemberToProto(m))
	}

	return &tenantV1.ListProjectMembersResponse{
		Members:    protoMembers,
		TotalCount: totalCount,
	}, nil
}

func (h *ProjectHandler) GetMember(ctx context.Context, req *tenantV1.GetProjectMemberRequest) (*tenantV1.ProjectMemberResponse, error) {
	projectID, _ := uuid.Parse(req.GetProjectId())
	memberID, _ := uuid.Parse(req.GetMemberId())

	member, err := h.projectMemberUC.GetMember(ctx, projectID, memberID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &tenantV1.ProjectMemberResponse{Member: mapProjectMemberToProto(member)}, nil
}

func (h *ProjectHandler) GetCurrentMember(ctx context.Context, req *tenantV1.ListProjectMembersRequest) (*tenantV1.ProjectMemberResponse, error) {
	callerID, err := h.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	projectID, _ := uuid.Parse(req.GetProjectId())

	member, err := h.projectMemberUC.GetCurrentMember(ctx, projectID, callerID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &tenantV1.ProjectMemberResponse{Member: mapProjectMemberToProto(member)}, nil
}

func (h *ProjectHandler) ChangeMemberRole(ctx context.Context, req *tenantV1.ChangeProjectMemberRoleRequest) (*tenantV1.TenantSuccessResponse, error) {
	projectID, _ := uuid.Parse(req.GetProjectId())
	memberID, _ := uuid.Parse(req.GetMemberId())

	if err := h.projectMemberUC.ChangeRole(ctx, projectID, memberID, req.GetRole()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &tenantV1.TenantSuccessResponse{Message: "Member role updated"}, nil
}

func (h *ProjectHandler) RemoveMember(ctx context.Context, req *tenantV1.RemoveProjectMemberRequest) (*tenantV1.TenantSuccessResponse, error) {
	callerID, err := h.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	projectID, _ := uuid.Parse(req.GetProjectId())
	memberID, _ := uuid.Parse(req.GetMemberId())

	if err := h.projectMemberUC.RemoveMember(ctx, projectID, memberID, callerID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &tenantV1.TenantSuccessResponse{Message: "Member removed from project"}, nil
}

func (h *ProjectHandler) LeaveProject(ctx context.Context, req *tenantV1.ProjectIdRequest) (*tenantV1.TenantSuccessResponse, error) {
	callerID, err := h.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	projectID, _ := uuid.Parse(req.GetProjectId())

	if err := h.projectMemberUC.LeaveProject(ctx, projectID, callerID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &tenantV1.TenantSuccessResponse{Message: "You have left the project"}, nil
}

// ==========================================
// DOMAIN TO PROTOBUF MAPPERS
// ==========================================

func mapProjectToProto(p *domain.Project) *tenantV1.Project {
	if p == nil {
		return nil
	}

	// Notice we DO NOT map the DatabaseURL here. It stays safely in the backend.
	protoProj := &tenantV1.Project{
		Id:             p.GetIDString(),
		OrganizationId: p.GetOrganizationIDString(),
		Name:           p.GetName(),
		Slug:           p.GetSlug(),
		Description:    p.GetDescription(),
		IsActive:       p.GetIsActive(),
		CreatedBy:      p.GetCreatedByString(),
		MaxApiKeys:     p.GetMaxAPIKeys(),
		MaxMembers:     p.GetMaxMembers(),
		MemberCount:    p.GetMemberCount(),
		CreatedAt:      p.GetCreatedAtString(),
	}

	// Map optional archived time pointer
	// if archived := p.GetArchivedAtString(); archived != "" {
	// 	protoProj.ArchivedAt = &archived
	// }

	return protoProj
}

func mapProjectMemberToProto(m *domain.ProjectMember) *tenantV1.ProjectMember {
	if m == nil {
		return nil
	}

	protoMem := &tenantV1.ProjectMember{
		Id:       m.GetIDString(),
		UserId:   m.GetUserIDString(),
		Email:    m.GetEmail(),
		FullName: m.GetFullName(),
		Role:     m.GetRole(),
		IsActive: m.GetIsActive(),
		JoinedAt: m.GetJoinedAtString(),
	}

	if removed := m.GetRemovedAtString(); removed != "" {
		protoMem.RemovedAt = &removed
	}

	return protoMem
}
