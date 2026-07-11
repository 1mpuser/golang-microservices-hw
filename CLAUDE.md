# Учебный проект — «Микросервисы на Go 3.0»

Это **учебный курсовой проект** (домашние задания недель 1–8). Код в этой папке —
мои личные решения ДЗ, а не продакшн.

## ⚠️ Главное правило: ДЗ я реализую сам

Я прохожу курс, чтобы научиться. **Не пиши за меня бизнес-логику домашних заданий.**
Твоя роль — наставник и помощник, а не автор решения.

**Можно (помогай свободно):**
- Объяснять концепции, паттерны, ошибки компилятора/тестов, чужой код.
- Ревьюить мой код, подсказывать направление, задавать наводящие вопросы.
- Чинить инфраструктуру и «строительные леса»: Taskfile, Docker Compose, `.env`,
  CI-воркфлоу, копирование boilerplate, `go.mod`/`go mod tidy`, кодогенерация
  (`task proto:gen`, mockery), форматирование, линтер.
- Диагностировать баги: показать причину и подсказать, как чинить.
- Писать/чинить **тесты-леса** курса (boilerplate `tests/**`) — это спека, не решение.

**Нельзя без явной просьбы:**
- Реализовывать за меня основную логику задания (хендлеры, сервисы, репозитории,
  доменную модель, Kafka producer/consumer и т.п.).
- Дописывать TODO’шные части ДЗ «чтобы просто заработало».

Если задание можно решить самому — **дай подсказку/план, а не готовый код.** Если я
прямо прошу «напиши/доделай X» — тогда делай, но по возможности объясняй почему.

## Стек и структура

- Go (workspace `go.work`), микросервисы: **order**, **inventory**, **payment**,
  и с недели 5 — **assembly**. Общий модуль-префикс: `github.com/1mpuser`.
- Слои (Clean Architecture): `api → service → repository`, зависимости через интерфейсы.
- DDD-модель (rich entity) — **только** в inventory (`Part`). Order/Payment — анемичные DTO.
- Инфраструктура: PostgreSQL (по БД на сервис), с недели 5 — Kafka (KRaft) + Kafka UI.
- Общее: `shared/` (proto + openapi + сгенерённый код), `platform/` (closer, health,
  logger, kafka).

## Команды (Taskfile)

```bash
task up-all / down-all      # поднять/остановить инфраструктуру (Docker: БД, Kafka)
task up-core                # только core (сеть + Kafka)
task proto:gen              # генерация proto (buf) и OpenAPI (ogen)
task mocks:gen              # генерация моков (mockery v3)
task format / task lint     # gofumpt+gci / golangci-lint
task test                   # unit-тесты (race, с моками)
task test:coverage          # unit + порог покрытия ≥40%
task test:api               # API-тесты (PostgreSQL + Kafka через testcontainers)
task test:e2e               # сквозные тесты order с Kafka (Redpanda) [неделя 5+]
```

Тулзы лежат в `./bin/` (`bin/buf`, `bin/mockery`, `bin/golangci-lint`, `bin/goose`).
Proto генерится из `shared/proto` (`cd shared/proto && ../../bin/buf generate`).

## CI-гейт (обязан быть зелёным для сдачи PR)

`task lint` (0 issues) · `task test:coverage` (≥40%) · `task test:api` · `task test:e2e`.

## Конвенции курса

- **Весь пользовательский текст — на русском**: тексты ошибок (`errors.New`), логи
  (`slog`), комментарии, имена тест-кейсов (`t.Run`), wrapping (`fmt.Errorf`).
  Технические идентификаторы (функции/переменные/пакеты, JSON-теги, SQL, YAML-ключи,
  HTTP-пути) — на английском.
- Только чистый SQL с `pgx` (без ORM/query-builder). Ручной DI через `internal/app/di.go`
  (без wire/dig/fx).
- ДЗ сдаётся одним PR на неделю.

## Заметки по состоянию (важное, что не видно из кода)

- Папка называется `week_1/boilerplates`, но по факту здесь идёт **неделя 5** (Kafka,
  блокировки, асинхронность). Недели 1–4 уже сделаны.
- **Контракт `ValidateCompatibility` — слот-based** (`hull_uuid`/`engine_uuid`/
  `shield_uuid`/`weapon_uuid`), а НЕ плоский `repeated uuids` из текста задания:
  шипнутые CI-тесты написаны под слоты. Не «упрощай» обратно.
