package entity

import (
	"gorm.io/gorm"
	"tiktok/global"
	"time"
)

type Model struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (model *Model) GenerateID() {
	model.ID = global.SnowFlakeNode.Generate().Int64()
}
