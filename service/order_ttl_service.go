package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type OrderTTLService interface {
	SetOrderTTL(ctx context.Context, orderId int, duration time.Duration) error 
	DeleteOrderTTL(ctx context.Context, orderId int) error
}

type orderTTLService struct {
	redis *redis.Client
}

func (ots *orderTTLService) SetOrderTTL(ctx context.Context, orderId int, duration time.Duration) error {
	key := fmt.Sprintf("order:ttl:%d", orderId)
	return ots.redis.Set(ctx, key, orderId, duration).Err()
}

func (ots *orderTTLService) DeleteOrderTTL(ctx context.Context, orderId int) error {
	key := fmt.Sprintf("order:ttl:%d", orderId)
	return ots.redis.Del(ctx, key).Err()
}