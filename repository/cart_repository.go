package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"github.com/redis/go-redis/v9"
)

type CartRepository interface {
	GetCart(userId int) (*model.Cart, error)
	SetCart(userId int, cart *model.Cart) (*model.Cart, error)
	ClearCart(userId int) error
}

type cartRepository struct {
	redis *redis.Client
}

func (cr *cartRepository) GetCart(userId int) (*model.Cart, error) {
	ctx := context.Background()
	key := fmt.Sprintf("cart: %d", userId)

	data, err := cr.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return &model.Cart{
			UserId: userId,
			Items: []model.Item{},
		}, nil
	} else if err != nil {
		return nil, err
	}

	var cart model.Cart
	json.Unmarshal([]byte(data), &cart)

	return &cart, nil
}

func (cr *cartRepository) SetCart(userId int, cart *model.Cart) (*model.Cart, error) {
	ctx := context.Background()
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

func (cr *cartRepository) ClearCart(userId int) error {
	ctx := context.Background()
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
