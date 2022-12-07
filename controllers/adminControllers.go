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

// ===================================================ADMIN ROUTES===========================================

//===================ADMIN LOGIN=====================

func AdminLogin(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	ok := adminLoggedStatus

	if ok {
		c.Redirect(303, "/adminProfile")
		return
	}
	c.HTML(http.StatusOK, "adminlogin.html", nil)

}

//===================POST LOIGN=====================

type User struct {
	ID       int
	UserName string
	Password string
}

func AdminPostLogin(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	//GEt email and password from request body

	adminNameFromForm := c.Request.FormValue("adminName")
	adminPasswordFromForm := c.Request.FormValue("password")

	var form struct {
		Name     string
		Password string
	}

	form.Name = adminNameFromForm
	form.Password = adminPasswordFromForm

	//Check if the admin exists in the database
	var admin models.Admin

	initializers.DB.First(&admin, "name = ?", form.Name)

	if admin.ID == 0 {
		c.Redirect(303, "/adminLogin")
		fmt.Println("invalid admin name")
		return
	}

	//Check if the password is correct

	// err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(form.Password))
	// if err != nil {
	// 	c.Redirect(303, "/adminLogin")
	// 	fmt.Println("invalid password")
	// 	return
	// }

	//Generate a jwt token
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": admin.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		fmt.Println("error signing token")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error signing admintoken",
		})
		return
	}

	//respond
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("admintoken", tokenString, 3600, "", "", false, true)
	c.Redirect(303, "/adminProfile")
	adminLoggedStatus = true
	fmt.Println("admin logged in")

}

// ===================ADMIN LOGOUT=====================
func AdminLogout(c *gin.Context) {

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("admintoken", "", -1, "", "", false, true)
	adminLoggedStatus = false
	c.Redirect(303, "/adminLogin")

}
// ===================ADMIN Status=====================
var adminLoggedStatus = false

func AdminLogged(c *gin.Context) {

	fmt.Println("admin logged status set to true")
	adminLoggedStatus = true
}

// ===================ADMIN PROFILE=====================

func AdminProfile(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	///////////////////////////////////////////////////////////////
	ok := adminLoggedStatus
	fmt.Println("admin logged status is ", ok)
	if ok {

		var user []models.User

		var us [50]string

		var id [50]uint
		initializers.DB.Raw("SELECT id,name FROM users").Scan(&user)
		for ind, i := range user {
			us[ind], id[ind] = i.Name, i.ID

		}

		c.HTML(http.StatusOK, "adminProfile.html", gin.H{

			"users": us,
			"id":    id,
		})

		fmt.Println("fetching users")
		return
	}
	c.Redirect(303, "/adminLogin")

}

// ===================CREATE USER FROM ADMIN PANEL=====================

func CreateUser(c *gin.Context) {

	usernameFromForm := c.Request.FormValue("username")
	passwordFromForm := c.Request.FormValue("password")

	var user models.User
	//GEt email and password from the form
	var form struct {
		Name     string
		Password string
	}

	form.Name = usernameFromForm
	form.Password = passwordFromForm
	//hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
	if err != nil {
		fmt.Println("error hashing password")
		c.Redirect(303, "/adminProfile")
		return
	}

	//Save the user in the database
	user = models.User{Name: form.Name, Password: string(hashedPassword)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.Redirect(303, "/adminProfile")
		fmt.Println("error saving user")
		return
	}

	//respond
	c.Redirect(303, "/adminProfile")
}


// ===================DELETE USER FROM ADMIN PANEL=====================

func DeleteUser(c *gin.Context) {
	
	var user models.User
	name := c.Param("name")
	
	initializers.DB.Where("name=?", name).Delete(&user)
	c.Redirect(303, "/adminProfile")
	fmt.Println("user deleted")

}