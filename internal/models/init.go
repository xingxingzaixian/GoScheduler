package models

import (
	"GoScheduler/internal/modules/global"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type CommonMap map[string]interface{}

var TablePrefix = "tbl_"

func InitDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
		})

	option := gorm.Config{
		Logger:                 newLogger,
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   TablePrefix,
			SingularTable: true,
		},
	}
	if global.Setting.DBInfo.Engine == "mysql" {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			global.Setting.DBInfo.User,
			global.Setting.DBInfo.Password,
			global.Setting.DBInfo.Host,
			global.Setting.DBInfo.Port,
			global.Setting.DBInfo.Database,
		)
		db, err := gorm.Open(
			mysql.Open(dsn),
			&option)
		if err != nil {
			zap.S().Fatalf("MySQL数据库连接失败:%s:%d", global.Setting.DBInfo.Host, global.Setting.DBInfo.Port)
		}

		global.DB.AutoMigrate(&Host{}, &User{})

		global.DB = db
	} else if global.Setting.DBInfo.Engine == "sqlite" {
		db, err := gorm.Open(
			mysql.Open(global.Setting.DBInfo.User),
			&option)
		if err != nil {
			zap.S().Fatalf("SQlite数据库连接失败:%s:%d", global.Setting.DBInfo.User)
		}

		global.DB.AutoMigrate(&Host{}, &User{})

		global.DB = db
	}
}
