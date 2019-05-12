package function

import (
	"github.com/prometheus/alertmanager/template"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

var discardLogger = log.New(ioutil.Discard, "", 0)

func TestDispatch(t *testing.T) {
	frl := &fakeRatelimit{}
	fc := &fakeConfig{&RemedyConfig{
		NodeFilter: ".*",
		Remedy:     "foo",
		RateLimit:  10,
	}}
	dc := &DispatcherConfig{frl, fc, discardLogger}
	d := NewDispatcher(dc)
	err := d.Dispatch(makeAlert("a", "n"))
	assert.Nil(t, err)
}

func makeAlert(alert, node string) template.Alert {
	return template.Alert{
		Status: "firing",
		Labels: map[string]string{
			"alertname": alert,
			"node":      node,
		},
	}
}
