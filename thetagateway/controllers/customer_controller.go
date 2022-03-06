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

type CustomerController struct {
	*mongo.Collection
}

func (customerController *CustomerController) GetCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.Customer
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		uid := c.MustGet("uid").(string)

		if uid == "" {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "No user id found"})
			return
		}

		if err := customerController.FindOne(ctx, bson.M{"_id": uid}).Decode(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "No user not found"})
			return
		}

		c.JSON(http.StatusAccepted, user)
	}
}

func (customerController *CustomerController) UpdateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		type updatableBody struct {
			FullName   *string   `bson:"fullName,omitempty" json:"fullName,omitempty"  validate:"min=3, max=150"`
			Email      *string   `bson:"email,omitempty" json:"email,omitempty" validate:"email"`
			Updated_at time.Time `bson:"updatedAt,omitempty" json:"-"`
		}
		var updatable updatableBody
		if err := c.BindJSON(&updatable); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fields"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		uid := c.MustGet("uid").(string)
		if uid == "" {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "No user id found"})
			return
		}
		updatable.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		filter := bson.M{"_id": uid}
		upsert := false
		opts := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := customerController.UpdateOne(ctx, filter, updatable, &opts)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fields"})
			return
		}
		var user models.Customer
		if err := customerController.FindOne(ctx, bson.M{"_id": result.UpsertedID}).Decode(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "No user not found"})
			return
		}
		c.JSON(http.StatusAccepted, user)
	}

}

func (customerController *CustomerController) CreateNewApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.MustGet("uid").(string)

		word, err := utilities.GenerateRandomString()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create key"})
			log.Panic(err.Error())
			return
		}
		hash, secret, err := utilities.GenerateApiKey(*word)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create key"})
			log.Panic(err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		updateTime, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		obj, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create key"})
			log.Panic(err)
			return
		}
		filter := bson.M{"_id": obj}
		upsert := false
		opts := options.UpdateOptions{
			Upsert: &upsert,
		}
		var updateObj bson.D = bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "updatedAt", Value: updateTime},
					{Key: "secretKey", Value: secret},
				},
			},
		}

		_, err = customerController.UpdateOne(ctx, filter, updateObj, &opts)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create key"})
			log.Panic(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"key": hash})
	}
}
