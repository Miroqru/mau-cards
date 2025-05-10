package card

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
)

const ASSETS_PATH = "assets/"

func LoadAsset(path string) (*image.NRGBA, error) {
	file, err := os.Open(ASSETS_PATH + path)
	if err != nil {
		return nil, fmt.Errorf("loader: %w", err)
	}
	defer file.Close()

	m, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("loader: %w", err)
	}
	return m.(*image.NRGBA), nil
}

func EncodeImage(m *image.NRGBA) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, m); err != nil {
		return nil, err
	}
	return &buf, nil
}

// Функции отрисовки
// =================

func addColor(base *image.NRGBA, color string) error {
	colorImage, err := LoadAsset("color/" + color + ".png")
	if err != nil {
		return err
	}

	draw.Draw(base, base.Bounds(), colorImage, image.Point{}, draw.Over)
	return nil
}

func addSym(base *image.NRGBA, sym []string) error {
	x, y := -48, -48
	for _, s := range sym {
		symImage, err := LoadAsset("sym/" + s + ".png")
		if err != nil {
			return err
		}

		bounds := symImage.Bounds()
		// TODO: Разобраться со сдвигом
		// yOffset := (64 - bounds.Max.Y/2)
		draw.Draw(base, base.Bounds(), symImage, image.Point{x, y}, draw.Over)

		x = x - 8 - bounds.Max.X
	}

	return nil
}

func addReverseSym(base *image.NRGBA, sym []string) error {
	x, y := -271, -463
	for _, s := range sym {
		symImage, err := LoadAsset("sym_reverse/" + s + ".png")
		if err != nil {
			return err
		}

		bounds := symImage.Bounds()
		// TODO: Разобраться со сдвигом
		px := x + bounds.Max.X
		py := y + bounds.Max.Y
		draw.Draw(base, base.Bounds(), symImage, image.Point{px, py}, draw.Over)

		x = x + 8 + bounds.Max.X
	}

	return nil
}

func addGlyph(base *image.NRGBA, color string, glyph string) error {
	glyphImage, err := LoadAsset("glyph/" + color + "/" + glyph + ".png")
	if err != nil {
		return err
	}

	bounds := base.Bounds()
	glyphBounds := glyphImage.Bounds()
	x := -(bounds.Max.X/2 - (glyphBounds.Max.X / 2))
	y := -(bounds.Max.Y/2 - (glyphBounds.Max.Y / 2))

	draw.Draw(base, bounds, glyphImage, image.Point{x, y}, draw.Over)
	return nil
}

func uncover(base *image.NRGBA) error {
	darknessFactor := 0.6
	bounds := base.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Получаем цвет пикселя из исходного изображения
			originalColor := base.At(x, y)
			r, g, b, a := originalColor.RGBA()

			// Применяем затемнение
			r = uint32(float64(r>>8) * darknessFactor)
			g = uint32(float64(g>>8) * darknessFactor)
			b = uint32(float64(b>>8) * darknessFactor)

			// Устанавливаем новый цвет в затемненное изображение
			base.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a >> 8)})
		}
	}

	return nil
}

type drawParams struct {
	sym   []string
	glyph string
}

func takeGlyph(n int) string {
	if n == 1 {
		return "take_1"
	}
	if n < 4 {
		return "take_2"

	}
	return "take_4"
}

func cardParams(card Card) drawParams {
	switch card.Behavior {
	case "turn":
		return drawParams{[]string{"block"}, "block"}

	case "reverse":
		return drawParams{[]string{"reverse"}, "reverse"}

	case "take":
		return drawParams{
			[]string{"plus", strconv.Itoa(card.Value)},
			takeGlyph(card.Value),
		}

	case "wild+color":
		return drawParams{[]string{}, "delta"}

	case "wild+take":
		return drawParams{
			[]string{"plus", strconv.Itoa(card.Value)},
			takeGlyph(card.Value),
		}

	default:
		return drawParams{
			[]string{strconv.Itoa(card.Value)},
			strconv.Itoa(card.Value),
		}
	}
}

// Главная функция
func RenderCard(card Card, cover string) (*image.NRGBA, error) {
	m, err := LoadAsset("base.png")
	if err != nil {
		return nil, err
	}
	colorStr := strconv.Itoa(int(card.Color))

	err = addColor(m, colorStr)
	if err != nil {
		return nil, err
	}

	params := cardParams(card)

	err = addSym(m, params.sym)
	if err != nil {
		return nil, err
	}

	err = addReverseSym(m, params.sym)
	if err != nil {
		return nil, err
	}

	err = addGlyph(m, colorStr, params.glyph)
	if err != nil {
		return nil, err
	}

	if cover == "false" {
		uncover(m)
	}

	return m, nil
}
