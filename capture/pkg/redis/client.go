package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func NewClient(addr string) error {
	Rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := Rdb.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}

func SetUrlStatus(ctx context.Context, domain string, slug string, value string) {
	key := fmt.Sprintf("domain:%s:%s", domain, slug)
	err := Rdb.Set(ctx, key, value, 7*24*time.Hour).Err()
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}
}
