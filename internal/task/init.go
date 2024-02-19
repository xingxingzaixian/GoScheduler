package task

import (
	"GoScheduler/internal/models"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	// 定时任务调度管理器
	serviceCron *cron.Cron
)

func Initialize() {
	zap.S().Info("开始初始化定时任务")

	serviceCron = cron.New(cron.WithSeconds())
	serviceCron.Start()

	// 获取定时任务列表
	taskModel := new(models.Task)
	taskList, err := taskModel.GetActiveList()
	if err != nil {
		zap.S().Fatalf("定时任务初始化#获取任务列表错误: %s", err)
	}
	for _, item := range taskList {
		AddTask(item)
	}
}
