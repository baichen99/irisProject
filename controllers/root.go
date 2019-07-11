package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type RootController struct {
	Context iris.Context
}

func (c *RootController) BeforeActivation (app mvc.BeforeActivation) {
	app.Handle("GET", "/", "Get")
}
func (c *RootController) Get() {
	hi := c.Context.Translate("Hello")

	language := c.Context.Values().GetString(c.Context.Application().ConfigurationReadOnly().GetTranslateLanguageContextKey())
	c.Context.Writef("%s, Your language is %s", hi, language)
}
