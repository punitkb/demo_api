package models

import (
    
    "sezzle_api/src/repository"
)

type OrderModel struct{}

func (m OrderModel) ListOrders(repository *repository.Repository) (interface{},error) {
  carts,err := repository.ListAllOrders()
    if err != nil {
        return nil, err
    }
    return carts,nil
}
