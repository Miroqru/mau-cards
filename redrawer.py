from pathlib import Path

from PIL import Image

COLORS = {
    (217, 54, 114): (204, 122, 163),
    (242, 61, 97): (230, 161, 162),
    (255, 89, 95): (255, 191, 179),
}


for i in Path("assets/glyph/0/").iterdir():
    image = Image.open(i)
    data = list(image.getdata())

    new_data = []
    for pixel in data:
        if len(pixel) < 4:
            new_data.append(pixel)
            continue

        new_color = COLORS.get(pixel[:3])
        if new_color is None:
            new_data.append(pixel)
            continue

        new_data.append(new_color + (pixel[3],))

    image.putdata(new_data)
    image.save(Path("out") / i.name)
