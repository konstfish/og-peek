package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RedisAddr       string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	RedisStreamName string `envconfig:"REDIS_STREAM_NAME" default:"screenshot-urls"`
	S3Endpoint      string `envconfig:"S3_BUCKET_ENDPOINT"`
	S3BucketName    string `envconfig:"S3_BUCKET_NAME"`
	S3AccessKeyId   string `envconfig:"S3_ACCESS_KEY_ID"`
	S3AccessKey     string `envconfig:"S3_SECRET_ACCESS_KEY"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
