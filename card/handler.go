package card

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type App struct {
	ctx context.Context
	rdb *redis.Client
}

func NewApp(ctx context.Context, rdb *redis.Client) App {
	return App{ctx, rdb}
}

func (a App) SendCard(c *fiber.Ctx) error {
	cardStr := c.Params("card")
	coverStr := c.Params("cover", "true")
	log.Info().Str("card", cardStr).Str("cover", coverStr).Msg("New request")

	unoCard, err := ParseCard(cardStr)
	if err != nil {
		return fiber.NewError(404, err.Error())
	}

	key := cardStr + "/" + coverStr
	val, err := a.rdb.Get(a.ctx, key).Result()
	if err == redis.Nil {
		log.Debug().Any("key", key).Msg("Not found")

	} else if err != nil {
		return err
	} else {
		c.Set("Content-Type", "image/png")
		return c.Send([]byte(val))
	}

	now := time.Now()
	log.Debug().Time("now", now).Msg("Start render")

	i, err := RenderCard(*unoCard, coverStr)
	if err != nil {
		return err
	}
	log.Info().Dur("took", time.Since(now)).Msg("Render complete")

	buf, err := EncodeImage(i)
	if err != nil {
		return err
	}
	log.Info().Dur("took", time.Since(now)).Msg("Card render complete")

	err = a.rdb.Set(a.ctx, key, buf.Bytes(), time.Duration(time.Duration.Hours(1))).Err()
	if err != nil {
		return err
	}

	c.Set("Content-Type", "image/png")
	return c.Send(buf.Bytes())
}
