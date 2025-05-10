package main

import (
	"mau/cards/card"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	app := fiber.New()
	app.Get("/card/:card/:cover", card.SendCard)
	app.Listen(":3112")
}
