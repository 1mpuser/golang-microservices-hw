package user

// Реализовать (неделя 6): UserRepository на PostgreSQL (pgx, чистый SQL).
// Конструктор NewRepository(db). Методы — в отдельных файлах: create, get_by_login, get_by_uuid.

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool, txGetter *trmpgx.CtxGetter) *repository {
	return &repository{
		pool,
	}
}

// New — конструктор для тестового харнесса: получает менеджер транзакций и использует
// стандартный trmpgx CtxGetter (репозиторий работает через транзакции trm).
func New(pool *pgxpool.Pool) *repository {
	return &repository{
		pool,
	}
}
