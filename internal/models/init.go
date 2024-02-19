package models

import (
	"GoScheduler/internal/modules/global"
	"fmt"
	"github.com/spf13/viper"
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
	if viper.GetString("db.engine") == "mysql" {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			viper.GetString("db.user"),
			viper.GetString("db.password"),
			viper.GetString("db.host"),
			viper.GetInt("db.port"),
			viper.GetString("db.name"),
		)
		db, err := gorm.Open(
			mysql.Open(dsn),
			&option)
		if err != nil {
			zap.S().Fatalf("MySQL数据库连接失败:%s:%d", viper.GetString("db.host"), viper.GetInt("db.port"))
		}

		global.DB.AutoMigrate(&Host{}, &User{})

		global.DB = db
	} else if viper.GetString("db.engine") == "sqlite" {
		db, err := gorm.Open(
			mysql.Open(viper.GetString("db.user")),
			&option)
		if err != nil {
			zap.S().Fatalf("SQlite数据库连接失败:%s:%d", viper.GetString("db.user"))
		}

		global.DB.AutoMigrate(&Host{}, &User{}, &TaskHostDetail{})

		global.DB = db
	}
}

func PageLimitOffset(db *gorm.DB, params CommonMap) *gorm.DB {
	page, ok := params["page"]
	if !ok || page.(int) <= 0 {
		page = global.Page
	}

	pageSize, ok := params["pageSize"]
	if !ok || pageSize.(int) <= 0 {
		pageSize = global.PageSize
	}

	return db.Limit(pageSize.(int)).Offset((page.(int) - 1) * pageSize.(int))
}
