package model

type Service struct {
	BaseModel
	Key string `gorm:"size:63;not null;uniqueIndex"`

	Name        string `gorm:"size:63;not null;uniqueIndex"`
	Description string `gorm:"type:text"`
	Category    string `gorm:"size:63"`
	Icon        string `gorm:"size:63"`

	Version string `gorm:"size:30"`
	Status  string `gorm:"size:10;default:'pending'"` // pending, ready, maintenance, disabled, beta

	IsPublic   bool `gorm:"default:false"`
	IsInternal bool `gorm:"default:false"`

	SortOrder uint32 `gorm:"default:0"`
}

func (Service) TableName() string {
	return "services"
}
