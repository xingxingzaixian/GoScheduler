package task

import (
	"GoScheduler/internal/models"
	rpcClient "GoScheduler/internal/modules/rpc/client"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}

// 批量添加任务
func BatchAddTask(tasks []models.Task) {
	for _, item := range tasks {
		RemoveAndAddTask(item)
	}
}

// 删除任务后添加
func RemoveAndAddTask(taskModel models.Task) {
	RemoveTask(taskModel.ID)
	AddTask(taskModel)
}

// 添加任务
func AddTask(taskModel models.Task) {
	if taskModel.Level == models.TaskLevelChild {
		zap.S().Errorf("添加任务失败#不允许添加子任务到调度器#任务Id-%d", taskModel.ID)
		return
	}
	taskFunc := createJob(taskModel)
	if taskFunc == nil {
		zap.S().Error("创建任务处理Job失败,不支持的任务协议#", taskModel.Protocol)
		return
	}

	if _, err := serviceCron.AddFunc(taskModel.Spec, taskFunc); err != nil {
		zap.S().Error("添加任务到调度器失败#", err)
	}
}

func StopTask(ip string, port int, id uint) {
	rpcClient.Stop(ip, port, id)
}

func RemoveTask(id uint) {
	serviceCron.Remove(cron.EntryID(id))
}

func RunTask(taskModel models.Task) {
	go createJob(taskModel)()
}
