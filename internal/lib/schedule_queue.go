package lib

import "github.com/google/uuid"

type ScheduleQueue interface {
	NextQueue() *Employee
	RotateEmployee(employeeId uuid.UUID)
	Initial()
	Size() int
	ReSetQueue()
}

type scheduleQueue struct {
	initQueue []*Employee
	queue     []*Employee
	index     int
	size      int
}

func (q *scheduleQueue) ReSetQueue() {
	q.index = -1
}

func (q *scheduleQueue) Initial() {
	q.queue = make([]*Employee, len(q.initQueue))
	copy(q.queue, q.initQueue)
}

func (q *scheduleQueue) NextQueue() *Employee {
	q.index = (q.index + 1) % q.size
	return q.queue[q.index]
}

func (q *scheduleQueue) RotateEmployee(employeeId uuid.UUID) {
	for i := range len(q.queue) {
		if q.queue[i].Id == employeeId {
			q.queue = append(q.queue[:i], append(q.queue[i+1:], q.queue[i])...)
			return
		}
	}
}

func (q *scheduleQueue) Size() int {
	return q.size
}

func NewScheduleQueue(employeesQueueOrdered []*Employee) ScheduleQueue {

	queue := make([]*Employee, len(employeesQueueOrdered))
	copy(queue, employeesQueueOrdered)

	return &scheduleQueue{
		initQueue: employeesQueueOrdered,
		queue:     queue,
		size:      len(employeesQueueOrdered),
		index:     -1,
	}
}
