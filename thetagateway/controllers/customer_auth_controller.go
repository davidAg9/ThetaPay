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

type CustomerAuthController struct {
	*mongo.Collection
}

func (controller *CustomerAuthController) LoginCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var customer models.Customer
		var foundCustomer models.Customer
		defer cancel()
		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Print(customer)
		err := controller.FindOne(ctx, bson.M{"email": customer.Email}).Decode(&foundCustomer)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := utilities.VerifyPassword(*customer.Password, *foundCustomer.Password)

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundCustomer.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Customer not found"})
		}
		//TODO:FIX GENERATE ALL TOKENS

		token, _ := utilities.GenerateAllTokens(*foundCustomer.Email, *foundCustomer.FullName, foundCustomer.ID.String())
		updateAllTokens(token, foundCustomer.ID.String(), controller)
		err = controller.FindOne(ctx, bson.M{"_id": foundCustomer.ID}).Decode(&foundCustomer)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundCustomer)
	}
}

func (controller *CustomerAuthController) SignUpCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func updateAllTokens(signedToken string, CustomerId string, CustomerCollection *CustomerAuthController) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updatedAt", Updated_at})

	upsert := false
	filter := bson.M{"_id": CustomerId}
	optns := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := CustomerCollection.UpdateOne(
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
