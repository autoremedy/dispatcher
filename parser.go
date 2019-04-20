package function

import (
	"encoding/json"

	"github.com/prometheus/alertmanager/template"
)

func parse(body []byte) (*template.Data, error) {
	var msg template.Data

	err := json.Unmarshal(body, &msg)
	return &msg, err
}
