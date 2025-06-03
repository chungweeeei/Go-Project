package routes

import (
	"example.com/restful-server/controllers"
	"example.com/restful-server/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	// add cors middleware
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5000"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	apiV1 := server.Group("/v1")
	apiV1.Use(middleware.Authenticate)
	{
		users := apiV1.Group("/users")
		{
			users.GET("/", controllers.GetUsers)
		}
		files := apiV1.Group("/files")
		{
			files.POST("/upload", controllers.Upload)
		}
	}

	// sign-up / login
	server.POST("/v1/signup", controllers.Signup)
	server.POST("/v1/login", controllers.Login)
}
