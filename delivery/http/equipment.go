package http

import (
	"CoalCompany/dto"
	appErrors "CoalCompany/errors"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

/*
pattern: /equipments?status=purchased
method:  GET
info:    query params

succeed:
  - status code:   200 Ok
  - response body: JSON represent found remunerations

failed:
  - status code:   400, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleGetEquipments(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	if status == "purchased" {
		h.sendEquipments(w, h.enterprise.GetPurchasedEquipments())
	} else {
		h.sendEquipments(w, h.enterprise.GetAllEquipments())
	}
}

func (h *HTTPHandlers) sendEquipments(w http.ResponseWriter, equipments dto.EquipmentResponse) {
	w.Header().Set("Content-Type", "application/json")

	b, err := json.MarshalIndent(equipments, "", "    ")
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
pattern: /equipments
method:  POST
info:    JSON in HTTP request body

succeed:
  - status code:   201 Created
  - response body: JSON represent found remunerations

failed:
  - status code:   400, 402, 404, 409, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleBuyEquipment(w http.ResponseWriter, r *http.Request) {
	var req dto.PurchaseEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errDTO := dto.Err{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := h.enterprise.BuyEquipment(req.Type); err != nil {
		errDTO := dto.Err{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, appErrors.ErrEquipmentNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else if errors.Is(err, appErrors.ErrInsufficientFunds) {
			http.Error(w, errDTO.ToString(), http.StatusPaymentRequired)
		} else if errors.Is(err, appErrors.ErrEquipmentAlreadyPurchased) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := dto.PurchaseEquipmentResponse{
		Type:    req.Type,
		Balance: h.enterprise.GetBalance(),
	}

	b, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		errDTO := dto.Err{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
