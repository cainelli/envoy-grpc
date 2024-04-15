package service

import (
	"context"
	"time"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

type ExtAuthz struct{}

var _ authv3.AuthorizationServer = &ExtAuthz{}

func (s *ExtAuthz) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	time.Sleep(10 * time.Millisecond)
	return &authv3.CheckResponse{
		Status: &status.Status{
			Code: int32(codes.OK),
		},
	}, nil
}
