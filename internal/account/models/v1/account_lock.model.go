package accountmodelv1

import (
	"time"

	"github.com/google/uuid"
)

type ServiceUserAccountLock struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ProjectID uuid.UUID `gorm:"type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid;index"`

	Reason string

	LockedAt  time.Time
	ExpiresAt *time.Time
}

func (ServiceUserAccountLock) TableName() string {
	return "account_service.account_lock"
}
