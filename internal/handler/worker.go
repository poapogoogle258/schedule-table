package handler

import (
	"schedule_table/internal/model/dao"
	"time"

	"github.com/google/uuid"
)

type Worker interface {
	GetId() uuid.UUID
	GetData() *dao.Members
	IsAvailable(when time.Time) bool
	AddTask(task *dao.Tasks, breakTimeSecond uint32)
}

type Member struct {
	Id       uuid.UUID
	Member   *dao.Members
	Leaves   []string
	RestTime time.Time
}

func (m *Member) GetData() *dao.Members {
	return m.Member
}

func (m *Member) GetId() uuid.UUID {
	return m.Id
}

func (m *Member) IsAvailable(when time.Time) bool {

	if m.Leaves != nil {
		date := when.Format(time.DateOnly)
		for _, v := range m.Leaves {
			if date == v {
				return false
			}
		}
	}

	return m.RestTime.Equal(when) || m.RestTime.Before(when)
}

func (m *Member) AddTask(task *dao.Tasks, breakTimeSecond uint32) {
	m.RestTime = task.End.Add(time.Second * time.Duration(breakTimeSecond))
}

func getBetweens(start time.Time, end time.Time) []string {
	betweens := make([]string, 0)

	now := start.Add(time.Second * 0)

	for now.Before(end) {
		betweens = append(betweens, now.Format(time.DateOnly))
	}

	return betweens
}

func selectOnlyDateLeaves(leaves *[]dao.Leaves) []string {
	leavesList := make([]string, 0)
	for i := 0; i < len(*leaves); i++ {
		leave := &(*leaves)[i]
		leavesList = append(leavesList, getBetweens(leave.Start, leave.End)...)
	}

	return leavesList
}

func NewWorker(member *dao.Members) Worker {

	var restTime time.Time
	if member.LastTimeTask != nil {
		restTime = *member.LastTimeTask
	} else {
		restTime = time.Date(2000, 12, 1, 0, 0, 0, 0, time.Local)
	}

	return &Member{
		Id:       member.Id,
		Member:   member,
		RestTime: restTime,
		Leaves:   selectOnlyDateLeaves(member.Leaves),
	}
}
