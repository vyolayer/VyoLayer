package apikey

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/platform/database/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type APIKeyInfo struct {
	ProjectID      uuid.UUID
	OrganizationID uuid.UUID
	KeyHash        string
	Mode           string
}

type APIKeyVerifier interface {
	Verify(rawKey string, projectID uuid.UUID) (*APIKeyInfo, error)
}

type apiKeyVerifier struct {
	client *gorm.DB
}

func NewAPIKeyVerifier(client *gorm.DB) APIKeyVerifier {
	return &apiKeyVerifier{
		client: client,
	}
}

func (v *apiKeyVerifier) Verify(rawKey string, projectID uuid.UUID) (*APIKeyInfo, error) {
	if rawKey == "" {
		return nil, status.Error(codes.Unauthenticated, "missing API key")
	}

	if projectID == uuid.Nil {
		return nil, status.Error(codes.Unauthenticated, "missing project ID")
	}

	keyHash := hashKey(rawKey)

	var apiKeyModel models.ApiKey
	err := v.client.Where("key_hash = ? AND project_id = ?", keyHash, projectID).First(&apiKeyModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.Unauthenticated, "invalid API key")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &APIKeyInfo{
		ProjectID:      apiKeyModel.ProjectID,
		OrganizationID: apiKeyModel.OrganizationID,
		KeyHash:        apiKeyModel.KeyHash,
		Mode:           apiKeyModel.Mode,
	}, nil
}

func hashKey(rawKey string) string {
	hash := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(hash[:])
}
