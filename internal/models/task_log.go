package models

import (
	"GoScheduler/internal/modules/global"
	"gorm.io/gorm"
)

type TaskLog struct {
	gorm.Model
	TaskId     uint         `gorm:"column:task_id;not null" json:"task_id"`                                // 任务id
	Name       string       `gorm:"column:name;not null" json:"name"`                                      // 任务名称
	Spec       string       `gorm:"column:spec;not null" json:"spec"`                                      // crontab
	Protocol   TaskProtocol `json:"protocol" gorm:"type:tinyint;not null"`                                 // 协议 1:http 2:RPC
	Command    string       `gorm:"column:command;not null" json:"command"`                                // URL地址或shell命令
	Timeout    int          `gorm:"column:timeout;not null;default:0" json:"timeout"`                      // 任务执行超时时间(单位秒),0不限制
	RetryTimes int8         `gorm:"column:retry_times;type:tinyint;not null;default:0" json:"retry_times"` // 任务重试次数
	HostName   string       `gorm:"column:host_name;not null;default:''" json:"host_name"`                 // RPC主机名，逗号分隔
	Result     string       `json:"result" gorm:"type:mediumtext;not null"`                                // 执行结果
	Status     Status       `json:"status" gorm:"type:tinyint;not null;index;default:1"`                   // 状态 0:执行失败 1:执行中  2:执行完毕 3:任务取消(上次任务未执行完成) 4:异步执行
	TotalTime  int          `json:"total_time" gorm:"-"`                                                   // 执行总时长
}

// Create 新增
func (taskLog *TaskLog) Create() (insertId uint, err error) {
	result := global.DB.Create(taskLog)
	if result.Error == nil {
		insertId = taskLog.Model.ID
	}

	return
}

func (taskLog *TaskLog) Update(id uint, data CommonMap) (int64, error) {
	result := global.DB.Model(&taskLog).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		// 更新数据时发生错误
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (taskLog *TaskLog) Clear() (int64, error) {
	result := global.DB.Delete(&TaskLog{})
	if result.Error != nil {
		// 更新数据时发生错误
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (taskLog *TaskLog) GetList(params CommonMap) (list []TaskLog, err error) {
	m := global.DB.Model(&taskLog)

	taskId, ok := params["task_id"]
	if ok {
		m = m.Where("task_id = ?", taskId)
	}

	protocol, ok := params["protocol"]
	if ok {
		m = m.Where("protocol = ?", protocol)
	}

	status, ok := params["status"]
	if ok {
		m = m.Where("status = ?", status)
	}

	err = PageLimitOffset(m, params).Find(&list).Error
	return
}
