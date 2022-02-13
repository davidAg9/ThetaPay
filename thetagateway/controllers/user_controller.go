package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/davidAg9/thetagateway/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	mongo.Collection
}

func (userController *UserController) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		if userId != "" {
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			var user models.User
			err := userController.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
			defer cancel()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, user)
		} else {
			err := errors.New("Invalid user identity")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	}
}

func (userController *UserController) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		userId := c.Param("userId")
		filter := bson.M{"_id": userId}

		defer cancel()
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D
		if user.FullName != nil {
			updateObj = append(updateObj, bson.E{"fullName", user.FullName})
		}
		if user.PhoneNumber != nil {
			updateObj = append(updateObj, bson.E{"phoneNumber", user.PhoneNumber})
		}

		if updateObj == nil {
			c.JSON(200, gin.H{"error": "Update body cannot be empty"})
			return
		}
		Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updatedAt", Updated_at})

		_, updateErr := userController.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
		)

		if updateErr != nil {
			msg := fmt.Sprintf("Update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(200, gin.H{"message": "Success"})
	}
}
