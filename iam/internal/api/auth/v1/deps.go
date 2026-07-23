package v1

import (
	"context"

	"github.com/1mpuser/iam/internal/model"
	"github.com/1mpuser/iam/internal/service/input"
)

// IAmService — то, что gRPC-хендлеры AuthService требуют от сервисного слоя.
// Только методы этого API (Login, Whoami, Logout); Register/GetUser — в интерфейсе
// пакета user/v1. Сигнатуры совпадают с service/iam, поэтому *service удовлетворяет
// интерфейс неявно.
type IAmService interface {
	Login(ctx context.Context, in input.LoginInput) (string, error)
	Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error)
	Logout(ctx context.Context, sessionUUID string) error
}
