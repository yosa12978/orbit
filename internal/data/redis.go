package data

import (
	"context"
	"orbit-app/internal/config"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	rdb       *redis.Client
	initRedis sync.Once
)

func Redis(ctx context.Context) *redis.Client {
	initRedis.Do(func() {
		conf := config.Get()
		client := redis.NewClient(&redis.Options{
			Addr:     conf.Redis.Addr,
			Password: conf.Redis.Password,
			DB:       conf.Redis.DB,
		})
		if err := client.Ping(ctx).Err(); err != nil {
			panic(err)
		}
		rdb = client
	})
	return rdb
}
