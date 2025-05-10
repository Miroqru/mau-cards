import requests
from mau.deck.presets import DeckGenerator
from tqdm import tqdm

card_gen = DeckGenerator.from_preset("classic")
for card in tqdm(card_gen._cards(), total=108):
    res = requests.get(f"http://127.0.0.1:3112/card/{card.pack()}/false")
