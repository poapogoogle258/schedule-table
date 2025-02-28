package repository

import (
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	Repository[*dao.Employee]
	InjectionTx(tx *gorm.DB) EmployeeRepository
}

type employeeRepository struct {
	Repository[*dao.Employee]
	db *gorm.DB
}

func (repo *employeeRepository) InjectionTx(tx *gorm.DB) EmployeeRepository {
	return &employeeRepository{
		Repository: NewRepository[*dao.Employee](tx),
		db:         tx,
	}
}

// func (repo *employeeRepository) Update(insert *dao.Employee) error {
// 	return repo.db.Scopes(allowColumnUpdate).Save(insert).Error
// }

// func allowColumnUpdate(db *gorm.DB) *gorm.DB {
// 	selectedField := []string{"Id", "ImageURL", "Name", "Nickname", "Color", "Description", "Position", "Email", "Telephone"}
// 	return db.Select(selectedField)
// }

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{
		Repository: NewRepository[*dao.Employee](db),
		db:         db,
	}
}
