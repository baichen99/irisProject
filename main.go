package main

import (
	"github.com/kataras/iris"

	"github.com/kataras/iris/mvc"
	"irisProject/controllers"
	"irisProject/service"
	"irisProject/utils"
)

func main() {
	app := iris.Default()
	db := utils.ConnectDB()
	utils.InitDB(db)
	mvc.Configure(app.Party("/user"), func(app *mvc.Application) {
		app.Register(service.NewUserService(db))
		app.Handle(new(controllers.UserController))
	})

	app.Run(iris.Addr(":8080"))
}
