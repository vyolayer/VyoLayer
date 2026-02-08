package repository

import (
	"log"
	"worklayer/internal/platform/database/models"
	"worklayer/internal/platform/database/types"
	"worklayer/pkg/errors"

	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (sr *sessionRepository) Save(session *models.UserSession) *errors.AppError {
	if err := sr.db.Create(session).Error; err != nil {
		return ConvertDBError(err, "saving session")
	}
	return nil
}

func (sr *sessionRepository) FindByTokenHash(hashedToken string) (*models.UserSession, *errors.AppError) {
	var session models.UserSession
	err := sr.db.Where("token_hash = ?", hashedToken).First(&session).Error
	if err != nil {
		log.Println("Error finding session:", err)
		return nil, ConvertDBError(err, "finding session by token hash")
	}
	return &session, nil
}

func (sr *sessionRepository) DeleteByTokenHash(hashedToken string) *errors.AppError {
	result := sr.db.Where("token_hash = ?", hashedToken).Delete(&models.UserSession{})
	if result.Error != nil {
		return ConvertDBError(result.Error, "deleting session")
	}

	if result.RowsAffected == 0 {
		return errors.NotFound("Session with token hash not found")
	}
	return nil
}

func (sr *sessionRepository) FindByUserId(userId types.UserID) (*models.UserSession, *errors.AppError) {
	var session models.UserSession
	err := sr.db.Where("user_id = ?", userId.InternalID).First(&session).Error
	if err != nil {
		return nil, ConvertDBError(err, "finding session by user ID")
	}
	return &session, nil
}

func (sr *sessionRepository) DeleteAllByUserId(userId types.UserID) *errors.AppError {
	if err := sr.db.Where("user_id = ?", userId.InternalID).Delete(&models.UserSession{}).Error; err != nil {
		return ConvertDBError(err, "deleting all sessions for user")
	}
	return nil
}

// RotateByTokenHash rotates a session by old token hash and new token hash
func (sr *sessionRepository) RotateByTokenHash(oldHashedToken string, newSession *models.UserSession) *errors.AppError {
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
		return TransactionError(err, "rotating session")
	}
	return nil
}
