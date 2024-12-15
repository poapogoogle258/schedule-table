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

	authSet = wire.NewSet(
		handler.NewAuthHandler,
		service.NewJWTAuthService,
		repository.NewUserRepository,
	)
)

func Injector() *router.Handlers {

	wire.Build(calendarSet, authSet, database.ConnectPostgresql, wire.Struct(new(router.Handlers), "*"))

	return &router.Handlers{}

}
