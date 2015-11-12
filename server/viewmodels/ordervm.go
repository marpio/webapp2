package viewmodels

import (
	"github.com/marpio/webapp2/server/constants/orderstate"
	"time"
)

type OrderDto struct {
	ID            string
	UserID        string
	State         int
	ServiceOrder  ServiceOrderDto   `json:"serviceOrder"`
	ProductOrders []ProductOrderDto `json:"productOrders"`
}

type ProductOrderDto struct {
	ID           string
	ProductID    string  `json:"productID"`
	ProductName  string  `json:"productName"`
	ProductPrice float64 `json:"productPrice"`
	Amount       int     `json:"amount"`
}

type ServiceOrderDto struct {
	ID                  string
	ServiceProviderID   string
	ServiceProviderName string
	ServiceDate         time.Time
}

func NewOrderDto() *OrderDto {
	var prodOrders = make([]ProductOrderDto, 0)
	order := OrderDto{
		State:         orderstate.NEW,
		ProductOrders: prodOrders,
	}
	return &order
}
