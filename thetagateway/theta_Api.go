package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/davidAg9/thetagateway/controllers"
	"github.com/davidAg9/thetagateway/databases"
	"github.com/davidAg9/thetagateway/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const TransactionCollection = "transactions"
const UserCollection = "thetaUsers"

//TODO:ENTER DATABASE NAME
const DatabaseName = "thetadb"

func main() {
	//load environment variables
	err := godotenv.Load(".env")

	if err != nil {

		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	mongoUrl := os.Getenv("MONGODB_URL")
	if port == "" {
		port = "8000"
	}

	// connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := databases.ConnnectDatabase(ctx, &mongoUrl)

	if err != nil {
		log.Fatal(err)
	}
	thetaDB := databases.ThetaDatabase{
		client.Database(DatabaseName),
	}

	// setup contollers
	authContoller := &controllers.AuthController{
		*thetaDB.Collection(UserCollection),
	}

	userContoller := &controllers.UserController{
		*thetaDB.Collection(UserCollection),
	}
	transactionContoller := &controllers.TransactionController{
		*thetaDB.Collection(TransactionCollection),
	}

	// start server
	server := gin.Default()

	//configure routes
	server.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})
	routes.AuthRoutes(server, authContoller)
	routes.UserRoutes(server, userContoller)
	routes.TransactionRoutes(server, transactionContoller)
	server.Run(":" + port)
	defer client.Disconnect(ctx)
}
