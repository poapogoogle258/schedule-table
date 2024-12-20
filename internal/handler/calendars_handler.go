package handler

import (
	"net/http"
	"schedule_table/internal/constant"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/now"
)

type CalendarsHandler interface {
	GetMyCalendar(c *gin.Context)
	GenerateTasks(c *gin.Context)
}

type CalendarsHandlerImpl struct {
	calRepo      repository.CalendarRepository
	scheduleRepo repository.SchedulesRepository
	recurService service.RecurrenceService
}

func (s *CalendarsHandlerImpl) GetMyCalendar(c *gin.Context) {
	defer pkg.PanicHandler(c)

	if userId, check := c.Keys["token_userId"]; check {
		calendar := s.calRepo.FindByOwnerId(userId.(string))
		c.JSON(200, pkg.BuildResponse(constant.Success, calendar))
	} else {
		pkg.PanicException(constant.DataNotFound)
	}
}

func (s *CalendarsHandlerImpl) GenerateTasks(c *gin.Context) {
	defer pkg.PanicHandler(c)

	userId, validUserId := c.Keys["token_userId"].(string)
	calendarId := c.Param("calendarId")

	if !validUserId {
		pkg.PanicException(constant.DataNotFound)
	}

	if calendarId != "default" && !s.calRepo.IsOwnerCalendar(userId, calendarId) {
		pkg.PanicException(constant.DataNotFound)
	}

	if calendarId == "default" {
		calendarId = s.calRepo.GetDefaultCalendarId(userId)
	}

	start, end := now.BeginningOfDay(), now.EndOfMonth()

	schedules := s.scheduleRepo.GetScheduleOfCalendar(calendarId, &start, &end)
	leaves := s.calRepo.GetLeavesOfCalendarId(calendarId, &start, &end)
	members := s.calRepo.GetMembersOfCalendarId(calendarId)

	workersTable := service.NewMapWorker(members, leaves)
	schedulesManager := make(map[uuid.UUID]*service.ScheduleManager)
	schedulesTasks := make([]dao.Tasks, 0)

	for i := 0; i < len(*schedules); i++ {
		schedule := &(*schedules)[i]
		tasksSchedule := s.recurService.GenerateScheduleTasks(schedule, &start, &end)
		responsiblePersons := s.scheduleRepo.GetResponsiblePersons(schedule.Id.String())
		schedulesManager[schedule.Id] = service.NewScheduleManager(schedule, responsiblePersons)
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
		workerId := schedulesManager[task.ScheduleId].Next(queue)
		for !workersTable.IsAvailable(workerId, task.Start) {
			queue++
			schedulesManager[task.ScheduleId].Skip()
			workerId = schedulesManager[task.ScheduleId].Next(queue)
		}
		schedulesManager[task.ScheduleId].Select(queue)
		workersTable.AddTask(workerId, task)
		task.MemberId = workerId
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, schedulesTasks))

}

func GetValueOrDefaultTime(t *time.Time, defaultValue *time.Time) *time.Time {
	if t == nil {
		return defaultValue
	}
	return t
}

func NewCalendarsHandler(
	calRepo repository.CalendarRepository,
	scheduleRepo repository.SchedulesRepository,
	recurService service.RecurrenceService) CalendarsHandler {
	return &CalendarsHandlerImpl{
		calRepo:      calRepo,
		scheduleRepo: scheduleRepo,
		recurService: recurService,
	}
}
