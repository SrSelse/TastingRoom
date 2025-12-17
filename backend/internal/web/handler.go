package web

import (
	"log/slog"
	"net/http"

	"github.com/rs/cors"
	"skafteresort.se/beers/internal/auth"
	"skafteresort.se/beers/internal/beers"
	"skafteresort.se/beers/internal/rooms"
)

const ContextUserKey = "userID"

func NewServer(
	logger *slog.Logger,
	allowedOrigins []string,
	jwtSecret string,
	centrifugoHmacKey string,
	userService *auth.UserService,
	roomService *rooms.RoomService,
	beerService *beers.BeerService,

) http.Handler {

	// panic(allowedOrigins)
	corsMw := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		Debug:          false,
		AllowedHeaders: []string{"*"},
		// AllowCredentials: true,
	})

	mux := http.NewServeMux()
	mux.Handle("/auth/",
		corsMw.Handler(
			loggingMiddleware(logger,
				addRoutes(logger, userService, roomService, beerService),
			),
		),
	)

	mux.Handle("/broadcasting/",
		corsMw.Handler(
			loggingMiddleware(logger,
				jwtMiddleware(
					addCentrifugoRoutes(logger, centrifugoHmacKey),
					jwtSecret,
					logger,
				),
			),
		),
	)

	mux.Handle("/api/",
		corsMw.Handler(
			loggingMiddleware(logger,
				jwtMiddleware(
					addApiRoutes(logger, userService, roomService, beerService),
					jwtSecret,
					logger,
				),
			),
		),
	)

	return mux
}

func addCentrifugoRoutes(
	logger *slog.Logger,
	centrifugoHmacKey string,
) *http.ServeMux {

	mux := http.NewServeMux()

	mux.Handle(
		"/broadcasting/connect",
		handleConnectionToken(logger, centrifugoHmacKey),
	)

	mux.Handle(
		"/broadcasting/auth",
		handleSubscriptionToken(logger, centrifugoHmacKey),
	)

	return mux
}

// Register service to service routes.
func addRoutes(
	logger *slog.Logger,
	userService *auth.UserService,
	roomService *rooms.RoomService,
	beerService *beers.BeerService,
) *http.ServeMux {

	mux := http.NewServeMux()

	mux.Handle(
		"/auth/login",
		handleLogin(userService, logger),
	)

	mux.Handle(
		"/auth/register",
		handleRegister(userService, logger),
	)

	return mux
}
