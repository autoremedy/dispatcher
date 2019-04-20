package function

import (
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {

	data, err := parse(req.Body)
	if err != nil {
		return handler.Response{
			Body:       []byte("failed to unmarshal json request"),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	for _, alert := range data.Alerts {
		processAlert(alert)
	}

	return handler.Response{
		Body:       []byte("ok"),
		StatusCode: http.StatusOK,
	}, nil
}
