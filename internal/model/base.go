package model

type Base struct {
	ID        int64 `gorm:"primaryKey" json:"id"`
	CreatedAt int64 `gorm:"autoCreateTime:milli" json:"createdAt"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli" json:"updatedAt"`
	DeletedAt int64 `json:"-"`
	//DeletedAt gorm.DeletedAt `gorm:"index"`
}
