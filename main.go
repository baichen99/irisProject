package main

import (
	"fmt"
	"irisProject/config"
	"irisProject/controllers"
	"irisProject/middlewares"
	"irisProject/service"
	"irisProject/utils"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func main() {
	app := initApp()
	db := utils.ConnectDB()
	utils.InitDB(db)
	utils.InitAdmin(db)
	mvc.Configure(app.Party("/user"), func(app *mvc.Application) {
		app.Register(service.NewUserService(db))
		app.Handle(new(controllers.UserController))
	})
	mvc.Configure(app.Party("/"), func(app *mvc.Application) {
		app.Handle(new(controllers.RootController))
	})

	test()

	app.Run(iris.Addr(":8080"))
}

func initApp() *iris.Application {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middlewares.NewI18nMiddleware(middlewares.I18nConf))
	app.Use(middlewares.CorsAllowAll)
	app.Use(middlewares.BeforeRequest)
	app.Done(middlewares.AfterRequest)

	return app
}

func test() {
	// get config
	fmt.Println(config.Conf)

}
