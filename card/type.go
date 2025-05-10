package card

import (
	"errors"
	"regexp"
	"strconv"
)

type CardColor int
type CardBehavior string

type Card struct {
	Color    CardColor
	Value    int
	Cost     int
	Behavior CardBehavior
}

var CARD_REGEX = regexp.MustCompile(`(\d)_(\d)_(\d+)_([a-z+]+)`)

// Собирает карт из строки
// Формат цвет_значение_стоимость_тип
func ParseCard(card string) (*Card, error) {
	cardMatch := CARD_REGEX.FindStringSubmatch(card)
	if cardMatch == nil {
		return nil, errors.New("incorrect card format")
	}

	colorInt, err := strconv.Atoi(cardMatch[1])
	if err != nil {
		return nil, err
	}

	value, err := strconv.Atoi(cardMatch[2])
	if err != nil {
		return nil, err
	}

	cost, err := strconv.Atoi(cardMatch[3])
	if err != nil {
		return nil, err
	}

	return &Card{
		Color:    CardColor(colorInt),
		Value:    value,
		Cost:     cost,
		Behavior: CardBehavior(cardMatch[4]),
	}, nil
}
