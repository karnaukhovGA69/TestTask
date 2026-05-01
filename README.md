# URL Shortener

HTTP-сервис для сокращения ссылок, написанный на Go. Поддерживает два хранилища: in-memory и PostgreSQL.

## Требования

- Go 1.22+
- Docker и Docker Compose (для запуска через контейнер)

## Локальный запуск

### In-memory хранилище

```bash
go run ./cmd/app dbelg
```

### PostgreSQL хранилище

Создайте файл `.env` на основе `.env.example` и запустите:

```bash
cp .env.example .env
go run ./cmd/app postgres
```

## Запуск через Docker Compose

```bash
cp .env.example .env
docker compose up --build
```

Сервис будет доступен на `http://localhost:8080`.

## API

### POST /url — сохранить ссылку

```bash
curl -X POST http://localhost:8080/url \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com"}'
```

Ответ:
```json
{"shortURL":"http://localhost:8080/aBcD123456"}
```

### GET /{shortURL} — получить оригинальную ссылку

```bash
curl http://localhost:8080/aBcD123456
```

Ответ:
```json
{"URL":"https://example.com"}
```

## Переменные окружения

| Переменная   | Описание              | Пример     |
|--------------|-----------------------|------------|
| DB_HOST      | Хост PostgreSQL       | postgres   |
| DB_PORT      | Порт PostgreSQL       | 5432       |
| DB_USER      | Пользователь БД       | admin      |
| DB_PASSWORD  | Пароль БД             | adminadmin |
| DB_NAME      | Имя базы данных       | shorturl   |
| DB_SSLMODE   | Режим SSL             | disable    |

## Тестирование

```bash
go test ./...
```
