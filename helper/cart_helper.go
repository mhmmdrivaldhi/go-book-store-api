package helper

import "github.com/mhmmmdrivaldhi/go-book-api/model"

func CalculateTotalPrice(cart *model.Cart) int {
	price := 0
	for _, item := range cart.Items {
		price += item.Price * item.Qty
	}

	return price
}