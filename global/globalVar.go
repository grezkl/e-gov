package global

import (
	"e-gov/config"

	"gorm.io/gorm"
)

var (
	Settings config.ServerConfig
	DB       *gorm.DB
)
