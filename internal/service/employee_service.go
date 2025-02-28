package service

import (
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/repository"
	"schedule_table/util"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type EmployeeService interface {
	GetEmployeesOfCalendar(calendarId string) ([]*dto.EmployeeInfo, error)
	GetPaginationEmployeesOfCalendarId(calendarId string, page int, limit int) ([]*dto.EmployeeInfo, error)
	CountEmployeesOfCalendarId(calendarId string) int64
	IsExist(employeeId string) bool
	GetEmployee(employeeId string) (*dto.EmployeeInfo, error)
	CreateEmployee(calendarId string, info *dto.EmployeeInfoRequest) (*dto.EmployeeInfo, error)
	UpdateEmployee(employeeId string, info *dto.EmployeeInfoRequest) (*dto.EmployeeInfo, error)
	DeleteEmployee(employeeId string) error
	EmployeesIsExistCalendar(calendarId string, employeeIds []string) bool
	InjectionTx(tx *gorm.DB) EmployeeService
}

type employeeService struct {
	employeeRepo repository.EmployeeRepository
}

func (s *employeeService) InjectionTx(tx *gorm.DB) EmployeeService {

	initRepo := s.employeeRepo.InjectionTx(tx)

	return &employeeService{
		employeeRepo: initRepo,
	}

}

func (s *employeeService) IsExist(employeeId string) bool {
	return s.employeeRepo.IsExist("id = ?", employeeId)
}

func (s *employeeService) GetEmployeesOfCalendar(calendarId string) ([]*dto.EmployeeInfo, error) {
	employees, err := s.employeeRepo.FindMany("calendar_id = ?", calendarId)
	if err != nil {
		return nil, err
	}
	result := []*dto.EmployeeInfo{}
	copier.Copy(&result, employees)

	return result, nil
}

func (s *employeeService) GetPaginationEmployeesOfCalendarId(calendarId string, page int, limit int) ([]*dto.EmployeeInfo, error) {
	employees, err := s.employeeRepo.FindManyPagination((page-1)*limit, page*limit, "calendar_id = ?", calendarId)
	if err != nil {
		return nil, err
	}

	result := []*dto.EmployeeInfo{}
	copier.Copy(&result, employees)

	return result, nil
}

func (s *employeeService) CountEmployeesOfCalendarId(calendarId string) int64 {
	return s.employeeRepo.Count("calendar_id = ?", calendarId)
}

func (s *employeeService) GetEmployee(employeeId string) (*dto.EmployeeInfo, error) {
	employee, err := s.employeeRepo.FindOne("id = ?", employeeId)
	if err != nil {
		return nil, err
	}

	result := dto.EmployeeInfo{}
	copier.Copy(&result, employee)

	return &result, nil
}

func (s *employeeService) CreateEmployee(calendarId string, info *dto.EmployeeInfoRequest) (*dto.EmployeeInfo, error) {

	insert := &dao.Employee{}
	copier.Copy(&insert, info)
	insert.Id = uuid.New()
	insert.CalendarId = uuid.MustParse(calendarId)

	if err := s.employeeRepo.Create(insert); err != nil {
		return nil, err
	}

	result := &dto.EmployeeInfo{}
	copier.Copy(&result, insert)

	return result, nil
}

func (s *employeeService) UpdateEmployee(employeeId string, info *dto.EmployeeInfoRequest) (*dto.EmployeeInfo, error) {

	employee := util.Must(s.employeeRepo.FindOne("id = ?", employeeId))
	copier.Copy(&employee, info)

	if err := s.employeeRepo.Save(employee); err != nil {
		return nil, err
	}

	result := &dto.EmployeeInfo{}
	copier.Copy(&result, employee)

	return result, nil
}

func (s *employeeService) DeleteEmployee(employeeId string) error {
	return s.employeeRepo.Delete("id = ?", employeeId)
}

func (s *employeeService) EmployeesIsExistCalendar(calendarId string, employeeIds []string) bool {
	for _, id := range employeeIds {
		if !s.employeeRepo.IsExist("id = ? AND calendar_id = ?", id, calendarId) {
			return false
		}
	}

	return true
}

func NewEmployeeService(employeeRepository repository.EmployeeRepository) EmployeeService {
	return &employeeService{
		employeeRepo: employeeRepository,
	}
}
