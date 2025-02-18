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
	worker "schedule_table/internal/workers/task"

	"github.com/google/wire"
)

var (
	repositorySet = wire.NewSet(
		repository.NewCalendarRepository,
		repository.NewLeaveRepository,
		repository.NewMemberRepository,
		repository.NewScheduleRepository,
		repository.NewTaskRepository,
		repository.NewUserRepository,
	)

	serviceSet = wire.NewSet(
		service.NewJWTAuthService,
		service.NewTaskService,
	)

	middlewareSet = wire.NewSet(
		middleware.NewAuthorizeJWTMiddleware,
		middleware.NewCalendarMiddleware,
		middleware.NewMemberMiddleware,
		middleware.NewScheduleMiddleware,
		middleware.NewTaskMiddleware,
	)

	handlerSet = wire.NewSet(
		handler.NewAuthHandler,
		handler.NewCalendarsHandler,
		handler.NewLeaveHandler,
		handler.NewScheduleHandler,
		handler.NewMemberHandler,
		handler.NewTasksHandler,
	)
)

func Injector() *router.Handlers {

	wire.Build(repositorySet, handlerSet, serviceSet, middlewareSet, database.ConnectPostgresql, wire.Struct(new(router.Handlers), "*"))

	return &router.Handlers{}

	// return &router.Handlers{}, &repository.CalendarRepository{}, &repository.ITaskRepository{}, &repository.ScheduleRepository{}, &repository.MembersRepository{}

}

func InjectorWorker() *worker.WorkerInit {

	wire.Build(repositorySet, database.ConnectPostgresql, wire.Struct(new(worker.WorkerInit), "*"))

	return &worker.WorkerInit{}

}

// go run github.com/google/wire/cmd/wire .
