package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"irisProject/middlewares"
	"irisProject/service"
	"irisProject/utils"
)

type ProfileController struct {
	Context iris.Context
	Service service.ProfileInterface
}

func (c *ProfileController) BeforeActivation(app mvc.BeforeActivation) {
	app.Router().Use(middlewares.CheckJWTToken)
	app.Handle("GET", "/", "GetProfileList")
	app.Handle("GET", "/{id:string}", "GetProfile")
	app.Handle("POST", "/", "CreateProfile")
	app.Handle("DELETE", "/{id:string}", "DeleteProfile")
	app.Handle("PUT", "/{id:string}", "UpdateProfile")
}

func (c *ProfileController) GetProfileList() {
	defer c.Context.Next()
	content := c.Context.URLParamDefault("content", "")
	listParams, err := utils.GetListParamsFromContext(c.Context, "profiles.id")
	if err != nil {
		return
	}
	listParameters := utils.GetProfileListParameters{
		GetListParameters: listParams,
		Content:           content,
	}

	profiles, count, err := c.Service.GetProfileList(listParameters)
	if err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "ProfileService:GetProfileList", err)
	}

	c.Context.JSON(iris.Map{
		"message": "success",
		"data": iris.Map{
			"profiles": profiles,
			"page":     listParams.Page,
			"limit":    listParams.Limit,
			"total":    count,
		},
	})

}

func (c *ProfileController) GetProfile() {
	defer c.Context.Next()
	id := c.Context.Params().Get("id")
	profile, err := c.Service.GetProfile(id)
	if err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "ProfileService::GetProfile", err)
		return
	}

	c.Context.JSON(iris.Map{
		"message": "success",
		"data":    profile,
	})

}

func (c *ProfileController) CreateProfile() {
	defer c.Context.Next()
	var form ProfileCreateForm
	err := utils.ReadValidateForm(c.Context, &form)
	if err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "paramKey", err)
	}
	profile := form.ConvertToModel()
	err = c.Service.CreateProfile(profile)
	if err != nil {
		utils.SetResponseError(c.Context, iris.StatusUnprocessableEntity, "ProfileController:CreateProfile", err)
	}

	c.Context.StatusCode(iris.StatusCreated)
	c.Context.JSON(iris.Map{
		"message": "success",
	})

}

func (c *ProfileController) DeleteProfile() {
	defer c.Context.Next()
	id := c.Context.Params().Get("id")
	if err := c.Service.DeleteProfile(id); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "ProfileService::DeleteProfile", err)
		return
	}

	c.Context.StatusCode(iris.StatusNoContent)
}

func (c *ProfileController) UpdateProfile() {
	defer c.Context.Next()
	id := c.Context.Params().Get("id")
	var form ProfileUpdateForm
	if err := utils.ReadValidateForm(c.Context, &form); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "ProfileService:UpdateProfile", err)
		return
	}

	if err := c.Service.UpdateProfile(id, form.ConvertToModel()); err != nil {
		utils.SetResponseError(c.Context, iris.StatusBadRequest, "ProfileService:UpdateProfile", err)
		return
	}

	c.Context.StatusCode(iris.StatusNoContent)
}
