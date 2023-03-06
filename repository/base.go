package repository

import (
	"tiktok/global"
)

type BaseRepo struct{}

var baseRepo BaseRepo

func (base *BaseRepo) Create(model any) error {
	return global.DB.Create(model).Error
}

func (base *BaseRepo) First(model any, where any) error {
	return global.DB.Where(where).First(model).Error
}

func (base *BaseRepo) Find(model any, where any) error {
	return global.DB.Where(where).Find(model).Error
}
func (base *BaseRepo) FindWhereOrderLimit(model any, where any, order any, limit int) error {
	return global.DB.Debug().Where(where).Order(order).Limit(limit).Find(model).Error
}