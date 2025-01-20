package handler

import (
	"errors"
	"net/http"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LeaveHandler interface {
	GetLeave(c *gin.Context) (*[]dao.Leaves, error)
	CreateNewLeave(c *gin.Context) (*dao.Leaves, error)
	Delete(c *gin.Context) error
}

type leaveHandler struct {
	CalRepo   repository.CalendarRepository
	LeaveRepo repository.LeaveRepository
}

type GetLeaveQueryString struct {
	Start string `form:"start"`
	End   string `form:"end"`
}

func (leaveHd *leaveHandler) GetLeave(c *gin.Context) (*[]dao.Leaves, error) {
	var query GetLeaveQueryString
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, pkg.NewErrorWithStatusCode(400, errors.New("ban request request start, end query string"))
	}

	calendarId := c.Param("calendarId")
	if err := leaveHd.CalRepo.CheckExist(calendarId); err != nil {
		return nil, err
	}

	start := util.Must(time.Parse(time.DateOnly, query.Start))
	end := util.Must(time.Parse(time.DateOnly, query.End))

	leaves, _ := leaveHd.LeaveRepo.Find("leaves.calendar_id = @calendarId AND leaves.date BETWEEN @start AND @end", map[string]interface{}{
		"calendarId": calendarId,
		"start":      start,
		"end":        end,
	})

	return leaves, nil
}

type CreateNewLeaveBody struct {
	MemberId string `json:"member_id"`
	Date     string `json:"date"` // format DateOnly exp. 2006-01-02
	Tzid     string `json:"tzid"`
}

func (leaveHd *leaveHandler) CreateNewLeave(c *gin.Context) (*dao.Leaves, error) {
	var body CreateNewLeaveBody
	if err := c.ShouldBind(&body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(400, errors.New("bad request body request date in body"))
	}

	calendarId := c.Param("calendarId")
	if err := leaveHd.CalRepo.CheckExist(calendarId); err != nil {
		return nil, err
	}

	userId := c.GetString("requestAuthUserId")
	dateOnly := NewDateOnlyFormat(body.Date)
	location, errLoadLocation := time.LoadLocation(body.Tzid)
	if errLoadLocation != nil {
		return nil, pkg.NewErrorWithStatusCode(400, errors.New("bad request request body tzid invalid"))
	}
	date := time.Date(dateOnly.Year, time.Month(dateOnly.Month), dateOnly.Date, 0, 0, 0, 0, location)

	insert := &dao.Leaves{
		CalendarId: uuid.MustParse(calendarId),
		MemberId:   uuid.MustParse(body.MemberId),
		UserId:     uuid.MustParse(userId),
		Date:       date,
		Tzid:       body.Tzid,
	}

	if err := leaveHd.LeaveRepo.Create(insert); err != nil {
		return nil, pkg.NewErrorWithStatusCode(500, err)
	}

	return insert, nil

}

type DateOnlyFormat struct {
	Date  int
	Month int
	Year  int
}

func NewDateOnlyFormat(s string) *DateOnlyFormat {
	spliced := util.MapStringToInt(strings.Split(s, "-"))
	return &DateOnlyFormat{
		Date:  spliced[2],
		Month: spliced[1],
		Year:  spliced[0],
	}
}

func (leaveHd *leaveHandler) Delete(c *gin.Context) error {
	leaveId := c.Param("leaveId")
	calendarId := c.Param("calendarId")

	if err := leaveHd.CalRepo.CheckExist(calendarId); err != nil {
		return err
	}

	if !leaveHd.LeaveRepo.Exits(leaveId) {
		return pkg.NewErrorWithStatusCode(http.StatusBadRequest, errors.New("not found leave id"))
	}

	if err := leaveHd.LeaveRepo.Delete(leaveId); err != nil {
		return err
	}

	return nil
}

func NewLeaveHandler(calRepo repository.CalendarRepository, leaveRepo repository.LeaveRepository) LeaveHandler {
	return &leaveHandler{
		CalRepo:   calRepo,
		LeaveRepo: leaveRepo,
	}
}
