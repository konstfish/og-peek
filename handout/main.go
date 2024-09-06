package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/og-peek/handout/pkg/config"
	"github.com/konstfish/og-peek/handout/pkg/controller"
	"github.com/konstfish/og-peek/handout/pkg/formatting"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/get", func(c *gin.Context) {
		ctx := c.Request.Context()

		// set gin timeout to 5 seconds
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// url from url parameters
		url := c.Query("url")
		if url == "" {
			c.String(http.StatusBadRequest, "missing url parameter")
			return
		}

		slug := formatting.UrlToSlug(url)

		// todo some sort of auth here

		// resultChan := make(chan string, 1)

		result := controller.UrlLifecycle(ctx, url, slug)

		c.String(http.StatusOK, result)
	})

	return r
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	controller.SetupCacheClient(cfg.RedisAddr)

	r := setupRouter()

	r.Run(":8080")
}
