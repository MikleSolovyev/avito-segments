package app

import (
	"avito-segments/config"
	v1 "avito-segments/internal/controller/http/v1"
	"avito-segments/internal/repository"
	"avito-segments/internal/service"
	"avito-segments/pkg/db/postgres"
	"avito-segments/pkg/httpserver"
	"avito-segments/pkg/logger"
	"avito-segments/pkg/validator"
	"github.com/ziflex/lecho/v3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title           Avito user segments service
// @version         1.0
// @description     This service provides adding users to segments, removing them from segments, and saving these actions in the history.

// @contact.name   Mikhail Solovyev

// @host      localhost:8080
// @BasePath  /api/v1

func Run(cfg *config.Config) {
	// setup logger
	appLogger, err := logger.NewZerolog(&cfg.Logger)
	if err != nil {
		log.Fatal(err)
	}

	appLogger.Info().Msgf("logger level: %s", cfg.Logger.Level)
	appLogger.Info().Msg("STARTING APP...")

	// open db
	appLogger.Info().Msg("Connecting to database...")
	storage, err := postgres.New(&cfg.Storage)
	if err != nil {
		appLogger.Fatal().Msgf("Cannot open DB: %s", err)
	}
	defer storage.Close()

	// init repositories
	appLogger.Info().Msg("Initializing repositories...")
	repo := repository.New(storage)

	// init services
	appLogger.Info().Msg("Initializing services...")
	services := service.New(repo)

	// init validator
	validate := validator.New()

	// init router
	appLogger.Info().Msg("Initializing router...")
	router := v1.New(services, lecho.From(*appLogger), validate)

	// start HTTP server
	appLogger.Info().Msgf("Starting HTTP server on %s...", cfg.HTTPServer.Address)
	server := httpserver.New(router, &cfg.HTTPServer)

	appLogger.Info().Msg("APP IS RUNNING")

	// graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		appLogger.Info().Msgf("INTERRUPTED BY %s", sig.String())
	case err = <-server.Notify():
		appLogger.Error().Msgf("Server error: %s", err)
	}

	appLogger.Info().Msg("Shutting down gracefully...")
	err = server.Shutdown()
	if err != nil {
		appLogger.Error().Msgf("%s", err)
	}

	appLogger.Info().Msg("BYE.")
}
