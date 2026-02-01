package mapper

import (
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/models"
	"worklayer/internal/platform/database/types"
)

// ToDomainUser converts a User model to a User domain object.
func ToDomainUser(userModel *models.User) *domain.User {
	if userModel == nil {
		return nil
	}

	userID, err := types.ReconstructUserID(userModel.ID.String())
	if err != nil {
		return nil
	}

	return domain.ReconstructUser(
		*userID,
		userModel.Email,
		userModel.PasswordHash,
		userModel.FullName,
		userModel.IsActive,
		userModel.IsEmailVerified,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	)
}

// ToDBUser converts a User domain object to a User model.
func ToDBUser(user *domain.User) *models.User {
	if user == nil {
		return nil
	}

	return &models.User{
		BaseModel: models.BaseModel{
			ID:        user.ID.InternalID().ID(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Email:           user.Email,
		PasswordHash:    user.HashedPassword,
		FullName:        user.FullName,
		IsActive:        user.IsActive,
		IsEmailVerified: user.IsEmailVerified,
	}
}
