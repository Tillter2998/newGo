package routes

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Tillter2998/truroMap/internal/handlers"
)

func AddRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
) {
	mux.Handle("GET /", healthCheck(logger))
	mux.Handle("POST /setTrip", handlers.GetDirections(logger))
}

func healthCheck(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.LogAttrs(context.Background(), slog.LevelInfo, "", slog.String("endpoint", "healthCheck"))

			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte("{\"Messaage\":\"Server is up\"}"))
		},
	)
}
