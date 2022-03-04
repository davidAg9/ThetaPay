package routes

import (
	"github.com/davidAg9/thetagateway/controllers"
	"github.com/davidAg9/thetagateway/middlewares"
	"github.com/gin-gonic/gin"
)

// TopUp() gin.HandlerFunc
// CheckBalance() gin.HandlerFunc
// ThetaTransfer() gin.HandlerFunc
// Refund() gin.HandlerFunc
// GetTransactions() gin.HandlerFunc
// AcceptPayment() gin.HandlerFunc
func TransactionRoutes(incomingRoutes *gin.Engine, transactionController *controllers.TransactionController) {
	incomingRoutes.Use(middlewares.AuhthenticateCustomer())
	incomingRoutes.POST("/transactions/users/topup", transactionController.TopUp())
	incomingRoutes.GET("transactions/balance", transactionController.CheckBalance())
	incomingRoutes.GET("/transactions/:accountId", transactionController.GetTransactions())
	incomingRoutes.Use(middlewares.VerifyApiKey())
	incomingRoutes.POST("/transactions/theta2theta", transactionController.ThetaTransfer())
	incomingRoutes.POST("/transactions/pay", transactionController.AcceptPayment())
	incomingRoutes.POST("transactions/refund/:txnId", transactionController.Refund())
}
