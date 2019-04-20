package function

import (
	"log"

	"github.com/prometheus/alertmanager/template"
)

func processAlert(alert template.Alert) {
	log.Println(alert)
	switch alert.Status {
	case "firing":
		log.Println("alert is firing")
	default:
		log.Printf("unknown alert status: %s", alert.Status)
		return
	}

	node, found := alert.Labels["node"]
	if !found {
		log.Println("alert doesn't have node label")
		return
	}

	alertName, found := alert.Labels["alertname"]
	if !found {
		log.Println("alert doesn't have alertname label")
		return
	}

	log.Printf("proccessing alert %s for node %s", alertName, node)
	// lookup remedies for alertName
	// lookup config for node
	// dispatch to remedy
}
