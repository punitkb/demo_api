package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	models "sezzle_api/src/models"
	"sezzle_api/src/repository"
	"net/http"
)

type OrderController struct{
	Repo *repository.Repository
}

var OrderModel = new(models.OrderModel)


func (orderController OrderController) ListOrders(c *gin.Context) {
	data,err:=  OrderModel.ListOrders(orderController.Repo)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while listing orders"})
		return
	}
	c.JSON(200, gin.H{"orders": data})
}



func NewOrderController(repo *repository.Repository) *OrderController {
	orderController :=  &OrderController{
		Repo : repo,
	}
	return orderController
}