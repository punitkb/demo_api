package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	models "sezzle_api/src/models"
	"sezzle_api/src/repository"
	"net/http"
	"strconv"
)

type CartController struct{
	Repo	*repository.Repository
}

var CartModel = new(models.CartModel)

func (cartController CartController) AddToCart(c *gin.Context) {

    token := c.GetHeader("token")
    
    //validating requested 
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
        c.Abort()
        return
    }

    redisClient := cartController.Repo.Redis
	val, err  := redisClient.Get(token).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please login to add item to cart."})
		return
	}

	if val == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	type Param struct {
		Name		string		`json:"item_name"`
	}

	var pjson Param
	if err := c.ShouldBindJSON(&pjson); err == nil {

		if pjson.Name == ""  {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide values for  item_name"})
			return
		}

		err := CartModel.AddToCart(pjson.Name,val,cartController.Repo)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while adding item to cart"})
			return
		}

	} else{
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid json"})
		return
	}

	c.JSON(200, gin.H{"message": "item added successfully"})
}


func (cartController CartController) CompleteOrder(c *gin.Context) {

    token := c.GetHeader("token")
    //validating requested 
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
        c.Abort()
        return
    }

	cartId,err := strconv.Atoi(c.Param("cartId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not a valid cartId"})
	}

	redisClient := cartController.Repo.Redis
	val, err  := redisClient.Get(token).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please login to add item to cart."})
		return
	}

	if val == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	err = CartModel.CompleteOrder(cartId,val,cartController.Repo)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while completing order"})
		return
	}

	c.JSON(200, gin.H{"message": "Order completed"})
}

func (cartController CartController) ListCarts(c *gin.Context) {
	data,err:=  CartModel.ListCarts(cartController.Repo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while listing users"})
		return
	}
	c.JSON(200, gin.H{"carts": data})
}


func NewCartController(repo *repository.Repository) *CartController {
	cartController :=  &CartController{
		Repo : repo,
	}
	return cartController
}