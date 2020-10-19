package models

import (
    "errors"
    _ "github.com/gin-gonic/gin"
     "sezzle_api/src/repository"
    "github.com/jinzhu/gorm"
    
)

type ItemModel struct{}

func (m ItemModel) AddItem(itemName  string, repository *repository.Repository) error {

    item, err := repository.GetItemByName(itemName)
    if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
        return err
    }
    if item != nil {
    	return errors.New("ItemAlreadyExists")
    }

    err =  repository.CreateItem(itemName)
    if err != nil {
        return err
    }
    return nil
}


func (m ItemModel) ListItems(repository *repository.Repository) (interface{},error) {
    items,err := repository.ListAllItems()
    if err != nil {
        return nil, err
    }

    return items,nil
}