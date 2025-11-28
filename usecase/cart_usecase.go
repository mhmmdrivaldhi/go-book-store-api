package usecase

import (
	"context"
	"errors"

	"github.com/mhmmmdrivaldhi/go-book-api/helper"
	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/repository"
)

type CartUsecase interface {
	AddToCart(ctx context.Context, userId int, item model.Item) (*model.Cart, error)
	GetCartFromUser(ctx context.Context, userId int) (*model.Cart, error)
	UpdateQtyFromItem(ctx context.Context, userId int, req dto.RequestUpdateQtyFromItem) (*model.Cart, error)
	UpdateItemFromCart(ctx context.Context, userId int, req dto.RequestUpdateItemFromCart) (*model.Cart, error)
	RemoveItemFromCart(ctx context.Context, userId, bookId int) (*model.Cart, error)
	ClearAllItemFromCart(ctx context.Context, userId int) error
}

type cartUsecase struct {
	cartRepo repository.CartRepository
	bookUsecase BookUsecase 
}

func (cu *cartUsecase) AddToCart(ctx context.Context, userId int, item model.Item) (*model.Cart, error) {
	qty := 0
	if qty <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	book, err := cu.bookUsecase.GetById(item.BookId)
	if err != nil {
		return nil, errors.New("book not found")
	}

	cart, err := cu.cartRepo.GetCart(ctx, userId)
	if err != nil {
		return nil, err
	}

	for i := range cart.Items {
		if cart.Items[i].BookId == item.BookId {
			cart.Items[i].Qty += item.Qty
			cart.Items[i].Price = book.Price
			break
		}
	}

	cart.TotalPrice = helper.CalculateTotalPrice(cart)

	_, err = cu.cartRepo.SetCart(ctx, userId, cart)
	if err != nil {
		return nil, errors.New("failed to added item to cart")
	}

	return cart, nil
}

func (cu *cartUsecase) GetCartFromUser(ctx context.Context, userId int) (*model.Cart, error) {
	cart, err := cu.cartRepo.GetCart(ctx, userId)
	if err != nil {
		return nil, errors.New("failed to get cart")
	}

	return cart, nil
}

func (cu *cartUsecase) UpdateQtyFromItem(ctx context.Context, userId int, req dto.RequestUpdateQtyFromItem) (*model.Cart, error) {
	if req.Qty != nil && *req.Qty <= 0{
		return nil, errors.New("quantity must be greater than 0")
	}

	cart, err := cu.cartRepo.GetCart(ctx, userId)
	if err != nil {
		return nil, errors.New("failed to get cart")
	}

	for i := range cart.Items {
		if cart.Items[i].BookId == *req.BookId {
			cart.Items[i].Qty = *req.Qty
			break
		}
	}

	cart.TotalPrice = helper.CalculateTotalPrice(cart)

	_, err = cu.cartRepo.SetCart(ctx, userId, cart)
	if err != nil {
		return nil, errors.New("failed to update qty from item")
	}

	return cart, nil
}

func (cu *cartUsecase) UpdateItemFromCart(ctx context.Context, userId int, req dto.RequestUpdateItemFromCart) (*model.Cart, error) {
	cart, err := cu.cartRepo.GetCart(ctx, userId)
	if err != nil {
		return nil, errors.New("failed to get cart")
	}

	book, err := cu.bookUsecase.GetById(*req.BookId)
	if err != nil {
		return nil, errors.New("book not found")
	}

	for i := range cart.Items {
		if cart.Items[i].BookId == *req.BookId {
			cart.Items[i] = model.Item{
				BookId: book.Id,
				Price: book.Price * (*req.Qty),
				Qty: *req.Qty,
			}
			break
		}
	}

	cart.TotalPrice = helper.CalculateTotalPrice(cart)
	_, err = cu.cartRepo.SetCart(ctx, userId, cart)
	if err != nil {
		return nil, errors.New("failed to update item from cart")
	}

	return cart, nil
}

func (cu *cartUsecase) RemoveItemFromCart(ctx context.Context, userId, bookId int) (*model.Cart, error) {
	cart, err := cu.cartRepo.GetCart(ctx, userId)
	if err != nil {
		return nil, errors.New("failed to get cart")
	} 

	items := []model.Item{}
	for _, item := range cart.Items {
		if item.BookId != bookId {
			items = append(items, item)
		}
	}

	cart.Items = items
	cart.TotalPrice = helper.CalculateTotalPrice(cart)

	_, err = cu.cartRepo.SetCart(ctx, userId, cart)
	if err != nil {
		return nil, errors.New("failed to remove item from cart")
	}
	return cart, nil
}

func (cu *cartUsecase) ClearAllItemFromCart(ctx context.Context, userId int) error {
	err := cu.cartRepo.ClearCart(ctx, userId)
	if err != nil {
		return errors.New("failed to clear all item from cart")
	}

	return nil
}

func NewCartUsecase(cartRepo repository.CartRepository, bookUsecase BookUsecase) CartUsecase {
	return &cartUsecase{
		cartRepo: cartRepo,
		bookUsecase: bookUsecase,
	}
}

	
