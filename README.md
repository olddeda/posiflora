# Posiflora — Telegram Integration MVP

Fullstack MVP: магазин подключает Telegram-бота и получает уведомления о новых заказах.

---

## Стек

### Backend (Go 1.25)

| Библиотека | Назначение |
|------------|-----------|
| `github.com/gin-gonic/gin` | HTTP-роутер и middleware |
| `gorm.io/gorm` | ORM |
| `gorm.io/driver/postgres` | PostgreSQL-драйвер для GORM |
| `github.com/swaggo/swag` | Генерация OpenAPI-спецификации |
| `github.com/swaggo/gin-swagger` | Serving swagger через Gin |
| `github.com/swaggo/files` | Статика для swagger UI |

### Frontend (Node 22)

| Библиотека | Назначение |
|------------|-----------|
| `react` + `react-dom` | UI-фреймворк |
| `react-router-dom` | Клиентский роутинг |
| `react-hook-form` | Управление формами |
| `zod` | Валидация схем форм |
| `@hookform/resolvers` | Интеграция zod с react-hook-form |
| `i18next` + `react-i18next` | Интернационализация (ru/en) |
| `i18next-browser-languagedetector` | Автоопределение языка |
| `tailwindcss` v4 | CSS utility-first |
| `vite` | Сборщик и dev-сервер |
| `typescript` | Типизация |
| `vitest` | Тест-раннер |
| `@testing-library/react` | Тестирование компонентов |
| `@testing-library/user-event` | Эмуляция действий пользователя |
| `jsdom` | DOM-среда для тестов |

### Инфраструктура

| Инструмент | Назначение |
|------------|-----------|
| PostgreSQL 16 | Основная БД |
| Docker + docker-compose | Локальный запуск всего стека |
| GitHub Actions | CI: build, lint, test, docker build |
| `act` | Локальный запуск GitHub Actions |

---

## Быстрый старт (Docker)

```bash
git clone https://github.com/olddeda/posiflora && cd posiflora
docker compose up --build
```

| Сервис | URL |
|--------|-----|
| Frontend | http://localhost:3000/shops/1/growth/telegram |
| Backend API | http://localhost:9090 |
| API Docs (Scalar) | http://localhost:9090/docs |
| PostgreSQL | localhost:5433 |

---

## Запуск без Docker

### Backend

```bash
# Требуется: Go 1.25+, PostgreSQL

createdb posiflora
psql posiflora -c "CREATE USER posiflora WITH PASSWORD 'posiflora';"
psql posiflora -c "GRANT ALL PRIVILEGES ON DATABASE posiflora TO posiflora;"

cp backend/.env.example backend/.env

cd backend
export $(cat .env | xargs)
go run ./cmd/server
```

Сервер стартует на `:8080`, автоматически применяет SQL-миграции и seed.

### Frontend

```bash
cd frontend
pnpm install
pnpm dev
```

Dev-сервер на http://localhost:5173, проксирует `/shops/*` на бэкенд.

---

## Переменные окружения

Файл `backend/.env.example`:

```env
DATABASE_URL=postgres://posiflora:posiflora@localhost:5432/posiflora?sslmode=disable
PORT=8080
ALLOWED_ORIGINS=http://localhost:5173
TELEGRAM_ENABLED=false
LOCALE=ru
LOCALES_DIR=locales
```

`TELEGRAM_ENABLED=false` — Telegram не вызывается, сообщения пишутся в stdout.
`TELEGRAM_ENABLED=true` — реальные вызовы Telegram Bot API.

---

## Миграции и seed

Миграции запускаются автоматически при старте backend через встроенный раннер.
SQL-файлы в `backend/migrations/`:

```
001_create_shops_table_up/down.sql
002_create_telegram_integrations_table_up/down.sql
003_create_orders_table_up/down.sql
004_create_telegram_send_logs_table_up/down.sql
005_seed_up/down.sql          — 1 магазин + 10 заказов
```

Применённые версии хранятся в таблице `schema_migrations`.

---

## API

Интерактивная документация: **http://localhost:8080/docs** (Scalar UI)

### POST `/shops/{shopId}/telegram/connect`

```json
{ "botToken": "123456:ABC...", "chatId": "987654321", "enabled": true }
```

### GET `/shops/{shopId}/telegram/status`

```json
{
  "enabled": true,
  "chatId": "****4321",
  "lastSentAt": "2024-01-15T10:30:00Z",
  "sentCount7d": 42,
  "failedCount7d": 1
}
```

### POST `/shops/{shopId}/orders`

```json
{ "number": "A-1005", "total": 2490, "customerName": "Анна" }
```

Ответ: `{ "order": {...}, "notifyStatus": "sent" | "failed" | "skipped" }`

---

## Тесты

### Backend (интеграционные — нужна PostgreSQL)

```bash
createdb posiflora_test
psql posiflora_test -c "CREATE USER posiflora WITH PASSWORD 'posiflora';"
psql posiflora_test -c "GRANT ALL PRIVILEGES ON DATABASE posiflora_test TO posiflora;"

cd backend
TEST_DATABASE_URL="postgres://posiflora:posiflora@localhost:5432/posiflora_test?sslmode=disable" \
  go test ./... -v
```

Или через Makefile:

```bash
make backend-test
```

Покрытие: `service`, `repository`, `handler` — 3 слоя.

### Frontend (unit + functional)

```bash
make frontend-test
make frontend-test-coverage
```

Покрытие: `apiFetch`, `cn`, UI-компоненты, формы с валидацией, i18n key parity.

---

## Swagger / API Docs

```bash
make backend-swagger   # регенерировать docs/ (нужен swag CLI)
```

Swagger CLI устанавливается через:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

После регенерации и запуска сервера: http://localhost:8080/docs

---

## Make-команды

```bash
make help               # список всех команд

make backend-build      # сборка бинарника
make backend-run        # запуск локально
make backend-test       # тесты
make backend-lint       # golangci-lint
make backend-fmt        # gofmt
make backend-swagger    # генерация OpenAPI-доков

make frontend-install   # pnpm install
make frontend-dev       # dev-сервер
make frontend-build     # production-сборка
make frontend-test      # vitest run
make frontend-test-coverage  # vitest с coverage

make docker-up          # запустить всё (db + backend + frontend)
make docker-build       # build + up
make docker-down        # остановить
make docker-clean       # остановить + удалить volumes
make docker-logs-backend
make docker-logs-frontend
make docker-logs-db
```

---

## Архитектура Backend

```
cmd/server/main.go
└── handler (Gin)
    ├── TelegramHandler  →  TelegramService  →  IntegrationRepository
    │                                        →  SendLogRepository
    └── OrderHandler     →  OrderService     →  OrderRepository
                                             →  TelegramClient (real / mock)
                                             →  i18n.Translator
```

Слои: `handler → service → repository`. Каждый слой изолирован через интерфейсы.

**Идемпотентность** — уникальный ключ `(shop_id, order_id)` в `telegram_send_logs` + `ON CONFLICT DO NOTHING`.

**Ошибка Telegram не роняет заказ** — ошибка логируется в `telegram_send_logs` со статусом `FAILED`, заказ возвращается клиенту.

---

## Допущения

| # | Что | Почему / Как улучшить в prod |
|---|-----|------------------------------|
| 1 | Авторизация отсутствует | MVP; добавить JWT/OAuth |
| 2 | Telegram отправка синхронная | Вынести в очередь (Redis Streams / RabbitMQ) с retry |
| 3 | Токен хранится в plaintext | Зашифровать AES-256 или вынести в Vault/KMS |
| 4 | Один язык уведомлений (LOCALE env) | Хранить locale per-shop |
| 5 | Тесты используют AutoMigrate вместо SQL-миграций | Перейти на testcontainers + реальные миграции |
