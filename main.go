package main

import (
	"irisProject/controllers"
	"irisProject/middlewares"
	"irisProject/service"
	"irisProject/utils"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"

	"github.com/kataras/iris/v12"
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
	mvc.Configure(app.Party("/profile"), func(app *mvc.Application) {
		app.Register(service.NewProfileService(db))
		app.Handle(new(controllers.ProfileController))
	})
	mvc.Configure(app.Party("/"), func(app *mvc.Application) {
		app.Handle(new(controllers.RootController))
	})

	app.Run(iris.Addr(":8080"))
}

func initApp() *iris.Application {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Status:             true,
		IP:                 true,
		Method:             true,
		Path:               true,
		Query:              true,
		MessageContextKeys: []string{"RequestID"},
		MessageHeaderKeys:  []string{"User-Agent"},
	}))
	app.OnAnyErrorCode(middlewares.ErrorHandler)
	app.I18n.Load("./locales/*.ini", "en-US", "zh-CN")
	// app.Use(middlewares.CorsAllowAll())
	app.UseGlobal(middlewares.BeforeHandleRequest)
	app.DoneGlobal(middlewares.AfterHandleRequest)

	return app
}
