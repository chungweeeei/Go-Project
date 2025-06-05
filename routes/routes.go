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
		tasks := apiV1.Group("/tasks")
		{
			tasks.POST("", controllers.CreateTask)
			tasks.GET("", controllers.GetTasks)
			tasks.GET("/status-counts", controllers.GetTaskStatusCounts)
			tasks.GET("/:id", controllers.GetTaskByID)
			tasks.PUT("/:id", controllers.UpdateTask)
			tasks.DELETE("/:id", controllers.DeleteTask)
		}
	}

	// sign-up / login
	server.POST("/v1/auth/register", controllers.Signup)
	server.POST("/v1/auth/login", controllers.Login)
}
