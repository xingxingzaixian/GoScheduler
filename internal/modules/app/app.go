package app

import (
	"GoScheduler/internal/modules/global"
	"GoScheduler/internal/modules/logger"
	"GoScheduler/internal/modules/utils"
	"go.uber.org/zap"
	"path/filepath"
)

func InitEnv() {
	logger.InitLogger()
	var err error
	global.AppDir, err = utils.GetWorkDir()
	if err != nil {
		zap.S().Fatal(err)
	}

	global.ConfDir = filepath.Join("/conf")
	utils.CreateDirIfNotExists(global.ConfDir)
	global.AppConfig = filepath.Join(global.ConfDir, "/app.yaml")
}
