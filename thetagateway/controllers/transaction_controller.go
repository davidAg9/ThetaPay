package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/davidAg9/thetagateway/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionsInterfaces interface {
	TopUp() gin.HandlerFunc
	CheckBalance() gin.HandlerFunc
	ThetaTransfer() gin.HandlerFunc
	Refund() gin.HandlerFunc
	//Retrieves all transactions made for a given user with [id] .
	// [pagination] is the number of transactions returns , 100 by default .
	GetTransactions() gin.HandlerFunc
	AcceptPayment() gin.HandlerFunc
}

type TransactionController struct {
	*mongo.Collection
}

func (controller *TransactionController) TopUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//retrieve transactions from json
		//retrieve the uid from context (taken from token and )

		//get user from database

		//verify if third party service has success

		//if success return status 201

		// if fail return error with status : could not deposit amount , insufficient balance , service timeout try again,

	}
}

func (controller *TransactionController) CheckBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve the uid from context (taken from token and )
		uid := c.Keys["uid"]
		if uid == nil && uid == "" {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "No user id found"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
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

func (controller *TransactionController) ThetaTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func (controller *TransactionController) AcceptPayment() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (controller *TransactionController) Refund() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (controller *TransactionController) GetTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		accid := c.Param("accountId")
		if accid == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "accountId was nil"})
			return
		}
		filter := bson.M{"accountId": accid}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		//get user from database
		var customer models.Customer
		err := controller.FindOne(ctx, filter).Decode(&customer)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "No user found"})
			return
		}

		c.JSON(http.StatusOK, customer.Transactions)

	}
}
