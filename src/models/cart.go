package models

import (
    _ "github.com/gin-gonic/gin"
    "sezzle_api/src/repository"
    "errors"
    "github.com/jinzhu/gorm"
    "strconv"
    
)


type CartModel struct{}

func (m CartModel) AddToCart(itemName, userid string, repository *repository.Repository) error {

    //if exist add to cart
    item, err := repository.GetItemByName(itemName)
    if err != nil && (errors.Is(err, gorm.ErrRecordNotFound)) {
        return errors.New("ErrItemNotExist")
    }
    itemId := item.Id

    user_id,err := strconv.Atoi(userid)
    if err != nil {
        return err
    }

    user, err := repository.GetUserById(int(user_id))
    if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)){
        return errors.New("ErrUserNotFound")
    }
    cartId := user.CartId

    
    err = repository.AddItemToCart(int(cartId),int(itemId))
    if err != nil {
        return err
    }
    
    //mark cart not purchased 
    isPurchased := false
    err = repository.UpdateCart(isPurchased,int(cartId))
    if err != nil {
        return err
    }

    return nil
}

func (m CartModel) CompleteOrder(cartId  int,userid string, repository *repository.Repository) error {

    user_id,err := strconv.Atoi(userid)
    if err != nil {
        return err
    }

    user, err := repository.GetUserById(int(user_id))
    if err != nil {
            return err
    }
    if user == nil {
            return errors.New("ErrInvalidToken")
    }

    if user.CartId != uint64(cartId) {
         return  errors.New("ErrCartIdNotBelongToUser")
    }

    // add to orders
    err = repository.CreateOrder(cartId,int(user.Id))
    if err != nil {
        return err
    }
    
    //mark cart purchased 
    isPurchased := true
    err = repository.UpdateCart(isPurchased,cartId)
    if err != nil {
        return err
    }

    // mark purchased in cart order relations
    err = repository.UpdateCartItemRelation(cartId,isPurchased)
    if err != nil {
        return err
    }

    return nil
}


func (m CartModel) ListCarts(repository *repository.Repository) (interface{},error) {
    carts,err := repository.ListAllCarts()
    if err != nil {
        return nil, err
    }
    return carts,nil
}
