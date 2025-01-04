package handler

import (
	"schedule_table/internal/model/dao"

	"github.com/google/uuid"
)

type Queue interface {
	Next(i int) uuid.UUID
	Select(i int)
	Skip()
}

type QueueMembers struct {
	MembersId []uuid.UUID
	SkipIndex int
}

func (q *QueueMembers) Next(i int) uuid.UUID {
	if i == len(q.MembersId) {
		panic("QueueMembers.Next: Skip More Member")
	}
	return q.MembersId[i]
}

func (s *QueueMembers) Select(i int) {
	s.SkipIndex = 0

	if i != len(s.MembersId)-1 {
		s.MembersId = append(s.MembersId[:i], append(s.MembersId[i+1:], s.MembersId[i])...)
	}
}

func (s *QueueMembers) Skip() {
	s.SkipIndex++
}

func selectOnlyMemberId(responsiblePersons *[]dao.Responsible) []uuid.UUID {
	MembersId := make([]uuid.UUID, len(*responsiblePersons))
	for i := 0; i < len(MembersId); i++ {
		MembersId[i] = (*responsiblePersons)[i].MemberId
	}

	return MembersId
}

func NewQueueMembers(responsiblePersons *[]dao.Responsible) Queue {
	return &QueueMembers{
		MembersId: selectOnlyMemberId(responsiblePersons),
		SkipIndex: 0,
	}
}

type ScheduleManager struct {
	Name         string
	Start        string
	Priority     int8
	BreakTime    uint32
	MembersQueue Queue
}

func NewScheduleManager(schedule *dao.Schedules, responsiblePersons *[]dao.Responsible) *ScheduleManager {
	return &ScheduleManager{
		Name:         schedule.Name,
		Start:        schedule.Hr_start,
		Priority:     schedule.Priority,
		BreakTime:    schedule.BreakTime,
		MembersQueue: NewQueueMembers(responsiblePersons),
	}
}
