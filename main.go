package main

import (
	"irisProject/config"
	"irisProject/controllers"
	"irisProject/middlewares"
	"irisProject/service"
	"irisProject/utils"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func main() {
	app := initApp()
	db := utils.ConnectDB()
	utils.InitDB(db)
	mvc.Configure(app.Party("/user"), func(app *mvc.Application) {
		app.Register(service.NewUserService(db))
		app.Handle(new(controllers.UserController))
	})
	mvc.Configure(app.Party("/"), func(app *mvc.Application) {
		app.Handle(new(controllers.RootController))
	})

	app.Run(iris.Addr(":8080"))
	//test()
}

func initApp() *iris.Application {
	app := iris.Default()
	app.Use(middlewares.NewI18nMiddleware(middlewares.I18nConf))
	app.Use(middlewares.CorsAllowAll)
	return app
}

func test() {
	// get config
	println(config.Conf.JWT.PrivateBytes)

}
