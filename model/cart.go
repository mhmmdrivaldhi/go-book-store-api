package model

type Item struct {
	BookId int `json:"book_id"`
	Price  int `json:"price"`
	Qty	   int `json:"qty"`
}

type Cart struct {
	UserId int `json:"user_id"`
	Items []Item `json:"items"`
	TotalQty int `json:"total_qty"`
	TotalPrice int `json:"total_price"`
}