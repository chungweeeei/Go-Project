package controllers

import (
	"errors"
	"net/http"

	"example.com/restful-server/database"
	"example.com/restful-server/models"
	"example.com/restful-server/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterUserRequest struct {
	Email    string `json:"email";binding:"required,email"`
	Username string `json:"username";binding:"required"`
	Password string `json:"password";binding:"required"`
}

func Signup(context *gin.Context) {
	var userReq RegisterUserRequest
	err := context.ShouldBindJSON(&userReq)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	hashedPassword, err := utils.HashPassword(userReq.Password)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not hash password."})
		return
	}

	user := models.User{
		Email:    userReq.Email,
		Username: userReq.Username,
		Password: hashedPassword,
		Role:     "guest",
	}

	// Save user to database
	result := database.DB.Create(&user)

	// check result for errors
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Successfully sign-up a new user."})
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

func Login(context *gin.Context) {
	var loginUser LoginUserRequest
	err := context.ShouldBindJSON(&loginUser)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = loginUser.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Could not authenticate user."})
		return
	}

	// generate jwt token
	// user => is from the request which only include email and password
	token, err := utils.GenerateToken(loginUser.Email)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successfully.", "email": loginUser.Email, "token": token})
}
