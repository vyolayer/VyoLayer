package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"vyolayer/internal/domain"
	"vyolayer/internal/platform/database/types"
	"vyolayer/internal/repository"
	"vyolayer/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

// ApiKeyService defines the interface for API key management.
type ApiKeyService interface {
	Generate(
		ctx *fiber.Ctx,
		userID types.UserID,
		projectID types.ProjectID,
		name, mode string,
	) (*domain.ApiKey, string, *errors.AppError) // returns domain key + raw key

	ListByProjectID(
		ctx *fiber.Ctx,
		userID types.UserID,
		projectID types.ProjectID,
	) ([]domain.ApiKey, *errors.AppError)

	GetByID(
		ctx *fiber.Ctx,
		userID types.UserID,
		projectID types.ProjectID,
		apiKeyID types.ApiKeyID,
	) (*domain.ApiKey, *errors.AppError)

	Revoke(
		ctx *fiber.Ctx,
		userID types.UserID,
		projectID types.ProjectID,
		apiKeyID types.ApiKeyID,
	) *errors.AppError

	ValidateKey(rawKey string) (*domain.ApiKey, *errors.AppError)
}

type apiKeyService struct {
	apiKeyRepo  repository.ApiKeyRepository
	memberRepo  repository.ProjectMemberRepository
	projectRepo repository.ProjectRepository
}

func NewApiKeyService(
	apiKeyRepo repository.ApiKeyRepository,
	memberRepo repository.ProjectMemberRepository,
	projectRepo repository.ProjectRepository,
) ApiKeyService {
	return &apiKeyService{
		apiKeyRepo:  apiKeyRepo,
		memberRepo:  memberRepo,
		projectRepo: projectRepo,
	}
}

// generateRawKey creates a crypto-random key with the given mode prefix.
// Format: wl_{mode}_{32 random hex chars}
func generateRawKey(mode string) (rawKey string, keyPrefix string, err error) {
	randomBytes := make([]byte, 32)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	randomHex := hex.EncodeToString(randomBytes)
	rawKey = fmt.Sprintf("wl_%s_%s", mode, randomHex)
	keyPrefix = rawKey[:16] // e.g. "wl_live_ab3f1234"
	return rawKey, keyPrefix, nil
}

// hashKey creates a SHA-256 hash of the raw API key.
func hashKey(rawKey string) string {
	hash := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(hash[:])
}

func (s *apiKeyService) Generate(
	ctx *fiber.Ctx,
	userID types.UserID,
	projectID types.ProjectID,
	name, mode string,
) (*domain.ApiKey, string, *errors.AppError) {
	// Validate mode
	if !domain.IsValidApiKeyMode(mode) {
		return nil, "", domain.ValidationError("API key mode must be 'dev' or 'live'")
	}

	// Verify user is project admin
	member, err := s.memberRepo.FindByUserAndProject(ctx.Context(), userID, projectID)
	if err != nil {
		return nil, "", err
	}
	if !member.IsAdmin() {
		return nil, "", errors.Forbidden("Only project admins can generate API keys")
	}

	// Check project API key limit
	project, err := s.projectRepo.FindByID(ctx.Context(), projectID)
	if err != nil {
		return nil, "", err
	}

	count, countErr := s.apiKeyRepo.CountByProjectID(ctx.Context(), projectID)
	if countErr != nil {
		return nil, "", countErr
	}
	if int(count) >= project.MaxApiKeys {
		return nil, "", domain.ApiKeyLimitReachedError()
	}

	// Generate the raw key
	rawKey, keyPrefix, genErr := generateRawKey(mode)
	if genErr != nil {
		return nil, "", errors.Internal("Failed to generate API key")
	}

	keyHash := hashKey(rawKey)

	// Create domain entity
	var expiresAt *time.Time
	if mode == domain.ApiKeyModeDev {
		// Dev keys expire after 90 days by default
		exp := time.Now().Add(90 * 24 * time.Hour)
		expiresAt = &exp
	}

	apiKey := domain.NewApiKey(
		projectID,
		project.OrganizationID,
		name,
		keyPrefix,
		keyHash,
		mode,
		userID,
		expiresAt,
	)
	if validErr := apiKey.Validate(); validErr != nil {
		return nil, "", validErr
	}

	created, createErr := s.apiKeyRepo.Create(ctx.Context(), apiKey)
	if createErr != nil {
		return nil, "", createErr
	}

	// Return both the domain key and the raw key (shown only once)
	return created, rawKey, nil
}

func (s *apiKeyService) ListByProjectID(
	ctx *fiber.Ctx,
	userID types.UserID,
	projectID types.ProjectID,
) ([]domain.ApiKey, *errors.AppError) {
	// Verify membership
	_, err := s.memberRepo.FindByUserAndProject(ctx.Context(), userID, projectID)
	if err != nil {
		return nil, errors.Forbidden("You are not a member of this project")
	}

	return s.apiKeyRepo.FindByProjectID(ctx.Context(), projectID)
}

func (s *apiKeyService) GetByID(
	ctx *fiber.Ctx,
	userID types.UserID,
	projectID types.ProjectID,
	apiKeyID types.ApiKeyID,
) (*domain.ApiKey, *errors.AppError) {
	// Verify membership
	_, err := s.memberRepo.FindByUserAndProject(ctx.Context(), userID, projectID)
	if err != nil {
		return nil, errors.Forbidden("You are not a member of this project")
	}

	return s.apiKeyRepo.FindByID(ctx.Context(), apiKeyID)
}

func (s *apiKeyService) Revoke(
	ctx *fiber.Ctx,
	userID types.UserID,
	projectID types.ProjectID,
	apiKeyID types.ApiKeyID,
) *errors.AppError {
	// Verify user is project admin
	member, err := s.memberRepo.FindByUserAndProject(ctx.Context(), userID, projectID)
	if err != nil {
		return err
	}
	if !member.IsAdmin() {
		return errors.Forbidden("Only project admins can revoke API keys")
	}

	// Verify the key exists and belongs to this project
	apiKey, err := s.apiKeyRepo.FindByID(ctx.Context(), apiKeyID)
	if err != nil {
		return err
	}
	if apiKey.ProjectID.String() != projectID.String() {
		return domain.ApiKeyNotFoundError(apiKeyID.String())
	}
	if apiKey.IsRevoked() {
		return domain.ApiKeyRevokedError(apiKeyID.String())
	}

	return s.apiKeyRepo.Revoke(ctx.Context(), apiKeyID, userID)
}

func (s *apiKeyService) ValidateKey(rawKey string) (*domain.ApiKey, *errors.AppError) {
	keyHash := hashKey(rawKey)

	apiKey, err := s.apiKeyRepo.FindByKeyHash(nil, keyHash)
	if err != nil {
		return nil, domain.ApiKeyInvalidError()
	}

	// Check if key is usable
	if usableErr := apiKey.IsUsable(); usableErr != nil {
		return nil, usableErr
	}

	return apiKey, nil
}
