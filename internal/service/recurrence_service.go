package service

import (
	"schedule_table/internal/constant"
	"schedule_table/internal/model/dao"
	"strconv"
	"strings"
	"time"

	rrule "github.com/teambition/rrule-go"
)

type RecurrenceService interface {
	NewScheduleRuleConfig(schedule *dao.Schedules, start *time.Time, end *time.Time) *rrule.ROption
}

type Recurrence struct {
}

func (s *Recurrence) NewScheduleRuleConfig(schedule *dao.Schedules, start *time.Time, end *time.Time) *rrule.ROption {

	var config rrule.ROption

	config.Freq = rrule.Frequency(schedule.Recurrence_freq)
	config.Interval = int(schedule.Recurrence_interval)
	config.Count = 31

	if schedule.Recurrence_freq == int8(constant.DAILY) && schedule.Recurrence_byweekday != "" {
		config.Byweekday = mapWeekdayFromStringSplit(schedule.Recurrence_byweekday, ",")
	}

	hr_start := mapIntFromStringSplit(schedule.Hr_start, ":")
	hr_end := mapIntFromStringSplit(schedule.Hr_end, ":")

	location, errLoadLocation := time.LoadLocation(schedule.Tzid)
	if errLoadLocation != nil {
		panic(errLoadLocation)
	}

	config.Dtstart = time.Date(start.Year(), start.Month(), start.Day(), hr_start[0], hr_start[1], 0, 0, location)
	config.Until = time.Date(end.Year(), end.Month(), end.Day(), hr_end[0], hr_end[1], 0, 0, location)

	return &config

}

func mapIntFromStringSplit(source string, s string) []int {
	splitSource := strings.Split(source, s)

	results := make([]int, len(splitSource))

	for i, _ := range results {
		results[i], _ = strconv.Atoi(splitSource[i])
	}

	return results
}

func mapWeekdayFromStringSplit(source string, s string) []rrule.Weekday {
	splitSource := strings.Split(source, s)

	MapStringWeekend := map[string]rrule.Weekday{
		"0": rrule.MO,
		"1": rrule.TU,
		"2": rrule.WE,
		"3": rrule.TH,
		"4": rrule.FR,
		"5": rrule.SA,
		"6": rrule.SU,
	}

	results := make([]rrule.Weekday, len(splitSource))

	for i, _ := range results {
		results[i] = MapStringWeekend[splitSource[i]]
	}

	return results
}

func NewRecurrenceService() RecurrenceService {
	return &Recurrence{}
}
