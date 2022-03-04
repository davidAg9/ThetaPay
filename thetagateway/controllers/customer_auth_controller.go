package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/davidAg9/thetagateway/models"
	"github.com/davidAg9/thetagateway/utilities"
	"github.com/dgrijalva/jwt-go"
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
		var ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
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

		foundId := foundCustomer.ID.String()
		creds := &utilities.ThetaCustomerCredentials{
			UserName: foundCustomer.Username,
			Uid:      &foundId,
			FullName: foundCustomer.FullName,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			},
		}

		token, _ := utilities.GenerateCustomerTokens(*creds)
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
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
		var customer models.Customer
		defer cancel()
		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		count, err := controller.CountDocuments(ctx, bson.M{"email": customer.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, "Error occured while signing in")
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, "User already Registered ")
			return
		}

		//verify password
		pass := *customer.Password

		hash := utilities.HashPassword(pass)

		customer.Password = &hash
		customer.ID = primitive.NewObjectID()
		customer.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		customer.Updated_at = customer.Created_at
		//Generate account information
		customer.AccountInfo, err = GenerateAccountInformation()
		if err != nil {
			c.String(http.StatusInternalServerError, "Could not create user")
			return
		}
		result, err := controller.InsertOne(ctx, customer)
		if err != nil {
			c.String(http.StatusInternalServerError, "Could not create user")
			return
		}
		c.JSON(http.StatusOK, result)

	}
}

func GenerateAccountInformation() (*models.AccountInfo, error) {
	var account models.AccountInfo
	accNo, err := utilities.GenerateAccountNumber()
	if err != nil {
		return nil, err
	}

	pin, err := utilities.GenerateTempPin()
	if err != nil {
		return nil, err
	}
	balance := 0.0
	account.AccountID = &accNo
	account.PinCode = pin
	account.Balance = &balance
	return &account, nil
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
