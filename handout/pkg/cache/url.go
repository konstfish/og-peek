package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	StatusNotFound = "not_found"
	StatusAwaiting = "awaiting"
	StatusCached   = "cached"
	StatusFailed   = "failed"
)

func GetUrlStatus(ctx context.Context, slug string) string {
	result, err := rdb.Get(ctx, slug).Result()
	if err != nil {
		return StatusNotFound
	}

	return result
}

func GetUrlTTL(ctx context.Context, slug string) (error, int) {
	ttl, err := rdb.TTL(ctx, slug).Result()
	if err != nil {
		return err, 0
	}

	return nil, int(ttl.Hours())
}

func SetUrlStatus(ctx context.Context, slug string, value string) {
	err := rdb.Set(ctx, slug, value, 96*time.Hour).Err()
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}
}

func SubmitUrlTask(ctx context.Context, url string, slug string) {
	rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "screenshot-urls",
		Values: map[string]interface{}{"url": url, "slug": slug},
	})
}
