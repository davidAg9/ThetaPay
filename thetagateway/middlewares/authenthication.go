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

		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		claims, err := utilities.ValidateUserToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("username", claims.UserName)
		c.Set("role", claims.Role)
		c.Set("uid", claims.Uid)
		c.Next()
	}

}

func AuhthenticateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		claims, err := utilities.ValidateCustomerToken(clientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("username", claims.UserName)
		c.Set("fullName", claims.FullName)
		c.Set("uid", claims.Uid)

		c.Next()
	}
}

func VerifyApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO:PROPER VALIDATION
		apikey := c.Request.Header.Get("key")
		if apikey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No Api key provided"})
			c.Abort()
			return
		}

		claims, err := utilities.ValidateAPIToken(apikey)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Issuer)

		c.Next()
	}
}
