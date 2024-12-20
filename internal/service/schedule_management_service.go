package service

import (
	"schedule_table/internal/model/dao"
	"time"

	"github.com/google/uuid"
)

// ######################################## Worker ###############################################

type Worker struct {
	Id       uuid.UUID
	Member   *dao.Members
	Leaves   *[]string
	RestTime time.Time
}

func (s *Worker) IsFreeAt(when time.Time) bool {

	if s.Leaves != nil {
		date := when.Format(time.DateOnly)
		for _, v := range *s.Leaves {
			if date == v {
				return false
			}
		}
	}

	return s.RestTime.Equal(when) || s.RestTime.Before(when)
}

func (s *Worker) AddTask(task *dao.Tasks) {
	s.RestTime = task.End.Add(time.Second * time.Duration(task.Description.RestTime))
}

type MapWorkers map[string]*Worker

func (workers *MapWorkers) CheckWorkerFree(id string, time time.Time) bool {
	worker, ok := (*workers)[id]
	if ok {
		return worker.IsFreeAt(time)
	} else {
		panic("CheckWorkerFree: not have this WorkerId in list")
	}
}

func GetBetweens(start *time.Time, end *time.Time) *[]string {
	betweens := make([]string, 0)

	now := start.Add(time.Second * 0)

	for now.Before(*end) {
		betweens = append(betweens, now.Format(time.DateOnly))
	}

	return &betweens
}

// ######################################## Tasks ###############################################

type TasksDaily map[string][]*dao.Tasks

func NewTasksDaily(tasks *[]dao.Tasks) *TasksDaily {
	tasksDaily := &TasksDaily{}
	tasksDaily.AddTasks(tasks)

	return tasksDaily
}

func (tasksDaily *TasksDaily) AddTasks(tasks *[]dao.Tasks) {

	for i := 0; i < len(*tasks); i++ {

		key := (*tasks)[i].Start.Format(time.DateOnly)

		if _, ok := (*tasksDaily)[key]; ok {
			(*tasksDaily)[key] = append((*tasksDaily)[key], &(*tasks)[i])
		} else {
			(*tasksDaily)[key] = append((*tasksDaily)[key], &(*tasks)[i])
		}

	}
}

// func (tasksDaily *TasksDaily) AddTask(task *dao.Tasks) {

// 	key := task.Start.Format(time.DateOnly)

// 	(*tasksDaily)[key] = []*dao.Tasks{task}

// 	// check task
// 	if _, ok := (*tasksDaily)[key]; ok {
// 		(*tasksDaily)[key] = append((*tasksDaily)[key], task)
// 	} else {
// 		(*tasksDaily)[key] = append((*tasksDaily)[key], task)
// 	}

// }

// ######################################## Schedule ###############################################

type ListResponsible struct {
	Members   *[]dao.Responsible
	SkipIndex int
}

func (s *ListResponsible) Next(i int) *dao.Responsible {

	if i == len(*s.Members) {
		panic("ListResponsible.Next : Skip More Member")
	}

	return &(*s.Members)[i]
}

func (s *ListResponsible) Select(i int) {

	*s.Members = append((*s.Members)[:s.SkipIndex+1], append((*s.Members)[i+1:], (*s.Members)[i])...)

	s.SkipIndex = 0

}

func (s *ListResponsible) Skip() {
	s.SkipIndex++
}

type Schedule struct {
	Name        *string
	Priority    int
	ListMembers ListResponsible
	TasksDaily  *TasksDaily
	Tasks       *[]dao.Tasks
}
