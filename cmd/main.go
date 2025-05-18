package main

import (
	"context"
	"mau/cards/card"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/redis/go-redis/v9"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	h := card.NewApp(ctx, rdb)

	app := fiber.New()
	app.Get("/card/:asset/:card/:filter", h.SendCard)
	app.Listen(":3112")
}
