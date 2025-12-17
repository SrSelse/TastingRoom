package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/centrifugal/gocent/v3"

	"skafteresort.se/beers/internal/auth"
	"skafteresort.se/beers/internal/beers"
	"skafteresort.se/beers/internal/providers"
	"skafteresort.se/beers/internal/rooms"
	"skafteresort.se/beers/internal/web"
)

var ServiceVersion string

type Server struct {
	db         *sql.DB
	logger     *slog.Logger
	config     ServerConfig
	httpServer *http.Server

	beerService *beers.BeerService
	roomService *rooms.RoomService
	userService *auth.UserService

	// beerRepo
}

func NewServer(config ServerConfig, logger *slog.Logger) *Server {
	return &Server{
		config: config,
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) {
	var err error

	s.logger.Info("Starting server", "version", ServiceVersion)

	if s.config.storageDriver != "mysql" {
		s.logger.Error("Invalid storage driver")
		return
	}

	s.db, err = sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			s.config.dbUser,
			s.config.dbPass,
			s.config.dbHost,
			s.config.dbPort,
			s.config.dbName,
		),
	)

	defer s.db.Close()

	if err != nil {
		s.logger.Error("Unable to connect to mysql", slog.String("error", err.Error()))
		return
	}

	if err = s.db.Ping(); err != nil {
		s.logger.Error("Unable to ping mysql", slog.String("error", err.Error()))
		return
	}

	gocentClient := gocent.New(gocent.Config{
		Addr: s.config.centrifugoApi,
		Key:  s.config.centrifugoKey,
	})

	centrifugoProvider := providers.NewCentrifugoProvider(gocentClient, s.logger)

	s.beerService = beers.NewBeerService(
		beers.NewBeerRepo(s.db),
		s.logger,
		centrifugoProvider,
	)

	s.roomService = rooms.NewRoomService(
		rooms.NewRoomRepo(s.db),
		s.logger,
		centrifugoProvider,
	)

	s.userService = auth.NewUserService(
		auth.NewUserRepo(s.db),
		s.logger,
	)

	s.serveHTTP()

	<-ctx.Done()
	s.Shutdown(ctx)
}

func (s *Server) serveHTTP() {
	handler := web.NewServer(
		s.logger,
		s.config.httpCorsAllowedOrigin,
		s.config.jwtSecret,
		s.config.centrifugoHmacKey,
		s.userService,
		s.roomService,
		s.beerService,
	)
	s.httpServer = &http.Server{
		Addr:         s.config.httpEndpointPort,
		Handler:      handler,
		ReadTimeout:  s.config.httpReadTimeout,
		WriteTimeout: s.config.httpWriteTimeout,
	}

	s.logger.Info("Starting webserver", slog.String("Listening", s.httpServer.Addr))

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Webserver error", slog.String("error", err.Error()))
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) {

	// Shutdown context/withTimeout
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("HTTP Server shutdown", "error", err)
		return
	}

	s.logger.Info("Graceful shutdown completed")
}
