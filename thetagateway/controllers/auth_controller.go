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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthController struct {
	mongo.Collection
}

func (controller *AuthController) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User
		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Print(user)
		err := controller.FindOne(ctx, bson.M{"phoneNumber": user.PhoneNumber}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "phoneNumber or password is incorrect"})
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
		//TODO:FIX GENERATE ALL TOKENS

		token, _ := utilities.GenerateAllTokens(*foundUser.PhoneNumber, *foundUser.FullName, *foundUser.Role, foundUser.ID.String())
		updateAllTokens(token, foundUser.ID.String(), controller)
		err = controller.FindOne(ctx, bson.M{"userId": foundUser.ID}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

func (controller *AuthController) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
func updateAllTokens(signedToken string, userId string, userCollection *AuthController) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updatedAt", Updated_at})

	upsert := false
	filter := bson.M{"_id": userId}
	optns := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&optns,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}

}
