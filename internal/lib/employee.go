package lib

import (
	"schedule_table/internal/model/dao"
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	Id       uuid.UUID
	Data     *dao.Employee
	TimeLine ITimeLine
}

// add work to timeline of employee
func (e *Employee) AddTask(task *dao.Task) {

	work := NewWork(task.Id, nil, task.Start, task.End, task.Description.Priority, Works)
	rest := NewWork(uuid.New(), &task.Id, task.End, task.End.Add(time.Minute*time.Duration(task.Description.BreakTime)), task.Description.Priority, Rest)

	e.TimeLine.AddWork(work)
	e.TimeLine.AddWork(rest)

	task.OriginalEmployeeId = &e.Id
	task.EmployeeId = &e.Id
	task.Person = e.Data
}

func (e *Employee) IsBusy(start time.Time, end time.Time) bool {
	return e.TimeLine.IsBusy(start, end)
}

func (e *Employee) GetWorksInRange(start time.Time, end time.Time) []*Work {
	return e.TimeLine.GetWorksInRange(start, end)
}

func NewEmployee(data *dao.Employee) *Employee {
	return &Employee{
		Id:       data.Id,
		Data:     data,
		TimeLine: NewTimeLine(),
	}
}
