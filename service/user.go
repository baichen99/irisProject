package service

import (
	"irisProject/models"
	"irisProject/utils"

	"github.com/jinzhu/gorm"
)

// UserService struct
type UserService struct {
	DB *gorm.DB
}

// NewUserService return a NewUserService
func NewUserService(db *gorm.DB) UserInterface {
	return &UserService{
		DB: db,
	}
}

// UserInterface  interface
type UserInterface interface {
	GetUserList(parameters utils.GetUserListParameters) (users []models.User, count int, err error)
	CreateUser(user models.User) (err error)
	GetUser(key string, value interface{}) (user models.User, err error)
	UpdateUser(id int, newUser models.User) (err error)
	DeleteUser(id int) (err error)
}

// GetUserList return users
func (s *UserService) GetUserList(parameters utils.GetUserListParameters) (users []models.User, count int, err error) {
	var user models.User
	db := s.DB.Scopes(
		utils.SearchByColumn("username", parameters.Username),
	)
	err = db.Model(&users).Count(&count).Error
	if err != nil {
		return
	}

	err = db.Model(&user).Offset((parameters.Page - 1) * parameters.Limit).Limit(parameters.Limit).Scan(&users).Error
	return
}

// CreateUser create a new user
func (s *UserService) CreateUser(user models.User) (err error) {
	user.Password, err = utils.HashPassword(user.Password)
	err = s.DB.Create(&user).Error
	if err != nil {
		return
	}
	return
}

// GetUser Get a user by id
func (s *UserService) GetUser(key string, value interface{}) (user models.User, err error) {
	switch key {
	case "id":
		err = s.DB.Where("id = ?", value).Take(&user).Error
	case "username":
		err = s.DB.Where("username = ?", value).Take(&user).Error
	}
	return
}

// UpdateUser update a user record
func (s *UserService) UpdateUser(id int, newUser models.User) (err error) {
	var user models.User
	if newUser.Password != "" {
		newUser.Password, err = utils.HashPassword(newUser.Password)
		if err != nil {
			return
		}
	}
	err = s.DB.Where("id = ?", id).Take(&user).Updates(&newUser).Error
	return
}

// DeleteUser delete a user record
func (s *UserService) DeleteUser(id int) (err error) {
	var user models.User
	err = s.DB.Where("id = ?", id).Delete(&user).Error
	return
}
