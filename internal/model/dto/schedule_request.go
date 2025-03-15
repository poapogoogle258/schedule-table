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
	Freq      int8  `json:"freq" binding:"required,min=0,max=6"`
	Interval  int32 `json:"interval" binding:"required,min=1"`
	Count     int32 `json:"count" binding:"required,min=0"`
	Byweekday []int `json:"byweekday" binding:"required,byweekday"`
	Bymonth   []int `json:"bymonth" binding:"required,bymonth"`
}

type Schedule struct {
	MasterScheduleId string     `json:"master_id" binding:"required"`
	Name             string     `json:"name" binding:"required"`
	Description      string     `json:"description" binding:"required"`
	ImageURL         string     `json:"imageURL" binding:"required,url"`
	Color            string     `json:"color" binding:"required,hexcolor"`
	Priority         int8       `json:"priority" binding:"required,min=1,max=99"`
	Start            time.Time  `json:"start" binding:"required" time_format:"2006-01-02"`
	End              *time.Time `json:"end" binding:"omitempty" time_format:"2006-01-02"`
	Hr_start         string     `json:"hr_start" binding:"required,hhmm"`
	Hr_end           string     `json:"hr_end" binding:"required,hhmm"`
	Tzid             string     `json:"tzid" binding:"required"`
	BreakTime        uint32     `json:"breaktime" binding:"required,min=0,max=1440"`
	RotationCycle    int        `json:"rotation_cycle" binding:"required,min=0"`
	UseNumberPeople  int8       `json:"use_number_people" binding:"required,gte=1"`
}

// schedule dto -> dao

type ScheduleInfoRequest struct {
	Schedule
	Recurrence Recurrence      `json:"recurrence" binding:"required"`
	Employees  []*EmployeeInfo `json:"employees" binding:"required,min=0"`
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
