package accountmodelv1

import (
	"time"

	"github.com/google/uuid"
)

type ServiceUserLoginAttempt struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ProjectID  uuid.UUID  `gorm:"type:uuid"`
	UserID     *uuid.UUID `gorm:"type:uuid"`
	Identifier string

	IPAddress string
	UserAgent string

	Success bool

	CreatedAt time.Time
}

func (ServiceUserLoginAttempt) TableName() string {
	return "login_attempt"
}
