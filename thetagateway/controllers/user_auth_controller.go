package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/davidAg9/thetagateway/models"
	"github.com/davidAg9/thetagateway/utilities"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserAuthController struct {
	*mongo.Collection
}

// login a user
// @Summary Users Login
// @Description Users Login with Email & Password
// @Tags Auth
// @Accept application/json
// @Produce json
// @Param Body body object true "Login Body"
// @Router /users/login [post]
// @Success 200 {object} models.User
func (controller *UserAuthController) LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
		var user models.User
		var foundUser models.User
		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Print(user)
		err := controller.FindOne(ctx, bson.M{"username": user.UserName}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username or password is incorrect"})
			return
		}

		passwordIsValid, msg := utilities.VerifyPassword(*user.Password, *foundUser.Password)

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.PhoneNumber == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}

		err = controller.FindOne(ctx, bson.M{"userId": foundUser.ID}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

func (controller *UserAuthController) SignUpUser() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
