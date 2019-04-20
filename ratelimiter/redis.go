package ratelimiter

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

type Rediser interface {
	Get(string) *redis.StringCmd
	Incr(string) *redis.IntCmd
	Expire(string, time.Duration) *redis.BoolCmd
}

type Client struct {
	client Rediser
	max    int64
}

func New(redis Rediser, max int64) *Client {
	return &Client{redis, max}
}

func (c *Client) Check(key string) (bool, error) {
	minute := time.Now().Format("04")
	k := key + ":" + minute

	s := c.client.Get(k).Val()

	if s == "" {
		log.Printf("%s is not set", k)
	} else {
		i, err := c.client.Get(k).Int64()
		if err != nil {
			return false, err
		}
		log.Printf("%s is %d", k, i)
		if i > c.max {
			return false, nil
		}
	}

	i, err := c.client.Incr(k).Result()
	if err != nil {
		return false, err
	}

	ok, err := c.client.Expire(k, time.Minute).Result()
	if err != nil {
		return false, err
	}
	if !ok {
		log.Printf("failed to set expiration on %s", k)
	}

	return i <= c.max, nil
}
