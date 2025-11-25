package dto

import "github.com/mhmmmdrivaldhi/go-book-api/model"

type RequestUpdateCartByMode struct {
	Qty  *int   `json:"qty" binding:"omitempty"`
	Mode string `json:"mode" binding:"omitempty"`
}

type RequestReplaceAllCart struct {
	Items []model.Item `json:"items" binding:"omitempty"`
}
