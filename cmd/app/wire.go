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

	generateTaskSet = wire.NewSet(
		handler.NewGenerateTaskHandler,
		repository.NewSchedulesRepository,
		service.NewRecurrenceService,
	)

	authSet = wire.NewSet(
		handler.NewAuthHandler,
		repository.NewUserRepository,
		service.NewJWTAuthService,
	)

	calSet = wire.NewSet(
		handler.NewMemberHandler,
		repository.NewMemberRepository,
	)
)

func Injector() *router.Handlers {

	wire.Build(generateTaskSet, calSet, calendarSet, authSet, database.ConnectPostgresql, wire.Struct(new(router.Handlers), "*"))

	return &router.Handlers{}

}

// go run github.com/google/wire/cmd/wire .
