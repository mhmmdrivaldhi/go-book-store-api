package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"github.com/redis/go-redis/v9"
)

type CartRepository interface {
	GetCart(ctx context.Context, userId int) (*model.Cart, error)
	SetCart(ctx context.Context, userId int, cart *model.Cart) (*model.Cart, error)
	ClearCart(ctx context.Context, userId int) error
}

type cartRepository struct {
	redis *redis.Client
}

func (cr *cartRepository) GetCart(ctx context.Context, userId int) (*model.Cart, error) {
	key := fmt.Sprintf("cart: %d", userId)

	data, err := cr.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return &model.Cart{
			UserId: userId,
			Items: []model.Item{},
			TotalQty: 0,
			TotalPrice: 0,
		}, nil
	} else if err != nil {
		return nil, err
	}

	var cart model.Cart
	err = json.Unmarshal([]byte(data), &cart)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (cr *cartRepository) SetCart(ctx context.Context, userId int, cart *model.Cart) (*model.Cart, error) {
	key := fmt.Sprintf("cart: %d", userId)

	data, err := json.Marshal(cart)
	if err != nil {
		return nil, err
	}

	err = cr.redis.Set(ctx, key, data, 0).Err()
	if err != nil {
		return nil, err
	}

	return cart, err
}

func (cr *cartRepository) ClearCart(ctx context.Context, userId int) error {
	key := fmt.Sprintf("cart: %d", userId)

	err := cr.redis.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func NewCartRepository(redis *redis.Client) *cartRepository {
	return &cartRepository{redis: redis}
}
