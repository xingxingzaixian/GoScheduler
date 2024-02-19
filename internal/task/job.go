package task

import (
	"GoScheduler/internal/models"
	taskModule "GoScheduler/internal/task/modules"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

func CreateJob(taskModel models.Task) cron.FuncJob {
	handler := CreateHandler(taskModel)
	if handler == nil {
		return nil
	}
	taskFunc := func() {
		taskCount.Add()
		defer taskCount.Done()

		taskLogId, err := beforeExecJob(taskModel)
		if taskLogId <= 0 {
			return
		}

		if taskModel.Multi == 0 {
			runInstance.add(taskModel.Id)
			defer runInstance.done(taskModel.Id)
		}

		concurrencyQueue.Add()
		defer concurrencyQueue.Done()

		zap.S().Infof("开始执行任务#%s#命令-%s", taskModel.Name, taskModel.Command)
		taskResult := execJob(handler, taskModel, taskLogId)
		zap.S().Infof("任务完成#%s#命令-%s", taskModel.Name, taskModel.Command)
		afterExecJob(taskModel, taskResult, taskLogId)
	}

	return taskFunc
}

func CreateHandler(taskModel models.Task) taskModule.Handler {
	var handler taskModule.Handler = nil
	switch taskModel.Protocol {
	case models.TaskHTTP:
		handler = new(taskModule.HTTPHandler)
	case models.TaskRPC:
		handler = new(taskModule.RPCHandler)
	}

	return handler
}

// 任务前置操作
func beforeExecJob(taskModel models.Task) (taskLogId uint, err error) {
	if taskModel.Multi == 0 && runInstance.has(taskModel.ID) {
		taskLogId, err = createTaskLog(taskModel, models.TaskCancel)
		return
	}
	taskLogId, err = createTaskLog(taskModel, models.TaskRunning)
	if err != nil {
		zap.S().Error("任务开始执行#写入任务日志失败-", err)
		return
	}
	zap.S().Debugf("任务命令-%s", taskModel.Command)

	return
}

// 任务执行后置操作
func afterExecJob(taskModel models.Task, taskResult TaskResult, taskLogId uint) {
	_, err := updateTaskLog(taskLogId, taskResult)
	if err != nil {
		zap.S().Error("任务结束#更新任务日志失败-", err)
	}

	// 发送邮件
	go SendNotification(taskModel, taskResult)
	// 执行依赖任务
	go execDependencyTask(taskModel, taskResult)
}

// 执行具体任务
func execJob(handler taskModule.Handler, taskModel models.Task, taskUniqueId uint) TaskResult {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Error("panic#service/task.go:execJob#", err)
		}
	}()
	// 默认只运行任务一次
	var execTimes int8 = 1
	if taskModel.RetryTimes > 0 {
		execTimes += taskModel.RetryTimes
	}
	var i int8 = 0
	var output string
	var err error
	for i < execTimes {
		output, err = handler.Run(taskModel, taskUniqueId)
		if err == nil {
			return TaskResult{Result: output, Err: err, RetryTimes: i}
		}
		i++
		if i < execTimes {
			zap.S().Warnf("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s", taskModel.ID, i, output, err.Error())
			if taskModel.RetryInterval > 0 {
				time.Sleep(time.Duration(taskModel.RetryInterval) * time.Second)
			} else {
				// 默认重试间隔时间，每次递增1分钟
				time.Sleep(time.Duration(i) * time.Minute)
			}
		}
	}

	return TaskResult{Result: output, Err: err, RetryTimes: taskModel.RetryTimes}
}

// 创建任务日志
func createTaskLog(taskModel models.Task, status models.Status) (uint, error) {
	taskLogModel := new(models.TaskLog)
	taskLogModel.TaskId = taskModel.ID
	taskLogModel.Name = taskModel.Name
	taskLogModel.Spec = taskModel.Spec
	taskLogModel.Protocol = taskModel.Protocol
	taskLogModel.Command = taskModel.Command
	taskLogModel.Timeout = taskModel.Timeout
	if taskModel.Protocol == models.TaskRPC {
		aggregationHost := ""
		for _, host := range taskModel.Hosts {
			aggregationHost += fmt.Sprintf("%s - %s<br>", host.Alias, host.Name)
		}
		taskLogModel.HostName = aggregationHost
	}
	taskLogModel.Status = status
	insertId, err := taskLogModel.Create()

	return insertId, err
}

// 更新任务日志
func updateTaskLog(taskLogId uint, taskResult TaskResult) (int64, error) {
	taskLogModel := new(models.TaskLog)
	var status models.Status
	result := taskResult.Result
	if taskResult.Err != nil {
		status = models.TaskFailure
	} else {
		status = models.TaskFinish
	}
	return taskLogModel.Update(taskLogId, models.CommonMap{
		"retry_times": taskResult.RetryTimes,
		"status":      status,
		"result":      result,
	})
}
