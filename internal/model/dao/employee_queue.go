package dao

import (
	"github.com/google/uuid"
)

type EmployeeQueue struct {
	EmployeeId uuid.UUID `gorm:"primaryKey;type:uuid;column:employee_id" json:"member_id"`
	ScheduleId uuid.UUID `gorm:"primaryKey;type:uuid;column:schedule_id" json:"schedule_id"`
	Queue      int8      `gorm:"column:queue" json:"queue"`
	Person     Employee  `gorm:"foreignKey:employee_id" json:"person"`
}
