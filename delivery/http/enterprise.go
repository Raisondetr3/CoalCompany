package http

import (
	"CoalCompany/dto"
	"encoding/json"
	"net/http"
	"time"
)

/*
pattern: /enterprise
method:  GET
info:    –

succeed:
  - status code:   200 Ok
  - response body: JSON represent found remunerations

failed:
  - status code:  500
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleGetEnterpriseStats(w http.ResponseWriter, r *http.Request) {
	stats := h.enterprise.GetEnterpriseStatsSafe()

	w.Header().Set("Content-Type", "application/json")

	b, err := json.MarshalIndent(stats, "", "    ")
	if err != nil {
		errDTO := dto.Err{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

/*
pattern: /enterprise
method:  POST
info:    –

succeed:
  - status code:   200 Ok
  - response body: JSON represent found remunerations

failed:
  - status code:  500
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleShutdownGame(w http.ResponseWriter, r *http.Request) {
	shutdownResponse, err := h.enterprise.ShutdownGame()
	if err != nil {
		errDTO := dto.Err{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	b, err := json.MarshalIndent(shutdownResponse, "", "    ")
	if err != nil {
		errDTO := dto.Err{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
