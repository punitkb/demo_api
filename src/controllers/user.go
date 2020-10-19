package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	models "sezzle_api/src/models"
	"net/http"
	"sezzle_api/src/repository"
	"github.com/go-redis/redis"
	"io"
	"time"
	"crypto/md5"
	"strconv"
)

type UserController struct{
	Repo	*repository.Repository
}

var UserModel = new(models.UserModel)

func (userController UserController) Login(c *gin.Context) {
	type Param struct {
		UserName	string 		`json:"user_name"`
		Password	string 		`json:"password"`
	}

	var pjson Param
	if err := c.ShouldBindJSON(&pjson); err == nil {

		if pjson.UserName == "" || pjson.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide user_name and password "})
			return
		}

		//check credentials
		user,err := userController.Repo.AuthByPassword(pjson.UserName,pjson.Password)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "401 errors"})
			return
		} 

		token := userController.GetRandomToken(pjson.UserName)

		redisClient := userController.Repo.Redis
		//check if user is already logged in
		val,err := redisClient.HGet("usermap_"+strconv.Itoa(int(user.Id)),"token").Result()
		if err != nil  && err != redis.Nil{
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Login failed! Please try again."})
			return
		}
		// if user is logged in
		if val != "" {
			//delete previous user token
			err  := redisClient.Del(val).Err()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Please try again."})
				return
			}
		}

		//set new user token
		err  = redisClient.HSet("usermap_"+strconv.Itoa(int(user.Id)),"token",token).Err()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Login failed! Please try again."})
			return
		}
		//set token user map
		err = redisClient.Set(token,user.Id,0).Err()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Login failed! Please try again."})
			return
		}

		//set expire for seesion/keys
		err = redisClient.Expire("usermap_"+strconv.Itoa(int(user.Id)),11*time.Minute).Err()
		if err != nil {
			fmt.Println("\n expire =",err)
		}
		err = redisClient.Expire(token,10*time.Minute).Err()
		if err != nil {
			fmt.Println("\n expire =",err)
		}


		c.JSON(http.StatusOK, gin.H{"message": "Singin successful, token :"+token})
		return
	} else{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid json"})
		return
	}

	return
	
}

func (userController UserController) Logout(c *gin.Context) {
	redisClient := userController.Repo.Redis
	token := c.GetHeader("token")
	val, err  := redisClient.Get(token).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	err  = redisClient.Del("usermap_"+val).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login failed! Please try againg."})
		return
	}	
	err = redisClient.Del(token).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please try again."})
		return
	}	
	c.JSON(http.StatusOK, gin.H{"message": "Logged out..."})
}

func (userController UserController) AddUser(c *gin.Context) {

	type Param struct {
		Name		string		`json:"name"`
		UserName	string 		`json:"user_name"`
		Password	string 		`json:"password"`
		ConfirmPassword	string	`json:"conf_password"`
	}

	var token string

	var pjson Param
	if err := c.ShouldBindJSON(&pjson); err == nil {

		if pjson.Name == "" || pjson.UserName == "" || pjson.Password == "" || pjson.ConfirmPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide values for all the parameters name,user_name,password,conf_password"})
			return
		}
		if pjson.Password != pjson.ConfirmPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords are not matching"})
			return
		}

		token := userController.GetRandomToken(pjson.UserName)

		err := UserModel.AddUser(pjson.Name,pjson.UserName,pjson.Password,token,userController.Repo)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while creating account"})
			return
		} 

	} else{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid json"})
		return
	}

	c.JSON(200, gin.H{"message": "Account created successufully. Your token is "+token})
}

func (userController UserController) ListUsers(c *gin.Context) {
	data,err:=  UserModel.ListUsers(userController.Repo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while listing users"})
		return
	}
	c.JSON(200, gin.H{"users": data})
}


func(userController UserController) GetRandomToken(userName string) string {
	
	now := time.Now()      // current local time
	sec := now.Unix()
	hash := md5.New()
	io.WriteString(hash, userName+"_"+strconv.Itoa(int(sec)))
	token := fmt.Sprintf("%x", hash.Sum(nil))

	return token
}

func NewUserController(repo *repository.Repository) *UserController {
	userController :=  &UserController{
		Repo : repo,
	}
	return userController
}