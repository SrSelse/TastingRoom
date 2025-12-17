package web

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func handleSubscriptionToken(logger *slog.Logger, hmacKey string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(ContextUserKey)

			if userId.(int) <= 0 {
				logger.Error("handleSubscriptionToken")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			data := struct {
				Channel string `json:"channel"`
			}{}

			err := json.NewDecoder(r.Body).Decode(&data)
			logger.Info("Datachannel", "data", data)
			if err != nil {
				logger.Error("handleSubscriptionToken", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub":     fmt.Sprintf("%d", userId.(int)),
				"channel": data.Channel,
			})

			// Sign and get the complete encoded token as a string using the secret
			tokenString, err := token.SignedString([]byte(hmacKey))
			if err != nil {
				logger.Error("handleSubscriptionToken", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
		},
	)
}

func handleConnectionToken(logger *slog.Logger, hmacKey string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(ContextUserKey)
			if userId.(int) <= 0 {
				logger.Error("handleConnectionToken")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return

			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": fmt.Sprintf("%d", userId.(int)),
				"exp": time.Now().Add(5 * time.Minute).Unix(),
			})

			// Sign and get the complete encoded token as a string using the secret
			tokenString, err := token.SignedString([]byte(hmacKey))
			if err != nil {
				logger.Error("handleConnectionToken", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
		},
	)
}
