package controllers

import (
	"fmt"
	"jwt/initializers"
	"jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func UserSignup(c *gin.Context) {
	c.HTML(http.StatusOK, "usersignup.html", nil)
}

func UserPostSignup(c *gin.Context) {

	usernameFromForm := c.Request.FormValue("username")
	passwordFromForm := c.Request.FormValue("password")

	var user models.User
	//GEt email and password from the form
	var form struct {
		Email    string
		Password string
	}

	form.Email = usernameFromForm
	form.Password = passwordFromForm
	//hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
	if err != nil {
		fmt.Println("error hashing password")
		c.Redirect(303, "/userSignup")
		return
	}

	//Save the user in the database
	user = models.User{Email: form.Email, Password: string(hashedPassword)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.Redirect(303, "/userSignup")
		fmt.Println("error saving user")
		return
	}

	//respond
	c.HTML(http.StatusOK, "userlogin.html", nil)
}




func UserLogin (c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	ok := welcomeUserStatus
	if ok {
		c.Redirect(303, "/welcomeUser")
		return
	}
	c.HTML(http.StatusOK, "userlogin.html", nil)

}
func UserPostLogin(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	//GEt email and password from request body

	usernameFromForm := c.Request.FormValue("username")
	passwordFromForm := c.Request.FormValue("password")


	var form struct {
		Email    string
		Password string
	}

	form.Email = usernameFromForm
	form.Password = passwordFromForm

	//Check if the user exists in the database
	var user models.User

	initializers.DB.First(&user, "email = ?", form.Email)

	if user.ID == 0 {
		c.Redirect(303, "/userLogin")
		fmt.Println("invalid user")
		return
	}

	//Check if the password is correct

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		c.Redirect(303, "/userLogin")
		fmt.Println("invalid user")
		return
	}

	//Generate a jwt token
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		fmt.Println("error signing token")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error signing token",
		})
		return
	}

	//respond
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenString, 3600, "", "", false, true)
	c.HTML(http.StatusOK, "welcomeUser.html", nil)

}

func UserLogout(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", "", -1, "", "", false, true)
	c.HTML(http.StatusOK, "userlogin.html", nil)

}


var welcomeUserStatus = false
func WelcomeUser(c *gin.Context){

	
	welcomeUserStatus = true
}

