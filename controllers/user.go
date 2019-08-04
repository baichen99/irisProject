package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"irisProject/middlewares"
	"irisProject/service"
	"irisProject/utils"
)

type UserController struct {
	Context iris.Context
	Service service.UserInterface
}

func (c *UserController) BeforeActivation(app mvc.BeforeActivation) {
	app.Handle("POST", "/login", "Login")
	app.Handle("GET", "/", "GetUserList", middlewares.CheckJWTToken)
	app.Handle("GET", "/{username:string}", "GetUser", middlewares.CheckJWTToken)
	app.Handle("POST", "/", "CreateUser",middlewares.CheckJWTToken, middlewares.CheckAdminRole)
	app.Handle("DELETE", "/{username:string}", "DeleteUser", middlewares.CheckJWTToken)
	app.Handle("PUT", "/{username:string}", "UpdateUser", middlewares.CheckJWTToken)
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
		"data":    token,
	})



}

func (c *UserController) GetUserList() {
	c.Context.Next()
	username := c.Context.URLParamDefault("username", "")
	listParams, err := utils.GetListParamsFromContext(c.Context, "username")
	if err != nil {
		return
	}
	listParameters := utils.GetUserListParameters{
		GetListParameters: listParams,
		Username: username,
	}

	users, count, err := c.Service.GetUserList(listParameters)
	if err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService:GetUserList", err)
	}

	c.Context.JSON(iris.Map{
			"message": "success",
			"data": iris.Map{
				"users":  users,
				"pageKey":  listParams.Page,
				"limitKey": listParams.Limit,
				"totalKey": count,
			},
		})

}


func (c *UserController) GetUser() {
	username := c.Context.Params().GetStringDefault("username", "")
	user, err := c.Service.GetUser("username", username)
	if err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService::GetUser", err)
		return
	}

	c.Context.JSON(iris.Map{
		"message": "success",
		"data": user,
	})

}

func (c *UserController) CreateUser() {
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
		"message" : "success",
	})

}


func (c *UserController) DeleteUser() {
	c.Context.Next()
	username := c.Context.Params().GetStringDefault("username", "")
	if err := c.Service.DeleteUser("username", username); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService::DeleteUser", err)
		return
	}

	c.Context.StatusCode(iris.StatusNoContent)
}


func (c *UserController) UpdateUser() {
	c.Context.Next()
	username := c.Context.Params().GetStringDefault("username", "")
	var form UserUpdateForm
	if err := utils.ReadValidateForm(c.Context, &form); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService:UpdateUser", err)
		return
	}

	if err := c.Service.UpdateUser("username", username, form.ConvertToModel()); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService:UpdateUser", err)
		return
	}

	c.Context.StatusCode(iris.StatusNoContent)
}
