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
func TransactionRoutes(incomingRoutes *gin.Engine, transactionController *controllers.TransactionController, customerContoler *controllers.CustomerController) {
	incomingRoutes.GET("/transactions/:merchantId", transactionController.GetTransactions())
	incomingRoutes.Use(middlewares.VerifyApiKey(customerContoler))
	incomingRoutes.POST("/transactions/refund/:txnId", transactionController.Refund())
	incomingRoutes.POST("/transactions/theta2theta", transactionController.ThetaTransfer())

}
