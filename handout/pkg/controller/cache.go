package controller

import (
	"context"
	"log"
	"time"

	"github.com/konstfish/og-peek/handout/pkg/redis"
	r "github.com/redis/go-redis/v9"
)

var client *redis.Client

func SetupCacheClient(address string) {
	var err error

	client, err = redis.NewClient(address)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func UrlLifecycle(ctx context.Context, url string, slug string) string {
	result := GetUrlFromCache(ctx, slug)
	log.Println(result)
	if result == "new" {
		SetUrlInCache(ctx, slug, "awaiting")
		SubmitUrlTask(ctx, url, slug)
	}

	// wait for url to be done
	for {
		result = GetUrlFromCache(ctx, slug)
		if result != "awaiting" {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	return result
}

func GetUrlFromCache(ctx context.Context, slug string) string {
	result, err := client.Get(ctx, slug)
	if err != nil {
		return "new"
	}

	return result
}

func SetUrlInCache(ctx context.Context, slug string, value string) {
	err := client.Set(ctx, slug, value)
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}
}

func SubmitUrlTask(ctx context.Context, url string, slug string) {
	client.Rdb.XAdd(ctx, &r.XAddArgs{
		Stream: "screenshot-urls",
		Values: map[string]interface{}{"url": url, "slug": slug},
	})
}
