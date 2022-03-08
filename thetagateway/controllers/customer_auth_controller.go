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

// login a Customer
// @Summary Customer Login
// @Description Cutomer Login with Email & Password
// @Tags Auth
// @Accept application/json
// @Produce json
// @Param Body body object true "Login Body"
// @Router /customers/login [post]
// @Success 200 {object} models.Customer
func (controller *CustomerAuthController) LoginCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		var customer models.Customer
		var foundCustomer models.Customer
		defer cancel()
		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := controller.FindOne(ctx, bson.M{"email": customer.Email}).Decode(&foundCustomer)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email might be incorrect"})
			return
		}

		passwordIsValid, msg := utilities.VerifyPassword(*customer.Password, *foundCustomer.Password)

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundCustomer.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Customer not found"})
			return
		}

		foundId := foundCustomer.ID.Hex()
		creds := &utilities.ThetaCustomerCredentials{
			Uid:      &foundId,
			FullName: foundCustomer.FullName,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			},
		}

		token, err := utilities.GenerateCustomerTokens(*creds)
		log.Println(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = updateAllTokens(token, foundCustomer.ID, controller)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = controller.FindOne(ctx, bson.M{"_id": foundCustomer.ID}).Decode(&foundCustomer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		foundCustomer.Password = nil
		log.Println(foundCustomer)
		c.JSON(http.StatusOK, gin.H{"token": foundCustomer.Token, "user": foundCustomer})
	}
}

// SignUp a customer
// @Summary Customer SignUp
// @Description Signup a new Customer
// @Tags Auth
// @Accept application/json
// @Produce json
// @Param Body body models.Customer true "Signup Body"
// @Router /customers/signup [post]
// @Success 200 {string} result
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

			c.JSON(http.StatusInternalServerError, "Error occured while signing in")
			log.Panic(err.Error())
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
		customer.Deleted_at = nil
		customer.Verified = false
		//Generate account information
		customer.AccountInfo, err = GenerateAccountInformation(*customer.AccountInfo.AccountType)
		if err != nil {
			log.Panic(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
			return
		}
		uid := customer.ID.Hex()
		creds := utilities.ThetaCustomerCredentials{
			Uid:      &uid,
			FullName: customer.FullName,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			},
		}
		token, err := utilities.GenerateCustomerTokens(creds)
		if err != nil {
			log.Panic(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
			return
		}
		customer.Token = &token
		if err != nil {
			c.String(http.StatusInternalServerError, "Could not create user")
			return
		}
		_, err = controller.InsertOne(ctx, customer)
		if err != nil {
			c.String(http.StatusInternalServerError, "Could not create user")
			return
		}

		customer.Password = nil
		customer.API_KEY = nil
		c.JSON(http.StatusOK, gin.H{"token": customer.Token, "user": customer})

	}
}

func GenerateAccountInformation(accounType models.AccountType) (*models.AccountInfo, error) {
	var account models.AccountInfo
	accNo, err := utilities.GenerateAccountNumber()
	if err != nil {
		return nil, err
	}

	balance := 0.0
	account.AccountID = &accNo
	account.Balance = &balance
	account.AccountType = &accounType
	return &account, nil
}

func updateAllTokens(signedToken string, CustomerId primitive.ObjectID, CustomerCollection *CustomerAuthController) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: updated_at})

	upsert := false
	filter := bson.M{"_id": CustomerId}
	optns := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := CustomerCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&optns,
	)

	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}
