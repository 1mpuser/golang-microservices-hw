package v1

import (
	userv1 "github.com/1mpuser/shared/pkg/proto/user/v1"
)

// Реализовать (неделя 6): реализация authv1.AuthServiceServer (gRPC). Конструктор NewAPI(service).

type api struct {
	userv1.UnimplementedUserServiceServer

	iAmService IAmService
}

func NewApi(iAmService IAmService) *api {
	return &api{
		iAmService: iAmService,
	}
}
