package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"irisProject/service"
	"irisProject/utils"
)

type UserController struct {
	Context iris.Context
	Service service.UserInterface
}

func (c *UserController) BeforeActivation(app mvc.BeforeActivation) {
	app.Handle("GET", "/", "GetUserList")
	app.Handle("GET", "/{id:int}", "GetUser")
	app.Handle("POST", "/", "CreateUser")
	app.Handle("DELETE", "/{id:int}", "DeleteUser")
	app.Handle("PUT", "/{id:int}", "UpdateUser")
}

func (c *UserController) GetUserList() {
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
	id := c.Context.Params().GetIntDefault("id", 0)
	user, err := c.Service.GetUser(id)
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
	id := c.Context.Params().GetIntDefault("id", 0)
	if err := c.Service.DeleteUser(id); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService::DeleteUser", err)
		return
	}

	c.Context.StatusCode(iris.StatusNoContent)
}


func (c *UserController) UpdateUser() {
	id := c.Context.Params().GetIntDefault("id", 0)
	var form UserUpdateForm
	if err := utils.ReadValidateForm(c.Context, &form); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService:UpdateUser", err)
		return
	}

	if err := c.Service.UpdateUser(id, form.ConvertToModel()); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "UserService:UpdateUser", err)
		return
	}

	c.Context.StatusCode(iris.StatusNoContent)
}
