package main

import (
	"github.com/kaasikodes/shop-ease/shared/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCClient(addr string, logger logger.Logger) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	if err != nil {
		logger.Fatal("Unable to connect %v", err)
	}
	logger.Info("Connected to grpc client", addr)
	return conn

}
