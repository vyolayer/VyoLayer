package repository

import (
	"log"
	"worklayer/internal/platform/database/models"
	"worklayer/internal/platform/database/types"

	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (sr *sessionRepository) Save(session *models.UserSession) RepositoryError {
	if err := sr.db.Create(session).Error; err != nil {
		return NewRepositoryError(500, err.Error())
	}
	return nil
}

func (sr *sessionRepository) FindByTokenHash(hashedToken string) (*models.UserSession, RepositoryError) {
	var session models.UserSession
	err := sr.db.Where("token_hash = ?", hashedToken).First(&session).Error
	if err != nil {
		log.Println("Error finding session:", err)
		return nil, NewRepositoryError(500, err.Error())
	}
	return &session, nil
}

func (sr *sessionRepository) DeleteByTokenHash(hashedToken string) RepositoryError {
	result := sr.db.Where("token_hash = ?", hashedToken).Delete(&models.UserSession{})
	if result.Error != nil {
		return NewRepositoryError(500, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return NewRepositoryError(404, "Session not found")
	}
	return nil
}

func (sr *sessionRepository) FindByUserId(userId types.UserID) (*models.UserSession, RepositoryError) {
	var session models.UserSession
	err := sr.db.Where("user_id = ?", userId.InternalID).First(&session).Error
	if err != nil {
		return nil, NewRepositoryError(500, err.Error())
	}
	return &session, nil
}

func (sr *sessionRepository) DeleteAllByUserId(userId types.UserID) RepositoryError {
	if err := sr.db.Where("user_id = ?", userId.InternalID).Delete(&models.UserSession{}).Error; err != nil {
		return NewRepositoryError(500, err.Error())
	}
	return nil
}

// RotateByTokenHash rotates a session by old token hash and new token hash
func (sr *sessionRepository) RotateByTokenHash(oldHashedToken string, newSession *models.UserSession) RepositoryError {
	err := sr.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&models.UserSession{}, "token_hash = ?", oldHashedToken)
		if result.Error != nil {
			log.Printf("SESSION REPOSITORY :: RotateByTokenHash : %v", result.Error)
			return result.Error
		}

		if result.RowsAffected == 0 {
			log.Printf("SESSION REPOSITORY :: RotateByTokenHash : %v", gorm.ErrRecordNotFound)
			return gorm.ErrRecordNotFound
		}

		if err := tx.Save(newSession).Error; err != nil {
			log.Printf("SESSION REPOSITORY :: RotateByTokenHash : %v", err)
			return err
		}
		return nil
	})

	if err != nil {
		return NewRepositoryError(500, err.Error())
	}
	return nil
}
