package function

import (
	"log"
	"net/http"

	"github.com/go-redis/redis"
	handler "github.com/openfaas-incubator/go-function-sdk"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	options := &redis.Options{}

	data, err := parse(req.Body)
	if err != nil {
		return handler.Response{
			Body:       []byte("failed to unmarshal json request"),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	// TODO get redis connnection string from environment
	r := &Redis{redis.NewClient(options)}
	dc := NewDispatcherConfig(r)
	d := NewDispatcher(dc)

	for _, alert := range data.Alerts {
		if err := d.Dispatch(alert); err != nil {
			log.Println(err)
		}
	}

	return handler.Response{
		Body:       []byte("ok"),
		StatusCode: http.StatusOK,
	}, nil
}
