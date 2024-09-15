package cache

import (
	"context"
	"fmt"
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

func GetUrlStatus(ctx context.Context, domain string, slug string) string {
	key := fmt.Sprintf("domain:%s:%s", domain, slug)
	result, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return StatusNotFound
	}

	return result
}

func GetUrlTTL(ctx context.Context, domain string, slug string) (error, int) {
	key := fmt.Sprintf("domain:%s:%s", domain, slug)
	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		return err, 0
	}

	return nil, int(ttl.Hours())
}

func SetUrlStatus(ctx context.Context, domain string, slug string, value string) {
	key := fmt.Sprintf("domain:%s:%s", domain, slug)
	err := rdb.Set(ctx, key, value, 96*time.Hour).Err()
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}
}

func SubmitUrlTask(ctx context.Context, domain string, slug string, url string) {
	rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "screenshot-urls",
		Values: map[string]interface{}{"url": url, "slug": slug, "domain": domain},
	})
}
