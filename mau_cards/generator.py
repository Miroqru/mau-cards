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


def _uncover(base: Image.Image) -> Image.Image:
    res = Image.new("RGBA", base.size)
    for py in range(base.size[1]):
        for px in range(base.size[0]):
            r, g, b, a = base.getpixel((px, py))  # type: ignore

            tr = round(r / 1.5)
            tg = round(g / 1.5)
            tb = round(b / 1.5)
            res.putpixel((px, py), (tr, tg, tb, a))

    return res


# Главная функция
# ===============


def to_image(card: UnoCard, uncover: bool = False) -> io.BytesIO:
    """Собирает изображение из карты."""
    base: ImageFile.ImageFile | Image.Image = Image.open(ASSETS_PATh / "base.png")
    _add_color(base, card.color)

    draw_layer = Image.new(mode="RGBA", size=base.size)

    if isinstance(card.behavior, behavior.TurnBehavior):
        sym: list[str | int] = ["block"]

    elif isinstance(card.behavior, behavior.ReverseBehavior):
        sym = ["reverse"]

    elif isinstance(card.behavior, behavior.TakeBehavior):
        sym = ["plus", card.value]

    elif isinstance(card.behavior, behavior.ColorBehavior):
        sym = []

    elif isinstance(card.behavior, behavior.ColorTakeBehavior):
        sym = []

    else:
        sym = [card.value]

    _add_sym(draw_layer, sym)
    _add_reverse_sym(draw_layer, sym)

    base.alpha_composite(draw_layer)

    if uncover:
        base = _uncover(base)

    buf = io.BytesIO()
    base.save(buf, format="PNG")
    buf.seek(0)

    return buf
