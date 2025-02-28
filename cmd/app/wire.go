//go:build wireinject
// +build wireinject

package main

import (
	"schedule_table/internal/database"
	"schedule_table/internal/handler"
	"schedule_table/internal/http/middleware"
	"schedule_table/internal/http/router"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"

	"github.com/google/wire"
)

var (
	repositorySet = wire.NewSet(
		repository.NewTransaction,
		repository.NewCalendarRepository,
		repository.NewUserRepository,
		repository.NewEmployeeRepository,
		repository.NewScheduleRepository,
	)

	serviceSet = wire.NewSet(
		service.NewJWTAuthService,
		service.NewUserService,
		service.NewEmployeeService,
		service.NewCalendarService,
		service.NewScheduleService,
	)

	middlewareSet = wire.NewSet(
		middleware.NewAuthorizeJWTMiddleware,
		middleware.NewCalendarMiddleware,
	)

	handlerSet = wire.NewSet(
		handler.NewAuthHandler,
		handler.NewCalendarsHandler,
		handler.NewEmployeeHandler,
		handler.NewScheduleHandler,
	)
)

func Injector() *router.Handlers {

	wire.Build(repositorySet, handlerSet, serviceSet, middlewareSet, database.ConnectPostgresql, wire.Struct(new(router.Handlers), "*"))

	return &router.Handlers{}

}

// go run github.com/google/wire/cmd/wire .
