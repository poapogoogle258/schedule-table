package lib

import (
	"errors"
	"schedule_table/internal/model/dao"
	"schedule_table/util"
	"slices"

	"github.com/google/uuid"
)

var (
	ErrAllWorkerAreCrossed = errors.New("all workers are crossed")
	ErrNotHaveWorker       = errors.New("not have workers available")
)

type ScheduleManageQueue interface {
	OrderQueue(task []*dao.Task) error
	Commit(task *dao.Task) error
	RollBack(task *dao.Task) error
}

type scheduleManageQueue struct {
	Id               uuid.UUID
	scheduleQueue    ScheduleQueue
	mapListCommitted map[uuid.UUID]ScheduleListCommitted
	orderShiftsCycle orderShiftsCycle
}

func (sm *scheduleManageQueue) OrderQueue(task []*dao.Task) error {
	return nil
}

var (
	ErrIncludeListCommitted = errors.New("include in schedule committed list")
	ErrIsBusy               = errors.New("is busy")
)

func clearAllListCommitted(mapListCommitted map[uuid.UUID]ScheduleListCommitted) error {
	for key := range mapListCommitted {
		mapListCommitted[key].Clear()
	}

	return nil

}

func (sm *scheduleManageQueue) Commit(task *dao.Task) error {

	for range sm.scheduleQueue.Size() {

		// check shift all committed
		if sm.mapListCommitted[task.ScheduleId].CheckAllListCommittedHaveLimit() {
			clearAllListCommitted(sm.mapListCommitted)
			sm.scheduleQueue.ReSetQueue()
		}

		employee := sm.scheduleQueue.NextQueue()

		if sm.mapListCommitted[task.ScheduleId].CheckCommittedLimit(employee.Id) {

			continue
		}

		// check employee busy
		if works := employee.GetWorksInRange(task.Start, task.End); len(works) > 0 {

			continue
		}

		sm.mapListCommitted[task.ScheduleId].Add(employee.Id)
		employee.AddTask(task)

		break

	}

	return nil
}

func (sm *scheduleManageQueue) RollBack(task *dao.Task) error {
	return nil
}

func NewScheduleManageQueue(schedule *dao.Schedule, shiftsSchedule []uuid.UUID, employeeQueueOrdered []*Employee) ScheduleManageQueue {

	employeeIds := util.Map(employeeQueueOrdered, func(e *Employee) uuid.UUID {
		return e.Id
	})

	mapListCommitted := make(map[uuid.UUID]ScheduleListCommitted, len(shiftsSchedule))
	for _, shift := range shiftsSchedule {
		mapListCommitted[shift] = NewScheduleListCommitted(shift, employeeIds, schedule.RotationCycle)
	}

	return &scheduleManageQueue{
		Id:               schedule.Id,
		scheduleQueue:    NewScheduleQueue(employeeQueueOrdered),
		mapListCommitted: mapListCommitted,
		orderShiftsCycle: orderShiftsCycle(shiftsSchedule),
	}
}

type orderShiftsCycle []uuid.UUID

func (sh orderShiftsCycle) Next(scheduleId uuid.UUID) uuid.UUID {
	i := slices.IndexFunc(sh, func(s uuid.UUID) bool {
		return s == scheduleId
	})

	if i == -1 {
		panic("orderShiftsCycle.Next: not fount scheduleId in orderShiftsCycle")
	}

	if i+1 == len(sh) {
		return sh[0]
	}

	return sh[i+1]
}
