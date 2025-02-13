package routers

import (
	"net/http"

	"github.com/Dor1ma/log-finder/logfinder/internal/server/handlers"
	"github.com/Dor1ma/log-finder/logfinder/internal/server/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(handler *handlers.LogHandler, rateLimit int) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/logs", middleware.RateLimit(middleware.LoggingMiddleware(handler.GetLogByTimestamp), rateLimit)).
		Methods("GET").
		Queries("timestamp", "{timestamp}")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	return r
}
