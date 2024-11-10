package handlers

import (
	"encoding/json"
	"net/http"
)

type Healthcheck struct {
	Healthy bool `json:"healthy"`
}

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	health, err := json.Marshal(Healthcheck{Healthy: true})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(health)
}
