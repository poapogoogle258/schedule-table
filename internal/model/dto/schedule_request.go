package dto

import (
	"schedule_table/internal/model/dao"
	"schedule_table/util"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
)

type Recurrence struct {
	Freq      int8  `json:"freq" validate:"required,min=0,max=6"`
	Interval  int32 `json:"interval" validate:"required,min=1"`
	Count     int32 `json:"count" validate:"required,min=0"`
	Byweekday []int `json:"byweekday" validate:"required,byweekday"`
	Bymonth   []int `json:"bymonth" validate:"required,bymonth"`
}

type Schedule struct {
	MasterScheduleId string     `json:"master_id"`
	Name             string     `json:"name" validate:"required"`
	Description      string     `json:"description"`
	ImageURL         string     `json:"imageURL" validate:"required,url"`
	Color            string     `json:"color" validate:"required,hexcolor"`
	Priority         int8       `json:"priority" validate:"required,min=1,max=99"`
	Start            time.Time  `json:"start" validate:"required"`
	End              *time.Time `json:"end"`
	Hr_start         string     `json:"hr_start" validate:"required,hhmm"`
	Hr_end           string     `json:"hr_end" validate:"required,hhmm"`
	Tzid             string     `json:"tzid" validate:"required"`
	BreakTime        uint32     `json:"breaktime" validate:"required,min=0,max=1440"`
	RotationCycle    int        `json:"rotation_cycle" validate:"required,min=0"`
	UseNumberPeople  int8       `json:"use_number_people" validate:"required,gte=1"`
}

// schedule dto -> dao

type ScheduleInfoRequest struct {
	Schedule
	Recurrence Recurrence      `json:"recurrence" validate:"required"`
	Employees  []*EmployeeInfo `json:"employees" validate:"required"`
}

func (s *ScheduleInfoRequest) Recurrence_freq() int8 {
	return s.Recurrence.Freq
}

func (s *ScheduleInfoRequest) Recurrence_count() int32 {
	return s.Recurrence.Count
}
func (s *ScheduleInfoRequest) Recurrence_interval() int32 {
	return s.Recurrence.Interval
}

func (s ScheduleInfoRequest) Recurrence_bymonth() string {
	if len(s.Recurrence.Bymonth) == 0 {
		return ""
	}
	strList := util.Map(s.Recurrence.Bymonth, func(i int) string {
		return strconv.Itoa(i)
	})

	return strings.Join(strList, ",")
}

func (s ScheduleInfoRequest) Recurrence_byweekday() string {
	if len(s.Recurrence.Byweekday) == 0 {
		return ""
	}
	strList := util.Map(s.Recurrence.Byweekday, func(i int) string {
		return strconv.Itoa(i)
	})

	return strings.Join(strList, ",")
}

// schedule dao -> dto

type ScheduleInfo struct {
	Id string `json:"id"`
	ScheduleInfoRequest
}

func (s *ScheduleInfo) EmployeeQueue(employeeQueue []*dao.EmployeeQueue) {
	s.Employees = make([]*EmployeeInfo, len(employeeQueue))
	for i := range employeeQueue {
		employeeInfo := &EmployeeInfo{}
		copier.Copy(&employeeInfo, employeeQueue[i].Person)
		s.Employees[i] = employeeInfo
	}
}

func (resSchedule *ScheduleInfo) Recurrence_freq(s int8) {
	resSchedule.Recurrence.Freq = s
}

func (resSchedule *ScheduleInfo) Recurrence_count(s int32) {
	resSchedule.Recurrence.Count = s
}
func (resSchedule *ScheduleInfo) Recurrence_interval(s int32) {
	resSchedule.Recurrence.Interval = s
}
func (resSchedule *ScheduleInfo) Recurrence_byweekday(s string) {
	if s == "" {
		resSchedule.Recurrence.Byweekday = []int{}
	} else {
		resSchedule.Recurrence.Byweekday = util.MapStringToInt(strings.Split(s, ","))
	}
}
func (resSchedule *ScheduleInfo) Recurrence_bymonth(s string) {
	if s == "" {
		resSchedule.Recurrence.Bymonth = []int{}
	} else {
		resSchedule.Recurrence.Bymonth = util.MapStringToInt(strings.Split(s, ","))
	}
}
