package routes

import (
	"example.com/restful-server/controllers"
	"example.com/restful-server/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	apiV1 := server.Group("/v1")
	apiV1.Use(middleware.Authenticate)
	{
		users := apiV1.Group("/users")
		{
			users.GET("/", controllers.GetUsers)
		}
	}

	// sign-up / login
	server.POST("/signup", controllers.Signup)
	server.POST("/login", controllers.Login)
}
