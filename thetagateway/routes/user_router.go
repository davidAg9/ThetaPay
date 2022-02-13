package routes

import (
	"github.com/davidAg9/thetagateway/controllers"
	"github.com/davidAg9/thetagateway/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine, controller *controllers.UserController) {
	incomingRoutes.Use(middlewares.AuhthenticateUser())
	incomingRoutes.GET("/users/:userId", controller.GetUser())
	incomingRoutes.PUT("/update/:userId", controller.UpdateUser())
}
