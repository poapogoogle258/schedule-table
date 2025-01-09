//go:build wireinject
// +build wireinject

package main

import (
	"schedule_table/internal/database"
	"schedule_table/internal/handler"
	"schedule_table/internal/http/router"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"

	"github.com/google/wire"
)

var (
	calendarSet = wire.NewSet(
		handler.NewCalendarsHandler,
		repository.NewCalendarRepository,
	)

	authSet = wire.NewSet(
		handler.NewAuthHandler,
		repository.NewUserRepository,
		service.NewJWTAuthService,
	)

	memberSet = wire.NewSet(
		handler.NewMemberHandler,
		repository.NewMemberRepository,
	)

	scheduleSet = wire.NewSet(
		handler.NewScheduleHandler,
		repository.NewScheduleRepository,
	)

	taskSet = wire.NewSet(
		handler.NewTasksHandler,
		service.NewScheduleService,
		service.NewManagerService,
		repository.NewTaskRepository,
	)

	leaveSet = wire.NewSet(
		handler.NewLeaveHandler,
		repository.NewLeaveRepository,
	)
)

func Injector() *router.Handlers {

	wire.Build(leaveSet, taskSet, scheduleSet, memberSet, calendarSet, authSet, database.ConnectPostgresql, wire.Struct(new(router.Handlers), "*"))

	return &router.Handlers{}

}

// go run github.com/google/wire/cmd/wire .
