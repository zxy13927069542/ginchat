package models

import (
	"ginchat/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var db *gorm.DB

//	Init() database连接初始化
func Init(c config.Config) *gorm.DB {
	//	自定义日志打印sql
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
		SlowThreshold: time.Second,
		LogLevel: logger.Info,
		Colorful: true,
	})

	var err error
	db, err = gorm.Open(mysql.Open(c.Mysql.Dns), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&UserBasic{})
	db.AutoMigrate(&Message{})
	db.AutoMigrate(&GroupBasic{})
	db.AutoMigrate(&Contact{})
	return db
}
