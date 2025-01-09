package service

import (
	"errors"
	"schedule_table/internal/model/dao"

	"github.com/google/uuid"
)

var ErrorSkipAllQueue = errors.New("skip all queue in list")

type IQueue interface {
	Next(i int) Worker
	Select(i int)
	Skip()
}

type QueueResponsible struct {
	ScheduleId uuid.UUID
	SkipIndex  int
	Members    []Worker
}

func (queue *QueueResponsible) Next(i int) Worker {
	if i >= len(queue.Members) {
		panic(ErrorSkipAllQueue)
	}

	return queue.Members[i]
}

func (queue *QueueResponsible) Select(i int) {
	if i != len(queue.Members)-1 {
		queue.Members = append(queue.Members[:i], append(queue.Members[i+1:], queue.Members[i])...) // move worker to last queue
	}
	// TO DO : skip queue in loop

	queue.SkipIndex = 0
}

func (queue *QueueResponsible) Skip() {
	queue.SkipIndex++
}

// HAVE TO: Sort Responsible By Queue.
func NewResponsibleQueue(scheduleId uuid.UUID, responsible *[]dao.Responsible) IQueue {

	len_responsible := len(*responsible)
	members := make([]Worker, len_responsible)

	for i := 0; i < len(members); i++ {
		members[i] = NewMemberWorker(&(*responsible)[i].Person)
	}

	return &QueueResponsible{
		ScheduleId: scheduleId,
		Members:    members,
		SkipIndex:  0,
	}
}
