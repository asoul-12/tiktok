package global

import (
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"tiktok/config"
)

var (
	Config        *config.Config
	DB            *gorm.DB
	SnowFlakeNode *snowflake.Node
)
