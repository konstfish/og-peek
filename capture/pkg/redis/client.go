package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	Rdb *redis.Client
}

func NewClient(addr string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &Client{Rdb: rdb}, nil
}

func (c *Client) Close() error {
	return c.Rdb.Close()
}

func (c *Client) Set(ctx context.Context, domain string, slug string, value string) error {
	key := fmt.Sprintf("domain:%s:%s", domain, slug)
	return c.Rdb.Set(ctx, key, value, 7*24*time.Hour).Err()
}
