package service

import (
	"time"
	// "github.com/teambition/rrule-go"
)

type Config struct {
}

type Recurrence interface {
	GetAll(start *time.Time, end *time.Time, rule *Config)
}
