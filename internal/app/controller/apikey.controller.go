package controller

import (
	"worklayer/internal/app/dto"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/service"
	"worklayer/pkg/errors"
	"worklayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ApiKeyController interface {
	GenerateKey(ctx *fiber.Ctx) error
	ListKeys(ctx *fiber.Ctx) error
	GetKeyByID(ctx *fiber.Ctx) error
	RevokeKey(ctx *fiber.Ctx) error
}

type apiKeyController struct {
	apiKeyService service.ApiKeyService
}

func NewApiKeyController(apiKeyService service.ApiKeyService) ApiKeyController {
	return &apiKeyController{apiKeyService: apiKeyService}
}

// @Summary Generate a new API key
// @Tags API Keys
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Param name body string true "API key name"
// @Param mode body string true "API key mode"
// @Success 201 {object} dto.ApiKeyCreatedDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId}/api-keys [post]
func (c *apiKeyController) GenerateKey(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	var req dto.CreateApiKeyRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	apiKey, rawKey, svcErr := c.apiKeyService.Generate(ctx, localUserID, projectID, req.Name, req.Mode)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.Created(ctx, dto.FromDomainApiKeyCreated(apiKey, rawKey))
}

// @Summary List all API keys for a project
// @Tags API Keys
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Success 200 {array} dto.ApiKeyDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId}/api-keys [get]
func (c *apiKeyController) ListKeys(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	keys, svcErr := c.apiKeyService.ListByProjectID(ctx, localUserID, projectID)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	dtos := make([]dto.ApiKeyDTO, 0, len(keys))
	for _, k := range keys {
		if d := dto.FromDomainApiKey(&k); d != nil {
			dtos = append(dtos, *d)
		}
	}

	return response.Success(ctx, dtos)
}

// @Summary Get a specific API key by ID
// @Tags API Keys
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Param apiKeyId path string true "API key ID"
// @Success 200 {object} dto.ApiKeyDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId}/api-keys/{apiKeyId} [get]
func (c *apiKeyController) GetKeyByID(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	apiKeyID, err := types.ReconstructApiKeyID(ctx.Params("apiKeyId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid API key ID format"))
	}

	apiKey, svcErr := c.apiKeyService.GetByID(ctx, localUserID, projectID, apiKeyID)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.Success(ctx, dto.FromDomainApiKey(apiKey))
}

// @Summary Revoke an API key
// @Tags API Keys
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Param apiKeyId path string true "API key ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId}/api-keys/{apiKeyId}/revoke [post]
func (c *apiKeyController) RevokeKey(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	apiKeyID, err := types.ReconstructApiKeyID(ctx.Params("apiKeyId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid API key ID format"))
	}

	if svcErr := c.apiKeyService.Revoke(ctx, localUserID, projectID, apiKeyID); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "API key revoked successfully")
}
