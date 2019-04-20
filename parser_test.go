package function

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	body, err := ioutil.ReadFile("testdata/alert.json")
	assert.Nil(t, err)
	msg, err := parse(body)
	assert.Nil(t, err)
	assert.Len(t, msg.Alerts, 1)
}
