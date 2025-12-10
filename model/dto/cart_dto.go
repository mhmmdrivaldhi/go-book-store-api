package dto

type RequestUpdateQtyFromItem struct {
	BookId *int   `json:"book_id" binding:"omitempty"`
	Qty  *int   `json:"qty" binding:"omitempty"`
}

type RequestUpdateItemFromCart struct {
	BookId *int   `json:"book_id" binding:"omitempty"`
	NewBookId *int `json:"new_book_id" binding:"omitempty"`
	Qty *int   `json:"qty" binding:"omitempty"`
}
