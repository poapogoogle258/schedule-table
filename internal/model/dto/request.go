package dto

import (
	"schedule_table/util"
	"strconv"
	"strings"
	"time"
)

type RequestMember struct {
	Id          string `json:"id"`
	ImageURL    string `json:"imageURL"`
	Name        string `json:"name"`
	Nickname    string `json:"nickname"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Position    string `json:"position"`
	Email       string `json:"email"`
	Telephone   string `json:"telephone"`
}

type RequestSchedule struct {
	MasterScheduleId *string    `json:"master_id"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	ImageURL         string     `json:"imageURL"`
	Priority         int8       `json:"priority"`
	Start            time.Time  `json:"start"`
	End              *time.Time `json:"end"`
	Hr_start         string     `json:"hr_start"`
	Hr_end           string     `json:"hr_end"`
	Tzid             string     `json:"tzid"`
	BreakTime        uint32     `json:"breaktime"`
	UseNumberPeople  int8       `json:"use_number_people"`
	Recurrence       Recurrence `json:"recurrence"`
	// Members          []ResponseMember `json:"members"`
}

func (reqSchedule *RequestSchedule) Validate() error {
	return nil
}

func (reqSchedule *RequestSchedule) Recurrence_freq() int8 {
	return reqSchedule.Recurrence.Freq
}
func (reqSchedule *RequestSchedule) Recurrence_interval() int32 {
	return reqSchedule.Recurrence.Interval
}
func (reqSchedule *RequestSchedule) Recurrence_count() int32 {
	return reqSchedule.Recurrence.Count
}
func (reqSchedule *RequestSchedule) Recurrence_bymonth() string {
	if len(reqSchedule.Recurrence.Bymonth) == 0 {
		return ""
	} else {
		sl := util.Map(reqSchedule.Recurrence.Bymonth, strconv.Itoa)
		return strings.Join(sl, ",")
	}

}
func (reqSchedule *RequestSchedule) Recurrence_byweekday() string {
	if len(reqSchedule.Recurrence.Byweekday) == 0 {
		return ""
	} else {
		sl := util.Map(reqSchedule.Recurrence.Byweekday, strconv.Itoa)
		return strings.Join(sl, ",")
	}
}

// func (reqSchedule *RequestSchedule) Responsibles() *[]dao.EmployeeQueue {

// 	if reqSchedule.MasterScheduleId != nil {
// 		return nil
// 	}

// 	responsibles := make([]dao.EmployeeQueue, 0, len(reqSchedule.Members))
// 	for i, member := range reqSchedule.Members {
// 		responsibles = append(responsibles, dao.EmployeeQueue{
// 			Queue:      int8(i),
// 			EmployeeId: uuid.MustParse(member.Id),
// 		})
// 	}
// 	return &responsibles

// }
