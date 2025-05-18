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

func NewApp(ctx context.Context, rdb *redis.Client) App {
	return App{ctx, rdb}
}

func (a App) SendCard(c *fiber.Ctx) error {
	assetParam := c.Params("asset")
	cardParam := c.Params("card")
	filterParam := c.Params("filter", "cover")
	log.Debug().Str("card", cardParam).Str("filter", filterParam).Str("asset", assetParam).Msg("New request")

	unoCard, err := ParseCard(cardParam)
	if err != nil {
		return fiber.NewError(404, fmt.Errorf("Card not found: %w", err).Error())
	}

	// TODO: Если редиска откажется работать
	key := assetParam + "/" + cardParam + "/" + filterParam
	val, err := a.rdb.Get(a.ctx, key).Result()
	if err == redis.Nil {
		log.Debug().Any("key", key).Msg("Not found")

	} else if err != nil {
		return err
	} else {
		c.Set("Content-Type", "image/png")
		return c.Send([]byte(val))
	}

	r, err := NewDrawer(assetParam, *unoCard)
	if err != nil {
		return err
	}

	i, err := r.Render(filterParam)
	if err != nil {
		return err
	}

	buf, err := EncodeImage(i)
	if err != nil {
		return err
	}
	log.Debug().Msg("Card render complete")

	err = a.rdb.Set(a.ctx, key, buf.Bytes(), time.Duration(time.Duration.Hours(1))).Err()
	if err != nil {
		return err
	}

	c.Set("Content-Type", "image/png")
	return c.Send(buf.Bytes())
}
