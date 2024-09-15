package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/og-peek/handout/pkg/cache"
	"github.com/konstfish/og-peek/handout/pkg/formatting"
	"github.com/konstfish/og-peek/handout/pkg/storage"
)

func GetUrl(c *gin.Context) {
	ctx := c.Request.Context()

	// set gin timeout to 5 seconds
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	// url from url parameters
	// TODO: make this middleware
	url := c.Query("url")
	if url == "" {
		c.String(http.StatusBadRequest, "missing url parameter")
		return
	}
	slug, err := formatting.UrlToSlug(url)
	domain, err := formatting.DomainFromUrl(url)
	if err != nil {
		c.String(http.StatusBadRequest, "issue parsing url")
		return
	}

	// minor tests here
	// - 404/other error
	// - whitelist

	stateChan := make(chan string, 1)
	stateChan <- cache.GetUrlStatus(ctx, domain, slug)

	for {
		select {
		case state := <-stateChan:
			log.Println(state)
			switch state {
			case cache.StatusNotFound:
				cache.SetUrlStatus(ctx, domain, slug, cache.StatusAwaiting)
				cache.SubmitUrlTask(ctx, domain, slug, url)
				// wait for avg. amount of time a screenshot usually takes
				time.Sleep(2000 * time.Millisecond)
				stateChan <- cache.GetUrlStatus(ctx, domain, slug)
			case cache.StatusAwaiting:
				time.Sleep(500 * time.Millisecond)
				stateChan <- cache.GetUrlStatus(ctx, domain, slug)
			case cache.StatusCached:
				// c.File(fmt.Sprintf("../screenshots/%s.png", slug))
				obj, err := storage.Download(ctx, domain, slug)
				if err != nil {
					c.String(http.StatusServiceUnavailable, "failed to retrieve screenshot from storage")
					// TODO: proper status here
					cache.SetUrlStatus(ctx, domain, slug, cache.StatusFailed)
					return
				}

				defer obj.Close()

				streamObject(c, obj)
				return
			case cache.StatusFailed:
				err, ttl := cache.GetUrlTTL(ctx, domain, slug)
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
}
