package service

import (
	"context"
	"time"

	ratelimit "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
)

type RateLimitService struct{}

var _ ratelimit.RateLimitServiceServer = &RateLimitService{}

func (rls *RateLimitService) ShouldRateLimit(ctx context.Context, req *ratelimit.RateLimitRequest) (*ratelimit.RateLimitResponse, error) {
	time.Sleep(10 * time.Millisecond)
	return &ratelimit.RateLimitResponse{
		OverallCode: ratelimit.RateLimitResponse_OK,
	}, nil
}
