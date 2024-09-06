package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dranikpg/gtrs"

	"github.com/konstfish/og-peek/capture/pkg/chrome"
	"github.com/konstfish/og-peek/capture/pkg/config"
	"github.com/konstfish/og-peek/capture/pkg/redis"
	"github.com/konstfish/og-peek/capture/pkg/screenshot"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chromeCancelFunc, err := chrome.Initialize(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Chrome: %v", err)
	}
	defer chromeCancelFunc()

	redisClient, err := redis.NewClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	consumerName, err := os.Hostname()
	if err != nil {
		consumerName = "capture"
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("goodbye!")
		cancel()
	}()

	cs := gtrs.NewGroupConsumer[redis.CaptureTask](ctx, redisClient.Rdb, "capture", consumerName, cfg.RedisStreamName, ">")

	log.Println("started capture service")

	for msg := range cs.Chan() {
		select {
		case <-ctx.Done():
			cs.Close()
			return
		default:
			taskCtx := context.Background()
			log.Println(msg)

			err := screenshot.Capture(taskCtx, msg.Data.Url)
			if err != nil {
				log.Printf("Failed to capture screenshot: %v", err)
				msg.Err = err
				log.Println(redisClient.Set(ctx, msg.Data.Slug, "failed"))
			} else {
				log.Println(redisClient.Set(ctx, msg.Data.Slug, "cached"))
			}

			cs.Ack(msg)
		}
	}
}
