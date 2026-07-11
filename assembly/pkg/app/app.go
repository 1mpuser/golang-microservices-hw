// Package app — тонкая обёртка над internal/app для e2e/интеграционных тестов
// (по аналогии с order/pkg/app и inventory/pkg/app).
package app

// TODO(неделя 5): экспортировать сборку зависимостей AssemblyService для тестов —
// например функцию запуска consumer'а в процессе теста с переданными brokers/топиками,
// чтобы e2e-тесты order могли поднять всю цепочку событий in-process.
