package main

import (
	"context"
	"os"
	"strings"
	"syscall"

	"github.com/andrewostroumov/lifecycle"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/subosito/gotenv"
	"go.uber.org/zap"

	"github.com/andrewostroumov/grpc-greeter/internal/app"
	grpcserver "github.com/andrewostroumov/grpc-greeter/internal/pkg/server/grpc"
)

const defaultServerAddr = ":8080"

const (
	serverAddrPath               = "grpc.server.addr"
	serverKeepAliveTime          = "grpc.server.keepalive.time"
	serverKeepAliveTimeoutPath   = "grpc.server.keepalive.timeout"
	lifecycleShutdownTimeoutPath = "lifecycle.shutdown.timeout"
)

func main() {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()

	k, err := newKoanf()
	if err != nil {
		logger.Fatal("koanf error", zap.Error(err))
	}

	lc := lifecycle.New(lifecycleOpts(k, logger)...)
	serv := grpcserver.NewServer(serverConfig(k, logger))

	ap := app.New(k, serv, logger)
	if err := ap.Run(ctx); err != nil {
		logger.Fatal("app error", zap.Error(err))
	}

	lc.Invoke(serv.Run)
	lc.Shutdown(serv.Stop)

	lc.Run(ctx)
}

func newKoanf() (*koanf.Koanf, error) {
	if err := gotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	var k = koanf.New(".")

	if err := k.Load(env.Provider("", "_", func(s string) string {
		return strings.ToLower(s)
	}), nil); err != nil {
		return nil, err
	}

	return k, nil
}

func lifecycleOpts(k *koanf.Koanf, logger *zap.Logger) []lifecycle.Option {
	opts := []lifecycle.Option{
		lifecycle.WithSignalContext(context.Background(), syscall.SIGINT, syscall.SIGTERM),
		lifecycle.WithLogger(logger),
	}

	if d := k.Duration(lifecycleShutdownTimeoutPath); d != 0 {
		opts = append(opts, lifecycle.WithShutdownTimeout(d))
	}

	return opts
}

func serverConfig(k *koanf.Koanf, logger *zap.Logger) grpcserver.Config {
	addr := k.String(serverAddrPath)
	if addr == "" {
		addr = defaultServerAddr
	}

	keepAliveTime := k.Duration(serverKeepAliveTime)
	keepAliveTimeout := k.Duration(serverKeepAliveTimeoutPath)

	return grpcserver.Config{
		Addr:   addr,
		Logger: logger,
		KeepAlive: grpcserver.KeepAlive{
			Time:    keepAliveTime,
			Timeout: keepAliveTimeout,
		},
	}
}
