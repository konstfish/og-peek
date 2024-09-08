package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func SetupCacheClient(address string) error {
	rdb = redis.NewClient(&redis.Options{
		Addr: address,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}
