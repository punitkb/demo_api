package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	models "sezzle_api/src/models"
	"sezzle_api/src/repository"
	"net/http"
)

type ItemController struct{
	Repo	*repository.Repository
}

var ItemModel = new(models.ItemModel)

func (itemController ItemController) AddItem(c *gin.Context) {
	type Param struct {
		Name		string		`json:"item_name"`
	}

	var pjson Param
	if err := c.ShouldBindJSON(&pjson); err == nil {
		if pjson.Name == ""  {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide values for  item_name"})
			return
		}

		err := ItemModel.AddItem(pjson.Name,itemController.Repo)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "err"})
			return
		}
	} else{
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid json"})
		return
	}

	c.JSON(200, gin.H{"message": "Item added successfully"})
}

func (itemController ItemController) ListItems(c *gin.Context) {
	data,err:=  ItemModel.ListItems(itemController.Repo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while listing users"})
		return
	}
	c.JSON(200, gin.H{"items": data})
}


func NewItemController(repo *repository.Repository) *ItemController {
	itemController :=  &ItemController{
		Repo : repo,
	}
	return itemController
}