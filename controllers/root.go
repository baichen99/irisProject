package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type RootController struct {
	Context iris.Context
}

func (c *RootController) BeforeActivation(app mvc.BeforeActivation) {
	app.Handle("GET", "/", "Get")
}
func (c *RootController) Get() {
	hi := c.Context.Tr("Hello")

	locale := c.Context.GetLocale()
	c.Context.Writef("%s, Your language is %s", hi, locale.Language())
}
