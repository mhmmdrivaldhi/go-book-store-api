package usecase

import (
	"context"
	"errors"
	"fmt"

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
	cartRepo    repository.CartRepository
	bookUsecase BookUsecase
}

func (cu *cartUsecase) AddToCart(ctx context.Context, userId int, item model.Item) (*model.Cart, error) {
	if item.Qty <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	book, err := cu.bookUsecase.GetById(item.BookId)
	if err != nil {
		return nil, fmt.Errorf("book not found for id %d", item.BookId)
	}

	cart, err := cu.cartRepo.GetCart(ctx, userId)
	if err != nil {
		return nil, errors.New("failed to retrived cart from repository")
	}

	itemFound := false
	for i := range cart.Items {
		if cart.Items[i].BookId == item.BookId {
			cart.Items[i].Qty += item.Qty
			cart.Items[i].Price = book.Price * cart.Items[i].Qty
			itemFound = true
			break
		}
	}

	if !itemFound {
		cart.Items = append(cart.Items, model.Item{
			BookId: item.BookId,
			Qty:    item.Qty,
			Price:  book.Price * item.Qty,
		})
	}

	cart.TotalPrice = helper.CalculateTotalPrice(cart)
	cart.TotalQty = helper.CalculateTotalQty(cart)

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
	if req.BookId == nil {
		return nil, errors.New("book id is required")
	}

	if req.Qty == nil {
		return nil, errors.New("qty is required")
	}

	if *req.Qty <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	cart, err := cu.cartRepo.GetCart(ctx, userId)
	if err != nil {
		return nil, errors.New("failed to get cart")
	}

	book, err := cu.bookUsecase.GetById(*req.BookId)
	if err != nil {
		return nil, errors.New("book not found")
	}

	itemFound := false
	for i := range cart.Items {
		if cart.Items[i].BookId == *req.BookId {
			cart.Items[i].Qty = *req.Qty
			cart.Items[i].Price = book.Price * (*req.Qty)
			itemFound = true
			break
		}
	}

	if !itemFound {
		return nil, errors.New("item not found in cart")
	}

	cart.TotalPrice = helper.CalculateTotalPrice(cart)
	cart.TotalQty = helper.CalculateTotalQty(cart)

	_, err = cu.cartRepo.SetCart(ctx, userId, cart)
	if err != nil {
		return nil, errors.New("failed to update qty from item")
	}

	return cart, nil
}

func (cu *cartUsecase) UpdateItemFromCart(ctx context.Context, userId int, req dto.RequestUpdateItemFromCart) (*model.Cart, error) {
	if req.BookId == nil {
		return nil, errors.New("book id is required")
	}

	if req.NewBookId == nil {
		return nil, errors.New("new book id is required")
	}

	if req.Qty == nil {
		return nil, errors.New("qty is required")
	}

	if *req.Qty <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	cart, err := cu.cartRepo.GetCart(ctx, userId)
	if err != nil {
		return nil, errors.New("failed to get cart")
	}

	newBook, err := cu.bookUsecase.GetById(*req.NewBookId)
	if err != nil {
		return nil, errors.New("book not found")
	}

	itemFound := false
	newItems := []model.Item{}

	for _, item := range cart.Items {
		if item.BookId == *req.BookId {
			newItems = append(newItems, model.Item{
				BookId: *req.NewBookId,
				Qty: *req.Qty,
				Price: newBook.Price * (*req.Qty),
			})
			itemFound = true
		} else {
			newItems = append(newItems, item)
		}
	}

	if !itemFound {
		return nil, errors.New("item not found in cart")
	}

	cart.Items = newItems
	cart.TotalPrice = helper.CalculateTotalPrice(cart)
	cart.TotalQty = helper.CalculateTotalQty(cart)

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
	cart.TotalQty = helper.CalculateTotalQty(cart)

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
		cartRepo:    cartRepo,
		bookUsecase: bookUsecase,
	}
}
