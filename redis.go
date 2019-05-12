package function

import (
	"time"

	"github.com/go-redis/redis"
)

type Rediser interface {
	GetInt64(string) (int64, error)
	Incr(string) (int64, error)
	Expire(string, time.Duration) (bool, error)
	LRange(string, int64, int64) ([]string, error)
}

type Redis struct {
	client *redis.Client
}

// GetInt64 wraps the redis call, returns -1 when key is not found
func (r *Redis) GetInt64(key string) (int64, error) {
	i := r.client.Get(key)
	if i.Val() == "" {
		return -1, nil
	}
	return i.Int64()
}

func (r *Redis) Incr(key string) (int64, error) {
	return r.client.Incr(key).Result()
}

func (r *Redis) Expire(key string, duration time.Duration) (bool, error) {
	return r.client.Expire(key, time.Minute).Result()
}

func (r *Redis) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(key, start, stop).Result()
}
