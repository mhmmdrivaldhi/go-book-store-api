package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/usecase"
)

type orderController struct {
	orderUsecase usecase.OrderUsecase
	cartUsecase usecase.CartUsecase
}

func (oc *orderController) OrderFromCart(ctx *gin.Context) { 
	userId := ctx.GetInt("user_id")

	cart, err := oc.cartUsecase.GetCartFromUser(ctx, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "failed to get cart from user"})
		return
	}

	order, err := oc.orderUsecase.OrderFromCart(ctx, userId, cart)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "failed to order from cart"})
		return
	}

	ctx.JSON(http.StatusCreated, dto.GeneralResponse{
		Message: "successfully to order book",
		Data: order,
	})
}

func (oc *orderController) GetAllOrders(ctx *gin.Context) {
	orders, err := oc.orderUsecase.GetAllOrders()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully get all orders",
		Data: orders,
	})
}

func (oc *orderController) GetOrderById(ctx *gin.Context) {
	orderId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid order id"})
		return
	}

	order, err := oc.orderUsecase.GetOrderById(orderId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get order by id"})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully get order by id",
		Data: order,
	})
}

func (oc *orderController) UpdatePaymentStatus(ctx *gin.Context) {
	orderId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid order id"})
		return
	}

	var request dto.UpdatePaymentStatusRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	update, err := oc.orderUsecase.UpdatePaymentStatus(orderId, request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully update payment status",
		Data: update,
	})
}

func (oc *orderController) CancelOrder(ctx *gin.Context) {
	orderId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid order id"})
		return
	}

	err = oc.orderUsecase.CancelOrder(orderId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully cancel order",
	})
}

func NewOrderController(orderUsecase usecase.OrderUsecase, cartUsecase usecase.CartUsecase, rg *gin.RouterGroup) *orderController {
	return &orderController{
		orderUsecase: orderUsecase,
		cartUsecase: cartUsecase,
	}

	controller := &orderController{
		orderUsecase: orderUsecase,
		cartUsecase: cartUsecase,
	}

	rg.POST("/order", controller.OrderFromCart)
	rg.GET("/order", controller.GetAllOrders)
	rg.GET("/order", controller.GetOrderById)
	rg.PUT("/order/:id/payment", controller.UpdatePaymentStatus)
	// rg.PUT("order/:id/status", controller.UpdateOrderStatus)
	rg.DELETE("/order/:id", controller.CancelOrder)

	return controller
}