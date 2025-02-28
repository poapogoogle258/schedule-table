package dao

import (
	"time"

	"github.com/google/uuid"
)

type Leave struct {
	Id          uuid.UUID `gorm:"primarykey;type:uuid;column:id" json:"id"`
	CalendarId  uuid.UUID `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	EmployeeId  uuid.UUID `gorm:"type:uuid;column:member_id" json:"member_id"`
	UserId      uuid.UUID `gorm:"type:uuid;column:user_id" json:"user_id"`
	Date        time.Time `gorm:"column:date" json:"date"`
	Tzid        string    `gorm:"column:tzid;default:Asia/Bangkok" json:"tzid"`
	Description string    `gorm:"column:description" json:"description"`
	Employee    Employee  `gorm:"references:id"`
	BaseModel
}
