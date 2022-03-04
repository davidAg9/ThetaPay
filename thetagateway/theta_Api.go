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
	"go.mongodb.org/mongo-driver/mongo/options"
)

const TransactionCollection = "transactions"
const UserCollection = "thetaUsers"
const CustomerCollection = "customers"
const AuditsCollection = "audits"

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
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(mongoUrl).SetServerAPIOptions(serverAPIOptions)
	// connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()
	client, err := databases.ConnnectDatabase(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	thetaDB := databases.ThetaDatabase{
		client.Database(DatabaseName),
	}

	// setup contollers
	userAuthContoller := &controllers.UserAuthController{
		thetaDB.Collection(UserCollection),
	}

	userContoller := &controllers.UserController{
		thetaDB.Collection(UserCollection),
	}
	customerAuthController := &controllers.CustomerAuthController{
		thetaDB.Collection(CustomerCollection),
	}
	customerController := &controllers.CustomerController{
		thetaDB.Collection(CustomerCollection),
	}
	transactionContoller := &controllers.TransactionController{
		thetaDB.Collection(TransactionCollection),
	}

	// start server
	server := gin.Default()

	//configure routes
	server.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})
	routes.AuthCustomerRoutes(server, customerAuthController)
	routes.CustomerRoutes(server, customerController)
	routes.AuthUserRoutes(server, userAuthContoller)
	routes.UserRoutes(server, userContoller)
	routes.TransactionRoutes(server, transactionContoller)
	server.Run(":" + port)

}
