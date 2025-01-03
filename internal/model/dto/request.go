package dto

import (
	"errors"
	"mime/multipart"
	"regexp"
	"schedule_table/internal/model/dao"
	"schedule_table/util"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RequestMember struct {
	ImageURL    string `json:"imageURL"`
	Name        string `json:"name"`
	Nickname    string `json:"nickname"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Position    string `json:"position"`
	Email       string `json:"email"`
	Telephone   string `json:"telephone"`
}

type RequestCreateNewMember struct {
	Name        string                `form:"name"`
	NickName    string                `form:"nickname"`
	Color       string                `form:"color"`
	Description string                `form:"description"`
	Position    string                `form:"position"`
	Email       string                `form:"email"`
	Telephone   string                `form:"telephone"`
	File        *multipart.FileHeader `form:"image"`
}

func (reqCreateMember *RequestCreateNewMember) ImageURL() string {
	if reqCreateMember.File != nil {
		return reqCreateMember.File.Filename
	}

	return ""
}

func (newMember *RequestCreateNewMember) Validate() error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if newMember.Email != "" && !emailRegex.MatchString(newMember.Email) {
		return errors.New("field 'Email' format validate failed")
	}

	colorRegex := regexp.MustCompile(`^#[a-fA-F0-9]{6}$`)
	if newMember.Color != "" && !colorRegex.MatchString(newMember.Color) {
		return errors.New("field 'Color' format validate failed")
	}

	telephoneRegex := regexp.MustCompile(`^[0-9]{10}$`)
	if newMember.Telephone != "" && !telephoneRegex.MatchString(newMember.Telephone) {
		return errors.New("field 'Telephone' format validate failed")
	}

	return nil
}

type RequestSchedule struct {
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
	sl := util.Map(reqSchedule.Recurrence.Bymonth, strconv.Itoa)
	return strings.Join(sl, ",")
}
func (reqSchedule *RequestSchedule) Recurrence_byweekday() string {
	sl := util.Map(reqSchedule.Recurrence.Byweekday, strconv.Itoa)
	return strings.Join(sl, ",")
}

func (reqSchedule *RequestSchedule) Responsibles() *[]dao.Responsible {

	if reqSchedule.MasterScheduleId != nil {
		return nil
	}

	responsibles := make([]dao.Responsible, 0, len(reqSchedule.Members))
	for i, member := range reqSchedule.Members {
		responsibles = append(responsibles, dao.Responsible{
			Queue:    int8(i),
			MemberId: uuid.MustParse(member.Id),
		})
	}
	return &responsibles

}
