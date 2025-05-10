"""Главный файл.

Настраивает сервер и запускает его.
"""

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse, StreamingResponse
from loguru import logger
from mau.deck.card import UnoCard
from redis.asyncio import Redis

from mau_cards.generator import to_image

app = FastAPI(title="mau:cards", version="v1.0", root_path="/card")
redis = Redis()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/{card}/{cover}")
async def get_card(card: str, cover: bool):
    # key = f"{card}/{cover}"
    # cache_image = await redis.get(key)
    # if cache_image is not None:
    #     logger.debug("From cache")
    #     return StreamingResponse(io.BytesIO(cache_image), media_type="image/png")

    uno_card = UnoCard.unpack(card)
    if uno_card is None:
        return FileResponse("assets/base.png")

    logger.info("Render {}", uno_card)
    image = to_image(uno_card, cover)
    # await redis.set(key, image.getvalue())
    return StreamingResponse(image, media_type="image/png")
