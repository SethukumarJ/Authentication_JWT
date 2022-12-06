package routes

import (
	"jwt/controllers"
	"jwt/middleware"

	"github.com/gin-gonic/gin"
)

func UserAuthRoutes(r *gin.Engine) {

	r.LoadHTMLGlob("templates/*.html")
	r.GET("/userSignup", controllers.UserSignup)
	r.POST("/userPostSignup", controllers.UserPostSignup)
	r.GET("/userLogin", controllers.UserLogin)
	r.POST("/userPostLogin", controllers.UserPostLogin)
	r.GET("/welocomeUser", middleware.RequireAuth, controllers.WelcomeUser)
	r.GET("/logout", controllers.UserLogout)
}
