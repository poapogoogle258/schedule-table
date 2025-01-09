package service

import (
	"fmt"
	"schedule_table/internal/model/dao"
	"time"

	"github.com/google/uuid"
)

type Id interface {
	uuid.UUID | string
}
type CountRecipe[T Id] struct {
	Id T
	Member
	Count int
}
type CountRecipes []CountRecipe[uuid.UUID]

func (list *CountRecipes) Count(Id uuid.UUID) int {
	for _, item := range *list {
		if item.Id == Id {
			return item.Count
		}
	}

	return 0
}

func (list *CountRecipes) Add(Id uuid.UUID) error {
	for i, item := range *list {
		if item.Id == Id {
			(*list)[i].Count++
			return nil
		}
	}

	*list = append(*list, CountRecipe[uuid.UUID]{Id: Id, Count: 1})
	return nil

}

type Manager struct {
	Id               uuid.UUID
	MasterScheduleId *uuid.UUID
	RestTime         time.Duration
	Queue            IQueue
	Count            *CountRecipes
	Tasks            *[]dao.Tasks
}

func NewManagerSchedule(schedule *dao.Schedules) *Manager {

	fmt.Println("Queue", NewResponsibleQueue(schedule.Id, schedule.Responsibles))

	return &Manager{
		Id:               schedule.Id,
		MasterScheduleId: schedule.MasterScheduleId,
		RestTime:         time.Duration(schedule.BreakTime * uint32(time.Second)),
		Count:            &CountRecipes{},
		Queue:            NewResponsibleQueue(schedule.Id, schedule.Responsibles),
	}
}

func NewManagerScheduleWithQueue(schedule *dao.Schedules, queue IQueue) *Manager {
	return &Manager{
		Id:               schedule.Id,
		MasterScheduleId: schedule.MasterScheduleId,
		RestTime:         time.Duration(schedule.BreakTime * uint32(time.Second)),
		Count:            &CountRecipes{},
		Queue:            queue,
	}
}

// export to ManagerService
type IManagerService interface {
	NewManagerSchedule(schedule *dao.Schedules) *Manager
	NewManagerScheduleWithQueue(schedule *dao.Schedules, queue IQueue) *Manager
}

type ManagerService struct{}

func (service *ManagerService) NewManagerSchedule(schedule *dao.Schedules) *Manager {
	return NewManagerSchedule(schedule)
}

func (service *ManagerService) NewManagerScheduleWithQueue(schedule *dao.Schedules, queue IQueue) *Manager {
	return NewManagerScheduleWithQueue(schedule, queue)
}

func NewManagerService() IManagerService {
	return &ManagerService{}
}
