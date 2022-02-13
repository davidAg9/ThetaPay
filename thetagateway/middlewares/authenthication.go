package middlewares

import (
	"net/http"

	"github.com/davidAg9/thetagateway/utilities"
	"github.com/gin-gonic/gin"
)

func AuhthenticateUser() gin.HandlerFunc {
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

		c.Set("email", claims.PhoneNumber)
		c.Set("fullName", claims.FullName)
		c.Set("uid", claims.Uid)
		c.Set("role", claims.Role)
		c.Next()
	}

}
