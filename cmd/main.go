package main

import (
	"context"
	"mau/cards/card"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Setup redis client
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	defer func() {
		if err := rdb.Close(); err != nil {
			log.Error().Err(err).Msg("Redis disconnection error")
		}
	}()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("Redis connection error")
	}

	// Create fiber app
	appHandler := card.NewApp(ctx, rdb)
	app := fiber.New()
	app.Get("/card/:asset/:card/:filter", appHandler.SendCard)

	// Run fiber app
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	serverErrors := make(chan error, 1)
	go func() {
		const addr = ":3112"
		log.Info().Msgf("Run fiber server on %s", addr)
		err := app.Listen(addr)
		if err != nil {
			serverErrors <- err
		}
	}()

	select {
	case sig := <-quit:
		log.Info().Msgf("Catch signal: %v. Shutdown...", sig)
	case err := <-serverErrors:
		log.Error().Err(err).Msg("Server error")
	}

	// Попытка graceful shutdown с таймаутом
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Fiber stopping error")
	}

	log.Info().Msg("Server stopped")
}
