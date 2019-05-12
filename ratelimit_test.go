package function

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeRatelimit struct{ limited bool }

func (f *fakeRatelimit) Limited(key string, max int64) (bool, error) { return f.limited, nil }

func TestRatelimit(t *testing.T) {
	r := &fakeRedis{}
	rl := Ratelimit{r, discardLogger}
	limited, err := rl.Limited("key", 10)
	assert.Nil(t, err)
	assert.True(t, limited)
}
