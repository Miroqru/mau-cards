"""Главный файл.

Настраивает сервер и запускает его.
"""

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse, StreamingResponse
from loguru import logger
from mau.deck.card import UnoCard

from mau_cards.generator import to_image

app = FastAPI(title="mau:cards", version="v1.0", root_path="/card")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/{card}/{cover}")
async def get_card(card: str, cover: bool):
    uno_card = UnoCard.unpack(card)
    if uno_card is None:
        return FileResponse("assets/base.png")

    image = to_image(uno_card, cover)
    logger.debug(image)
    return StreamingResponse(image, media_type="image/png")
