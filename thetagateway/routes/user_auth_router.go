package routes

import (
	"github.com/davidAg9/thetagateway/controllers"
	"github.com/gin-gonic/gin"
)

func AuthUserRoutes(incomingRoutes *gin.Engine, controller *controllers.UserAuthController) {
	incomingRoutes.POST("users/signup", controller.SignUpUser())//localhost:8080/users/signup
	incomingRoutes.POST("users/login", controller.LoginUser())
}
