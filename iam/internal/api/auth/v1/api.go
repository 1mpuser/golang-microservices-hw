package v1

import authv1 "github.com/1mpuser/shared/pkg/proto/auth/v1"

// Реализовать (неделя 6): реализация authv1.AuthServiceServer (gRPC). Конструктор NewAPI(service).

type api struct {
	authv1.UnimplementedAuthServiceServer

	iAmService IAmService
}

func NewApi(iAmService IAmService) *api {
	return &api{
		iAmService: iAmService,
	}
}
