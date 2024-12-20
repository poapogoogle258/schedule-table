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
	Leaves   []string
	RestTime time.Time
}

func (s *Worker) IsFreeAt(when time.Time) bool {

	if s.Leaves != nil {
		date := when.Format(time.DateOnly)
		for _, v := range s.Leaves {
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

type MapWorkers map[uuid.UUID]*Worker

func (workers *MapWorkers) IsAvailable(memberId uuid.UUID, when time.Time) bool {
	worker, ok := (*workers)[memberId]
	if ok {
		return worker.IsFreeAt(when)
	} else {
		panic("CheckWorkerFree: not have this WorkerId in list")
	}
}

func (workers *MapWorkers) AddTask(memberId uuid.UUID, task *dao.Tasks) {

	(*workers)[memberId].RestTime = task.End.Add(time.Second * time.Duration(task.Description.RestTime))

}

func NewMapWorker(members *[]dao.Members, leaves *[]dao.Leaves) MapWorkers {
	workersMap := MapWorkers{}
	now := time.Now()
	startOfDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	for i := 0; i < len(*members); i++ {
		member := &(*members)[i]
		workersMap[member.Id] = &Worker{
			Id:       member.Id,
			Member:   member,
			RestTime: *getValueOrDefaultTime(member.LastTimeTask, &startOfDate),
		}
	}
	for i := 0; i < len(*leaves); i++ {
		workersMap[(*leaves)[i].MemberId].Leaves = getBetweens(&(*leaves)[i].Start, &(*leaves)[i].End)
	}

	return workersMap

}

func getValueOrDefaultTime(t *time.Time, defaultValue *time.Time) *time.Time {
	if t == nil {
		return defaultValue
	}
	return t
}

func getBetweens(start *time.Time, end *time.Time) []string {
	betweens := make([]string, 0)

	now := start.Add(time.Second * 0)

	for now.Before(*end) {
		betweens = append(betweens, now.Format(time.DateOnly))
	}

	return betweens
}

// ######################################## ScheduleManager ###############################################

type ScheduleManager struct {
	Name      string
	Start     string
	Priority  int
	MembersId []uuid.UUID
	SkipIndex int
}

func (s *ScheduleManager) Next(i int) uuid.UUID {

	if i == len(s.MembersId) {
		panic("ScheduleManager.Next : Skip More Member")
	}

	return s.MembersId[i]
}

func (s *ScheduleManager) Select(i int) {

	s.SkipIndex = 0

	if i != len(s.MembersId)-1 {
		s.MembersId = append(s.MembersId[:i], append(s.MembersId[i+1:], s.MembersId[i])...)
	}

}

func (s *ScheduleManager) Skip() {
	s.SkipIndex++
}

func NewScheduleManager(schedule *dao.Schedules, responsiblePersons *[]dao.Responsible) *ScheduleManager {

	scheduleManager := &ScheduleManager{}

	scheduleManager.Name = schedule.Name
	scheduleManager.Start = schedule.Hr_start
	scheduleManager.Priority = int(schedule.Priority)
	scheduleManager.SkipIndex = 0

	MembersId := make([]uuid.UUID, len(*responsiblePersons))
	for i := 0; i < len(MembersId); i++ {
		MembersId[i] = (*responsiblePersons)[i].MemberId
	}

	scheduleManager.MembersId = MembersId

	return scheduleManager

}
