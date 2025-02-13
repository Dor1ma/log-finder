package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dor1ma/log-finder/logfinder/internal/service"
)

var timeFormat = "2006-01-02T15:04:05.000"

type LogHandler struct {
	useCase *service.LogService
}

func NewLogHandler(uc *service.LogService) *LogHandler {
	return &LogHandler{useCase: uc}
}

func (h *LogHandler) GetLogByTimestamp(w http.ResponseWriter, r *http.Request) {
	timestampParam := r.URL.Query().Get("timestamp")
	if timestampParam == "" {
		http.Error(w, "timestamp parameter is required", http.StatusBadRequest)
		return
	}

	timestamp, err := time.Parse(timeFormat, timestampParam)
	if err != nil {
		http.Error(w, "invalid timestamp format", http.StatusBadRequest)
		return
	}

	result, err := h.useCase.FindLog(r.Context(), timestamp)
	if err != nil {
		if err == service.ErrNotFound {
			http.Error(w, "log entry not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Timestamp time.Time `json:"timestamp"`
		Message   string    `json:"message"`
	}{
		Timestamp: timestamp,
		Message:   result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
