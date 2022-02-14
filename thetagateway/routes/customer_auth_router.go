package routes

import (
	"github.com/davidAg9/thetagateway/controllers"
	"github.com/gin-gonic/gin"
)

func AuthCustomerRoutes(incomingRoutes *gin.Engine, controller *controllers.CustomerAuthController) {
	incomingRoutes.POST("customers/signup", controller.SignUpCustomer())
	incomingRoutes.POST("customers/login", controller.LoginCustomer())
}
