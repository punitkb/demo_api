package models

import (
    "errors"
    _ "github.com/gin-gonic/gin"
    "sezzle_api/src/repository"
    "github.com/jinzhu/gorm"
)

type UserModel struct{}

func (m UserModel) AddUser(name string, username string, passwd string, token string , repository *repository.Repository) error {

    // check parser name exist
    user, err := repository.GetUserByName(username)
    if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)){
        return err
    }
    
    if user != nil {
            return  errors.New("ErrUserNameAlreadyExits")
    }

    //create cart for user
    cart_id,err := repository.CreateCart()
    if err != nil {
        return err
    }

    //create user
    err =  repository.CreateUser(name, username, passwd, token, cart_id)
    if err != nil {
            return err
    }

    return nil
}


func (m UserModel) ListUsers(repository *repository.Repository) (interface{},error) {
    users,err := repository.ListAllUsers()
    if err != nil {
        return nil, err
    }
    return users,nil
}   