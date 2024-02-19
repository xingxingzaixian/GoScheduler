package models

import (
	"GoScheduler/internal/modules/global"
	"gorm.io/gorm"
)

type TaskHostDetail struct {
	gorm.Model
	TaskId uint   `gorm:"column:task_id;not null" json:"task_id"`
	HostId int16  `gorm:"column:host_id;not null" json:"host_id"`
	Name   string `gorm:"column:name" json:"name"`
	Port   int    `gorm:"column:port" json:"port"`
	Alias  string `gorm:"column:alias" json:"alias"`
}

func (TaskHostDetail) TableName() string {
	return "task_host"
}

func (t *TaskHostDetail) Remove(taskId uint) error {
	result := global.DB.Where("task_id = ?", taskId).Delete(&TaskHostDetail{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (t *TaskHostDetail) Add(taskId uint, hostIds []int) error {
	err := t.Remove(taskId)
	if err != nil {
		return err
	}

	taskHosts := make([]TaskHostDetail, len(hostIds))
	for i, value := range hostIds {
		taskHosts[i].TaskId = taskId
		taskHosts[i].HostId = int16(value)
	}
	result := global.DB.Create(&taskHosts)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
