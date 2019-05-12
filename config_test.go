package function

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeConfig struct{ conf *RemedyConfig }

func (f *fakeConfig) Fetch(alert, node string) (*RemedyConfig, error) {
	return f.conf, nil
}

func TestLookupRemedyConfig(t *testing.T) {
	fr := &fakeRedis{
		lrange: []string{
			`{"filter":"foo.*","remedy":"foo"}`,
			`{"filter":"bar.*","remedy":"bar"}`,
		},
	}
	conf := ConfigFetcher{fr, discardLogger}

	c, err := conf.Fetch("alert", "bar")
	assert.Nil(t, err)
	assert.Equal(t, "bar", c.Remedy)
}
