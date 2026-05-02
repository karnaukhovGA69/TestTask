# Project Review: URL Shortener

## 1. Краткий вывод

Проект частично соответствует ТЗ. Основная функциональность реализована: есть HTTP API, генератор коротких ссылок, in-memory storage, PostgreSQL storage, Dockerfile, docker-compose.yml, README и часть unit-тестов.

Отправлять как тестовое задание пока рискованно. Главные проблемы перед сдачей:

- Docker Compose не готов к запуску из чистого checkout: `README.md` требует `.env.example`, но файла `.env.example` нет, а `docker compose config` подставляет пустые значения для `DB_USER`, `DB_PASSWORD`, `DB_NAME`.
- Тесты есть, но покрытие неполное: нет тестов PostgreSQL storage, а tests для HTTP handler проверяют тестовую копию handler-логики, а не production `Handler`.
- HTTP API слабо различает ошибки: любая ошибка `POST` возвращается как `400`, любая ошибка `GET` возвращается как `404`.
- Конфигурация и Docker требуют доработки перед публикацией: `.dockerignore` не исключает `.env`, README ссылается на отсутствующий `.env.example`, модуль называется `main`.

### Команды проверки

- `go fmt ./...`: не запускался, потому что команда может менять файлы, а условие ревью запрещает менять код. Вместо этого выполнен dry-run: `find . -name '*.go' -print0 | xargs -0 gofmt -l`, вывод пустой, неформатированных Go-файлов не найдено.
- `go vet ./...`: выполнено успешно.
- `go test ./...`: выполнено успешно. Пакеты `main/cmd/app` и `main/internal/storage/postgres` без test files.
- `go mod tidy`: не запускался, потому что команда может менять `go.mod`/`go.sum`. Вместо этого выполнено `go mod tidy -diff`, diff отсутствует.
- `docker compose config`: команда завершилась с кодом 0, но вывела предупреждения, что `DB_PASSWORD`, `DB_USER`, `DB_NAME` не заданы. В итоговом config эти значения пустые.
- `docker build -t url-shortener .`: выполнено успешно.
- `docker compose up --build`: реальный запуск не выполнялся, потому что в проекте нет `.env` и `.env.example`, а compose-файл монтирует `./.env`. Выполнен `docker compose --dry-run up --build`: dry-run прошел, но повторил предупреждения о незаданных env-переменных.
- Дополнительно проверен Docker-образ без PostgreSQL: `docker run --rm --name review-url-shortener-dbelg -p 18080:8080 url-shortener dbelg` запустил HTTP-сервис, `POST /url` вернул `200 OK`.

## 2. Проверка соответствия ТЗ

| Требование | Статус | Комментарий | Файлы |
|---|---|---|---|
| 1. Сокращенная ссылка должна быть уникальной | Выполнено | Генератор делает случайный short code, in-memory storage проверяет `shortToLong`, PostgreSQL имеет `UNIQUE` на `shortURL` и повторяет генерацию при `23505`. | `shorturl/shorturl.go:8`, `shorturl/shorturl.go:11`, `internal/storage/dbelg/dbelg.go:49`, `internal/storage/postgres/postgres.go:71`, `migrations/init.sql:4` |
| 2. На один оригинальный URL должна ссылаться только одна сокращенная ссылка | Выполнено | In-memory storage сначала проверяет `longToShort`; PostgreSQL сначала делает `SELECT`, а при `UNIQUE` по `longURL` повторно читает существующий short URL. | `internal/storage/dbelg/dbelg.go:45`, `internal/storage/postgres/postgres.go:63`, `internal/storage/postgres/postgres.go:81` |
| 3. Short code длиной ровно 10 символов | Выполнено | `length = 10`, результат создается как `make([]byte, length)`. | `shorturl/shorturl.go:9`, `shorturl/shorturl.go:12` |
| 4. Short code состоит только из `a-z`, `A-Z`, `0-9`, `_` | Выполнено | Алфавит соответствует ТЗ. | `shorturl/shorturl.go:8`, `shorturl/shorturl.go:15` |
| 5. HTTP POST сохраняет URL и возвращает сокращенный, GET возвращает оригинальный URL | Частично выполнено | Роуты есть и логика работает, но API слабо различает ошибки и всегда возвращает полный URL с `http://`. | `cmd/app/main.go:39`, `cmd/app/main.go:40`, `internal/handler/handlers.go:19`, `internal/handler/handlers.go:40` |
| 6. Сервис написан на Go | Выполнено | Основной код сервиса написан на Go. | `cmd/app/main.go`, `internal/...`, `shorturl/shorturl.go` |
| 7. Сервис распространяется в виде Docker-образа | Частично выполнено | Dockerfile собирается и образ стартует в `dbelg`-режиме, но Docker Compose из чистого checkout не готов из-за отсутствующего `.env.example`/`.env`. | `Dockerfile:1`, `Dockerfile:19`, `docker-compose.yml:1`, `README.md:20` |
| 8. Есть PostgreSQL и in-memory storage | Выполнено | Обе реализации есть и соответствуют интерфейсу `storage.DB`. | `internal/storage/storage.go:3`, `internal/storage/dbelg/dbelg.go:11`, `internal/storage/postgres/postgres.go:15` |
| 9. Выбор storage задается параметром запуска | Выполнено | `os.Args[1]` передается в `storage.MakeDB`, где выбираются `dbelg` или `postgres`. | `cmd/app/main.go:24`, `cmd/app/main.go:29`, `internal/storage/factory.go:11` |
| 10. Функциональность покрыта unit-тестами | Частично выполнено | Тесты добавлены для генератора, dbelg, service, factory и handler-сценариев. Нет тестов PostgreSQL storage, а handler tests не вызывают production `Handler`. | `shorturl/shorturl_test.go`, `internal/storage/dbelg/dbelg_test.go`, `internal/service/service_test.go`, `internal/handler/handlers_test.go`, `internal/storage/factory_test.go` |
| 11. Проект готов к публикации в публичном GitHub | Частично выполнено | README есть, `.gitignore` есть, но отсутствует `.env.example`, `.dockerignore` не исключает `.env`, модуль называется `main`, Docker Compose требует локальный env-файл. | `README.md:20`, `.gitignore:1`, `.dockerignore:1`, `go.mod:1`, `docker-compose.yml:30` |

## 3. Критические ошибки

### Ошибка 1: Docker Compose не запускается корректно из чистого checkout

- Где находится: `docker-compose.yml:7`, `docker-compose.yml:8`, `docker-compose.yml:9`, `docker-compose.yml:30`, `README.md:20`
- Почему это критично: требование про Docker-развертывание закрыто неполно. Пользователь следует README, но `.env.example` в проекте отсутствует, а без env-переменных Compose подставляет пустые значения для PostgreSQL.
- Как проявится: `docker compose config` выводит предупреждения `The "DB_PASSWORD" variable is not set`, `The "DB_USER" variable is not set`, `The "DB_NAME" variable is not set`; в итоговой конфигурации `POSTGRES_DB`, `POSTGRES_PASSWORD`, `POSTGRES_USER` пустые.
- Как исправить: добавить `.env.example` с `DB_HOST=postgres`, `DB_PORT=5432`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`; либо явно передавать env для `app` и `postgres` через `env_file`.

### Ошибка 2: HTTP handler tests не тестируют production `Handler`

- Где находится: `internal/handler/handlers_test.go:30`, `internal/handler/handlers_test.go:47`
- Почему это критично: требование про unit-тесты формально частично выполнено, но важный слой HTTP API проверяется через `testHandler`, который дублирует production-логику вместо вызова `Handler.PostHandler` и `Handler.GetHandler`.
- Как проявится: production `Handler` можно сломать, а handler tests останутся зелеными, потому что тестируют отдельную копию методов `post` и `get`.
- Как исправить: сделать handler зависимым от интерфейса service-слоя и тестировать настоящий `NewHandler(...).PostHandler` / `GetHandler`, либо использовать реальный `service.Service` с mock storage.

### Ошибка 3: PostgreSQL storage не покрыт unit/integration tests

- Где находится: `internal/storage/postgres/postgres.go:19`, отсутствие `internal/storage/postgres/*_test.go`
- Почему это критично: PostgreSQL storage содержит наиболее рискованную логику проекта: подключение, SQL, `UNIQUE`-конфликты, `sql.ErrNoRows`, но она не проверяется автоматическими тестами.
- Как проявится: регрессии в SQL-запросах, схеме или обработке конфликтов могут пройти `go test ./...`.
- Как исправить: добавить integration tests через testcontainers или отдельную тестовую БД; минимум проверить `AddURL`, `GetShortURL`, `GetLongURL`, `ErrNotFound`, повторный original URL и concurrent insert одного URL.

### Ошибка 4: `.env` может попасть в Docker build context

- Где находится: `.dockerignore:1`, `.gitignore:1`
- Почему это критично: `.env` игнорируется git, но не игнорируется Docker. Если разработчик создаст `.env`, он попадет в build context при `COPY . .`.
- Как проявится: секреты могут попасть в build context, Docker cache или слои builder stage.
- Как исправить: добавить `.env`, `.env.*`, кроме безопасного `.env.example`, в `.dockerignore`.

## 4. Места, где могут возникнуть ошибки

### Потенциальная проблема 1: `POST` возвращает `400` для любой ошибки service/storage

- Где находится: `internal/handler/handlers.go:48`
- В каком сценарии возникнет: PostgreSQL недоступен, генератор вернул ошибку, storage вернул внутреннюю ошибку.
- Чем опасно: клиент получит `400 Bad Request`, хотя проблема на стороне сервера.
- Как лучше исправить: ввести sentinel errors для validation/not found и возвращать `500` для инфраструктурных ошибок.

### Потенциальная проблема 2: `GET` возвращает `404` для любой ошибки service/storage

- Где находится: `internal/handler/handlers.go:26`
- В каком сценарии возникнет: ошибка соединения с PostgreSQL, SQL-ошибка, внутренняя ошибка storage.
- Чем опасно: реальные ошибки сервера маскируются под “не найдено”.
- Как лучше исправить: использовать `errors.Is(err, storage.ErrNotFound)` и различать `404` и `500`.

### Потенциальная проблема 3: DSN собирается из env без валидации

- Где находится: `internal/storage/factory.go:17`
- В каком сценарии возникнет: нет `.env`, не задан `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_NAME`, `DB_SSLMODE`.
- Чем опасно: приложение падает на `db.Ping()`, а лог в `cmd/app/main.go:31` не показывает исходную ошибку.
- Как лучше исправить: валидировать обязательные переменные до сборки DSN и логировать `zap.Error(err)`.

### Потенциальная проблема 4: `godotenv.Load()` игнорирует ошибку

- Где находится: `cmd/app/main.go:18`
- В каком сценарии возникнет: локальный запуск PostgreSQL без `.env`.
- Чем опасно: приложение продолжит стартовать с пустыми env-переменными и упадет позже на подключении к БД.
- Как лучше исправить: для `postgres`-режима явно проверять обязательные env; для `dbelg`-режима `.env` не нужен.

### Потенциальная проблема 5: полный short URL всегда строится с `http://`

- Где находится: `internal/handler/handlers.go:53`
- В каком сценарии возникнет: запуск за HTTPS reverse proxy или ingress.
- Чем опасно: клиент получит `http://...`, хотя внешний сервис доступен по HTTPS.
- Как лучше исправить: возвращать только short code или настраивать base URL через конфиг.

### Потенциальная проблема 6: фиксированный лимит `69` попыток генерации

- Где находится: `internal/storage/dbelg/dbelg.go:9`, `internal/storage/postgres/postgres.go:13`
- В каком сценарии возникнет: большое количество ссылок или частые коллизии.
- Чем опасно: сервис может вернуть ошибку, хотя свободные short codes еще есть.
- Как лучше исправить: сделать лимит конфигурируемым и добавить метрики/логирование отказов по коллизиям.

### Потенциальная проблема 7: PostgreSQL connection не закрывается при остановке сервиса

- Где находится: `internal/storage/postgres/postgres.go:108`, `cmd/app/main.go:17`
- В каком сценарии возникнет: graceful shutdown, тесты с поднятой БД, долгоживущий сервис.
- Чем опасно: нет контролируемого освобождения ресурсов.
- Как лучше исправить: добавить graceful shutdown и вызывать `Close()` для storage, если оно поддерживает закрытие.

### Потенциальная проблема 8: README требует `.env.example`, которого нет

- Где находится: `README.md:20`, `README.md:30`
- В каком сценарии возникнет: первый запуск проекта новым пользователем.
- Чем опасно: инструкция запуска не работает.
- Как лучше исправить: добавить `.env.example` или переписать инструкцию под явные env-переменные.

### Потенциальная проблема 9: Go version в README не совпадает с `go.mod`

- Где находится: `README.md:7`, `go.mod:3`
- В каком сценарии возникнет: пользователь с Go 1.22 попытается запустить проект.
- Чем опасно: README обещает Go 1.22+, но `go.mod` требует `go 1.25.1`.
- Как лучше исправить: синхронизировать версию Go в README, Dockerfile и `go.mod`.

## 5. Проверка генерации short URL

- Длина: соответствует ТЗ. `length = 10`, результат создается длиной 10 байт.
- Алфавит: соответствует ТЗ. Используются `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_`.
- Префиксы: внутри short code нет `http://`; префикс добавляется только в HTTP-ответе.
- Алгоритм: надежный для учебного сервиса, используется `crypto/rand`.
- Уникальность: сам генератор не проверяет глобальную уникальность, это правильно делается на уровне storage.
- Коллизии: обработаны в `dbelg.AddURL` и `postgres.AddURL`.
- Бесконечный цикл: отсутствует, циклы ограничены `maxGenerateAttempts`.
- Тесты: есть тесты длины, алфавита и вероятностный тест уникальности в `shorturl/shorturl_test.go`.

## 6. Проверка in-memory storage

- `longToShort` и `shortToLong` используются корректно.
- `AddURL` сначала проверяет `longToShort`, поэтому повторный original URL возвращает тот же short URL.
- `GetShortURL` читает из `longToShort`, `GetLongURL` читает из `shortToLong`.
- Потокобезопасность есть: используется `sync.RWMutex`.
- Запись защищена `Lock`, чтение защищено `RLock`.
- Параллельный сценарий частично покрыт тестом `TestDBelg_Concurrent`.
- Ошибки not found есть, но представлены обычным `errors.New`, без общего sentinel error.
- Неиспользуемые файлы `longToShort.json` и `shortToLong.json` лежат рядом с in-memory storage, но код их не читает и не пишет.

## 7. Проверка PostgreSQL storage

- Подключение реализовано через `database/sql` и `github.com/lib/pq`.
- `NewPostgresDB` вызывает `sql.Open`, затем `Ping`, при ошибке `Ping` закрывает соединение.
- DSN собирается в `storage.MakeDB` из env-переменных.
- SQL-запросы параметризованы через `$1`, `$2`, SQL-инъекций в найденном коде нет.
- `sql.ErrNoRows` обрабатывается в `GetShortURL`, `GetLongURL`, `AddURL`.
- `UNIQUE`-конфликты обрабатываются через `pq.Error.Code == "23505"`.
- При `UNIQUE` по `longURL` код повторно читает существующий `shortURL`.
- Таблица соответствует базовым требованиям: `longURL TEXT NOT NULL UNIQUE`, `shortURL VARCHAR(10) NOT NULL UNIQUE`.
- Недостаток: нет PostgreSQL tests.
- Недостаток: логика завязана на имя constraint `urls_longurl_key`; это имя генерируется PostgreSQL из текущей схемы, но при изменении схемы может сломаться.
- Недостаток: `Close()` есть, но на уровне приложения не вызывается.

## 8. Проверка HTTP API

- `POST /url` зарегистрирован в `cmd/app/main.go:40`.
- `GET /{shortURL}` зарегистрирован в `cmd/app/main.go:39`.
- `POST` принимает JSON вида `{"url":"https://example.com"}`.
- Плохой JSON возвращает `400 Bad Request`.
- Пустой URL отсекается через service layer и возвращает `400 Bad Request`.
- Успешный `POST` возвращает JSON с ключом `shortURL`.
- Успешный `GET` возвращает JSON с ключом `URL`.
- Несуществующая короткая ссылка возвращает `404 Not Found`.
- Проблема: ошибки storage не различаются по типам, поэтому возможны неверные HTTP status codes.
- Проблема: API возвращает полный URL с `http://`, а не только short code; для HTTPS/deploy это может быть неудобно.
- Проблема: JSON-ключи `shortURL` и `URL` стилистически несогласованы.

## 9. Проверка Docker

- Dockerfile собирается успешно.
- Собранный образ запускается в `dbelg`-режиме и отвечает на `POST /url`.
- Docker Compose config формально валиден, но без `.env` содержит пустые переменные PostgreSQL.
- В репозитории нет `.env.example`, хотя README требует его создать через `cp .env.example .env`.
- `docker-compose.yml` монтирует `./.env` в app-контейнер, но сам файл отсутствует в чистом checkout.
- PostgreSQL volume настроен: `postgres_data:/var/lib/postgresql/data`.
- `init.sql` подключается в `/docker-entrypoint-initdb.d/init.sql`.
- `init.sql` создает таблицу `urls`; важно, что он выполняется только при первом создании volume.
- Внутри Compose для PostgreSQL должен использоваться host `postgres`; README это документирует в таблице env.
- `localhost` для DB внутри Docker Compose в найденных файлах не используется.
- `.dockerignore` не исключает `.env`; это нужно исправить перед публикацией.

## 10. Проверка тестов

- Уже есть:
  - `shorturl/shorturl_test.go`: длина, алфавит, вероятностная уникальность.
  - `internal/storage/dbelg/dbelg_test.go`: add/get/not found/idempotency/concurrency.
  - `internal/service/service_test.go`: пустые значения, valid path, storage error.
  - `internal/storage/factory_test.go`: `dbelg`, case-insensitive names, unknown storage.
  - `internal/handler/handlers_test.go`: handler-like scenarios, но через тестовую копию handler-кода.
- Не хватает:
  - PostgreSQL storage tests.
  - Tests для production `Handler`.
  - Tests для `storage.MakeDB("postgres")` с env validation или test DB.
  - Tests для Docker/compose можно оставить как e2e smoke-check в CI.
  - Tests для error mapping: storage internal error должен давать `500`, not found должен давать `404`.
- Обязательно добавить:
  - `postgres.AddURL` повторный original URL.
  - `postgres.AddURL` конкурентный одинаковый original URL.
  - `postgres.GetLongURL` / `GetShortURL` not found.
  - Production HTTP handler tests без дублирования логики.
  - Test, что bad JSON и пустой URL не вызывают storage.

## 11. Рекомендации по улучшению

### Обязательно исправить

- Добавить `.env.example` и синхронизировать его с `docker-compose.yml`, README и `storage.MakeDB`.
- Добавить `.env` в `.dockerignore`.
- Переписать handler tests так, чтобы они тестировали production `Handler`.
- Добавить PostgreSQL storage tests.
- Разделить ошибки validation/not found/internal error и возвращать корректные HTTP status codes.
- Синхронизировать Go version в `README.md`, `go.mod`, `Dockerfile`.

### Желательно исправить

- Логировать исходную ошибку в `cmd/app/main.go:31` через `zap.Error(err)`.
- Валидировать env-переменные до подключения к PostgreSQL.
- Возвращать только short code или настраивать base URL через env.
- Добавить graceful shutdown HTTP-сервера.
- Вызывать `Close()` у PostgreSQL storage при завершении приложения.
- Переименовать `DBelg` во что-то понятное, например `MemoryStorage`.
- Удалить или объяснить `longToShort.json` и `shortToLong.json`.

### Можно улучшить позже

- Переименовать модуль из `main` в имя GitHub-модуля.
- Добавить CI: `go fmt` check, `go vet`, `go test`, `docker build`.
- Сделать OpenAPI/README-описание API более строгим.
- Использовать единый JSON-формат ошибок.
- Добавить health endpoint для приложения.
- Добавить миграционный инструмент вместо одного `init.sql`.

## 12. Итоговый чеклист перед сдачей

- [ ] go fmt ./...
- [x] go vet ./...
- [x] go test ./...
- [x] docker build
- [ ] docker compose up
- [x] POST работает
- [ ] GET работает
- [ ] PostgreSQL работает
- [x] memory storage работает
- [x] README написан
- [ ] проект залит на GitHub

Примечания к чеклисту:

- `go fmt ./...` не запускался из-за запрета менять код; dry-run форматирования чистый.
- `docker compose up` не запускался как реальный старт из-за отсутствующего `.env`/`.env.example`; dry-run показывает, что env-переменные PostgreSQL пустые.
- `POST работает` подтвержден через Docker-образ в `dbelg`-режиме на временном порту `18080`.
- `GET работает` для текущего запуска отдельно не проверялся.
- `PostgreSQL работает` через реальный Compose-старт не проверялся из-за отсутствующего env-файла в текущей ветке.
