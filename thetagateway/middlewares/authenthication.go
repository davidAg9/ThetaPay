package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/davidAg9/thetagateway/controllers"
	"github.com/davidAg9/thetagateway/models"
	"github.com/davidAg9/thetagateway/utilities"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

		c.Set("username", *claims.UserName)
		c.Set("role", *claims.Role)
		c.Set("uid", *claims.Uid)
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

		c.Set("fullName", *claims.FullName)
		c.Set("uid", *claims.Uid)

		c.Next()
	}
}

func VerifyApiKey(controller *controllers.CustomerController) gin.HandlerFunc {
	return func(c *gin.Context) {

		var customer models.Customer
		var transaction bson.M
		apikey := c.Request.Header.Get("key")
		if apikey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No Api key provided"})
			c.Abort()
			return
		}

		err := c.BindJSON(&transaction)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err = controller.FindOne(ctx, bson.M{"accountInfo.accountId": transaction["merchantId"]}).Decode(&customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		valid, err := utilities.ValidateAPIToken(&apikey, customer.API_KEY)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if !valid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid api key"})
			c.Abort()
			return
		}
		// c.Set("accNo", *customer.AccountInfo.AccountID)
		c.Next()
	}
}
