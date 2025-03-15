package service

import (
	"schedule_table/internal/lib"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg/logger"
	"schedule_table/internal/repository"
	"schedule_table/util"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type TaskService interface {
	GetTasksOfCalendars(calendarId string, start time.Time, end time.Time) ([]*dto.TaskInfo, error)
	UpdateTask(taskId string, data *dto.TaskRequest) (*dto.TaskInfo, error)
	RefreshTasksOfCalendar(calendar *dao.Calendar, start time.Time, end time.Time) error
	IsExist(calendarId string, taskId string) bool
}

type taskService struct {
	taskRepo    repository.TaskRepository
	transaction repository.Transaction
}

func (s *taskService) IsExist(calendarId string, taskId string) bool {
	return s.taskRepo.IsExist("id = ? AND calendar_id = ?", taskId, calendarId)
}

func (s *taskService) GetTasksOfCalendars(calendarId string, start time.Time, end time.Time) ([]*dto.TaskInfo, error) {

	tasks, err := s.taskRepo.FindManyWithJoin("tasks.calendar_id = ? AND tasks.start BETWEEN ? AND ?", calendarId, start, end)
	if err != nil {
		return nil, err
	}

	result := []*dto.TaskInfo{}
	copier.Copy(&result, tasks)

	return result, nil

}

func (s *taskService) RefreshTasksOfCalendar(calendar *dao.Calendar, start time.Time, end time.Time) error {

	mapEmployee := make(map[uuid.UUID]*lib.Employee)
	for _, employee := range calendar.Employees {
		mapEmployee[employee.Id] = lib.NewEmployee(employee)
	}

	// load leave_day to employees
	for _, leave := range calendar.Leaves {
		if leave.Status == dao.Accept {
			if _, ok := mapEmployee[leave.EmployeeId]; ok {
				leaveTask := util.Must(lib.FactoryWork(leave, lib.Leave))
				mapEmployee[leave.EmployeeId].TimeLine.AddWork(leaveTask)
			}
		}
	}

	// init scheduleManage queue
	scheduleManage := make(map[uuid.UUID]lib.ScheduleManageQueue)
	masterSchedule := make(map[uuid.UUID][]uuid.UUID)
	mapSchedule := make(map[uuid.UUID]*dao.Schedule)
	for _, schedule := range calendar.Schedules {
		mapSchedule[schedule.Id] = schedule
		masterScheduleId := schedule.MasterScheduleId
		if _, ok := masterSchedule[masterScheduleId]; !ok {
			masterSchedule[masterScheduleId] = []uuid.UUID{masterScheduleId}
		} else {
			masterSchedule[masterScheduleId] = append(masterSchedule[masterScheduleId], schedule.Id)
		}
	}

	for id := range masterSchedule {
		orderQueueEmployee := OrderQueueEmployee(mapEmployee, mapSchedule[id].EmployeeQueue)
		scheduleManage[id] = lib.NewScheduleManageQueue(mapSchedule[id], masterSchedule[id], orderQueueEmployee)
	}

	// init generate task from schedule config
	tasksCreated := util.Must(s.taskRepo.FindMany("calendar_id = ? AND extra = false AND start BETWEEN ? AND ?", calendar.Id.String(), start, end))

	tasks := make([]*dao.Task, 0)
	for _, schedule := range calendar.Schedules {
		tasksGenerate := lib.CreateRecurrenceTasks(schedule, start, end)
		tasks = append(tasks, tasksGenerate...)
	}

	// init task before commit employee to task
	OverrideTasks(tasks, tasksCreated)
	slices.SortFunc(tasks, SortTasksWithPriorityAndStartTime)

	for i := range tasks {
		scheduleMasterId := mapSchedule[tasks[i].ScheduleId].MasterScheduleId
		scheduleManage[scheduleMasterId].Commit(tasks[i])
	}

	// Refresh Tasks
	db := s.transaction.Begin()
	defer db.Rollback()

	taskRepo := s.taskRepo.InjectionTx(db)

	tasksExcept := ExpectTasks(tasks, tasksCreated)
	tasksIdDelete := make([]string, len(tasksExcept))
	for i := range tasksExcept {
		tasksIdDelete[i] = tasksExcept[i].Id.String()
	}
	if len(tasksIdDelete) > 0 {
		if err := taskRepo.Delete("id IN ?", tasksIdDelete); err != nil {
			return err
		}
	}

	if err := taskRepo.Upsert(tasks); err != nil {
		return err
	}

	db.Commit()
	return nil

}

func (s *taskService) UpdateTask(taskId string, data *dto.TaskRequest) (*dto.TaskInfo, error) {

	insert := util.Must(s.taskRepo.FindOne("id = ?", taskId))
	insert.Start = data.Start
	insert.End = data.End
	insert.Status = data.Status
	insert.EmployeeId = nil
	if data.Person != nil {
		insert.EmployeeId = &data.Person.Id
	}

	if err := s.taskRepo.Save(insert); err != nil {
		return nil, err
	}

	task := util.Must(s.taskRepo.FindOneWithJoin("id = ?", taskId))
	result := &dto.TaskInfo{}
	copier.Copy(result, task)

	return result, nil

}

func NewTaskService(taskRepo repository.TaskRepository, transaction repository.Transaction) TaskService {
	return &taskService{
		taskRepo:    taskRepo,
		transaction: transaction,
	}
}

func SortTasksWithPriorityAndStartTime(a, b *dao.Task) int {
	if a.Priority > b.Priority {
		return -1
	} else if a.Priority < b.Priority {
		return 1
	} else {
		if a.Start.Before(b.Start) {
			return -1
		} else if a.Start.After(b.Start) {
			return 1
		} else {
			return 0
		}
	}
}

func OrderQueueEmployee(mapEmployee map[uuid.UUID]*lib.Employee, queue []*dao.EmployeeQueue) []*lib.Employee {
	result := make([]*lib.Employee, len(queue))
	for i := range queue {
		if _, ok := mapEmployee[queue[i].EmployeeId]; !ok {
			logger.Warn("orderQueueEmployee: employee not found")
			continue
		}
		result[i] = mapEmployee[queue[i].EmployeeId]
	}

	return result
}

func OverrideTasks(tasks []*dao.Task, override []*dao.Task) {

	mapOverride := make(map[string]*dao.Task)
	for i := range override {
		mapOverride[override[i].RecurrenceId] = override[i]
	}

	for i := range tasks {
		if _, ok := mapOverride[tasks[i].RecurrenceId]; ok {
			tasks[i] = mapOverride[tasks[i].RecurrenceId]
		}
	}

}

func ExpectTasks(tasks []*dao.Task, expect []*dao.Task) []*dao.Task {
	result := make([]*dao.Task, 0)
	mapExpect := make(map[string]*dao.Task)
	for _, task := range expect {
		mapExpect[task.RecurrenceId] = task
	}

	for i := range tasks {
		if _, ok := mapExpect[tasks[i].RecurrenceId]; !ok {
			result = append(result, tasks[i])
		}
	}

	return result

}
