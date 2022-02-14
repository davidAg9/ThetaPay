package routes

import (
	"github.com/davidAg9/thetagateway/controllers"
	"github.com/davidAg9/thetagateway/middlewares"
	"github.com/gin-gonic/gin"
)

func CustomerRoutes(incomingRoutes *gin.Engine, controller *controllers.CustomerController) {
	incomingRoutes.Use(middlewares.AuhthenticateCustomer())
	incomingRoutes.GET("/customers/:userId", controller.GetCustomer())
	incomingRoutes.PUT("/customers/update/:userId", controller.UpdateCustomer())
}
