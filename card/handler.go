package card

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SendCard(c *fiber.Ctx) error {
	cardStr := c.Params("card")
	coverStr := c.Params("cover", "true")
	log.Info().Str("card", cardStr).Str("cover", coverStr).Msg("New request")

	unoCard, err := ParseCard(cardStr)
	if err != nil {
		return fiber.NewError(404, err.Error())
	}
	now := time.Now()
	log.Debug().Time("now", now).Msg("Start render")

	i, err := RenderCard(*unoCard, coverStr)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	log.Info().Dur("took", time.Since(now)).Msg("Render complete")

	buf, err := EncodeImage(i)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	log.Info().Dur("took", time.Since(now)).Msg("Card render complete")

	c.Set("Content-Type", "image/png")
	return c.SendStream(buf)
}
