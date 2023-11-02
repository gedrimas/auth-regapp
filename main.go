package main

import (
	"auth-regapp/database"
	"auth-regapp/helpers"
	"github.com/gin-gonic/gin"
	"auth-regapp/routes"
)

func main() {

	router := gin.Default()

	database.Start()

	router.Use(gin.Logger())

	routes.AuthRoutes(router)

	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"response": "api started successfuly",
		})
	})

	port := helpers.EnvFileVal("PORT")
	router.Run(":" + port)
}
