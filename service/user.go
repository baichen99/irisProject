package service

import (
	"errors"
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
	GetUser(key string, value string) (user models.User, err error)
	UpdateUser(id string, newUser models.User) (err error)
	DeleteUser(id string) (err error)
}

// GetUserList return users
func (s *UserService) GetUserList(parameters utils.GetUserListParameters) (users []models.User, count int, err error) {
	db := s.DB.Scopes(
		utils.SearchByColumn("username", parameters.Username),
	)
	queryExp := db.Joins("JOIN profiles ON profiles.id = users.profile_id")
	err = queryExp.Find(&users).Count(&count).Error

	if err != nil {
		return
	}

	err = queryExp.
		Offset((parameters.Page - 1) * parameters.Limit).
		Limit(parameters.Limit).
		Find(&users).
		Error
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
func (s *UserService) GetUser(key string, value string) (user models.User, err error) {
	switch key {
	case "id":
		err = s.DB.Where("id = ?", value).Take(&user).Error
	case "username":
		err = s.DB.Where("username = ?", value).Take(&user).Error
	default:
		err = errors.New("Unsupported key")
	}
	return
}

// UpdateUser update a user record
func (s *UserService) UpdateUser(id string, newUser models.User) (err error) {
	if newUser.Password != "" {
		newUser.Password, err = utils.HashPassword(newUser.Password)
		if err != nil {
			return
		}
	}
	err = s.DB.
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(newUser).
		Error
	return
}

// DeleteUser delete a user record
func (s *UserService) DeleteUser(id string) (err error) {
	var user models.User
	err = s.DB.
		Where("id = ?", id).
		Delete(&user).
		Error
	return
}
