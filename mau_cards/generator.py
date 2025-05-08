"""Собирает изображение карты по её описанию."""

import io
from collections.abc import Iterable
from pathlib import Path

from mau.deck import behavior
from mau.deck.card import UnoCard
from mau.enums import CardColor
from PIL import Image, ImageFile

# Путь к частям карты
ASSETS_PATh = Path("assets/")


def _add_color(base: Image.Image, color: CardColor) -> None:
    color_image = Image.open(ASSETS_PATh / "color" / f"{color.value}.png")
    base.alpha_composite(color_image)


def _add_sym(base: Image.Image, sym: Iterable[str | int]) -> None:
    paste_to = (48, 48)
    for s in sym:
        sym_image = Image.open(ASSETS_PATh / "sym" / f"{s}.png")
        y_offset = (64 - sym_image.size[1]) // 2
        base.paste(sym_image, (paste_to[0], paste_to[1] + y_offset))
        paste_to = (paste_to[0] + 8 + sym_image.size[0], paste_to[1])


def _add_reverse_sym(base: Image.Image, sym: Iterable[str | int]) -> None:
    paste_to = (271, 463)
    offset = 0
    for s in sym:
        sym_image = Image.open(ASSETS_PATh / "sym_reverse" / f"{s}.png")
        y_offset = (64 - sym_image.size[1]) // 2
        paste_coords = (
            paste_to[0] - sym_image.size[0] - offset,
            paste_to[1] - sym_image.size[1] - y_offset,
        )
        base.paste(sym_image, paste_coords)
        offset += 8 + sym_image.size[0]


def _add_glyph(base: Image.Image, color: str, glyph: str) -> None:
    glyph_image = Image.open(ASSETS_PATh / "glyph" / color / f"{glyph}.png")

    x = base.size[0] // 2 - (glyph_image.size[0] // 2)
    y = base.size[1] // 2 - (glyph_image.size[1] // 2)
    base.alpha_composite(glyph_image, (x, y))


# Главная функция
# ===============


def to_image(card: UnoCard, uncover: bool = False) -> io.BytesIO:
    """Собирает изображение из карты."""
    base: ImageFile.ImageFile | Image.Image = Image.open(ASSETS_PATh / "base.png")
    _add_color(base, card.color)

    draw_layer = Image.new(mode="RGBA", size=base.size)

    if isinstance(card.behavior, behavior.TurnBehavior):
        sym: list[str | int] = ["block"]
        glyph = "block"

    elif isinstance(card.behavior, behavior.ReverseBehavior):
        sym = ["reverse"]
        glyph = "reverse"

    elif isinstance(card.behavior, behavior.TakeBehavior):
        sym = ["plus", card.value]
        glyph = "take_2"

    elif isinstance(card.behavior, behavior.WildColorBehavior):
        sym = []
        glyph = None

    elif isinstance(card.behavior, behavior.WildTakeBehavior):
        sym = []
        glyph = "take_4"

    else:
        sym = [card.value]
        glyph = str(card.value)

    _add_sym(draw_layer, sym)
    _add_reverse_sym(draw_layer, sym)
    if glyph is not None:
        _add_glyph(draw_layer, str(card.color.value), glyph)

    base.alpha_composite(draw_layer)

    if uncover:
        base = base.point(lambda p: p // 1.5)

    buf = io.BytesIO()
    base.save(buf, format="PNG")
    buf.seek(0)

    return buf
