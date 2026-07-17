package order

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool     *pgxpool.Pool
	txGetter *trmpgx.CtxGetter
}

func NewRepository(pool *pgxpool.Pool, txGetter *trmpgx.CtxGetter) *repository {
	return &repository{
		pool,
		txGetter,
	}
}

// New — конструктор для тестового харнесса: получает менеджер транзакций и использует
// стандартный trmpgx CtxGetter (репозиторий работает через транзакции trm).
func New(pool *pgxpool.Pool, _ *trmmanager.Manager) *repository {
	return &repository{
		pool,
		trmpgx.DefaultCtxGetter,
	}
}
