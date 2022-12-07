package routes

import (
	"jwt/controllers"
	"jwt/middleware"

	"github.com/gin-gonic/gin"
)

func AdminAuthRoutes(r *gin.Engine) {

	r.LoadHTMLGlob("templates/*.html")
	r.GET("/adminLogin", controllers.AdminLogin)
	r.POST("/adminPostLogin", controllers.AdminPostLogin)
	r.GET("/adminProfile", middleware.RequireAuthAdmin, controllers.AdminLogged, controllers.AdminProfile)
	r.GET("/logoutadmin", controllers.AdminLogout)
	r.POST("/createUser", controllers.CreateUser)
	r.GET("/deleteUser/:name", controllers.DeleteUser)
	
}
