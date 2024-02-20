package service

import (
	"GoScheduler/internal/models"
	"GoScheduler/internal/modules/global"
	setting2 "GoScheduler/internal/modules/setting"
	"GoScheduler/internal/task"
	"go.uber.org/zap"
)

func Start() {
	setting2.InitConfig(global.AppConfig)

	models.InitDB()

	// 3. 初始化任务调度
	task.Initialize()

	zap.S().Info("开始初始化定时任务")
}
