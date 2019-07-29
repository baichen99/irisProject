package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"irisProject/config"
	"irisProject/controllers"
	"irisProject/middlewares"
	"irisProject/service"
	"irisProject/utils"
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
	return app
}

func test() {
	// get config
	println(config.Conf.JWT.PrivateBytes)

}


