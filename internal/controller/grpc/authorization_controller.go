package grpc

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	v31 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	pb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthorizationController struct {
	logger *zap.Logger
	pb.UnimplementedAuthorizationServer
}

func NewAuthorizationController(logger *zap.Logger) *AuthorizationController {
	return &AuthorizationController{
		logger: logger,
	}
}

func (c *AuthorizationController) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	recv, _ := metadata.FromIncomingContext(ctx)
	c.logger.Info("request",
		zap.String("body", req.Attributes.Request.Http.Body),
		zap.Any("metadata", recv),
	)

	spew.Dump(req)

	res := pb.OkHttpResponse{
		Headers: []*v31.HeaderValueOption{
			{
				Header: &v31.HeaderValue{
					Key:   "x-user-id",
					Value: "123",
				},
			},
		},
	}

	return &pb.CheckResponse{
		Status: status.New(codes.OK, "ok").Proto(),
		HttpResponse: &pb.CheckResponse_OkResponse{
			OkResponse: &res,
		},
	}, nil
}

func (c *AuthorizationController) RegisterService(s grpc.ServiceRegistrar) {
	serv := s.(interface {
		Server() *grpc.Server
	})

	pb.RegisterAuthorizationServer(serv.Server(), c)
}
