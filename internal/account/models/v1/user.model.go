package accountmodelv1

import (
	"time"

	"github.com/google/uuid"
)

type ServiceUser struct {
	UUID
	TimeStamps

	ProjectID uuid.UUID `gorm:"type:uuid;not null;"`

	Email         string `gorm:"type:varchar(255);not null;"`
	EmailVerified bool   `gorm:"default:false;"`
	Password      string `gorm:"type:text;"`
	Username      string `gorm:"type:varchar(50);"`
	FirstName     string `gorm:"type:varchar(50);not null;"`
	LastName      string `gorm:"type:varchar(50);"`

	AvatarID uuid.UUID         `gorm:"type:uuid;"`
	Avatar   ServiceUserAvatar `gorm:"foreignKey:AvatarID;references:ID;constraint:OnDelete:SET NULL"`

	LastLoginAt *time.Time
	Status      string `gorm:"size:20;not null;default:'pending';"`
}

func (ServiceUser) TableName() string {
	return "user"
}

// IsActive checks if the user is active
func (u ServiceUser) IsActive() bool {
	return !u.TimeStamps.DeletedAt.Valid && u.Status == "active"
}

// IsVerified checks if the user is verified
func (u ServiceUser) IsVerified() bool {
	return u.EmailVerified
}
