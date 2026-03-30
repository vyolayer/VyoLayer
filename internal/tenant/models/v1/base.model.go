package tenantmodelv1

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimeStamps struct {
	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"<-:update;type:timestamp;default:CURRENT_TIMESTAMP"`
}

type BaseModel struct {
	ID uuid.UUID `gorm:"<-:create;type:uuid;primaryKey"`
	TimeStamps
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
