package service

import (
	"fmt"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/repository"
	"time"

	"github.com/jinzhu/copier"
)

type CalendarService interface {
	GetCalendarsOfUser(userId string) ([]*dto.CalendarInfo, error)
	FindCalendarIdWithFullAggregate(calendarId string, start time.Time, end time.Time) (*dao.Calendar, error)
	IsExist(calendarId string) bool
}

type calendarService struct {
	calRepo repository.CalendarRepository
}

func (s *calendarService) FindCalendarIdWithFullAggregate(calendarId string, start time.Time, end time.Time) (*dao.Calendar, error) {
	return s.calRepo.FindCalendarIdWithFullAggregate(calendarId, start, end)

}

func (s *calendarService) GetCalendarsOfUser(userId string) ([]*dto.CalendarInfo, error) {
	calendars, err := s.calRepo.FindMany("user_id = ?", userId)
	if err != nil {
		return nil, err
	}

	result := []*dto.CalendarInfo{}
	if err := copier.Copy(&result, &calendars); err != nil {
		return nil, err
	}

	fmt.Printf("test")

	return result, nil

}

func (s *calendarService) IsExist(calendarId string) bool {
	return s.calRepo.IsExist("id = ?", calendarId)
}

func NewCalendarService(calRepo repository.CalendarRepository) CalendarService {
	return &calendarService{
		calRepo: calRepo,
	}
}
