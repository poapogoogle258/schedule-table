//go:build wireinject
// +build wireinject

package main

import (
	"schedule_table/internal/database"
	"schedule_table/internal/handler"
	"schedule_table/internal/repository"
	"schedule_table/internal/router"
	"schedule_table/internal/service"

	"github.com/google/wire"
)

var (
	calendarSet = wire.NewSet(
		handler.NewCalendarsHandler,
		repository.NewCalendarRepository,
	)

	scheduleSet = wire.NewSet(
		repository.NewSchedulesRepository,
	)

	authSet = wire.NewSet(
		handler.NewAuthHandler,
		service.NewJWTAuthService,
		repository.NewUserRepository,
	)
)

func Injector() *router.Handlers {

	wire.Build(scheduleSet, calendarSet, authSet, database.ConnectPostgresql, wire.Struct(new(router.Handlers), "*"))

	return &router.Handlers{}

}

// go run github.com/google/wire/cmd/wire .
