package apikey

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/response"
	pb "github.com/vyolayer/vyolayer/proto/apikey/v1"
)

const (
	grpcTimeout = 10 * time.Second
)

type ApiKeyHandler struct {
	logger *logger.AppLogger
	client pb.APIKeyServiceClient
	iamJWT jwt.IamJWT
}

func NewHandler(
	logger *logger.AppLogger,
	client pb.APIKeyServiceClient,
	iamJWT jwt.IamJWT,
) *ApiKeyHandler {
	return &ApiKeyHandler{
		logger: logger.WithContext("APIKeyHandler"),
		client: client,
		iamJWT: iamJWT,
	}
}

func (h *ApiKeyHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(grpcTimeout)

	apiKey := router.Group("/api-keys")
	apiKey.Use(grpcCtxMiddleware.Handler())
	apiKey.Use(middleware.IamJWTVerify(h.iamJWT))

	apiKey.Get("/", h.List)
	apiKey.Post("/", h.Create)

	apiKey.Get("/:apiKeyID", h.Get)
	apiKey.Patch("/:apiKeyID/rotate", h.Rotate)
	apiKey.Delete("/:apiKeyID/revoke", h.Revoke)
	apiKey.Get("/:apiKeyID/validate", h.Validate)

	h.logger.Info("APIKey routes registered", "")
}

func (h *ApiKeyHandler) Create(c *fiber.Ctx) error {
	var req pb.CreateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	var (
		organizationID, projectID, actorID string
	)

	// organization and project from query params
	organizationID = c.Query("org-id")
	if organizationID == "" {
		return response.Error(c, errors.BadRequest("organization id is required"))
	}
	projectID = c.Query("project-id")
	if projectID == "" {
		return response.Error(c, errors.BadRequest("project id is required"))
	}

	//
	actorID, _ = ctxutil.ExtractIAMUserID(c.UserContext())
	if actorID == "" {
		return response.Error(c, errors.Unauthorized("unauthorized"))
	}

	req.OrganizationId = organizationID
	req.ProjectId = projectID
	req.ActorId = actorID

	resp, err := h.client.CreateAPIKey(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.Success(c, resp)
}

func (h *ApiKeyHandler) List(c *fiber.Ctx) error {

	var (
		req            pb.ListAPIKeysRequest
		organizationID string
		projectID      string
	)

	// organization and project from query params
	organizationID = c.Query("org-id")
	if organizationID == "" {
		return response.Error(c, errors.BadRequest("organization id is required"))
	}
	projectID = c.Query("project-id")
	if projectID == "" {
		return response.Error(c, errors.BadRequest("project id is required"))
	}

	req.OrganizationId = organizationID
	req.ProjectId = projectID

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 5)
	status := c.Query("status")
	req.Page = int32(page)
	req.Limit = int32(limit)

	if status != "" {
		req.Status = status
	}

	resp, err := h.client.ListAPIKeys(c.UserContext(), &req)

	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.Success(c, resp)
}

func (h *ApiKeyHandler) Get(c *fiber.Ctx) error {

	var (
		req            pb.GetAPIKeyRequest
		organizationID string
		projectID      string
		apiKeyID       string
	)

	// organization and project from query params
	organizationID = c.Query("org-id")
	if organizationID == "" {
		return response.Error(c, errors.BadRequest("organization id is required"))
	}
	projectID = c.Query("project-id")
	if projectID == "" {
		return response.Error(c, errors.BadRequest("project id is required"))
	}

	apiKeyID = c.Params("apiKeyID")
	if apiKeyID == "" {
		return response.Error(c, errors.BadRequest("api key id is required"))
	}

	req.OrganizationId = organizationID
	req.ProjectId = projectID
	req.Id = apiKeyID

	resp, err := h.client.GetAPIKey(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.Success(c, resp)
}

func (h *ApiKeyHandler) Revoke(c *fiber.Ctx) error {

	var (
		req            pb.RevokeAPIKeyRequest
		organizationID string
		projectID      string
		actorID        string
		apiKeyID       string
	)

	organizationID = c.Query("org-id")
	if organizationID == "" {
		return response.Error(c, errors.BadRequest("organization id is required"))
	}
	projectID = c.Query("project-id")
	if projectID == "" {
		return response.Error(c, errors.BadRequest("project id is required"))
	}

	apiKeyID = c.Params("apiKeyID")
	if apiKeyID == "" {
		return response.Error(c, errors.BadRequest("api key id is required"))
	}

	actorID, _ = ctxutil.ExtractIAMUserID(c.UserContext())
	if actorID == "" {
		return response.Error(c, errors.Unauthorized("unauthorized"))
	}

	req.OrganizationId = organizationID
	req.ProjectId = projectID
	req.ActorId = actorID
	req.Id = apiKeyID

	resp, err := h.client.RevokeAPIKey(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.Success(c, resp)
}

func (h *ApiKeyHandler) Rotate(c *fiber.Ctx) error {

	var (
		req            pb.RotateAPIKeyRequest
		organizationID string
		projectID      string
		actorID        string
		apiKeyID       string
	)

	organizationID = c.Query("org-id")
	if organizationID == "" {
		return response.Error(c, errors.BadRequest("organization id is required"))
	}
	projectID = c.Query("project-id")
	if projectID == "" {
		return response.Error(c, errors.BadRequest("project id is required"))
	}

	apiKeyID = c.Params("apiKeyID")
	if apiKeyID == "" {
		return response.Error(c, errors.BadRequest("api key id is required"))
	}

	actorID, _ = ctxutil.ExtractIAMUserID(c.UserContext())
	if actorID == "" {
		return response.Error(c, errors.Unauthorized("unauthorized"))
	}

	req.OrganizationId = organizationID
	req.ProjectId = projectID
	req.ActorId = actorID
	req.Id = apiKeyID

	resp, err := h.client.RotateAPIKey(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.Success(c, resp)
}

func (h *ApiKeyHandler) Validate(c *fiber.Ctx) error {
	var req pb.ValidateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	resp, err := h.client.ValidateAPIKey(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.Success(c, resp)
}
