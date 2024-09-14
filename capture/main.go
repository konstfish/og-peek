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
	"github.com/konstfish/og-peek/capture/pkg/storage"
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

	// setup storage
	s3Client, err := storage.NewClient(cfg.S3Endpoint, cfg.S3BucketName, cfg.S3AccessKeyId, cfg.S3AccessKey, true)
	if err != nil {
		log.Fatalf("Failed to set up S3 Client: %v", err)
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

			content, err := screenshot.Capture(taskCtx, msg.Data.Url)
			if err != nil {
				log.Printf("Failed to capture screenshot: %v", err)
				msg.Err = err

				err := redisClient.Set(ctx, msg.Data.Slug, "failed")
				if err != nil {
					log.Printf("Failed to set failed status: %v", err)
				}
			} else {
				// upload to s3
				err := storage.Upload(taskCtx, s3Client, content, msg.Data.Slug)
				if err != nil {
					log.Printf("Failed to upload to S3: %v", err)
				}

				err = redisClient.Set(ctx, msg.Data.Slug, "cached")
				if err != nil {
					log.Printf("Failed to set failed status: %v", err)
				}
			}

			cs.Ack(msg)
		}
	}
}
