package ratelimiter

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

type fakeRedis struct{}

func (f *fakeRedis) Get(string) *redis.StringCmd {
	return redis.NewStringCmd("asd")
}

func (f *fakeRedis) Incr(string) *redis.IntCmd {
	return redis.NewIntCmd()
}

func (f *fakeRedis) Expire(string, time.Duration) *redis.BoolCmd {
	return redis.NewBoolCmd()
}

func TestRedis(t *testing.T) {
	fake := &fakeRedis{}
	r := New(fake, 2)
	ok, err := r.Check("key")
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestRedisIntegration(t *testing.T) {
	t.Skip("integration")
	r := redis.NewClient(&redis.Options{
		Addr: "192.168.174.133:6379",
	})
	rl := New(r, 2)

	key := "bar"
	ok, err := rl.Check(key)
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = rl.Check(key)
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = rl.Check(key)
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = rl.Check(key)
	assert.Nil(t, err)
	assert.False(t, ok)
}
