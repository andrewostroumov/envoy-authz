package health

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func Register(s grpc.ServiceRegistrar) {
	srv := health.NewServer()
	healthgrpc.RegisterHealthServer(s, srv)
}
