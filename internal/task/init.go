package task

import (
	"GoScheduler/internal/models"
	"GoScheduler/internal/modules/global"
	concurrencyQueue2 "GoScheduler/lib/ConcurrencyQueue"
	"GoScheduler/lib/instance"
	taskCout "GoScheduler/lib/task_cout"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"strings"
)

var (
	// 定时任务调度管理器
	serviceCron *cron.Cron

	// 同一任务是否有实例处于运行中
	runInstance instance.Instance

	// 任务计数-正在运行的任务
	taskCount taskCout.TaskCount

	// 并发队列, 限制同时运行的任务数量
	ccQueue concurrencyQueue2.ConcurrencyQueue
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

	taskNum := 0
	for _, item := range taskList {
		AddTask(item)
		taskNum++
	}

	zap.S().Infof("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
}

// 发送任务结果通知
func SendNotification(taskModel models.Task, taskResult global.TaskResult) {
	var statusName string
	// 未开启通知
	if taskModel.NotifyStatus == 0 {
		return
	}
	if taskModel.NotifyStatus == 3 {
		// 关键字匹配通知
		if !strings.Contains(taskResult.Result, taskModel.NotifyKeyword) {
			return
		}
	}
	if taskModel.NotifyStatus == 1 && taskResult.Err == nil {
		// 执行失败才发送通知
		return
	}
	if taskModel.NotifyType != 3 && taskModel.NotifyReceiverId == "" {
		return
	}
	if taskResult.Err != nil {
		statusName = "失败"
	} else {
		statusName = "成功"
	}
	// 发送通知
	fmt.Println(statusName)
}
