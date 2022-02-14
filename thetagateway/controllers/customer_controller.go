package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerController struct {
	*mongo.Collection
}

func (customerController *CustomerController) GetCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (customerController *CustomerController) UpdateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
