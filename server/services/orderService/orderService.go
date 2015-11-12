package orderService

import (
	vm "github.com/marpio/webapp2/server/viewmodels"
)

func AddProductOrder(order *vm.OrderDto, po vm.ProductOrderDto) {
	prods := append(order.ProductOrders, po)
	order.ProductOrders = prods
}

func RemoveProductOrder(order *vm.OrderDto, productID string) {
	var index = -1
	for i, v := range order.ProductOrders {
		if v.ProductID == productID {
			index = i
		}
	}
	if index >= 0 {
		order.ProductOrders = append(order.ProductOrders[:index], order.ProductOrders[index+1:]...)
	}
}

func GetTotalPrice(order *vm.OrderDto) float64 {
	if order == nil {
		return 0
	} else {
		var price float64
		for _, po := range order.ProductOrders {
			price = price + (po.ProductPrice * float64(po.Amount))
		}
		return price
	}
}
func GetItemsCount(order *vm.OrderDto) int {
	if order == nil {
		return 0
	} else {
		counter := 0

		if order.ServiceOrder != (vm.ServiceOrderDto{}) {
			counter = counter + 1
		}
		return counter + len(order.ProductOrders)
	}
}
