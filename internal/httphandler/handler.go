package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kfc-forecast/internal/core"
)

type Handler struct {
	svc core.ForecastService
}

func New(svc core.ForecastService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stores", h.listStores)
	mux.HandleFunc("GET /api/forecast", h.getForecast)
}

func (h *Handler) listStores(w http.ResponseWriter, r *http.Request) {
	stores, err := h.svc.ListStores()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, stores)
}

func (h *Handler) getForecast(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.URL.Query().Get("store_id")
	dateStr := r.URL.Query().Get("date")

	if storeIDStr == "" || dateStr == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("store_id and date are required"))
		return
	}

	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid store_id"))
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("date must be YYYY-MM-DD"))
		return
	}

	entries, err := h.svc.GetForecast(storeID, date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if entries == nil {
		entries = []core.ForecastEntry{}
	}
	writeJSON(w, http.StatusOK, entries)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) 
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
