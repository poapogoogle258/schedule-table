package service

import (
	"errors"
	"schedule_table/internal/constant"
	"schedule_table/internal/model/dao"
	"schedule_table/util"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	rrule "github.com/teambition/rrule-go"
)

type IScheduleService interface {
	NewSchedule(schedule *dao.Schedules) ISchedule
}

type ScheduleService struct{}

type ISchedule interface {
	GetId() uuid.UUID
	GetCalendarId() uuid.UUID
	GetPriority() int8
	GenerateTasks(start time.Time, end time.Time) *[]dao.Tasks
}

type Schedule struct {
	Id                  uuid.UUID
	CalendarId          uuid.UUID
	Priority            int8
	TimeStart           hrTime
	TimeEnd             hrTime
	Tzid                *time.Location
	UseNumberPeople     int
	RecurrenceFreq      int
	RecurrenceInterval  int
	RecurrenceByMonth   []int8
	RecurrenceByWeekday []rrule.Weekday
}

type hrTime struct {
	H int
	M int
}

func (schedule *Schedule) GetId() uuid.UUID {
	return schedule.Id
}

func (schedule *Schedule) GetCalendarId() uuid.UUID {
	return schedule.CalendarId
}

func (schedule *Schedule) GetPriority() int8 {
	return schedule.Priority
}

func (schedule *Schedule) GenerateTasks(start time.Time, end time.Time) *[]dao.Tasks {
	var option rrule.ROption

	option.Freq = rrule.Frequency(schedule.RecurrenceFreq)
	option.Interval = schedule.RecurrenceInterval

	if schedule.RecurrenceFreq == constant.DAILY && len(schedule.RecurrenceByWeekday) != 0 {
		option.Byweekday = schedule.RecurrenceByWeekday
	}

	option.Dtstart = time.Date(start.Year(), start.Month(), start.Day(), schedule.TimeStart.H, schedule.TimeStart.M, 0, 0, schedule.Tzid)
	option.Until = time.Date(end.Year(), end.Month(), end.Day(), schedule.TimeEnd.H, schedule.TimeEnd.M, 0, 0, schedule.Tzid)

	scheduleRule, ErrRule := rrule.NewRRule(option)
	if ErrRule != nil {
		panic(ErrRule)
	}

	duration := hrTimeDuration(schedule.TimeStart, schedule.TimeEnd)

	recurrenceSchedule := scheduleRule.All()
	tasks := make([]dao.Tasks, 0, len(recurrenceSchedule)*schedule.UseNumberPeople)

	for i := 0; i < len(recurrenceSchedule); i++ {
		for number := 1; number <= schedule.UseNumberPeople; number++ {
			tasks = append(tasks, dao.Tasks{
				Id:         uuid.New(),
				ScheduleId: schedule.Id,
				CalendarId: schedule.CalendarId,
				Priority:   schedule.Priority,
				Start:      recurrenceSchedule[i],
				End:        recurrenceSchedule[i].Add(duration),
			})
		}
	}

	return &tasks

}

func (scheService *ScheduleService) NewSchedule(schedule *dao.Schedules) ISchedule {
	service := &Schedule{}

	service.Id = schedule.Id
	service.CalendarId = schedule.CalendarId
	service.Priority = schedule.Priority
	service.TimeStart = parseHrTime(schedule.Hr_start)
	service.TimeEnd = parseHrTime(schedule.Hr_end)

	location, errLoadLocation := time.LoadLocation(schedule.Tzid)
	if errLoadLocation != nil {
		panic(errLoadLocation)
	}
	service.Tzid = location
	service.UseNumberPeople = int(schedule.UseNumberPeople)

	service.RecurrenceFreq = int(schedule.Recurrence_freq)
	service.RecurrenceInterval = int(schedule.Recurrence_count)
	service.RecurrenceByMonth = util.Map(strings.Split(schedule.Recurrence_bymonth, ","), getInt8)
	service.RecurrenceByWeekday = util.Map(strings.Split(schedule.Recurrence_byweekday, ","), getWeekDay)

	return service
}

func parseHrTime(t string) hrTime {
	splitTime := strings.Split(t, ":")
	_hrTime, errHrtime := strconv.Atoi(splitTime[0])
	if errHrtime != nil {
		panic(errHrtime)
	}
	_mnTime, errMntime := strconv.Atoi(splitTime[1])
	if errMntime != nil {
		panic(errMntime)
	}

	return hrTime{
		H: _hrTime,
		M: _mnTime,
	}
}

func getWeekDay(d string) rrule.Weekday {
	switch d {
	case "0":
		return rrule.MO
	case "1":
		return rrule.TU
	case "2":
		return rrule.WE
	case "3":
		return rrule.TH
	case "4":
		return rrule.FR
	case "5":
		return rrule.SA
	case "6":
		return rrule.SU
	default:
		panic(errors.New("getWeekDay: not exits weekday"))
	}
}

func getInt8(s string) int8 {
	i, _ := strconv.Atoi(s)
	return int8(i)
}

func hrTimeDuration(a, b hrTime) time.Duration {
	nanosecondStart := (time.Duration(a.H) * time.Hour) + (time.Duration(a.M) * time.Minute)
	nanosecondEnd := (time.Duration(b.H) * time.Hour) + (time.Duration(b.M) * time.Minute)

	if nanosecondStart > nanosecondEnd {
		nanosecondEnd = nanosecondEnd + (time.Duration(24) * time.Hour)
	}

	return time.Duration(nanosecondEnd - nanosecondStart)
}

func NewScheduleService() IScheduleService {
	return &ScheduleService{}
}
