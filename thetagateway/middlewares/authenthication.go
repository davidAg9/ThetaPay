package middlewares

import (
	"net/http"

	"github.com/davidAg9/thetagateway/utilities"
	"github.com/gin-gonic/gin"
)

func AuhthenticateSystemUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		//verify username and password

		//verify role

		//verify request authorization

		// 	c.Set("phoneNumber", claims.PhoneNumber)
		// 	c.Set("userName", claims.FullName)
		// 	c.Set("uid", claims.Uid)
		// 	c.Set("role", claims.Role)
		// 	c.Next()
	}

}

func AuhthenticateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		claims, err := utilities.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("fullName", claims.FullName)
		c.Set("uid", claims.Uid)

		c.Next()
	}
}
