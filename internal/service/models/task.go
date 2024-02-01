package models

import (
	"GoScheduler/internal/models"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type Task struct{}

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}

// 批量添加任务
func (task Task) BatchAdd(tasks []models.Task) {
	for _, item := range tasks {
		task.RemoveAndAdd(item)
	}
}

// 删除任务后添加
func (task Task) RemoveAndAdd(taskModel models.Task) {
	task.Remove(taskModel.Id)
	task.Add(taskModel)
}

// 添加任务
func (task Task) Add(taskModel models.Task) {
	if taskModel.Level == models.TaskLevelChild {
		zap.S().Errorf("添加任务失败#不允许添加子任务到调度器#任务Id-%d", taskModel.Id)
		return
	}
	taskFunc := createJob(taskModel)
	if taskFunc == nil {
		zap.S().Error("创建任务处理Job失败,不支持的任务协议#", taskModel.Protocol)
		return
	}

	cronName := strconv.Itoa(taskModel.Id)
	err := goutil.PanicToError(func() {
		serviceCron.AddFunc(taskModel.Spec, taskFunc, cronName)
	})
	if err != nil {
		zap.S().Error("添加任务到调度器失败#", err)
	}
}

func (task Task) NextRunTime(taskModel models.Task) time.Time {
	if taskModel.Level != models.TaskLevelParent ||
		taskModel.Status != models.Enabled {
		return time.Time{}
	}
	entries := serviceCron.Entries()
	taskName := strconv.Itoa(taskModel.Id)
	for _, item := range entries {
		if item.Name == taskName {
			return item.Next
		}
	}

	return time.Time{}
}

// 停止运行中的任务
func (task Task) Stop(ip string, port int, id int64) {
	rpcClient.Stop(ip, port, id)
}

func (task Task) Remove(id int) {
	serviceCron.RemoveJob(strconv.Itoa(id))
}

// 等待所有任务结束后退出
func (task Task) WaitAndExit() {
	serviceCron.Stop()
	taskCount.Exit()
}

// 直接运行任务
func (task Task) Run(taskModel models.Task) {
	go createJob(taskModel)()
}
