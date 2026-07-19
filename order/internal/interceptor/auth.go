package interceptor

// Реализовать (неделя 6): gRPC client (outgoing) interceptor OrderService.
//   auth.SessionUUIDFromContext(ctx) → metadata.AppendToOutgoingContext(ctx, "session-uuid", uuid)
//   подключается к gRPC-клиенту InventoryService в app/di.go.
//   Без него Inventory не получит токен сессии и упадёт на Unauthenticated.
