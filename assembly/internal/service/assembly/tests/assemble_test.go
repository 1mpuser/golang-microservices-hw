package assembly_test

// Заметка (неделя 5): unit-тесты сервиса сборки —
//  - успешная сборка публикует ShipAssembled с корректными полями (user_uuid проброшен);
//  - ошибка producer'а возвращается наверх.
// Для мока ShipAssembledProducer добавь интерфейс в .mockery.yaml и сгенерируй (task mocks:gen).
