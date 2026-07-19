package middleware

// Реализовать (неделя 6): HTTP auth-middleware OrderService (контракт «HTTP Middleware»).
//   читает "Authorization: Bearer <session_uuid>"
//   нет заголовка → 401 "отсутствует заголовок Authorization"
//   не Bearer     → 401 "неверный формат Authorization"
//   AuthService.Whoami(session_uuid); ошибка → 401 "недействительная сессия"
//   успех → auth.WithUserUUID + auth.WithSessionUUID (platform/pkg/auth) → next
//   Защищает все эндпоинты /api/v1/orders*.
