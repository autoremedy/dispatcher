package function

import (
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type fakeRedis struct {
	lrange []string
}

func (r *fakeRedis) GetInt64(key string) (int64, error) {
	return 0, nil
}

func (r *fakeRedis) Incr(key string) (int64, error) {
	return 0, nil
}

func (r *fakeRedis) Expire(key string, duration time.Duration) (bool, error) {
	return true, nil
}

func (r *fakeRedis) LRange(key string, start, stop int64) ([]string, error) {
	return r.lrange, nil
}

func TestRedis(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)
	defer mr.Close()

	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	r := Redis{rc}
	i, err := r.Incr("key")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), i)
}
