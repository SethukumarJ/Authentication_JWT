package controllers

import (
	"jwt/initializers"
	"jwt/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {

	var user models.User
	//GEt email and password from request body
	var body struct {
		Email string
		Password string
	}

	if c.BindJSON(&user) != nil {
		c.JSON(400, gin.H{
			"message": "invalid json",
		})
		return
	}
	//hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "error hashing password",
		})
		return
	}
	

	//Save the user in the database
	user = models.User{Email:body.Email,Password: string(hashedPassword)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "failed to create user",
		})
	}

	//respond
	c.JSON(http.StatusOK, gin.H{})
}
