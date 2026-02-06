# URL Shortener

Сервис сокращения ссылок на Go. Поддерживает два хранилища: in-memory и PostgreSQL. Возвращает короткий URL для одного и того же `long_url`. Сервис полностью покрыт тестами.

## Возможности
- REST API: создание коротких ссылок, получение оригинальной, редирект по алиасу
- Хранилища: `memory` и `postgresql`
- Генерация алиасов фиксированной длины (10 символов)
- Запуск через Docker

## Быстрый старт (Docker)

### Вариант 1: in-memory

1) Подготовьте окружение:

- Создать .env файл, можно скопировать .env.example
- Прописать в .env
```bash
STORAGE=memory
```

2) Запустите приложение (запускает только сервис не поднимая postgres):
```bash
docker compose -f docker-compose-memory.yml up --build   
```

### Вариант 2: PostgreSQL


1) Подготовьте окружение:

- Создать .env файл, можно скопировать .env.example
- Прописать в .env
```bash
STORAGE=postgresql
```

2) Запустите приложение:

```bash
docker compose up -d --build
```

Сервис будет доступен по http://localhost:8081 (или по значению HTTP_PORT в .env).

## Локальный запуск без Docker

### Вариант 1: in-memory

1) Подготовьте окружение:

- Создать .env файл, можно скопировать .env.example
- Прописать в .env
```bash
STORAGE=memory
```

2) Запустите приложение:

```bash
go run ./cmd/url-shortener
```

### Вариант 2: PostgreSQL

1) Подготовьте окружение:

- Создать .env файл, можно скопировать .env.example
- Прописать в .env
```bash
STORAGE=postgresql
```

2) Поднимите Postgres (можно через compose):

```bash
docker compose up -d postgres
```

3) Примените миграции:

```bash
docker compose up migrate
```

4) Запустите приложение локально:

```bash
go run ./cmd/url-shortener
```

## Конфигурация

Все настройки читаются из `.env` (пример — `.env.example`).

Ключевые параметры:
- `BASE_URL` — базовый URL, который будет возвращаться в поле `short_url`.
- `STORAGE` — `postgresql` или `memory`.
- `HTTP_HOST`, `HTTP_PORT` — адрес и порт HTTP-сервера.
- `DB_*` — параметры подключения к Postgres (нужны только при `STORAGE=postgresql`).
- `ALIAS_SECRET` — секрет для генератора алиасов (смешивается с ID).

## API

Базовый URL: `http://localhost:8081` (или ваш `HTTP_HOST:HTTP_PORT`).

### Создать короткую ссылку

`POST /api`

```bash
curl -X POST http://localhost:8081/api \
  -H 'Content-Type: application/json' \
  -d '{"long_url": "https://example.com"}'
```

Ответ:

```json
{"short_url":"http://localhost:8081/aaacy0kMHk"}
```

Особенности:
- `long_url` должен быть валидным URL с схемой (`http://` или `https://`).
- Для уже существующего `long_url` вернется тот же алиас.

### Получить оригинальную ссылку по алиасу

`GET /api/:alias`

```bash
curl http://localhost:8081/api/aaacy0kMHk
```

Ответ:

```json
{"long_url":"https://example.com"}
```

### Редирект

`GET /:alias` — ответ 302 и редирект на оригинальный URL.

```bash
curl -i http://localhost:8081/aaacy0kMHk
```

### Ошибки

Формат ошибок:

```json
{"error":"invalid input"}
```

Коды:
- `400` - некорректный ввод
- `404` - алиас не найден
- `409` - конфликт алиаса
- `500` - внутренняя ошибка

## Тесты

Полный прогон тестов с тестовой БД:

```bash
make test
```

Что делает `make test`:
- поднимает тестовый Postgres через `docker-compose-test.yml`
- применяет миграции
- запускает `go test ./...`
- останавливает и удаляет контейнеры

Тестовая БД поднимается на `localhost:54329`.
