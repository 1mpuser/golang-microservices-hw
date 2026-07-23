package v1

import (
	"context"

	"github.com/1mpuser/iam/internal/model"
	"github.com/1mpuser/iam/internal/service/input"
)

// IAmService — то, что gRPC-хендлеры UserService требуют от сервисного слоя.
// Только методы, нужные этому API (Register, GetUser); Login/Whoami/Logout —
// в интерфейсе пакета auth/v1. Сигнатуры совпадают с service/iam, поэтому
// *service удовлетворяет интерфейс неявно.
type IAmService interface {
	Register(ctx context.Context, in input.RegisterInput) (string, error)
	GetUser(ctx context.Context, uid string) (model.User, error)
}
