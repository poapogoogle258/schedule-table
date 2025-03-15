package main

import (
	"fmt"
	"os"
	"schedule_table/internal/http/router"
	"schedule_table/internal/pkg/logger"
	vcustom "schedule_table/internal/pkg/validator"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil { // default load is .env
		panic(err)
	}

	handlers := Injector()
	server := router.NewRouter(handlers)

	// initial logger
	logger.InitLogger()
	defer logger.Sync()

	// initial add custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("telephone", vcustom.TelephoneFormat)
		v.RegisterValidation("hexcolor", vcustom.HexColorFormat)
		v.RegisterValidation("hhmm", vcustom.HHMMTimeFormat)
		v.RegisterValidation("byweekday", vcustom.Weekday)
		v.RegisterValidation("bymonth", vcustom.Month)
		logger.Info("init validator success")
	}

	// start http servers
	addr := fmt.Sprintf("%s:%s", os.Getenv("IP"), os.Getenv("PORT"))
	server.Run(addr) // listen and serve on 0.0.0.0:8080

}
