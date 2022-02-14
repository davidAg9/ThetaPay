package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionsInterfaces interface {
	TopUp() gin.HandlerFunc
	CheckBalance() gin.HandlerFunc
	ThetaTransfer() gin.HandlerFunc
	Refund() gin.HandlerFunc
	GetTransactions() gin.HandlerFunc
	AcceptPayment() gin.HandlerFunc
}

type TransactionController struct {
	*mongo.Collection
}

func (controller *TransactionController) TopUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (controller *TransactionController) CheckBalance() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (controller *TransactionController) ThetaTransfer() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
func (controller *TransactionController) AcceptPayment() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (controller *TransactionController) Refund() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (controller *TransactionController) GetTransactions() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
