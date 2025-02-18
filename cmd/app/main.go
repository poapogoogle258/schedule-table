package main

import (
	"fmt"
	"os"
	"schedule_table/internal/http/router"
	"schedule_table/internal/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	handlers := Injector()
	server := router.NewRouter(handlers)

	// initial logger
	logger.InitLogger()
	defer logger.Sync()

	// start CalendarTaskUpdate workers run on background
	// ctxBg := context.Background()
	// initWorker := InjectorWorker()
	// jobQueue := worker.NewJobQueue(100)
	// for i := 0; i < 5; i++ {
	// 	taskWorker := worker.NewTaskWorker(i, jobQueue, initWorker)
	// 	go taskWorker.Start(ctxBg)
	// }

	// // check tasks already
	// listCalendarIdScheduleChanged, errGetListIdOfScheduleChanged := initWorker.CalRepo.GetListIdOfScheduleChanged()
	// if errGetListIdOfScheduleChanged != nil {
	// 	logger.Error("errGetListIdOfScheduleChangedError", zap.Error(errGetListIdOfScheduleChanged))
	// 	panic(errGetListIdOfScheduleChanged)
	// }
	// for i := range listCalendarIdScheduleChanged {
	// 	jobQueue.Queue <- worker.NewJob(listCalendarIdScheduleChanged[i], "schedule")
	// }

	// // add hook RecurrentScheduleChanged to logger Observer
	// logger.AddHook("RecurrentScheduleChanged", func(data map[string]interface{}) error {
	// 	calendarId, _ := data["calendarId"].(string)
	// 	jobQueue.Queue <- worker.NewJob(calendarId, "schedule")
	// 	return nil
	// })

	// start http servers
	addr := fmt.Sprintf("%s:%s", os.Getenv("IP"), os.Getenv("PORT"))
	server.Run(addr) // listen and serve on 0.0.0.0:8080

}
