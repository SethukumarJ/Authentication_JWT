package middleware

import (
	"fmt"
	"jwt/initializers"
	"jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireAuthAdmin(c *gin.Context) {

	fmt.Println("Middleware called middle were called")
	

	//Get the cookie of request

	tokenString, err := c.Cookie("admintoken")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		
	}

	// Parse takes the token string and a function for looking up the key. The latter is especially
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		//CHECK IF THE TOKEN IS EXPIRED
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			
		}

		//FIND THE USER WITH TOKEN SUB
		var admin models.Admin
		initializers.DB.First(&admin, claims["sub"])

		if admin.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			
		}

		//ATTACH TO REQUEST
		c.Set("admin", admin)

		//CALL NEXT
		c.Next()
		fmt.Println(claims["foo"], claims["nbf"])
	
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		
	}

}
