package service

import (
	"errors"
	models "hamkaran_system/bootcamp/final/project/model"

	"golang.org/x/crypto/bcrypt"
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
	err := db.Where("username = ?", user.Username).First(&existingUser).Error
	if err == nil {
		return errors.New("duplicate user exists with this username")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func CheckPassword(plainPassword, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
