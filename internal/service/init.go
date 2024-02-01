package service

import (
	"GoScheduler/internal/service/models"
	"github.com/gogf/gf/v2/os/gcron"
	"go.uber.org/zap"
)

var (
	// 定时任务调度管理器
	serviceCron *gcron.Cron

	// 同一任务是否有实例处于运行中
	runInstance models.Instance

	// 任务计数-正在运行的任务
	taskCount *models.TaskCount

	// 并发队列, 限制同时运行的任务数量
	concurrencyQueue *models.ConcurrencyQueue
)

func Initialize() {
	serviceCron := gcron.New()
	serviceCron.Start()
	concurrencyQueue = models.NewConcurrencyQueue()
	taskCount = models.NewTaskCount()
	go taskCount.Wait()

	zap.S().Info("开始初始化定时任务")
	taskModel := new(models.Task)
	taskNum := 0
	page := 1
	pageSize := 1000
	maxPage := 1000
	for page < maxPage {
		taskList, err := taskModel.ActiveList(page, pageSize)
		if err != nil {
			zap.S().Fatalf("定时任务初始化#获取任务列表错误: %s", err)
		}
		if len(taskList) == 0 {
			break
		}
		for _, item := range taskList {
			task.Add(item)
			taskNum++
		}
		page++
	}
	zap.S().Infof("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
}
