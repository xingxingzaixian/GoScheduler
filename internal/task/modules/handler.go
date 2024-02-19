package modules

import "GoScheduler/internal/models"

type Handler interface {
	Run(taskModel models.Task, taskUniqueId uint) (string, error)
}
