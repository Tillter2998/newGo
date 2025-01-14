package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type tripOptions struct {
	StartingPoint string `json:"startingPoint"`
	Destination   string `json:"destination"`
}

func GetDirections(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.LogAttrs(context.Background(), slog.LevelInfo, "", slog.String("endpoint", "setTrip"))

			var tripOptions tripOptions
			if err := json.NewDecoder(r.Body).Decode(&tripOptions); err != nil {
				logger.LogAttrs(context.Background(), slog.LevelError, err.Error())
			}
			defer r.Body.Close()

			// Get weighted digraph and find shortest path to destination from startingPoint

			w.Header().Add("Content-Type", "application/json")
			w.Write()
		},
	)
}
