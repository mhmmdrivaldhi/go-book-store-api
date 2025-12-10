package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)


func NewRedisClient(cfg *Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: cfg.RedisPassword,
		DB: 0,
		DialTimeout: 3 * time.Second,
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize: 10,
		MinIdleConns: 3,
	})
	
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("failed connecting to Redis at %s: %v", addr, err))
	}

	fmt.Printf("successfully connected to Redis at %s (DB %d)\n", addr, cfg.RedisDB)
	return rdb
}
