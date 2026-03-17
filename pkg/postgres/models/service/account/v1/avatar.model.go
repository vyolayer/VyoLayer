package servicemodelv1

type ServiceUserAvatar struct {
	UUID
	TimeStamps

	URL           string `gorm:"type:text;"`
	FallbackChar  string `gorm:"type:varchar(1);"`
	FallbackColor string `gorm:"type:varchar(7);"`
}

func (ServiceUserAvatar) TableName() string {
	return "account_service.user_avatar"
}
