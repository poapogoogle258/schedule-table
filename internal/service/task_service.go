package service

import (
	"schedule_table/internal/model/dao"
	"schedule_table/internal/repository"
	"time"
)

type TaskService interface {
	FindRangeTasksOfCalendarId(calendarId string, start time.Time, end time.Time) ([]*dao.Tasks, error)
	EditTaskId(taskId string, value interface{}) (*dao.Tasks, error)
}

type taskService struct {
	TaskRepo repository.ITaskRepository
}

func (taskService *taskService) FindRangeTasksOfCalendarId(calendarId string, start time.Time, end time.Time) ([]*dao.Tasks, error) {
	tasks, errFindWithAssociation := taskService.TaskRepo.FindWithAssociation("calendar_id = ? AND start BETWEEN ? AND ?", calendarId, start, end)
	if errFindWithAssociation != nil {
		return nil, errFindWithAssociation
	}

	return tasks, nil
}

func (taskService *taskService) EditTaskId(taskId string, value interface{}) (*dao.Tasks, error) {
	task, errUpdateTask := taskService.TaskRepo.UpdatesAndFind(taskId, value)
	if errUpdateTask != nil {
		return nil, errUpdateTask
	}

	return task, nil
}

func NewTaskService(taskRepo repository.ITaskRepository) TaskService {
	return &taskService{
		TaskRepo: taskRepo,
	}
}
