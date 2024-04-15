package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cainelli/envoy-grpc/pkg/service"
	"golang.org/x/sync/errgroup"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	extproc "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	ratelimit "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"

	"google.golang.org/grpc"

	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	if err := run(); err != nil {
		slog.Error("oops", "error", err)
		os.Exit(1)
	}
}

func run() error {
	tracer.Start()
	defer tracer.Stop()

	si := grpctrace.StreamServerInterceptor(
		grpctrace.WithMetadataTags(),
		grpctrace.WithErrorDetailTags(),
		grpctrace.WithRequestTags(),
		grpctrace.WithStreamMessages(true),
	)

	ui := grpctrace.UnaryServerInterceptor(
		grpctrace.WithMetadataTags(),
		grpctrace.WithErrorDetailTags(),
		grpctrace.WithRequestTags(),
	)

	server := grpc.NewServer(
		grpc.StreamInterceptor(si),
		grpc.UnaryInterceptor(ui),
	)
	extprocService := &service.ExtProcessor{}
	extproc.RegisterExternalProcessorServer(server, extprocService)

	ratelimitService := &service.RateLimitService{}
	ratelimit.RegisterRateLimitServiceServer(server, ratelimitService)

	extauthzService := &service.ExtAuthz{}
	authv3.RegisterAuthorizationServer(server, extauthzService)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("cannot listen: %w", err)
	}

	httpsrv := httptrace.NewServeMux()
	httpsrv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	var eg errgroup.Group
	eg.Go(func() error {
		slog.Info("listening", "addr", listener.Addr().String(), "type", "grpc")
		return server.Serve(listener)
	})

	eg.Go(func() error {
		slog.Info("listening", "addr", ":8081", "type", "http")
		return http.ListenAndServe(":8081", httpsrv)
	})

	return eg.Wait()
}
