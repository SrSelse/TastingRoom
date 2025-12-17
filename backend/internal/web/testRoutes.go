package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func handleTestToken(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(ContextUserKey)
			logger.Info("handleTestToken", "userId", userId)
			if userId.(int) < 1 {
				logger.Error("handleTestToken", "err", "userId < 1")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"userId": userId})
		},
	)
}
