package app

// Реализовать (неделя 6): ручной DI-контейнер (без wire/fx).
// pgxpool → Redis-клиент (platform/pkg/redis) → repository(user, session) →
// service(iam) → api(auth, user) → grpc.NewServer(error-interceptor из internal/interceptor).
