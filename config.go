package function

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"regexp"
)

type Configurator interface {
	Fetch(string, string) (*RemedyConfig, error)
}

type ConfigFetcher struct {
	redis Rediser
	log   *log.Logger
}

type RemedyConfig struct {
	Name       string `json:"name"`
	NodeFilter string `json:"filter"`    // regexp
	Remedy     string `json:"remedy"`    // name of the faas function to invoke
	RateLimit  int64  `json:"ratelimit"` // max invocations per minute
	// TODO should contain a vault app role name that is used to pass a token to the remedy
}

// Fetch returns the first remedy configuration for alert that matches node,
// or nil if no config matched
func (cf *ConfigFetcher) Fetch(alert, node string) (*RemedyConfig, error) {
	// get all configurations and see if we have one that matches the alert and node
	configs, err := cf.redis.LRange(alert, 0, -1)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch configurations for %s", alert)
	}

	var conf RemedyConfig
	for _, c := range configs {
		err := json.Unmarshal([]byte(c), &conf)
		if err != nil {
			cf.log.Printf(`failed to unmarshal config "%s": %+v`, c, err)
			continue
		}
		re, err := regexp.Compile(conf.NodeFilter)
		if err != nil {
			cf.log.Printf(`failed to compile regexp "%s": %+v`, conf.NodeFilter, err)
			continue
		}
		if re.MatchString(node) {
			return &conf, nil
		}
	}

	return nil, nil
}
