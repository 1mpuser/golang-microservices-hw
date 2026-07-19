package session

// Реализовать (неделя 6): SessionRepository на Redis (go-redis v9, platform/pkg/redis).
//   Create: HSET session:{uuid} (поля из redisview) + EXPIRE (TTL 24h)
//   Get:    HGETALL session:{uuid}; пусто → errs.ErrSessionNotFound
//   Delete: DEL session:{uuid} (идемпотентно)
