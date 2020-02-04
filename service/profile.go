package service

import (
	"irisProject/models"
	"irisProject/utils"

	"github.com/jinzhu/gorm"
)

// ProfileService struct
type ProfileService struct {
	DB *gorm.DB
}

// NewProfileService return a New ProfileService
func NewProfileService(db *gorm.DB) ProfileInterface {
	return &ProfileService{
		DB: db,
	}
}

// ProfileInterface  interface
type ProfileInterface interface {
	GetProfileList(parameters utils.GetProfileListParameters) (profiles []models.Profile, count int, err error)
	CreateProfile(profile models.Profile) (err error)
	GetProfile(id string) (profile models.Profile, err error)
	UpdateProfile(id string, newProfile models.Profile) (err error)
	DeleteProfile(id string) (err error)
}

// GetProfileList return profiles
func (s *ProfileService) GetProfileList(parameters utils.GetProfileListParameters) (profiles []models.Profile, count int, err error) {
	db := s.DB.Scopes(
		utils.SearchByColumn("content", parameters.Content),
	)
	err = db.Model(&models.Profile{}).
		Count(&count).
		Error

	if err != nil {
		return
	}

	err = db.
		Offset((parameters.Page - 1) * parameters.Limit).
		Limit(parameters.Limit).
		Find(&profiles).
		Error
	return
}

// CreateProfile create a new Profile
func (s *ProfileService) CreateProfile(profile models.Profile) (err error) {
	err = s.DB.
		Create(&profile).
		Error
	return
}

// GetProfile Get a profile by id
func (s *ProfileService) GetProfile(id string) (profile models.Profile, err error) {
	err = s.DB.
		Preload("Creator").
		Where("id = ?", id).
		Take(&profile).
		Error
	return
}

// UpdateProfile update a profile record
func (s *ProfileService) UpdateProfile(id string, profile models.Profile) (err error) {
	err = s.DB.
		Model(&models.Profile{}).
		Where("id = ?", id).
		Updates(profile).
		Error
	return
}

// DeleteProfile delete a profile record
func (s *ProfileService) DeleteProfile(id string) (err error) {
	var profile models.Profile
	err = s.DB.
		Where("id = ?", id).
		Delete(&profile).
		Error
	return
}
