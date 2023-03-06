package model

import (
	"gorm.io/gorm"
	"tiktok/tools"
	"time"
)

type Model struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (model *Model) GenerateID() {
	model.ID = tools.NewSnowFlake(0).GenSnowID()
}
