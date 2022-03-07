package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/davidAg9/thetagateway/databases"
	"github.com/davidAg9/thetagateway/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var customerCollection = databases.ThetaClient.Database("thetadb").Collection("customers")

type TransactionsInterfaces interface {
	ThetaTransfer() gin.HandlerFunc
	Refund() gin.HandlerFunc
	//Retrieves all transactions made for a given user with [id] .
	// [pagination] is the number of transactions returns , 100 by default .
	GetTransactions() gin.HandlerFunc
}

type TransactionController struct {
	*mongo.Collection
}

func (controller *TransactionController) GetTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		accid := c.Param("accountId")
		if accid == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "accountId was nil"})
			return
		}
		filter := bson.M{"merchantId": accid}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		//get user from database
		var transactions []bson.M
		cursor, err := controller.Find(ctx, filter)
		cursor.All(ctx, transactions)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "No user found"})
			return
		}

		c.JSON(http.StatusOK, transactions)

	}
}

func (controller *TransactionController) Refund() gin.HandlerFunc {
	return func(c *gin.Context) {
		txnid := c.Param("txnId")
		if txnid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No transaction id was specified"})
			return
		}
		var transaction models.ThetaTransaction
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		objId, err := primitive.ObjectIDFromHex(txnid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = controller.FindOne(ctx, bson.M{"_id": objId}).Decode(&transaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve transaction"})
			return
		}

		err = performAccreditation(&transaction, controller, true)
		if err != nil {
			c.JSON(http.StatusPreconditionFailed, gin.H{"error": "Transaction failed"})
			log.Panic(err)
			return
		}

	}
}

func (controller *TransactionController) ThetaTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {

		var transaction models.ThetaTransaction
		err := c.BindJSON(transaction)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "reason": "invalid theta payload"})
			return
		}

		err = performAccreditation(&transaction, controller, false)
		if err != nil {
			c.JSON(http.StatusPreconditionFailed, gin.H{"error": "Transaction failed"})
			log.Panic(err)
			return
		}

	}
}

func performAccreditation(tt *models.ThetaTransaction, controller *TransactionController, refund bool) error {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	session, err := controller.Database().Client().StartSession()
	if err != nil {
		log.Panic(err.Error())
		return errors.New("could not start transaction")
	}
	defer session.EndSession(context.Background())
	var callback func(sessionContext mongo.SessionContext) (inter interface{}, err error)
	if refund {
		callback = controller.thetaTranferRefund(tt)
	} else {
		callback = controller.thetaTranferCreditation(tt)
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)
	if err != nil {
		return err
	}
	return nil
}

func (trans *TransactionController) thetaTranferCreditation(tt *models.ThetaTransaction) func(sessionContext mongo.SessionContext) (interface{}, error) {
	var debitCustomer models.Customer
	var creditCustomer models.Customer
	debitFilter := bson.M{"accountInfo.accountId": *tt.AccountId}
	creditFilter := bson.M{"accountInfo.accountId": *tt.MerchantId}

	return func(sessionContext mongo.SessionContext) (inter interface{}, err error) {
		//find debit account and check balance
		err = customerCollection.FindOne(sessionContext, debitFilter).Decode(&debitCustomer)
		if err != nil {
			return nil, err
		}
		if *debitCustomer.AccountInfo.Balance < *tt.Amount {
			return nil, errors.New(" Account balance is insufficient")
		}
		//substract balance
		newBalance := *debitCustomer.AccountInfo.Balance - *tt.Amount
		updateDeb := bson.D{
			{
				Key: "$set", Value: bson.E{
					Key: "accountInfo.balance", Value: newBalance,
				},
			},
		}
		//update balance in debit account
		_, err = customerCollection.UpdateOne(sessionContext, debitFilter, updateDeb)
		if err != nil {
			return nil, err
		}
		//find creditor
		err = customerCollection.FindOne(sessionContext, creditFilter).Decode(&creditCustomer)
		if err != nil {
			return nil, err
		}
		// calculate credit balance
		newCredBalance := *creditCustomer.AccountInfo.Balance + *tt.Amount
		updateCred := bson.D{
			{
				Key: "$set", Value: bson.E{
					Key: "accountInfo.balance", Value: newCredBalance,
				},
			},
		}
		//update creditor balance
		_, err = customerCollection.UpdateOne(sessionContext, creditFilter, updateCred)
		if err != nil {
			return nil, err
		}
		tt.TxnID = primitive.NewObjectID()
		tt.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		tt.Updated_at = tt.Created_at
		tt.Refund = false
		tt.PayMethod = models.ThetaMethod
		tt.Trans_Type = models.Tranfer
		//	persistent transaction to transaction collections

		result, err := trans.InsertOne(sessionContext, tt)
		if err != nil {
			return nil, err
		}

		return result, err

	}

}

func (trans *TransactionController) thetaTranferRefund(tt *models.ThetaTransaction) func(sessionContext mongo.SessionContext) (interface{}, error) {
	var debitCustomer models.Customer
	var creditCustomer models.Customer
	debitFilter := bson.M{"accountInfo.accountId": *tt.MerchantId}
	creditFilter := bson.M{"accountInfo.accountId": *tt.AccountId}

	return func(sessionContext mongo.SessionContext) (inter interface{}, err error) {
		//find debit account and check balance
		err = customerCollection.FindOne(sessionContext, debitFilter).Decode(&debitCustomer)
		if err != nil {
			return nil, err
		}
		if *debitCustomer.AccountInfo.Balance < *tt.Amount {
			return nil, errors.New(" Account balance is insufficient")
		}
		//substract balance
		newBalance := *debitCustomer.AccountInfo.Balance - *tt.Amount
		updateDeb := bson.D{
			{
				Key: "$set", Value: bson.E{
					Key: "accountInfo.balance", Value: newBalance,
				},
			},
		}
		//update balance in debit account
		res, err := customerCollection.UpdateOne(sessionContext, debitFilter, updateDeb)
		if err != nil {
			return nil, err
		}
		//find creditor
		err = customerCollection.FindOne(sessionContext, creditFilter).Decode(&creditCustomer)
		if err != nil {
			return nil, err
		}
		// calculate credit balance
		newCredBalance := *creditCustomer.AccountInfo.Balance + *tt.Amount
		updateCred := bson.D{
			{
				Key: "$set", Value: bson.E{
					Key: "accountInfo.balance", Value: newCredBalance,
				},
			},
		}
		//update creditor balance
		res, err = customerCollection.UpdateOne(sessionContext, creditFilter, updateCred)
		if err != nil {
			return nil, err
		}
		tt.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		tt.Refund = true
		//	persistent transaction to transaction collections
		upsert := false
		opts := options.UpdateOptions{
			Upsert: &upsert,
		}
		res, err = trans.UpdateByID(sessionContext, tt.TxnID, tt, &opts)
		if err != nil {
			return nil, err
		}

		return res, err
	}
}
