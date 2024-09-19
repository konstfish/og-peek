package queue

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dranikpg/gtrs"
	"github.com/konstfish/og-peek/capture/pkg/config"
	"github.com/konstfish/og-peek/capture/pkg/redis"
	"github.com/konstfish/og-peek/capture/pkg/screenshot"
	"github.com/konstfish/og-peek/capture/pkg/storage"
)

var cs *gtrs.GroupConsumer[CaptureTask]
var sigChan chan os.Signal

// TODO: error here
func NewClient(ctx context.Context, cfg *config.Config) {
	cs = gtrs.NewGroupConsumer[CaptureTask](
		ctx,
		redis.Rdb,
		cfg.RedisConsumerGroupName,
		cfg.RedisConsumerName,
		cfg.RedisStreamName,
		">",
	)
}

func Close() {
	cs.Close()
}

func SignalHandler(cancel context.CancelFunc) {
	sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("goodbye!")
		cancel()
	}()
}

func Listen(ctx context.Context) {
	log.Println("starting consumer")

	for msg := range cs.Chan() {
		select {
		case <-ctx.Done():
			cs.Close()
			return
		default:
			taskCtx := context.Background()
			log.Println(msg.Data)

			content, err := screenshot.Capture(taskCtx, msg.Data.Url)
			if err != nil {
				log.Printf("Failed to capture screenshot: %v", err)
				msg.Err = err

				redis.SetUrlStatus(ctx, msg.Data.Domain, msg.Data.Slug, "failed")

				continue
			}
			// upload to s3
			err = storage.Upload(taskCtx, content, msg.Data.Domain, msg.Data.Slug)
			if err != nil {
				log.Printf("Failed to upload to S3: %v", err)
				redis.SetUrlStatus(ctx, msg.Data.Domain, msg.Data.Slug, "failed")

				continue
			}

			redis.SetUrlStatus(ctx, msg.Data.Domain, msg.Data.Slug, "cached")

			cs.Ack(msg)
		}
	}
}
