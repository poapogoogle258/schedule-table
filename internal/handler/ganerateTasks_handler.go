package handler

import (
	"fmt"
	"net/http"
	"schedule_table/internal/constant"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"
	"schedule_table/util"
	"slices"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ######################################## Worker ###############################################

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

// ######################################## ScheduleManager ###############################################

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

// ######################################## GenerateTaskHandler ###############################################

type GenerateTaskHandler interface {
	GenerateTasks(c *gin.Context)
}

type generateTaskHandler struct {
	calRepo      repository.CalendarRepository
	scheduleRepo repository.SchedulesRepository
	recurService service.RecurrenceService
}

type RequestGenerateTasksBody struct {
	Start string `json:"start" binding:"required" ` // format : RFC3339 "2006-01-02T15:04:05+07:00"
	End   string `json:"end" binding:"required"`
}

func (gt *generateTaskHandler) GenerateTasks(c *gin.Context) {
	defer pkg.PanicHandler(c)

	body := &RequestGenerateTasksBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		pkg.PanicException(constant.InvalidRequest)
	}

	start := util.Must(time.Parse(time.RFC3339, body.Start))
	end := util.Must(time.Parse(time.RFC3339, body.End))

	userId, validUserId := c.Keys["token_userId"].(string)
	calendarId := c.Param("calendarId")

	if !validUserId {
		pkg.PanicException(constant.DataNotFound)
	}

	if calendarId != "default" && !gt.calRepo.IsOwnerCalendar(userId, calendarId) {
		pkg.PanicException(constant.DataNotFound)
	}

	if calendarId == "default" {
		calendarId = gt.calRepo.GetDefaultCalendarId(userId)
	}

	schedules := gt.scheduleRepo.GetScheduleOfCalendar(calendarId, &start, &end)
	members := gt.calRepo.GetMembersOfCalendarId(calendarId)

	workers := make(map[uuid.UUID]Worker)
	for i := 0; i < len(*members); i++ {
		member := &(*members)[i]
		workers[member.Id] = NewWorker(member)
	}

	schedulesManager := make(map[uuid.UUID]*ScheduleManager)
	schedulesTasks := make([]dao.Tasks, 0)
	for i := 0; i < len(*schedules); i++ {
		schedule := &(*schedules)[i]
		tasksSchedule := gt.recurService.GenerateScheduleTasks(schedule, &start, &end)

		fmt.Println("schedule.MasterScheduleId", schedule.MasterScheduleId, schedule.MasterScheduleId != nil)

		if schedule.MasterScheduleId == nil {
			responsiblePersons := gt.scheduleRepo.GetResponsiblePersons(schedule.Id.String())
			schedulesManager[schedule.Id] = NewScheduleManager(schedule, responsiblePersons)
		} else {
			schedulesManager[schedule.Id] = schedulesManager[*schedule.MasterScheduleId]
		}
		schedulesTasks = append(schedulesTasks, (*tasksSchedule)...)
	}

	// soft schedulesManager by Start, Priority
	slices.SortFunc(schedulesTasks, func(a, b dao.Tasks) int {
		if c := a.Start.Compare(b.Start); c == 0 {
			if a.Priority > b.Priority {
				return 1
			} else {
				return -1
			}
		} else {
			return c
		}
	})

	for i := 0; i < len(schedulesTasks); i++ {
		task := &schedulesTasks[i]

		queue := 0
		workerId := schedulesManager[task.ScheduleId].MembersQueue.Next(queue)
		for !workers[workerId].IsAvailable(task.Start) {
			queue++
			schedulesManager[task.ScheduleId].MembersQueue.Skip()
			workerId = schedulesManager[task.ScheduleId].MembersQueue.Next(queue)
		}
		schedulesManager[task.ScheduleId].MembersQueue.Select(queue)
		task.Person = *workers[workerId].GetData()
		workers[workerId].AddTask(task, schedulesManager[task.ScheduleId].BreakTime)
		task.MemberId = workerId
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, schedulesTasks))

}

func NewGenerateTaskHandler(
	calRepo repository.CalendarRepository,
	scheduleRepo repository.SchedulesRepository,
	recurService service.RecurrenceService) GenerateTaskHandler {
	return &generateTaskHandler{
		calRepo:      calRepo,
		scheduleRepo: scheduleRepo,
		recurService: recurService,
	}
}
