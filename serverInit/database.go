package serverInit

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"tiktok/global"
)

func InitDatabase() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/db_tiktok?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	global.DB = db
}

func InitDatabaseTest() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/db_tiktok?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	global.DB = db
}
