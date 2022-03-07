package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/davidAg9/thetagateway/databases"
	"github.com/davidAg9/thetagateway/models"
	"github.com/davidAg9/thetagateway/utilities"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type CustomerController struct {
	*mongo.Collection
}

var transactionCollection = databases.ThetaClient.Database("thetadb").Collection("transactions")

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

func (controller *CustomerController) TopUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve transactions from json
		var balance float64
		method := c.Query("paymentmethod")
		if method == "momo" {
			var transaction models.MomoTransaction
			err := c.BindJSON(transaction)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "reason": "invalid momo payload"})
				return
			}
			balance, err = controller.momoTranferCreditation(&transaction)
			if err != nil {
				c.JSON(http.StatusPreconditionFailed, gin.H{"error": "Transaction failed"})
				return
			}
		} else if method == "visa" {
			var transaction models.VisaTransaction
			err := c.BindJSON(transaction)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "reason": "invalid visa payload"})
				return
			}
			balance, err = controller.visaTranferCreditation(&transaction)
			if err != nil {
				c.JSON(http.StatusPreconditionFailed, gin.H{"error": "Transaction failed"})
				return
			}

		}
		c.JSON(http.StatusOK, gin.H{"balance": balance})

		//verify if third party service has success

		//if success return status 201

		// if fail return error with status : could not deposit amount , insufficient balance , service timeout try again,

	}
}

func (controller *CustomerController) momoTranferCreditation(tt *models.MomoTransaction) (float64, error) {

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	session, err := controller.Database().Client().StartSession()
	if err != nil {
		log.Panic(err.Error())
		return 0, errors.New("could not start transaction")
	}
	defer session.EndSession(context.Background())
	var newCredBalance float64
	callback := func(sctx mongo.SessionContext) (interface{}, error) {
		//TODO:VALIDATE MOMO SUCCESS

		//CREDIT ACCOUNT
		credFilter := bson.M{"accountInfo.accountId": tt.MerchantId}
		var customer models.Customer
		err = controller.FindOne(sctx, credFilter).Decode(&customer)
		if err != nil {
			return nil, err
		}
		// calculate credit balance
		newCredBalance = *customer.AccountInfo.Balance + *tt.Amount
		updateCred := bson.D{
			{
				Key: "$set", Value: bson.E{
					Key: "accountInfo.balance", Value: newCredBalance,
				},
			},
		}
		//update creditor balance
		_, err = controller.UpdateOne(sctx, credFilter, updateCred)
		if err != nil {
			return nil, err
		}
		tt.TxnID = primitive.NewObjectID()
		tt.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		tt.Updated_at = tt.Created_at
		tt.Refund = false
		tt.PayMethod = models.MomoMethod
		tt.Status = models.TxnSuccess
		res, err := transactionCollection.InsertOne(sctx, tt)
		if err != nil {
			return nil, err
		}
		return res, err
	}
	_, err = session.WithTransaction(context.Background(), callback, txnOpts)
	if err != nil {
		return 0, err
	}
	return newCredBalance, nil

}

func (controller *CustomerController) visaTranferCreditation(tt *models.VisaTransaction) (float64, error) {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	session, err := controller.Database().Client().StartSession()
	if err != nil {
		log.Panic(err.Error())
		return 0, errors.New("could not start transaction")
	}
	defer session.EndSession(context.Background())
	var newCredBalance float64
	callback := func(sctx mongo.SessionContext) (interface{}, error) {
		//TODO:VALIDATE VISA SUCCESS

		// CREDIT ACCOUNT
		credFilter := bson.M{"accountInfo.accountId": tt.MerchantId}
		var customer models.Customer
		err = controller.FindOne(sctx, credFilter).Decode(&customer)
		if err != nil {
			return nil, err
		}
		// calculate credit balance
		newCredBalance = *customer.AccountInfo.Balance + *tt.Amount
		updateCred := bson.D{
			{
				Key: "$set", Value: bson.E{
					Key: "accountInfo.balance", Value: newCredBalance,
				},
			},
		}
		//update creditor balance
		_, err = controller.UpdateOne(sctx, credFilter, updateCred)
		if err != nil {
			return nil, err
		}
		tt.TxnID = primitive.NewObjectID()
		tt.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		tt.Updated_at = tt.Created_at
		tt.Refund = false
		tt.PayMethod = models.VisaMethod
		tt.Status = models.TxnSuccess
		res, err := transactionCollection.InsertOne(sctx, tt)
		if err != nil {
			return nil, err
		}

		return res, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)
	if err != nil {
		return 0, err
	}
	return newCredBalance, nil

}

func (controller *CustomerController) CheckBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve the uid from context (taken from token and )
		uid := c.MustGet("uid").(string)
		if uid == "" {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "No user id found"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		//get user from database
		var customer models.Customer
		err := controller.FindOne(ctx, bson.M{"_id": uid}).Decode(&customer)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "No user found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"balance": customer.AccountInfo.Balance})

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
