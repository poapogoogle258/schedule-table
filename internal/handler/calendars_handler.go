package handler

import (
	"fmt"
	"net/http"
	"schedule_table/internal/constant"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"
	"time"

	"github.com/gin-gonic/gin"
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
	// defer pkg.PanicHandler(c)

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

	// create members part
	workersMap := service.MapWorkers{}
	for _, member := range *members {
		workersMap[member.Id.String()] = &service.Worker{
			Id:       member.Id,
			Member:   &member,
			RestTime: *GetValueOrDefaultTime(member.LastTimeTask),
		}
	}
	for _, leaves := range *leaves {
		workersMap[leaves.MemberId.String()].Leaves = service.GetBetweens(&leaves.Start, &leaves.End)
	}

	// create schedule part
	scheduleTasks := make([]service.Schedule, len(*schedules))
	for i, schedule := range *schedules {

		tasks_schedule := s.recurService.GenerateScheduleTasks(&schedule, &start, &end)
		responsible_persons := s.scheduleRepo.GetResponsiblePersons(schedule.Id.String())

		scheduleTasks[i].Name = &schedule.Name
		scheduleTasks[i].Priority = int(schedule.Priority)
		scheduleTasks[i].Tasks = tasks_schedule
		scheduleTasks[i].TasksDaily = service.NewTasksDaily(tasks_schedule)
		scheduleTasks[i].ListMembers = service.ListResponsible{
			Members: responsible_persons,
		}
	}

	// processing part

	// TO DO CHECK INTINITY LOOP
	// TO DO Sort Tasks in TasksDaily

	// start, end := now.BeginningOfDay(), now.EndOfMonth()

	for start.Before(end) {
		date := start.Format(time.DateOnly)
		fmt.Printf("Day : %s\n", date)

		for i := 0; i < len(scheduleTasks); i++ {
			if tasksDaily, ok := (*scheduleTasks[i].TasksDaily)[date]; ok {
				fmt.Printf("TasksDaily of schedule \"%s\"\n", *scheduleTasks[i].Name)

				for j := 0; j < len(tasksDaily); j++ {
					fmt.Printf("task[%d] : start(%s) , end(%s) \n", j+1, tasksDaily[j].Start.Format(time.DateOnly), tasksDaily[j].End.Format(time.DateOnly))

					index := 0
					queue := scheduleTasks[i].ListMembers.Next(index)
					fmt.Printf("Select Member : %s\n", queue.MemberId.String())

					for {
						if workersMap.CheckWorkerFree(queue.MemberId.String(), tasksDaily[j].Start) {
							tasksDaily[j].MemberId = queue.MemberId
							fmt.Printf("Set Member : %s into task: %s\n", queue.MemberId.String(), tasksDaily[j].Id.String())
							break
						} else {
							index++
							scheduleTasks[i].ListMembers.Skip()
							queue = scheduleTasks[i].ListMembers.Next(index)
							fmt.Printf("Skip Member: %s\n", queue.MemberId.String())
						}
					}
				}
			}
		}

		start = start.Add(time.Hour * 24)
	}

	// -------------------------------------------------------------

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, scheduleTasks[0].Tasks))

}

func GetValueOrDefaultTime(t *time.Time) *time.Time {
	if t == nil {
		now := time.Now()
		startOfDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

		return &startOfDate
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
