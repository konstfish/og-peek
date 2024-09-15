package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/og-peek/handout/pkg/cache"
	"github.com/konstfish/og-peek/handout/pkg/config"
	"github.com/konstfish/og-peek/handout/pkg/controllers"
	"github.com/konstfish/og-peek/handout/pkg/formatting"
	"github.com/konstfish/og-peek/handout/pkg/storage"
)

var s3Client storage.S3Client

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(controllers.Cors())

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/get", func(c *gin.Context) {
		ctx := c.Request.Context()

		// set gin timeout to 5 seconds
		ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
		defer cancel()

		// url from url parameters
		url := c.Query("url")
		if url == "" {
			c.String(http.StatusBadRequest, "missing url parameter")
			return
		}
		slug := formatting.UrlToSlug(url)

		// minor tests here
		// - 404/other error
		// - whitelist

		stateChan := make(chan string, 1)
		stateChan <- cache.GetUrlStatus(ctx, slug)

		for {
			select {
			case state := <-stateChan:
				log.Println(state)
				switch state {
				case cache.StatusNotFound:
					cache.SetUrlStatus(ctx, slug, cache.StatusAwaiting)
					cache.SubmitUrlTask(ctx, url, slug)
					// wait for avg. amount of time a screenshot usually takes
					time.Sleep(2000 * time.Millisecond)
					stateChan <- cache.GetUrlStatus(ctx, slug)
				case cache.StatusAwaiting:
					time.Sleep(500 * time.Millisecond)
					stateChan <- cache.GetUrlStatus(ctx, slug)
				case cache.StatusCached:
					// c.File(fmt.Sprintf("../screenshots/%s.png", slug))
					obj, err := storage.Download(ctx, s3Client, slug)
					if err != nil {
						c.String(http.StatusServiceUnavailable, "failed to retrieve screenshot from storage")
						cache.SetUrlStatus(ctx, slug, cache.StatusFailed)
						return
					}

					defer obj.Close()

					controllers.StreamObject(c, obj)
					return
				case cache.StatusFailed:
					err, ttl := cache.GetUrlTTL(ctx, slug)
					if err == nil {
						c.Header("Next-Attempt-In", fmt.Sprintf("%dh", ttl))
					}
					c.String(http.StatusInternalServerError, "failed to capture")
					return
				}
			case <-ctx.Done():
				c.Header("Retry-After", "1")
				c.String(http.StatusServiceUnavailable, "failed to generate screenshot within timeout, retrying...")
				return
			}
		}
	})

	return r
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	err = cache.SetupCacheClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to setup cache client: %v", err)
	}

	// setup storage
	s3Client, err = storage.NewClient(cfg.S3Endpoint, cfg.S3BucketName, cfg.S3AccessKeyId, cfg.S3AccessKey, true)
	if err != nil {
		log.Fatalf("Failed to set up S3 Client: %v", err)
	}

	r := setupRouter()

	r.Run(":8080")
}
