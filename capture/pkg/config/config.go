package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RedisAddr              string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	RedisStreamName        string `envconfig:"REDIS_STREAM_NAME" default:"screenshot-urls"`
	RedisConsumerGroupName string `envconfig:"REDIS_CONSUMER_GROUP_NAME" default:"capture"`
	RedisConsumerName      string `envconfig:"REDIS_CONSUMER_NAME" default:"capture"`
	S3Endpoint             string `envconfig:"S3_BUCKET_ENDPOINT"`
	S3BucketName           string `envconfig:"S3_BUCKET_NAME"`
	S3AccessKeyId          string `envconfig:"S3_ACCESS_KEY_ID"`
	S3AccessKey            string `envconfig:"S3_SECRET_ACCESS_KEY"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	// set consumer name dynamically
	if cfg.RedisConsumerName == "capture" {
		consumerName, err := os.Hostname()
		if err != nil {
			cfg.RedisConsumerName = consumerName
		}
	}

	return &cfg, nil
}
