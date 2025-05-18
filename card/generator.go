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

// Генератор карт
type CardDrawer struct {
	asset string
	card  Card
	image *image.NRGBA
}

func LoadAsset(asset string, path string) (*image.NRGBA, error) {
	file, err := os.Open(ASSETS_PATH + asset + "/" + path)
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

func (d CardDrawer) addColor(color string) error {
	colorImage, err := LoadAsset(d.asset, "color/"+color+".png")
	if err != nil {
		return err
	}

	draw.Draw(d.image, d.image.Bounds(), colorImage, image.Point{}, draw.Over)
	return nil
}

func (d CardDrawer) addSym(sym []string) error {
	x, y := -48, -48
	for _, s := range sym {
		symImage, err := LoadAsset(d.asset, "sym/"+s+".png")
		if err != nil {
			return err
		}

		bounds := symImage.Bounds()
		// TODO: Разобраться со сдвигом
		// yOffset := (64 - bounds.Max.Y/2)
		draw.Draw(d.image, d.image.Bounds(), symImage, image.Point{x, y}, draw.Over)

		x = x - 8 - bounds.Max.X
	}

	return nil
}

func (d CardDrawer) addReverseSym(sym []string) error {
	x, y := -271, -463
	for _, s := range sym {
		symImage, err := LoadAsset(d.asset, "sym_reverse/"+s+".png")
		if err != nil {
			return err
		}

		bounds := symImage.Bounds()
		// TODO: Разобраться со сдвигом
		px := x + bounds.Max.X
		py := y + bounds.Max.Y
		draw.Draw(d.image, d.image.Bounds(), symImage, image.Point{px, py}, draw.Over)

		x = x + 8 + bounds.Max.X
	}

	return nil
}

func (d CardDrawer) addGlyph(color string, glyph string) error {
	assetPath := fmt.Sprintf("glyph/%s/%s.png", color, glyph)
	if d.card.Wild {
		assetPath = fmt.Sprintf("glyph/%s.png", glyph)
	}

	glyphImage, err := LoadAsset(d.asset, assetPath)
	if err != nil {
		return err
	}

	bounds := d.image.Bounds()
	glyphBounds := glyphImage.Bounds()
	x := -(bounds.Max.X/2 - (glyphBounds.Max.X / 2))
	y := -(bounds.Max.Y/2 - (glyphBounds.Max.Y / 2))

	draw.Draw(d.image, bounds, glyphImage, image.Point{x, y}, draw.Over)
	return nil
}

func uncover(base *image.NRGBA) error {
	darknessFactor := 0.6
	bounds := base.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := base.At(x, y)
			r, g, b, a := originalColor.RGBA()

			r = uint32(float64(r>>8) * darknessFactor)
			g = uint32(float64(g>>8) * darknessFactor)
			b = uint32(float64(b>>8) * darknessFactor)
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
			[]string{"plus", strconv.Itoa(card.Value)}, takeGlyph(card.Value),
		}

	case "twist":
		return drawParams{
			[]string{strconv.Itoa(card.Value), "reverse"},
			strconv.Itoa(card.Value),
		}

	case "put":
		return drawParams{
			[]string{"minus", strconv.Itoa(card.Value)}, takeGlyph(card.Value),
		}

	case "color":
		return drawParams{[]string{"color"}, "/color"}

	case "delta":
		return drawParams{[]string{"delta"}, "/delta"}

	case "wild+color":
		return drawParams{[]string{"color"}, "color"}

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

func NewDrawer(asset string, card Card) (*CardDrawer, error) {
	m, err := LoadAsset(asset, "base.png")
	if err != nil {
		return nil, err
	}
	return &CardDrawer{asset: asset, card: card, image: m}, nil
}

// Главная функция
func (d CardDrawer) Render(filter string) (*image.NRGBA, error) {
	colorStr := strconv.Itoa(int(d.card.Color))

	err := d.addColor(colorStr)
	if err != nil {
		return nil, err
	}

	params := cardParams(d.card)

	err = d.addSym(params.sym)
	if err != nil {
		return nil, err
	}

	err = d.addReverseSym(params.sym)
	if err != nil {
		return nil, err
	}

	err = d.addGlyph(colorStr, params.glyph)
	if err != nil {
		return nil, err
	}

	if filter == "uncover" {
		uncover(d.image)
	}

	return d.image, nil
}
