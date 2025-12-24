package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/repository"
	"github.com/mhmmmdrivaldhi/go-book-api/service"
)

type OrderUsecase interface {
	OrderFromCart(ctx context.Context, userId int, cart *model.Cart) (*model.Order, error)
	GetAllOrders() ([]dto.OrderResponse, error)
	GetOrderById(orderId int) (*dto.OrderResponse, error)
	UpdatePaymentStatus(orderId int, req dto.UpdatePaymentStatusRequest) (*model.Order, error)
	UpdateOrderStatus(orderId int, req dto.UpdateOrderStatusRequest) (*model.Order, error)
	CancelOrder(orderId int) error
}

type orderUsecase struct {
	orderRepo repository.OrderRepository
	ttlService service.OrderTTLService
}

const (
	PaymentUnpaid = "Unpaid"
	PaymentPaid = "Paid"
	OrderCreated = "Created"
	OrderProccesing = "Processing"
	OrderShipped = "Shipping"
	OrderDelivered = "Delivered"
	OrderCanceled = "Canceled"
	OrderExpired = "Expired"
)

func (ou *orderUsecase) OrderFromCart(ctx context.Context, userId int, cart *model.Cart) (*model.Order, error) {
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	order := &model.Order{
		UserID: userId,
		PaymentStatus: PaymentUnpaid,
		OrderStatus: OrderCreated,
		TotalAmount: cart.TotalPrice,
		ExpiredAt: time.Now().Add(time.Duration(60 * time.Minute)),
	}

	for _, item := range cart.Items {
		order.OrderItems = append(order.OrderItems, model.OrderItem{
			BookID: item.BookId,
			Qty: item.Qty,
			Price: item.Price,
		})
	}

	create, err := ou.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, errors.New("failed to create order")
	}

	err = ou.ttlService.SetOrderTTL(ctx, create.ID, time.Minute * 60)
	if err != nil {
		return nil, errors.New("failed to set order ttl")
	}

	return create, nil
}

func (ou *orderUsecase) GetAllOrders() ([]dto.OrderResponse, error) {
	var resp []dto.OrderResponse

	orders, err := ou.orderRepo.FindAll()
	if err != nil {
		return nil, errors.New("failed to get all orders")
	}

	for _, order := range orders {
		var userResponse *dto.UserResponse
		if order.User.ID != 0 {
			userResponse = &dto.UserResponse{
				ID: order.User.ID,
				Name: order.User.Name,
				Email: order.User.Email,
				Address: order.User.Address,
				BankName: order.User.BankName,
				BankNumber: order.User.BankNumber,
			}
		}

		orderResponse := dto.OrderResponse{
			ID: order.ID,
			UserId: order.UserID,
			User: *userResponse,
			TotalQty: order.TotalQty,
			TotalAmount: order.TotalAmount,
			OrderStatus: order.OrderStatus,
			PaymentStatus: order.PaymentStatus,
			ShippingStatus: order.ShippingStatus,
			ExpiredAt: order.ExpiredAt,
		}
		resp = append(resp, orderResponse)
	}
	return resp, nil
}

func (ou *orderUsecase) GetOrderById(orderId int) (*dto.OrderResponse, error) {
	order, err := ou.orderRepo.FindById(orderId)
	if err != nil {
		return nil, errors.New("order not found")
	}

	var userResponse *dto.UserResponse
	if order.User.ID != 0 {
		userResponse = &dto.UserResponse{
			ID: order.User.ID,
				Name: order.User.Name,
				Email: order.User.Email,
				Address: order.User.Address,
				BankName: order.User.BankName,
				BankNumber: order.User.BankNumber,
		}
	}
	orderResponse := &dto.OrderResponse{
		ID: order.ID,
		UserId: order.UserID,
		User: *userResponse,
		TotalQty: order.TotalQty,
		TotalAmount: order.TotalAmount,
		OrderStatus: order.OrderStatus,
		PaymentStatus: order.PaymentStatus,
		ShippingStatus: order.ShippingStatus,
		ExpiredAt: order.ExpiredAt,
	}

	if time.Now().After(order.ExpiredAt) && order.PaymentStatus == PaymentUnpaid {
		orderResponse.OrderStatus = OrderExpired
	}

	return orderResponse, nil
}

func (ou *orderUsecase) UpdatePaymentStatus(orderId int, req dto.UpdatePaymentStatusRequest) (*model.Order, error) {
	var ctx context.Context

	if req.PaymentStatus != PaymentPaid {
		return nil, errors.New("invalid payment status")
	}

	updatePayment := &model.Order{
		PaymentStatus: PaymentPaid,
		OrderStatus: OrderProccesing,
	}

	_, err := ou.orderRepo.UpdateOrder(orderId, updatePayment)
	if err != nil {
		return nil, errors.New("failed to update payment status")
	}

	err = ou.ttlService.DeleteOrderTTL(ctx, orderId)

	return updatePayment, nil
}

func (ou *orderUsecase) UpdateOrderStatus(orderId int, req dto.UpdateOrderStatusRequest) (*model.Order, error) {
	paidOrder := map[string]bool{
		OrderProccesing: true,
		OrderShipped: true,
		OrderDelivered: true,
	}

	if !paidOrder[req.OrderStatus] {
		return nil, errors.New("invalid order status")
	}

	updateOrder := &model.Order{
		OrderStatus: req.OrderStatus,
	}

	_, err := ou.orderRepo.UpdateOrder(orderId, updateOrder)
	if err != nil {
		return nil, errors.New("failed to update order status")
	}
	return updateOrder, nil
}

func (ou *orderUsecase) CancelOrder(orderId int) error {
	order, err := ou.orderRepo.FindById(orderId)
	if err != nil {
		return errors.New("order not found")
	}

	if order.PaymentStatus == PaymentPaid {
		return errors.New("cannot cancel paid order")
	}

	err = ou.orderRepo.DeleteOrder(orderId)
	if err != nil {
		return errors.New("failed to cancel order")
	}

	return nil
}

func NewOrderUsecase(orderRepo repository.OrderRepository) OrderUsecase {
	return &orderUsecase{
		orderRepo: orderRepo,
	}
}