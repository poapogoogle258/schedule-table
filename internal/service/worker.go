package service

import (
	"cmp"
	"errors"
	"schedule_table/internal/model/dao"
	"time"

	"github.com/google/uuid"
)

type Worker interface {
	GetId() uuid.UUID
	GetData() *dao.Members
	AddBooking(task *dao.Tasks, restTime time.Duration) error
	AddLeaven(day time.Time) error
	AddTask(task *dao.Tasks, restTime time.Duration, options ...func(*AddTaskOption)) error
}

var (
	ErrMemberNotAvailable = errors.New("worker is not available")
	ErrMemberNotReady     = errors.New("worker in rest time")
	ErrMemberIsLeaven     = errors.New("worker is leaven")
	ErrMemberReserved     = errors.New("worker was reserved")
)

type Member struct {
	Id        uuid.UUID
	Info      *dao.Members
	Available bool
	ReadyTime time.Time
	Leaves    []time.Time
	Booking   []*dao.Tasks
}

func (member *Member) GetId() uuid.UUID {
	return member.Id
}

func (member *Member) GetData() *dao.Members {
	return member.Info
}

type AddTaskOption struct {
	ForceAvailable bool
	ForceLeaves    bool
	ForceReadyTime bool
	ForceReserved  bool
}

func (member *Member) AddTask(task *dao.Tasks, restTime time.Duration, options ...func(*AddTaskOption)) error {
	addTaskOption := &AddTaskOption{}
	for i, _ := range options {
		options[i](addTaskOption)
	}

	if addTaskOption.ForceAvailable || !member.isAvailable() {
		return ErrMemberNotAvailable
	}
	if addTaskOption.ForceLeaves || member.isLeaves(task.Start) {
		return ErrMemberIsLeaven
	}
	if addTaskOption.ForceReadyTime || !member.isReadyTime(task.Start) {
		return ErrMemberNotReady
	}
	if addTaskOption.ForceReserved || member.isReserved(task, restTime) {
		return ErrMemberReserved
	}

	member.ReadyTime = task.End.Add(restTime)
	return nil
}

func (member *Member) AddBooking(task *dao.Tasks, restTime time.Duration) error {

	if !member.isAvailable() {
		return ErrMemberReserved
	} else if member.isReserved(task, restTime) {
		return ErrMemberNotAvailable
	} else if member.isLeaves(task.Start) {
		return ErrMemberIsLeaven
	}

	member.Booking = append(member.Booking, task)
	return nil
}

func (member *Member) AddLeaven(day time.Time) error {
	if !member.isAvailable() {
		return ErrMemberNotAvailable
	} else if member.isReservedDay(day) {
		return ErrMemberReserved
	}

	member.Leaves = append(member.Leaves, day)
	return nil
}

func (member *Member) isAvailable() bool {
	return member.Available
}

func (member *Member) isReadyTime(t time.Time) bool {
	return t.Equal(member.ReadyTime) || t.After(member.ReadyTime)
}

func (member *Member) isLeaves(day time.Time) bool {
	for i := 0; i < len(member.Leaves); i++ {
		if day.Format(time.DateOnly) == member.Leaves[i].Format(time.DateOnly) {
			return true
		}
	}

	return false
}

func (member *Member) isReserved(task *dao.Tasks, restTime time.Duration) bool {
	start := task.Start
	end := task.End.Add(restTime)

	for i := 0; i < len(member.Booking); i++ {
		if member.Booking[i].Id == task.Id {
			return false
		}
	}

	for i := 0; i < len(member.Booking); i++ {
		if sameOrAfter(start, member.Booking[i].Start) && sameOrBefore(start, member.Booking[i].End) {
			return true
		} else if sameOrAfter(end, member.Booking[i].Start) && sameOrBefore(end, member.Booking[i].End) {
			return true
		}
	}

	return false
}

func (member *Member) isReservedDay(day time.Time) bool {
	for i := 0; i < len(member.Booking); i++ {
		if member.Booking[i].Start.Format(time.DateOnly) == day.Format(time.DateOnly) ||
			member.Booking[i].End.Format(time.DateOnly) == day.Format(time.DateOnly) {
			return true
		}
	}

	return false
}

func sameOrBefore(t1, t2 time.Time) bool {
	return t1.Equal(t2) || t1.Before(t2)
}

func sameOrAfter(t1, t2 time.Time) bool {
	return t1.Equal(t2) || t1.After(t2)
}

func NewMemberWorker(member *dao.Members) Worker {

	return &Member{
		Id:        member.Id,
		Info:      member,
		Available: true,
		ReadyTime: cmp.Or(*member.LastTimeTask, member.CreatedAt),
	}
}
