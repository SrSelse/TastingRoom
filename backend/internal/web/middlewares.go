package web

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

type TokenClaim struct {
	jwt.RegisteredClaims

	Subject  int    `json:"sub,omitempty"`
	Username string `json:"username"`
}

func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func(start time.Time) {
			if req.Method == http.MethodOptions {
				next.ServeHTTP(w, req)

				return
			}

			// attributes := []slog.Attr{}
			// if userID, ok := req.Context().Value(ContextUserKey).(int); ok {
			// 	attributes = append(attributes, slog.Int("context.userID", userID))
			// }

			// attributes = append(
			// 	attributes,
			// 	slog.Attr{
			// 		Key: "request",
			// 		Value: slog.GroupValue(
			// 			slog.String("method", req.Method),
			// 			slog.String("path", req.URL.Path),
			// 			slog.String("ip", req.RemoteAddr),
			// 			slog.Duration("time", time.Since(start)),
			// 		),
			// 	},
			// )
			// logger.LogAttrs(req.Context(), slog.LevelInfo, "http", attributes...)
		}(time.Now())

		next.ServeHTTP(w, req)
	})
}
func jwtMiddleware(handler http.Handler, jwtSecret string, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, err := request.ParseFromRequest(req, request.BearerExtractor{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		}, request.WithClaims(&TokenClaim{}))

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
			// error
		}

		claims, ok := token.Claims.(*TokenClaim)
		if !ok || claims.Subject == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(
			req.Context(),
			ContextUserKey,
			token.Claims.(*TokenClaim).Subject,
		)

		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}

// func jwtMiddleware(handler http.Handler, jwtSecret string, logger *slog.Logger) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		ct, err := req.Cookie("token")
// 		cookies := req.Cookies()
//
// 		logger.Info("jwtMiddleware", "ct", ct, "err", err, "cookies", cookies)
// 		token, err := jwt.ParseWithClaims(
// 			ct.Value,
// 			&TokenClaim{},
// 			func(t *jwt.Token) (interface{}, error) {
// 				return []byte(jwtSecret), nil
// 			},
// 		)
//
// 		if err != nil {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 			// error
// 		}
//
// 		claims, ok := token.Claims.(*TokenClaim)
// 		if !ok || claims.Subject == 0 {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
//
// 			return
// 		}
//
// 		ctx := context.WithValue(
// 			req.Context(),
// 			ContextUserKey,
// 			token.Claims.(*TokenClaim).Subject,
// 		)
//
// 		handler.ServeHTTP(w, req.WithContext(ctx))
// 	})
// }
