package app

import (
	"context"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"

	grpccontroller "github.com/andrewostroumov/grpc-greeter/internal/controller/grpc"
	grpcserver "github.com/andrewostroumov/grpc-greeter/internal/pkg/server/grpc"
)

type Application struct {
	k      *koanf.Koanf
	serv   *grpcserver.Server
	logger *zap.Logger
}

func New(k *koanf.Koanf, serv *grpcserver.Server, logger *zap.Logger) *Application {
	return &Application{
		k:      k,
		serv:   serv,
		logger: logger,
	}
}

func (a *Application) Run(ctx context.Context) error {
	grpccontroller.NewAuthorizationController(a.logger).RegisterService(a.serv)

	return nil
}
