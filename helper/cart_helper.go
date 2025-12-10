package helper

import "github.com/mhmmmdrivaldhi/go-book-api/model"

func CalculateTotalPrice(cart *model.Cart) int {
	price := 0
	for _, item := range cart.Items {
		price += item.Price * item.Qty
	}

	return price
}

func CalculateTotalQty(cart *model.Cart) int {
	qty := 0
	for _, item := range cart.Items {
		qty += item.Qty
	}

	return qty
}