package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/handlers/dto"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/response"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

// ProjectHandler handles all /organizations/:organizationID/projects/* routes.
type ProjectHandler struct {
	logger *logger.AppLogger
	client tenantV1.ProjectServiceClient
	iamJWT jwt.IamJWT
}

func NewProjectHandler(
	logger *logger.AppLogger,
	client tenantV1.ProjectServiceClient,
	iamJWT jwt.IamJWT,
) *ProjectHandler {
	return &ProjectHandler{
		logger: logger.WithContext("Project Handler"),
		client: client,
		iamJWT: iamJWT,
	}
}

// RegisterRoutes mounts all project and project-member routes.
//
// Route hierarchy:
//
//	GET    /organizations/:organizationID/projects
//	POST   /organizations/:organizationID/projects
//	GET    /organizations/:organizationID/projects/:projectID
//	PATCH  /organizations/:organizationID/projects/:projectID
//	DELETE /organizations/:organizationID/projects/:projectID
//
//	GET    /organizations/:organizationID/projects/:projectID/members
//	GET    /organizations/:organizationID/projects/:projectID/members/me
//	GET    /organizations/:organizationID/projects/:projectID/members/:memberID
//	POST   /organizations/:organizationID/projects/:projectID/members
//	PATCH  /organizations/:organizationID/projects/:projectID/members/:memberID/role
//	DELETE /organizations/:organizationID/projects/:projectID/members/:memberID
//	DELETE /organizations/:organizationID/projects/:projectID/members/leave
func (h *ProjectHandler) RegisterRoutes(router fiber.Router) {
	router.Use(grpcCtxMiddleware(tenantGRPCTimeout))
	router.Use(middleware.IamJWTVerify(h.iamJWT))

	// /organizations/:organizationID/projects
	projects := router.Group("/organizations/:organizationID/projects")
	projects.Use(middleware.ValidateOrganizationID())

	// Collection routes
	projects.Get("/", h.listProjects)
	projects.Post("/", h.createProject)

	// Single-project sub-group — requires a valid :projectID
	project := projects.Group("/:projectID", middleware.ValidateProjectID())
	project.Get("/", h.getProject)
	project.Patch("/", h.updateProject)
	project.Delete("/", h.deleteProject)

	// ── Members sub-group ────────────────────────────────────────────────────
	members := project.Group("/members")
	members.Get("/", h.listMembers)
	members.Post("/", h.addMember)
	members.Get("/me", h.getCurrentMember)
	// members.Delete("/leave", h.leaveProject)
	members.Get("/:memberID", h.getMember)
	members.Post("/:memberID/role", h.changeMemberRole)
	members.Delete("/:memberID", h.removeMember)

	h.logger.Info("Project routes registered", "")
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func getProjectIDFromLocals(c *fiber.Ctx) string {
	id, _ := c.Locals("project_id").(string)
	return id
}

// ─── Project CRUD ─────────────────────────────────────────────────────────────

func (h *ProjectHandler) createProject(c *fiber.Ctx) error {
	var req tenantV1.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.CreateProject(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusCreated,
		"project created successfully",
		protoProjectResponseToDTO(resp),
	)
}

func (h *ProjectHandler) getProject(c *fiber.Ctx) error {
	req := tenantV1.GetProjectRequest{
		OrganizationId: getOrgIDFromLocals(c),
		ProjectId:      getProjectIDFromLocals(c),
	}

	resp, err := h.client.GetProject(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"project fetched successfully",
		protoProjectResponseToDTO(resp),
	)
}

func (h *ProjectHandler) listProjects(c *fiber.Ctx) error {
	req := tenantV1.ListProjectsRequest{
		OrganizationId: getOrgIDFromLocals(c),
		PageSize:       int32(c.QueryInt("page_size", 0)),
		PageToken:      c.Query("page_token", ""),
	}

	resp, err := h.client.ListProjects(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	projects := make([]*dto.TProject, len(resp.GetProjects()))
	for i, p := range resp.GetProjects() {
		projects[i] = protoProjectToDTO(p)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"projects fetched successfully",
		&dto.ListProjects{
			Projects:      projects,
			TotalCount:    resp.GetTotalCount(),
			NextPageToken: resp.GetNextPageToken(),
		},
	)
}

func (h *ProjectHandler) updateProject(c *fiber.Ctx) error {
	var req tenantV1.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)
	req.ProjectId = getProjectIDFromLocals(c)

	resp, err := h.client.UpdateProject(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"project updated successfully",
		protoProjectToDTO(resp.GetProject()),
	)
}

func (h *ProjectHandler) deleteProject(c *fiber.Ctx) error {
	var req tenantV1.DeleteProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)
	req.ProjectId = getProjectIDFromLocals(c)

	resp, err := h.client.DeleteProject(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

// ─── Project Member operations ────────────────────────────────────────────────

func (h *ProjectHandler) listMembers(c *fiber.Ctx) error {
	req := tenantV1.ListProjectMembersRequest{
		OrganizationId: getOrgIDFromLocals(c),
		ProjectId:      getProjectIDFromLocals(c),
		PageSize:       int32(c.QueryInt("page_size", 0)),
		PageToken:      c.Query("page_token", ""),
	}

	resp, err := h.client.ListMembers(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	members := make([]*dto.TProjectMember, len(resp.GetMembers()))
	for i, m := range resp.GetMembers() {
		members[i] = protoProjectMemberToDTO(m)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"members fetched successfully",
		&dto.ListProjectMembers{
			Members:       members,
			TotalCount:    resp.GetTotalCount(),
			NextPageToken: resp.GetNextPageToken(),
		},
	)
}

func (h *ProjectHandler) getMember(c *fiber.Ctx) error {
	req := tenantV1.GetProjectMemberRequest{
		OrganizationId: getOrgIDFromLocals(c),
		ProjectId:      getProjectIDFromLocals(c),
		MemberId:       c.Params("memberID"),
	}

	resp, err := h.client.GetMember(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"member fetched successfully",
		protoProjectMemberToDTO(resp.GetMember()),
	)
}

func (h *ProjectHandler) getCurrentMember(c *fiber.Ctx) error {
	req := tenantV1.ListProjectMembersRequest{
		OrganizationId: getOrgIDFromLocals(c),
		ProjectId:      getProjectIDFromLocals(c),
	}

	resp, err := h.client.GetCurrentMember(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"current member fetched successfully",
		protoProjectMemberToDTO(resp.GetMember()),
	)
}

func (h *ProjectHandler) addMember(c *fiber.Ctx) error {
	var req tenantV1.AddProjectMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)
	req.ProjectId = getProjectIDFromLocals(c)

	resp, err := h.client.AddMember(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusCreated,
		"member added successfully",
		protoProjectMemberToDTO(resp.GetMember()),
	)
}

func (h *ProjectHandler) changeMemberRole(c *fiber.Ctx) error {
	var req tenantV1.ChangeProjectMemberRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)
	req.ProjectId = getProjectIDFromLocals(c)
	req.MemberId = c.Params("memberID")

	resp, err := h.client.ChangeMemberRole(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

func (h *ProjectHandler) removeMember(c *fiber.Ctx) error {
	req := tenantV1.RemoveProjectMemberRequest{
		OrganizationId: getOrgIDFromLocals(c),
		ProjectId:      getProjectIDFromLocals(c),
		MemberId:       c.Params("memberID"),
	}

	resp, err := h.client.RemoveMember(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

func (h *ProjectHandler) leaveProject(c *fiber.Ctx) error {
	req := tenantV1.ProjectIdRequest{
		OrganizationId: getOrgIDFromLocals(c),
		ProjectId:      getProjectIDFromLocals(c),
	}

	resp, err := h.client.LeaveProject(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}
