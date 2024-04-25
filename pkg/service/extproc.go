package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	extproc "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type ExtProcessor struct{}

var _ extproc.ExternalProcessorServer = &ExtProcessor{}

const (
	TraceMessageOperationName = "grpc.message"
)

var (
	RequestHeadersResourceName   = tracer.ResourceName("RequestHeaders")
	RequestBodyResourceName      = tracer.ResourceName("RequestBody")
	RequestTrailersResourceName  = tracer.ResourceName("RequestTrailers")
	ResponseHeadersResourceName  = tracer.ResourceName("ResponseHeaders")
	ResponseBodyResourceName     = tracer.ResourceName("ResponseBody")
	ResponseTrailersResourceName = tracer.ResourceName("ResponseTrailers")
)

func (svc *ExtProcessor) Process(procsrv extproc.ExternalProcessor_ProcessServer) error {
	ctx := procsrv.Context()

	for {
		procreq, err := procsrv.Recv()
		switch {
		case errors.Is(err, io.EOF), errors.Is(err, status.Error(codes.Canceled, context.Canceled.Error())):
			return nil
		case err != nil:
			return status.Errorf(codes.Unknown, "cannot receive stream request: %v", err)
		}

		switch procreq.Request.(type) {
		case *extproc.ProcessingRequest_RequestHeaders:
			span, _ := tracer.StartSpanFromContext(ctx, TraceMessageOperationName, RequestHeadersResourceName)
			time.Sleep(10 * time.Millisecond)
			if err := procsrv.Send(&extproc.ProcessingResponse{
				Response: &extproc.ProcessingResponse_RequestHeaders{
					RequestHeaders: &extproc.HeadersResponse{
						Response: &extproc.CommonResponse{
							HeaderMutation: &extproc.HeaderMutation{
								SetHeaders: []*corev3.HeaderValueOption{{
									Header: &corev3.HeaderValue{
										Key:      "x-req-header",
										RawValue: []byte("ok"),
									},
								}},
							},
						},
					},
				},
			}); err != nil {
				span.Finish(tracer.WithError(err))
				return fmt.Errorf("ResponseBody: failed sending response: %w", err)
			}
			span.Finish()
		case *extproc.ProcessingRequest_RequestBody:
			span, _ := tracer.StartSpanFromContext(ctx, TraceMessageOperationName, RequestBodyResourceName)
			time.Sleep(10 * time.Microsecond)
			if err := procsrv.Send(&extproc.ProcessingResponse{
				Response: &extproc.ProcessingResponse_RequestBody{},
			}); err != nil {
				span.Finish(tracer.WithError(err))
				return fmt.Errorf("ResponseBody: failed sending response: %w", err)
			}
			span.Finish()
		case *extproc.ProcessingRequest_RequestTrailers:
			span, _ := tracer.StartSpanFromContext(ctx, TraceMessageOperationName, RequestTrailersResourceName)
			time.Sleep(10 * time.Microsecond)
			if err := procsrv.Send(&extproc.ProcessingResponse{
				Response: &extproc.ProcessingResponse_RequestTrailers{},
			}); err != nil {
				span.Finish(tracer.WithError(err))
				return fmt.Errorf("ResponseBody: failed sending response: %w", err)
			}
			span.Finish()
		case *extproc.ProcessingRequest_ResponseHeaders:
			span, _ := tracer.StartSpanFromContext(ctx, TraceMessageOperationName, ResponseHeadersResourceName)
			time.Sleep(10 * time.Millisecond)
			if err := procsrv.Send(&extproc.ProcessingResponse{
				Response: &extproc.ProcessingResponse_ResponseHeaders{
					ResponseHeaders: &extproc.HeadersResponse{
						Response: &extproc.CommonResponse{
							HeaderMutation: &extproc.HeaderMutation{
								SetHeaders: []*corev3.HeaderValueOption{{
									Header: &corev3.HeaderValue{
										Key:      "x-res-header",
										RawValue: []byte("ok"),
									},
								}},
							},
						},
					},
				},
			}); err != nil {
				span.Finish(tracer.WithError(err))
				return fmt.Errorf("ResponseBody: failed sending response: %w", err)
			}
			span.Finish()
		case *extproc.ProcessingRequest_ResponseBody:
			span, _ := tracer.StartSpanFromContext(ctx, TraceMessageOperationName, ResponseBodyResourceName)
			time.Sleep(10 * time.Microsecond)
			if err := procsrv.Send(&extproc.ProcessingResponse{
				Response: &extproc.ProcessingResponse_ResponseBody{},
			}); err != nil {
				span.Finish(tracer.WithError(err))
				return fmt.Errorf("ResponseBody: failed sending response: %w", err)
			}
			span.Finish()
		case *extproc.ProcessingRequest_ResponseTrailers:
			span, _ := tracer.StartSpanFromContext(ctx, TraceMessageOperationName, ResponseTrailersResourceName)
			time.Sleep(10 * time.Microsecond)
			if err := procsrv.Send(&extproc.ProcessingResponse{
				Response: &extproc.ProcessingResponse_ResponseTrailers{},
			}); err != nil {
				span.Finish(tracer.WithError(err))
				return fmt.Errorf("ResponseTrailers: failed sending response: %w", err)
			}
			span.Finish()
		default:
			return fmt.Errorf("unknown request type: %T", procreq.Request)
		}
	}
}
