package interceptor

// Реализовать (неделя 6): gRPC server (incoming) interceptor InventoryService (контракт «gRPC Interceptor»).
//   metadata.FromIncomingContext → нет metadata → Unauthenticated "отсутствует metadata"
//   ключ "session-uuid" отсутствует → Unauthenticated "отсутствует session-uuid"
//   AuthService.Whoami(session_uuid); ошибка → Unauthenticated "недействительная сессия"
//   успех → auth.WithUserUUID → handler. Защищает все 6 методов Inventory.
