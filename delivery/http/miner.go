package http

import (
	"CoalCompany/domain/miner"
	"CoalCompany/dto"
	appErrors "CoalCompany/errors"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

/*
Я изменил 1-ый endpoint и сделал так,
чтобы просто выводилась инфа о характеристиках каждого типа шахтера,
мне кажется, что так более логично.

pattern: /miners?info=types|active|not_active
method:  GET
info:    query params – info = types (classes) | active | not_active; class = small | normal | strong

succeed:
  - status code:   200 Ok
  - response body: JSON represent found remunerations

failed:
  - status code:   400, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleGetMiners(w http.ResponseWriter, r *http.Request) {
	infoParam := r.URL.Query().Get("info")
	class := r.URL.Query().Get("class")

	switch infoParam {
	case "types":
		h.sendMinerTypes(w)
	case "active":
		h.sendMinersResponse(w, h.enterprise.FindHiredMiners(boolPtr(true), class))
	case "not_active":
		h.sendMinersResponse(w, h.enterprise.FindHiredMiners(boolPtr(false), class))
	default:
		h.sendMinersResponse(w, h.enterprise.FindHiredMiners(nil, class))
	}
}

func (h *HTTPHandlers) sendMinerTypes(w http.ResponseWriter) {
	minerTypes := dto.MinerTypesResponse{
		Types: []dto.MinerTypeInfo{
			dto.MapMinerToTypeInfo("small", miner.NewSmallMiner()),
			dto.MapMinerToTypeInfo("normal", miner.NewNormalMiner()),
			dto.MapMinerToTypeInfo("strong", miner.NewStrongMiner()),
		},
	}

	w.Header().Set("Content-Type", "application/json")

	b, err := json.MarshalIndent(minerTypes, "", "    ")
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

func (h *HTTPHandlers) sendMinersResponse(w http.ResponseWriter, miners dto.HiredMinersResponse) {
	w.Header().Set("Content-Type", "application/json")

	b, err := json.MarshalIndent(miners, "", "    ")
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
pattern: /miners
method:  POST
info:    JSON in HTTP request body

succeed:
  - status code:   201 Created
  - response body: JSON represent found remunerations

failed:
  - status code:   400, 402, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleCreateMiner(w http.ResponseWriter, r *http.Request) {
	var minerType dto.MinerType
	if err := json.NewDecoder(r.Body).Decode(&minerType); err != nil {
		errDTO := dto.Err{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	hireMinerInfo, err := h.enterprise.HireMiner(minerType.Type)
	if err != nil {
		errDTO := dto.Err{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, appErrors.ErrMinerTypeNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else if errors.Is(err, appErrors.ErrInsufficientFunds) {
			http.Error(w, errDTO.ToString(), http.StatusPaymentRequired)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(hireMinerInfo, "", "    ")
	if err != nil {
		errDTO := dto.Err{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func boolPtr(b bool) *bool { return &b }
