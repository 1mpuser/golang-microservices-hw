<!-- TODO: Замени YOUR_USERNAME на свой GitHub username и YOUR_GIST_ID на ID своего gist для coverage badge -->
![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/YOUR_USERNAME/YOUR_GIST_ID/raw/coverage.json)

# Микросервисы на Go — Week 1/2

Три микросервиса: `order`, `inventory`, `payment`. Архитектура слоёв — Clean Architecture
(api → service → repository, входные/выходные адаптеры через интерфейсы).

## Команды

```bash
task build          # Собрать все модули
task lint           # Линтер
task mocks:gen      # Генерация моков (mockery)
task test           # Unit-тесты с race-детектором
task test:coverage  # Unit-тесты + покрытие (порог 40%)
task coverage:html  # HTML-отчёт покрытия
task test:api       # API-тесты
```
