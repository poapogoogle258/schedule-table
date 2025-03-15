package dto

import (
	"schedule_table/internal/constant"
	"time"

	"github.com/google/uuid"
)

type EmployeeShort struct {
	Id       uuid.UUID `json:"id" binding:"required,uuid4"`
	Name     string    `json:"name" binding:"required"`
	NickName string    `json:"nickname" binding:"required"`
	Position string    `json:"position" binding:"required"`
	ImageURL string    `json:"imageURL" binding:"required,url"`
}

type ScheduleShort struct {
	Id       uuid.UUID `json:"id" binding:"required,uuid4"`
	Name     string    `json:"name" binding:"required"`
	ImageURL string    `json:"imageURL" binding:"required,url"`
	Color    string    `json:"color" binding:"required,hexcolor"`
}

type TaskRequest struct {
	Start       time.Time           `json:"start" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	End         time.Time           `json:"end" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	Status      constant.TaskStatus `json:"status" binding:"required,oneof=1 2 3"`
	Person      *EmployeeShort      `json:"person" binding:"omitempty"`
	Description ScheduleShort       `json:"description" binding:"required"`
}

type TaskInfo struct {
	Id uuid.UUID `json:"id"`
	TaskRequest
}
