package middleware

import (
	"net/http"

	"example.com/restful-server/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {

	// check token
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	userEmail, err := utils.VerifyToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	context.Set("userEmail", userEmail)
	context.Next()
}
