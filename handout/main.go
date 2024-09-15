package main

import (
	"log"

	"github.com/konstfish/og-peek/handout/pkg/cache"
	"github.com/konstfish/og-peek/handout/pkg/config"
	"github.com/konstfish/og-peek/handout/pkg/mappings"
	"github.com/konstfish/og-peek/handout/pkg/storage"
)

func main() {
	log.Println("loading config")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Println("setting up redis")
	err = cache.SetupCacheClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to setup cache client: %v", err)
	}

	log.Println("setting up s3 client")
	err = storage.NewClient(cfg.S3Endpoint, cfg.S3BucketName, cfg.S3AccessKeyId, cfg.S3AccessKey, true)
	if err != nil {
		log.Fatalf("Failed to set up S3 Client: %v", err)
	}

	log.Println("setting up gin router")
	mappings.CreateUrlMappings()
	mappings.Router.Run(":8080")
}
