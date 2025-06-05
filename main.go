package main

import (
	"example.com/restful-server/database"
	"example.com/restful-server/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	database.InitDB() // initialize database connection

	// create a gin default server instance
	server := gin.Default()

	// http server register routes
	routes.RegisterRoutes(server)

	// // register publish job
	// job := jobs.NewPublishJob("Hello, Go!")
	// go job.Process()

	// // select listening port for gin server
	server.Run(":3000")
}
