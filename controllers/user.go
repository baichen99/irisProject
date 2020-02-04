package controllers

import (
	"irisProject/middlewares"
	"irisProject/service"
	"irisProject/utils"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type UserController struct {
	Context iris.Context
	Service service.UserInterface
}

func (c *UserController) BeforeActivation(app mvc.BeforeActivation) {
	app.Handle("POST", "/login", "Login")
	app.Router().Use(middlewares.CheckJWTToken)
	app.Handle("GET", "/", "GetUserList")
	app.Handle("GET", "/{id:string}", "GetUser")
	app.Router().Use(middlewares.CheckSuper)
	app.Handle("POST", "/", "CreateUser")
	app.Handle("DELETE", "/{id:string}", "DeleteUser")
	app.Handle("PUT", "/{id:string}", "UpdateUser")
}

func (c *UserController) Login() {
	defer c.Context.Next()
	var form UserLoginForm
	if err := utils.ReadValidateForm(c.Context, &form); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService:Login", err)
		return
	}
	// get user from database
	user, err := c.Service.GetUser("username", form.Username)
	if err != nil {
		utils.SetResponseError(c.Context, iris.StatusUnauthorized, "UserService::GetUser", err)
		return
	}
	password, _ := utils.HashPassword(form.Password)
	if !utils.ComparePassword(password, form.Password) {
		utils.SetResponseError(c.Context, iris.StatusUnauthorized, "UserService::Login", err)
		return
	}

	token, _ := middlewares.SignJWTToken(user.ID, user.Role)
	c.Context.JSON(iris.Map{
		"message": "success",
		"token":   token,
	})

}

func (c *UserController) GetUserList() {
	defer c.Context.Next()
	username := c.Context.URLParamDefault("username", "")
	listParams, err := utils.GetListParamsFromContext(c.Context, "users.username")
	if err != nil {
		return
	}
	listParameters := utils.GetUserListParameters{
		GetListParameters: listParams,
		Username:          username,
	}

	users, count, err := c.Service.GetUserList(listParameters)
	if err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService:GetUserList", err)
	}

	c.Context.JSON(iris.Map{
		"message": "success",
		"data": iris.Map{
			"users": users,
			"page":  listParams.Page,
			"limit": listParams.Limit,
			"total": count,
		},
	})
}

func (c *UserController) GetUser() {
	defer c.Context.Next()
	id := c.Context.Params().Get("id")
	user, err := c.Service.GetUser("id", id)
	if err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService::GetUser", err)
		return
	}

	c.Context.JSON(iris.Map{
		"message": "success",
		"data":    user,
	})

}

func (c *UserController) CreateUser() {
	defer c.Context.Next()
	var form UserCreateForm
	err := utils.ReadValidateForm(c.Context, &form)
	if err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "paramKey", err)
	}
	user := form.ConvertToModel()
	err = c.Service.CreateUser(user)
	if err != nil {
		utils.SetResponseError(c.Context, iris.StatusUnprocessableEntity, "UserController:CreateUser", err)
	}

	c.Context.StatusCode(iris.StatusCreated)
	c.Context.JSON(iris.Map{
		"message": "success",
	})

}

func (c *UserController) DeleteUser() {
	defer c.Context.Next()
	id := c.Context.Params().Get("id")
	if err := c.Service.DeleteUser(id); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService::DeleteUser", err)
		return
	}

	c.Context.StatusCode(iris.StatusNoContent)
}

func (c *UserController) UpdateUser() {
	defer c.Context.Next()
	id := c.Context.Params().Get("id")
	var form UserUpdateForm
	if err := utils.ReadValidateForm(c.Context, &form); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "paramKey", err)
		return
	}

	if err := c.Service.UpdateUser(id, form.ConvertToModel()); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService:UpdateUser", err)
		return
	}

	c.Context.StatusCode(iris.StatusNoContent)
}
