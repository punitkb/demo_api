package api

import (
	"sezzle_api/src/config"
	"github.com/gin-gonic/gin"
	"sezzle_api/src/controllers"
	"sezzle_api/src/repository"
	"fmt"

)

func RunServer() {
	
	//initialize DB connection
	db, err := config.InitDb()
	defer db.Close()
	if err != nil {
		panic(err)
	}

	//initialize redis connection
	redisDb,err := config.InitRedisdb()

	connectionRepo := repository.NewRepository(db,redisDb)
    if err := connectionRepo.Init(); err != nil {
        panic(err)
    }

	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.Use()
	{
		
		user := controllers.NewUserController(connectionRepo)
		v1.POST("/login", user.Login)
		v1.GET("/logout", user.Logout)
		v1.POST("/user/create", user.AddUser)
		v1.GET("/user/list", user.ListUsers)

		item := controllers.NewItemController(connectionRepo)
		v1.POST("/item/create", item.AddItem)
		v1.GET("/item/list", item.ListItems)

		cart := controllers.NewCartController(connectionRepo)
		v1.GET("/cart/list",cart.ListCarts)
		v1.POST("/cart/add",cart.AddToCart)
		v1.POST("cart/complete/:cartId", cart.CompleteOrder)
		
		order := controllers.NewOrderController(connectionRepo)
		v1.GET("/order/list",order.ListOrders)

	}

	//Run api 
	router.Run(":9001")
	fmt.Println("Server listening on port 9001")
}
