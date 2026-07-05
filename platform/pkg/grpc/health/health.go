package health

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Server реализует gRPC Health Checking Protocol (https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
type Server struct {
	grpc_health_v1.UnimplementedHealthServer
}

// Check — unary RPC для проверки здоровья сервиса
func (s *Server) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch — server-streaming RPC для подписки на изменения статуса здоровья
func (s *Server) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	return stream.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

// RegisterService регистрирует health-сервис на gRPC-сервере
func RegisterService(s *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(s, &Server{})
}
