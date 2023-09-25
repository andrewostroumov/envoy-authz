package grpc

import (
	"time"

	"go.uber.org/zap"
)

const defaultServerAddr = ":8080"

var defaultConfig = Config{
	Addr:   defaultServerAddr,
	Logger: zap.NewNop(),
}

type KeepAlive struct {
	Time    time.Duration
	Timeout time.Duration
}

type Config struct {
	Addr      string
	Logger    *zap.Logger
	KeepAlive KeepAlive
}

func (c *Config) apply(conf Config) {
	if conf.Addr != "" {
		c.Addr = conf.Addr
	}

	if conf.Logger != nil {
		c.Logger = conf.Logger
	}

	if conf.KeepAlive.Time != 0 {
		c.KeepAlive.Time = conf.KeepAlive.Time
	}

	if conf.KeepAlive.Timeout != 0 {
		c.KeepAlive.Timeout = conf.KeepAlive.Timeout
	}
}
