package service

import (
	"errors"
	models "hamkaran_system/bootcamp/final/project/model"

	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) error {

	if user.Password == "" {
		return errors.New("please enter password")
	}
	if user.Username == "" {
		return errors.New("please enter password")
	}
	var existingUser models.User
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return errors.New("duplicate user exists with this username")
	}
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
