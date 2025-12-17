package web

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"skafteresort.se/beers/internal/auth"
)

func handleLogin(
	us *auth.UserService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var u auth.LoginAttempt
			json.NewDecoder(r.Body).Decode(&u)
			user, err := us.SignIn(r.Context(), u.Username, u.Password)
			if err != nil {
				if errors.As(err, &auth.UnauthenticatedError{}) {
					logger.Error("handleLogin", "err", err)
					http.Error(w, "Username or Password incorrect", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				logger.Error("handleLogin", "err", err)
				return
			}

			jwtToken, err := auth.CreateToken(user.Username, user.Id)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				logger.Error("handleLogin/jwt", "err", err)
				return
			}
			cookie := http.Cookie{
				Name:     "token",
				Value:    jwtToken,
				Path:     "/",
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				HttpOnly: true,
				MaxAge:   86400, // Lives for 1 hour
			}
			http.SetCookie(w, &cookie)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"token": jwtToken, "user": user})
			return

		},
	)
}

func handleRegister(
	us *auth.UserService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var u auth.SignupAttempt
			json.NewDecoder(r.Body).Decode(&u)
			if exists, err := us.UsernameInUse(r.Context(), u.Username); exists || err != nil {
				logger.Error("handleRegister/db", "err", err)
				http.Error(w, "Username already taken", http.StatusUnprocessableEntity)
				return
			}
			// username := r.PostForm.Get("username")
			// password := r.PostForm.Get("password")
			err := us.SignUp(r.Context(), u.Username, u.Password, u.Name)
			if err != nil {
				if errors.As(err, &auth.DatabaseError{}) {
					logger.Error("handleRegister/db", "err", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				logger.Error("handleLogin", "err", err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Success")
			return

		},
	)
}
