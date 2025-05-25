package card

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type App struct {
	ctx context.Context
	rdb *redis.Client
}

func NewApp(ctx context.Context, rdb *redis.Client) *App {
	return &App{ctx, rdb}
}

func (a App) SendCard(c *fiber.Ctx) error {
	assetParam := c.Params("asset")
	cardParam := c.Params("card")
	filterParam := c.Params("filter", "cover")

	unoCard, err := ParseCard(cardParam)
	if err != nil {
		errMsg := fmt.Errorf("parse card error: %w", err)
		log.Error().Err(err).Msg("Parse card error")
		return fiber.NewError(404, errMsg.Error())
	}

	// TODO: Если редиска откажется работать
	cacheKey := fmt.Sprintf("%s/%s/%s", assetParam, cardParam, filterParam)

	val, err := a.rdb.Get(a.ctx, cacheKey).Result()
	if err == nil {
		log.Debug().Str("key", cacheKey).Msg("Use card from cache")
		c.Set("Content-Type", "image/png")
		return c.Send([]byte(val))
	} else if err == redis.Nil {
		log.Debug().Any("key", cacheKey).Msg("Not found")
	} else {
		log.Error().Err(err).Msg("Redis error")
		return fiber.NewError(500, err.Error())
	}

	drawer, err := NewDrawer(assetParam, *unoCard)
	if err != nil {
		log.Error().Err(err).Msg("Create drawer")
		return err
	}

	i, err := drawer.Render(filterParam)
	if err != nil {
		log.Error().Err(err).Msg("Card render error")
		return err
	}

	buf, err := EncodeImage(i)
	if err != nil {
		log.Error().Err(err).Msg("Encode image error")
		return err
	}
	log.Debug().Msg("Card render complete")

	err = a.rdb.Set(a.ctx, cacheKey, buf.Bytes(), time.Duration(time.Hour)).Err()
	if err != nil {
		log.Error().Err(err).Msg("Redis error")
		return err
	}

	c.Set("Content-Type", "image/png")
	return c.Send(buf.Bytes())
}
