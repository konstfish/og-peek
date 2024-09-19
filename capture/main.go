package main

import (
	"context"
	"log"

	"github.com/konstfish/og-peek/capture/pkg/chrome"
	"github.com/konstfish/og-peek/capture/pkg/config"
	"github.com/konstfish/og-peek/capture/pkg/queue"
	"github.com/konstfish/og-peek/capture/pkg/redis"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("loading config")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Println("setting up chrome")
	chromeCancelFunc, err := chrome.Initialize(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Chrome: %v", err)
	}
	defer chromeCancelFunc()

	log.Println("setting up redis")
	err = redis.NewClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Rdb.Close()

	/*log.Println("setting up storage")
	err = storage.NewClient(cfg.S3Endpoint, cfg.S3BucketName, cfg.S3AccessKeyId, cfg.S3AccessKey, true)
	if err != nil {
		log.Fatalf("Failed to set up S3 client: %v", err)
	}*/

	log.Println("setting up message queue")
	queue.SignalHandler(cancel)
	queue.NewClient(ctx, cfg)
	defer queue.Close()

	queue.Listen(ctx)
}
