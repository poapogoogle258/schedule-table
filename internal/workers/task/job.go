package worker

import (
	"time"

	"github.com/google/uuid"
)

type Table string

const (
	TableSchedule Table = "schedule"
	TableMember   Table = "member"
	TableLeave    Table = "leave"
)

type Job struct {
	Id         uuid.UUID
	CalendarId string
	Table      Table
	UpdatedAt  time.Time
}

type JobQueue struct {
	Queue chan Job
}

func NewJob(calendarId string, table string) Job {
	return Job{
		Id:         uuid.New(),
		CalendarId: calendarId,
		Table:      Table(table),
		UpdatedAt:  time.Now(),
	}
}

func NewJobQueue(size int) *JobQueue {
	return &JobQueue{
		Queue: make(chan Job, size),
	}
}
