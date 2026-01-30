package repository

import (
	"log"
	"worklayer/internal/platform/database/models"

	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (sr *sessionRepository) Save(session *models.UserSession) error {
	return sr.db.Create(session).Error
}

func (sr *sessionRepository) FindByTokenHash(hashedToken string) (*models.UserSession, error) {
	var session models.UserSession
	err := sr.db.Where("token_hash = ?", hashedToken).First(&session).Error
	if err != nil {
		log.Println("Error finding session:", err)
		return nil, err
	}
	return &session, nil
}

func (sr *sessionRepository) DeleteByTokenHash(hashedToken string) error {
	result := sr.db.Where("token_hash = ?", hashedToken).Delete(&models.UserSession{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (sr *sessionRepository) FindByUserId(userId uint) (*models.UserSession, error) {
	var session models.UserSession
	err := sr.db.Where("user_id = ?", userId).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (sr *sessionRepository) DeleteAllByUserId(userId uint) error {
	return sr.db.Where("user_id = ?", userId).Delete(&models.UserSession{}).Error
}

// RotateByTokenHash rotates a session by old token hash and new token hash
func (sr *sessionRepository) RotateByTokenHash(oldHashedToken string, newSession *models.UserSession) error {
	return sr.db.Transaction(func(tx *gorm.DB) error {
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
}
