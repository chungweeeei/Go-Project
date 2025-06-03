package controllers

import (
	"errors"

	"example.com/restful-server/database"
	"example.com/restful-server/models"
	"example.com/restful-server/utils"
	"gorm.io/gorm"
)

type RegisterUserRequest struct {
	Email    string `json:"email";binding:"required,email"`
	Username string `json:"username";binding:"required"`
	Password string `json:"password";binding:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email";binding:"required,email"`
	Password string `json:"password";binding:"required"`
}

func (loginUser *LoginUserRequest) ValidateCredentials() error {

	var user models.User
	result := database.DB.Where("email = ?", loginUser.Email).First(&user)

	// Check for record not found error
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("Could not find user.")
	}

	// Check other database errors
	if result.Error != nil {
		return errors.New("Unexpected error occurred while fetching user.")
	}

	passwordIsValid := utils.CheckPasswordHash(loginUser.Password, user.Password)

	if !passwordIsValid {
		return errors.New("Credentials invalid")
	}

	return nil
}

type UserInfo struct {
	Email     string `json:"email";bind:"required"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
