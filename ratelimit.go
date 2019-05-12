package function

import (
	"log"
	"time"
)

type Ratelimiter interface {
	Limited(key string, max int64) (bool, error)
}

// Ratelimit is a distributed rate limiter, which uses this pattern:
// https://redislabs.com/redis-best-practices/basic-rate-limiting/
type Ratelimit struct {
	redis Rediser
	log   *log.Logger
}

// Limited checks if key has reached the max rate limit per hour
func (r *Ratelimit) Limited(key string, max int64) (bool, error) {
	// TODO make it possible to use a prate limit per minute, hour or day
	// minute is 04
	hour := time.Now().Format("15")
	k := key + ":" + hour

	s, err := r.redis.GetInt64(k)
	if err != nil {
		// fail closed so if redis is down we don't run without rate limit enforcement
		return true, err
	}

	if s == -1 {
		r.log.Printf("%s is not set", k)
	} else {
		r.log.Printf("%s is %d", k, s)
		if s > max {
			return false, nil
		}
	}

	i, err := r.redis.Incr(k)
	if err != nil {
		return false, err
	}

	ok, err := r.redis.Expire(k, time.Minute)
	if err != nil {
		return false, err
	}
	if !ok {
		r.log.Printf("failed to set expiration on %s", k)
	}

	return i <= max, nil
}
