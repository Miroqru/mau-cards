# Mau;cards

<img src="./assets/logo.png" width="256"></img>

[![License](https://img.shields.io/badge/License-AGPL%20v3-red?style=flat&labelColor=%23B38B74&color=%23FF595F)](./LICENSE)
![Mau version](https://img.shields.io/badge/dynamic/toml?url=https%3A%2F%2Fcodeberg.org%2FSalormoon%2Fmauno%2Fraw%2Fbranch%2Fmain%2Fpyproject.toml&query=project.version&prefix=v&style=flat&label=Mau&labelColor=%23B38B74&color=%2373FFAD)
[![Docs](https://img.shields.io/badge/docs-miroq-%2300cc99?style=flat&labelColor=%23805959&color=%2330BFB3&link=https%3A%2F%2Fmau.miroq.ru%2Fdocs%2F)](https://mau.miroq.ru/docs/)
![GitHub stars](https://img.shields.io/github/stars/miroqru/mau-cards?style=flat&logo=github&logoColor=%23E6D0A1&label=Stars&labelColor=%23805959&color=%23FFF766)

**Mau;cards** - Небольшой сервис для сборки изображений карт для игры Mau.

**Особенности**:

- Быстрая **скорость** работы
- Поддержка нескольких **стилей**
- **Простота** использования
- **Кеширование** карт

<img src="https://mau.miroq.ru/card/next/0_1_1_number/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/next/2_2_1_take/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/next/3_2_1_reverse/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/next/4_2_1_turn/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/next/6_2_1_wild+color/cover" width="128"></img>

## Использование

Ссылка на изображение карты выглядит следующим образом: https://mau.miroq.ru/card/next/0_1_1_number/cover.

Её можно представить в виде шаблона:
`https://mau.miroq.ru/card/{style}/{card}/{filter}` где.
- `{style}`: Стиль карт, `progressive` или `next`.
- `{card}`: Представление карты в Mau. `{color}_{value}_{cost}_{type}`.
- `{filter}`: Фильтр, `cover` или `uncover`.


> Более подробно смотрите в [документации](https://mau.miroq.ru/docs/use/card_generator/).

## Стили карт

Генератор поддерживает несколько стилей карт, которые определяются в директории
`assets/`.
Сейчас существует два стиля карт:

### `next`
Появился после обновления Mau v2.0. Основной стиль для карт.

<img src="https://mau.miroq.ru/card/next/0_1_1_number/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/next/2_2_1_take/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/next/3_2_1_reverse/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/next/4_2_1_turn/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/next/6_2_1_wild+color/cover" width="128"></img>

Особые карты:

<img src="https://mau.miroq.ru/card/next/1_1_1_put/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/next/3_2_1_delta/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/next/5_2_1_twist/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/next/6_2_1_color/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/next/7_5_1_take/cover" width="128"></img>

### `progressive`

Появились после обновления Mau v1.5.
Самый первый вариант нарисованных карт.

<img src="https://mau.miroq.ru/card/progressive/0_1_1_number/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/progressive/2_2_1_take/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/progressive/3_2_1_reverse/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/progressive/4_2_1_turn/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/progressive/6_2_1_wild+color/cover" width="128"></img>

Особые карты:

<img src="https://mau.miroq.ru/card/progressive/1_1_1_put/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/progressive/3_2_1_delta/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/progressive/5_2_1_twist/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/progressive/6_2_1_color/cover" width="128"></img>
<img src="https://mau.miroq.ru/card/progressive/7_5_1_take/cover" width="128"></img>

## Установка

Если вы желаете развернуть локальный генератор карт:

- Клонируем репозиторий:

```sh
git clone https://github.com/miroqru/mau-cards
```

- Собираем сервер:

```sh
go build ./cmd/main.go
```

> Обратите внимание что директория `assets/` должна находиться рабом с
> исполняемым файлом.
