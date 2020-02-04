package controllers

import "irisProject/models"

type UserCreateForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (f UserCreateForm) ConvertToModel() models.User {
	return models.User{
		Username: f.Username,
		Password: f.Password,
	}
}

type UserLoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (f UserLoginForm) ConvertToModel() models.User {
	return models.User{
		Username: f.Username,
		Password: f.Password,
	}
}

type UserUpdateForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (f UserUpdateForm) ConvertToModel() models.User {
	return models.User{
		Username: f.Username,
		Password: f.Password,
	}
}

type ProfileCreateForm struct {
	Content string
}

func (f ProfileCreateForm) ConvertToModel() models.Profile {
	return models.Profile{
		Content: f.Content,
	}
}

type ProfileUpdateForm struct {
	Content string
}

func (f ProfileUpdateForm) ConvertToModel() models.Profile {
	return models.Profile{
		Content: f.Content,
	}
}
