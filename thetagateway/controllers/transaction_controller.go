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

var customerCollection = databases.ThetaClient.Database("thetaDb").Collection("customers")

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

func (controller *TransactionController) AcceptPayment() gin.HandlerFunc {
	return func(c *gin.Context) {

		method := c.Query("paymentmethod")
		if method == "momo" {
			var transaction models.MomoTransaction
			err := c.BindJSON(transaction)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "reason": "invalid momo payload"})
				return
			}
			err = performAccreditation(&method, transaction, controller)
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
			err = performAccreditation(&method, transaction, controller)
			if err != nil {
				c.JSON(http.StatusPreconditionFailed, gin.H{"error": "Transaction failed"})
				return
			}

		} else if method == "theta" {
			var transaction models.ThetaTransaction
			err := c.BindJSON(transaction)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "reason": "invalid theta payload"})
				return
			}
			err = performAccreditation(&method, transaction, controller)
			if err != nil {
				c.JSON(http.StatusPreconditionFailed, gin.H{"error": "Transaction failed"})
				return
			}
		}
	}
}

func (controller *TransactionController) Refund() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
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

func performAccreditation(method *string, tt interface{}, controller *TransactionController) error {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	session, err := controller.Database().Client().StartSession()
	if err != nil {
		log.Panic(err.Error())
		return errors.New("could not start transaction")
	}
	defer session.EndSession(context.Background())
	var callback func(sctx mongo.SessionContext) (interface{}, error)
	if *method == "theta" {
		callback = controller.thetaTranferCreditation(tt.(*models.ThetaTransaction))

	} else if *method == "momo" {
		callback = controller.momoTranferCreditation(tt.(*models.MomoTransaction))

	} else if *method == "visa" {
		callback = controller.visaTranferCreditation(tt.(*models.VisaTransaction))

	}
	_, err = session.WithTransaction(context.Background(), callback, txnOpts)
	if err != nil {
		return err
	}
	return nil
}
func (trans *TransactionController) momoTranferCreditation(tt *models.MomoTransaction) func(sessionContext mongo.SessionContext) (interface{}, error) {
	return func(sctx mongo.SessionContext) (interface{}, error) {
		//TODO:VALIDATE MOMO SUCCESS

		//TODO: CREDIT ACCOUNT

		tt.TxnID = primitive.NewObjectID()
		tt.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		tt.Updated_at = tt.Created_at
		tt.Refund = false
		tt.PayMethod = models.MomoMethod
		tt.Status = models.TxnSuccess
		res, err := trans.InsertOne(sctx, tt)
		if err != nil {
			return nil, err
		}
		return res, err
	}
}

func (trans *TransactionController) visaTranferCreditation(tt *models.VisaTransaction) func(sessionContext mongo.SessionContext) (interface{}, error) {
	return func(sctx mongo.SessionContext) (interface{}, error) {
		//TODO:VALIDATE VISA SUCCESS

		//TODO: CREDIT ACCOUNT

		tt.TxnID = primitive.NewObjectID()
		tt.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		tt.Updated_at = tt.Created_at
		tt.Refund = false
		tt.PayMethod = models.VisaMethod
		res, err := trans.InsertOne(sctx, tt)
		if err != nil {
			return nil, err
		}

		return res, err
	}
}

func (trans *TransactionController) thetaTranferCreditation(tt *models.ThetaTransaction) func(sessionContext mongo.SessionContext) (interface{}, error) {
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
		//	persistent transaction to transaction collections

		result, err := trans.InsertOne(sessionContext, tt)
		if err != nil {
			return nil, err
		}

		return result, err

	}

}

// func (controller *TransactionController) ThetaTransfer() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		var transaction models.ThetaTransaction
// 		err := c.BindJSON(transaction)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "reason": "invalid momo payload"})
// 			return
// 		}

// 		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()
// 		transaction.TxnID = primitive.NewObjectID()
// 		transaction.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// 		transaction.Updated_at = transaction.Created_at
// 		transaction.Refund = false
// 		_, err = controller.InsertOne(ctx, transaction)
// 		if err != nil {
// 			log.Panic(err)
// 			c.JSON(http.StatusPreconditionFailed, gin.H{"error": "Transaction failed"})
// 			return
// 		}

// 	}
// }
