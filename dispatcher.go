package function

import (
	"fmt"
	"github.com/prometheus/alertmanager/template"
	"log"
	"os"
)

type DispatcherConfig struct {
	rl     Ratelimiter
	config Configurator
	log    *log.Logger
}

func NewDispatcherConfig(r Rediser) *DispatcherConfig {
	l := log.New(os.Stderr, "", log.LstdFlags)
	return &DispatcherConfig{
		&Ratelimit{r, l},
		&ConfigFetcher{r, l},
		l,
	}
}

type Dispatcher struct {
	*DispatcherConfig
}

func NewDispatcher(conf *DispatcherConfig) *Dispatcher {
	return &Dispatcher{conf}
}

func (d *Dispatcher) Dispatch(alert template.Alert) error {
	if alert.Status != "firing" {
		return fmt.Errorf("unknown alert status: %s", alert.Status)
	}

	alertName, found := alert.Labels["alertname"]
	if !found {
		return fmt.Errorf("alert doesn't have alertname label")
	}

	node, found := alert.Labels["node"]
	if !found {
		return fmt.Errorf("alert %s doesn't have node label", alertName)
	}

	d.log.Printf("proccessing alert %s for node %s", alertName, node)

	conf, err := d.config.Fetch(alertName, node)
	if err != nil {
		// TODO: log failure metric
		return fmt.Errorf("failed to lookup config %s and %s: %+v", alertName, node, err)
	}
	if conf == nil {
		return fmt.Errorf("no remedy available for %s and %s", alertName, node)

	}

	limited, err := d.rl.Limited(alertName, conf.RateLimit)
	if err != nil {
		// TODO: log failure metric
		return fmt.Errorf("failed to check rate limit for %s/%s: %+v", alertName, node, err)
	}

	if limited {
		return fmt.Errorf("rate limit reached")
	}

	// async dispatch to remedy function
	// TODO: log metric
	return nil
}
