package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/usecase"
)

type cartController struct {
	cartUsecase usecase.CartUsecase
}

func (cc *cartController) AddItem(ctx *gin.Context) {
	userId := ctx.GetInt("user_id")

	var req model.Item
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	cart, err := cc.cartUsecase.AddToCart(ctx, userId, req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.GeneralResponse{
		Message: "item added to cart",
		Data: cart,
	})
}

func (cc *cartController) GetCart(ctx *gin.Context) {
	userId := ctx.GetInt("user_id")

	cart, err := cc.cartUsecase.GetCartFromUser(ctx, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully get cart from user",
		Data: cart,
	})
}

func (cc *cartController) UpdateQty(ctx *gin.Context) {
	userId := ctx.GetInt("user_id")

	bookIdParam := ctx.Param("book_id")
	bookId, err := strconv.Atoi(bookIdParam)
	if err != nil { 
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid book id"})
		return
	}

	var book struct {
		Qty int `json:"qty"`
	}

	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	update := dto.RequestUpdateQtyFromItem{
		BookId: &bookId,
		Qty: &book.Qty,
	}

	cart, err := cc.cartUsecase.UpdateQtyFromItem(ctx, userId, update)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(200, dto.GeneralResponse{
		Message: "successfully update qty from item",
		Data: cart,
	})
}

func (cc *cartController) UpdateItem(ctx *gin.Context) {
	userId := ctx.GetInt("user_id")
	
	bookIdParam := ctx.Param("book_id")
	bookId, err := strconv.Atoi(bookIdParam)
	if err != nil { 
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid book id"})
		return
	}

	var update dto.RequestUpdateItemFromCart
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	update.BookId = &bookId

	cart, err := cc.cartUsecase.UpdateItemFromCart(ctx, userId, update)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(200, dto.GeneralResponse{
		Message: "successfully updated item",
		Data: cart,
	})
}

func (cc *cartController) RemoveItem(ctx *gin.Context) {
	userId := ctx.GetInt("user_id")
	bookId := ctx.GetInt("book_id")

	item, err := cc.cartUsecase.RemoveItemFromCart(ctx, userId, bookId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(200, dto.GeneralResponse{
		Message: "item removed",
		Data: item,
	})
}

func (cc *cartController) ClearAllItems(ctx *gin.Context) {
	userId := ctx.GetInt("user_id")

	if err := cc.cartUsecase.ClearAllItemFromCart(ctx, userId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(200, dto.GeneralResponse{
		Message: "successfully clear all items",
	})
}

func NewCartController(cartUsecase usecase.CartUsecase, rg *gin.RouterGroup) *cartController {
	controller := &cartController{cartUsecase: cartUsecase}

	rg.POST("/cart", controller.AddItem)
	rg.GET("/cart", controller.GetCart)
	rg.PUT("/cart/qty/:book_id", controller.UpdateQty)
	rg.PUT("/cart/item/:book_id", controller.UpdateItem)
	rg.DELETE("/cart/:book_id", controller.RemoveItem)
	rg.DELETE("/cart", controller.ClearAllItems)

	return controller
}