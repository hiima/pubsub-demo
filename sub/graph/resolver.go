package graph

import (
	"context"
	"sub/graph/model"
	"sync"

	"github.com/go-redis/redis/v8"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	redisClient *redis.Client
	observers   map[string]map[string]chan<- *model.Notification // 接続中のクライアントの状態を格納
	mutex       sync.Mutex
}

func NewResolver(ctx context.Context, redisURL string) *Resolver {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic("failed to ping redis server")
	}

	return &Resolver{
		redisClient: redisClient,
		observers:   map[string]map[string]chan<- *model.Notification{},
		mutex:       sync.Mutex{},
	}
}
