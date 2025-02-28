package lib

import (
	"github.com/google/uuid"
)

type ScheduleListCommitted interface {
	Add(employeeId uuid.UUID)
	Clear()
	includes(employeeId uuid.UUID) bool
	GetCommitted(employeeId uuid.UUID) int
	CheckCommittedLimit(employeeId uuid.UUID) bool
	GetRotationCycle() int
	CheckAllListCommittedHaveLimit() bool
}

type scheduleListCommitted struct {
	Id                   uuid.UUID
	rotationCycle        int
	mapEmployeeCommitted map[uuid.UUID]int
}

func (sm *scheduleListCommitted) Add(employeeId uuid.UUID) {
	sm.mapEmployeeCommitted[employeeId]++
}

func (sm *scheduleListCommitted) Clear() {
	for id := range sm.mapEmployeeCommitted {
		sm.mapEmployeeCommitted[id] = 0
	}
}

func (sm *scheduleListCommitted) includes(employeeId uuid.UUID) bool {
	_, ok := sm.mapEmployeeCommitted[employeeId]
	return ok
}

func (sm *scheduleListCommitted) GetCommitted(employeeId uuid.UUID) int {
	if committed, ok := sm.mapEmployeeCommitted[employeeId]; ok {
		return committed
	} else {
		return -1
	}
}

func (sm *scheduleListCommitted) CheckCommittedLimit(employeeId uuid.UUID) bool {
	committed, ok := sm.mapEmployeeCommitted[employeeId]
	if !ok {
		return false
	}

	return committed == sm.rotationCycle
}

func (sm *scheduleListCommitted) GetRotationCycle() int {
	return sm.rotationCycle
}

func (sm *scheduleListCommitted) CheckAllListCommittedHaveLimit() bool {
	for id := range sm.mapEmployeeCommitted {
		if sm.mapEmployeeCommitted[id] < sm.rotationCycle {
			return false
		}
	}

	return true
}

func NewScheduleListCommitted(scheduleId uuid.UUID, memberIds []uuid.UUID, rotationCycle int) ScheduleListCommitted {

	mapEmployeeCommitted := make(map[uuid.UUID]int, len(memberIds))
	for _, memberId := range memberIds {
		mapEmployeeCommitted[memberId] = 0
	}

	return &scheduleListCommitted{
		Id:                   scheduleId,
		rotationCycle:        rotationCycle,
		mapEmployeeCommitted: mapEmployeeCommitted,
	}
}
