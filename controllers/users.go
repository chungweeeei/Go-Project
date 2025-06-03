package controllers

import (
	"net/http"

	"example.com/restful-server/database"
	"example.com/restful-server/helper"
	"example.com/restful-server/models"
	"github.com/gin-gonic/gin"
)

func GetUsers(context *gin.Context) {
	var users []models.User
	result := database.DB.Find(&users)

	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error when fetching users"})
		return
	}

	var usersInfo []UserInfo
	for _, user := range users {
		usersInfo = append(usersInfo, UserInfo{Email: user.Email, Username: user.Username, Role: user.Role, CreatedAt: helper.ToISO(user.CreatedAt)})
	}

	context.JSON(http.StatusOK, usersInfo)
}
