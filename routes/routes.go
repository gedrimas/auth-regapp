package routes

import (
	"auth-regapp/controllers"
	"github.com/gin-gonic/gin"
	"auth-regapp/middleware"
)

func AuthRoutes(router *gin.Engine) {
	router.Use(middleware.CORSMiddleware())
	router.POST("users/signup", controllers.Signup())
	router.POST("users/login", controllers.Login())
}