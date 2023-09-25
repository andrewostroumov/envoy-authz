package grpc

import (
	"context"
	"net"
	"sync/atomic"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/andrewostroumov/grpc-greeter/internal/pkg/server/grpc/health"
)

type Server struct {
	addr     string
	serv     *grpc.Server
	logger   *zap.Logger
	stopping atomic.Bool
}

func NewServer(conf Config) *Server {
	co := defaultConfig
	co.apply(conf)

	opts := serverOpts(co)
	serv := grpc.NewServer(opts...)

	health.Register(serv)
	reflection.Register(serv)

	return &Server{
		addr:   co.Addr,
		serv:   serv,
		logger: co.Logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting gRPC server", zap.String("addr", s.addr))
	lis, err := listen(s.addr)
	if err != nil {
		return err
	}

	return s.serv.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) error {
	if !s.stopping.CompareAndSwap(false, true) {
		return nil
	}

	s.logger.Info("Stopping gRPC server...")
	s.serv.GracefulStop()
	s.logger.Info("gRPC server gracefully stopped")

	return nil
}

func (s *Server) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	s.logger.Debug("Registering gRPC service", zap.String("service_name", sd.ServiceName))
	s.serv.RegisterService(sd, ss)
}

func (s *Server) Server() *grpc.Server {
	return s.serv
}

func listen(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}

func serverOpts(conf Config) []grpc.ServerOption {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	}

	if conf.KeepAlive.Time != 0 && conf.KeepAlive.Timeout != 0 {
		opts = append(opts, grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    conf.KeepAlive.Time,
				Timeout: conf.KeepAlive.Timeout,
			},
		))
	}

	return opts
}
