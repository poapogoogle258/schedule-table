package dto

import (
	"schedule_table/internal/model/dao"
	"schedule_table/util"
	"strings"
	"time"

	"github.com/jinzhu/copier"
)

type ResponseCalendar struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ImageURL    string `json:"imageUrl"`
	Description string `json:"description"`
}

type ResponseMember struct {
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

type Recurrence struct {
	Freq      int8  `json:"freq"`
	Count     int32 `json:"count"`
	Interval  int32 `json:"interval"`
	Byweekday []int `json:"byweekday"`
	Bymonth   []int `json:"bymonth"`
}

type ResponseSchedule struct {
	Id               string           `json:"id"`
	MasterScheduleId *string          `json:"master_id"`
	CalendarId       string           `json:"calendar_id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	ImageURL         string           `json:"imageURL"`
	Priority         int8             `json:"priority"`
	Start            time.Time        `json:"start"`
	End              time.Time        `json:"end"`
	Hr_start         string           `json:"hr_start"`
	Hr_end           string           `json:"hr_end"`
	Tzid             string           `json:"tzid"`
	BreakTime        uint32           `json:"breaktime"`
	Recurrence       Recurrence       `json:"recurrence"`
	Members          []ResponseMember `json:"members"`
}

func (resSchedule *ResponseSchedule) Responsibles(responsibles *[]dao.Responsible) {

	for _, responsible := range *responsibles {
		member := &ResponseMember{}
		if err := copier.Copy(&member, &responsible.Person); err != nil {
			panic(err)
		}
		resSchedule.Members = append(resSchedule.Members, *member)
	}
}

func (resSchedule *ResponseSchedule) Recurrence_freq(s int8) {
	resSchedule.Recurrence.Freq = s
}
func (resSchedule *ResponseSchedule) Recurrence_count(s int32) {
	resSchedule.Recurrence.Count = s
}
func (resSchedule *ResponseSchedule) Recurrence_interval(s int32) {
	resSchedule.Recurrence.Interval = s
}
func (resSchedule *ResponseSchedule) Recurrence_byweekday(s string) {
	resSchedule.Recurrence.Byweekday = util.MapStringToInt(strings.Split(s, ","))
}
func (resSchedule *ResponseSchedule) Recurrence_bymonth(s string) {
	resSchedule.Recurrence.Bymonth = util.MapStringToInt(strings.Split(s, ","))
}

type ResponseTask struct {
	Id         string        `json:"id"`
	CalendarId string        `json:"calendar_id"`
	ScheduleId string        `json:"schedule_id"`
	Start      time.Time     `json:"start"`
	End        time.Time     `json:"end"`
	Status     int8          `json:"status"`
	Person     RequestMember `json:"person"`
}

type ResponseProfile struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	ImageURL    string  `json:"imageURL"`
	Description string  `json:"description"`
	CalendarId  *string `json:"calendar_id"`
}

func (responseProfile *ResponseProfile) Calendar(cal *dao.Calendars) {
	if cal != nil {
		id := cal.Id.String()
		responseProfile.CalendarId = &id
	} else {
		responseProfile.CalendarId = nil
	}
}

type Pagination struct {
	Total       int64 `json:"total_records"`
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
}

type ResponseMembersTable struct {
	Data       *[]ResponseMember `json:"data"`
	Pagination *Pagination       `json:"pagination"`
}
