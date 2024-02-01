package models

import (
	"GoScheduler/internal/modules/global"
	"gorm.io/gorm"
	"time"
)

type TaskProtocol int8

const (
	TaskHTTP TaskProtocol = iota + 1 // HTTP协议
	TaskRPC                          // RPC方式执行命令
)

type TaskLevel int8

const (
	TaskLevelParent TaskLevel = 1 // 父任务
	TaskLevelChild  TaskLevel = 2 // 子任务(依赖任务)
)

type TaskDependencyStatus int8

const (
	TaskDependencyStatusStrong TaskDependencyStatus = 1 // 强依赖
	TaskDependencyStatusWeak   TaskDependencyStatus = 2 // 弱依赖
)

type TaskHTTPMethod int8

const (
	TaskHTTPMethodGet  TaskHTTPMethod = 1
	TaskHttpMethodPost TaskHTTPMethod = 2
)

type Task struct {
	gorm.Model
	Name             string               `gorm:"size:32;not null" json:"name"`
	Level            TaskLevel            `gorm:"type:tinyint;not null;default:1" json:"level"`
	DependencyTaskId string               `gorm:"size:64;not null;default:''" json:"dependency_task_id"`
	DependencyStatus TaskDependencyStatus `gorm:"type:tinyint;not null;default:1" json:"dependency_status"`
	Spec             string               `gorm:"size:64;not null" json:"spec"`
	Protocol         TaskProtocol         `gorm:"type:tinyint;not null;index" json:"protocol"`
	Command          string               `gorm:"size:256;not null" json:"command"`
	HttpMethod       TaskHTTPMethod       `gorm:"type:tinyint;not null;default:1" json:"http_method"`
	Timeout          int                  `gorm:"type:mediumint;not null;default:0" json:"timeout"`
	Multi            int8                 `gorm:"type:tinyint;not null;default:1" json:"multi"`
	RetryTimes       int8                 `gorm:"type:tinyint;not null;default:0" json:"retry_times"`
	RetryInterval    int16                `gorm:"type:smallint;not null;default:0" json:"retry_interval"`
	NotifyStatus     int8                 `gorm:"type:tinyint;not null;default:1" json:"notify_status"`
	NotifyType       int8                 `gorm:"type:tinyint;not null;default:0" json:"notify_type"`
	NotifyReceiverId string               `gorm:"size:256;not null;default:''" json:"notify_receiver_id"`
	NotifyKeyword    string               `gorm:"size:128;not null;default:''" json:"notify_keyword"`
	Tag              string               `gorm:"size:32;not null;default:''" json:"tag"`
	Remark           string               `gorm:"size:100;not null;default:''" json:"remark"`
	Status           global.Status        `gorm:"type:tinyint;not null;index;default:0" json:"status"`
	Hosts            []TaskHostDetail     `json:"hosts" gorm:"-"`
	NextRunTime      time.Time            `json:"next_run_time" gorm:"-"`
}

// 新增
func (task *Task) Create() (insertId uint, err error) {
	result := global.DB.Create(task)
	if result.Error == nil {
		insertId = task.Model.ID
	}

	return
}

func (task *Task) UpdateBean(id uint) (int64, error) {
	result := global.DB.Model(&Task{}).Where("id = ?", id).
		Select("name", "spec", "protocol", "command", "timeout", "multi",
			"retry_times", "retry_interval", "remark", "notify_status",
			"notify_type", "notify_receiver_id", "dependency_task_id", "dependency_status", "tag", "http_method", "notify_keyword").
		Updates(task)
	if result.Error != nil {
		// 更新数据时发生错误
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (task *Task) Update(id uint, data CommonMap) (int64, error) {
	result := global.DB.Model(&task).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		// 更新数据时发生错误
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (task *Task) Delete(id int) (int64, error) {
	result := global.DB.Delete(&task, id)

	if result.Error != nil {
		// 删除数据时发生错误
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (task *Task) Disable(id uint) (int64, error) {
	return task.Update(id, CommonMap{"status": global.Disabled})
}

func (task *Task) Enable(id uint) (int64, error) {
	return task.Update(id, CommonMap{"status": global.Enabled})
}
