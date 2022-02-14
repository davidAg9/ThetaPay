package routes

import (
	"github.com/davidAg9/thetagateway/controllers"
	"github.com/gin-gonic/gin"
)

func AuthUserRoutes(incomingRoutes *gin.Engine, controller *controllers.UserAuthController) {
	incomingRoutes.POST("users/signup", controller.SignUpUser())
	incomingRoutes.POST("users/login", controller.LoginUser())
}
