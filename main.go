package main

import (
	"jwt/initializers"
	"jwt/routes"

	"github.com/gin-gonic/gin"
)

func init() {

	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()

}

func main() {
	

	router := gin.Default()
	router.Use(gin.Logger())

	routes.UserAuthRoutes(router)

	router.Run()

}
