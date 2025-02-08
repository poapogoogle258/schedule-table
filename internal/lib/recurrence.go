package lib

import (
	"errors"
	"fmt"
	"schedule_table/internal/constant"
	"schedule_table/internal/model/dao"
	"schedule_table/util"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"time"

	rrule "github.com/teambition/rrule-go"
)

func CreateRecurrenceTasks(schedule *dao.Schedules, start time.Time, end time.Time) []*dao.Tasks {
	var option rrule.ROption

	location, errLoadLocation := time.LoadLocation(schedule.Tzid)
	if errLoadLocation != nil {
		panic(errLoadLocation)
	}
	timeStart := parseHrTime(schedule.Hr_start)
	timeEnd := parseHrTime(schedule.Hr_end)
	useNumberPeople := int(schedule.UseNumberPeople)
	recurrenceFreq := int(schedule.Recurrence_freq)
	recurrenceInterval := int(schedule.Recurrence_count)

	recurrenceByMonth := []int{}
	if schedule.Recurrence_bymonth != "" {
		recurrenceByMonth = util.Map(strings.Split(schedule.Recurrence_bymonth, ","), getInt)
	}
	recurrenceByWeekday := []rrule.Weekday{}
	if schedule.Recurrence_byweekday != "" {
		recurrenceByWeekday = util.Map(strings.Split(schedule.Recurrence_byweekday, ","), getWeekDay)
	}

	option.Freq = rrule.Frequency(recurrenceFreq)
	option.Interval = recurrenceInterval

	if len(recurrenceByWeekday) > 0 {
		option.Byweekday = recurrenceByWeekday
	}

	if len(recurrenceByMonth) > 0 {
		option.Bymonth = recurrenceByMonth
	}

	if schedule.Start.Before(start) {
		option.Dtstart = schedule.Start
	} else {
		option.Dtstart = time.Date(start.Year(), start.Month(), start.Day(), timeStart.H, timeStart.M, 0, 0, location)
	}

	if schedule.End != nil && schedule.End.After(end) {
		option.Until = *schedule.End
	} else {
		option.Until = time.Date(end.Year(), end.Month(), end.Day(), timeEnd.H, timeEnd.M, 0, 0, location)
	}

	scheduleRule, ErrRule := rrule.NewRRule(option)
	if ErrRule != nil {
		panic(ErrRule)
	}

	duration := hrTimeDuration(timeStart, timeEnd)

	recurrenceSchedule := scheduleRule.All()
	tasks := make([]*dao.Tasks, 0, len(recurrenceSchedule)*useNumberPeople)

	for i := 0; i < len(recurrenceSchedule); i++ {
		for number := 1; number <= useNumberPeople; number++ {
			tasksStart := recurrenceSchedule[i]
			tasksEnd := recurrenceSchedule[i].Add(duration)
			resttime := tasksEnd.Add(time.Duration(schedule.BreakTime) * time.Minute)
			recurrenceId := fmt.Sprint(schedule.Id.String(), "-", tasksStart.Format(time.DateOnly), "#", number)

			tasks = append(tasks, &dao.Tasks{
				Id:           uuid.New(),
				ScheduleId:   schedule.Id,
				CalendarId:   schedule.CalendarId,
				RecurrenceId: recurrenceId,
				Status:       constant.TaskStatus_Created,
				Priority:     schedule.Priority,
				Start:        tasksStart,
				End:          tasksEnd,
				RestTime:     resttime,
				Description:  *schedule,
			})
		}
	}

	return tasks
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

func getInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
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

type hrTime struct {
	H int
	M int
}

func hrTimeDuration(a, b hrTime) time.Duration {
	nanosecondStart := (time.Duration(a.H) * time.Hour) + (time.Duration(a.M) * time.Minute)
	nanosecondEnd := (time.Duration(b.H) * time.Hour) + (time.Duration(b.M) * time.Minute)

	if nanosecondStart > nanosecondEnd {
		nanosecondEnd = nanosecondEnd + (time.Duration(24) * time.Hour)
	}

	return time.Duration(nanosecondEnd - nanosecondStart)
}
