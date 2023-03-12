package serverInit

import (
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"tiktok/global"
	"tiktok/router"
	"time"
)

func ServerInitAndStart() {
	InitDatabase()
	InitLogConfig()
	GenerateSnowFlakeNode()
	InitHertz()
}

func InitDatabase() {
	dsn := global.Config.Mysql.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("链接数据库出错了,%s", err)
	}
	_db, err := db.DB()
	if err != nil {
		log.Fatalf("链接数据库出错了,%s", err)
	}
	_db.SetMaxIdleConns(20)
	_db.SetMaxOpenConns(100)
	global.DB = db
}

func InitHertz() {
	hertz := server.Default(
		server.WithReadTimeout(global.Config.Server.ReadTimeOut),
		server.WithWriteTimeout(global.Config.Server.WriteTimeOut*time.Second),
		server.WithHostPorts(global.Config.Server.Addr))

	dir, _ := os.Getwd()
	hertz.Static("/assets", dir)
	router.RegisterRouter(hertz)
	hertz.Spin()
}

func InitLogConfig() {
	logrus.SetReportCaller(true)
}
func GenerateSnowFlakeNode() {
	node, err := snowflake.NewNode(0)
	global.SnowFlakeNode = node
	if err != nil {
		log.Fatalf("生成雪花节点失败了，%s", err)
	}
}
