package rerpc

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/protobuf/proto"

	healthpb "github.com/akshayjshah/rerpc/internal/healthpb/v1"
)

// HealthStatus describes the health of a service.
type HealthStatus int32

// Constants representing the known health states.
//
// These correspond to the ServingStatus enum in gRPC's health.proto. Because
// reRPC doesn't support watching health, SERVICE_UNKNOWN isn't exposed here.
//
// For details, see the protobuf schema:
//   https://github.com/grpc/grpc/blob/master/src/proto/grpc/health/v1/health.proto
const (
	HealthUnknown    HealthStatus = 0 // health state indeterminate
	HealthServing    HealthStatus = 1 // ready to accept requests
	HealthNotServing HealthStatus = 2 // process healthy but service not accepting requests
)

// DefaultCheckFunc returns a health-checking function that always returns
// HealthServing for the process and all registered services.
func DefaultCheckFunc(reg *Registrar) func(context.Context, string) (HealthStatus, error) {
	return func(_ context.Context, service string) (HealthStatus, error) {
		if service == "" {
			return HealthServing, nil
		}
		if reg.IsRegistered(service) {
			return HealthServing, nil
		}
		return HealthUnknown, errorf(CodeNotFound, "unknown service %s", service)
	}
}

// NewHealthHandler wraps the supplied function to build an HTTP handler for
// gRPC's health-checking API. It returns the HTTP handler and the correct path
// on which to mount it. The health-checking function will be called with a
// fully-qualified protobuf service name (e.g., "acme.ping.v0.Ping").
//
// The supplied health-checking function should:
//   * Return HealthUnknown, HealthServing, or HealthNotServing.
//   * Return the health status of the whole process when called with an empty
//     string.
//   * Return a CodeNotFound error when called with an unknown service.
//
// Note that the returned handler only supports the unary Check method, not the
// streaming Watch. As suggested in gRPC's health schema, reRPC returns
// CodeUnimplemented for the Watch method. For more details on gRPC's health
// checking protocol, see:
//   https://github.com/grpc/grpc/blob/master/doc/health-checking.md
//   https://github.com/grpc/grpc/blob/master/src/proto/grpc/health/v1/health.proto
func NewHealthHandler(
	check func(context.Context, string) (HealthStatus, error),
	opts ...HandlerOption,
) (string, http.Handler) {
	const serviceFQN = "grpc.health.v1.Health"
	const checkFQN = serviceFQN + ".Check"
	const watchFQN = serviceFQN + ".Watch"

	mux := http.NewServeMux()
	checkHandler := NewHandler(
		checkFQN,
		func(ctx context.Context, req proto.Message) (proto.Message, error) {
			typed, ok := req.(*healthpb.HealthCheckRequest)
			if !ok {
				return nil, errorf(
					CodeInternal,
					"can't call %s/Check with a %v",
					serviceFQN,
					req.ProtoReflect().Descriptor().FullName(),
				)
			}
			status, err := check(ctx, typed.Service)
			if err != nil {
				return nil, err
			}
			return &healthpb.HealthCheckResponse{
				Status: healthpb.HealthCheckResponse_ServingStatus(status),
			}, nil
		},
		opts...,
	)
	mux.HandleFunc(fmt.Sprintf("/%s/Check", serviceFQN), func(w http.ResponseWriter, r *http.Request) {
		checkHandler.Serve(w, r, &healthpb.HealthCheckRequest{})
	})

	watch := NewHandler(
		watchFQN,
		func(ctx context.Context, req proto.Message) (proto.Message, error) {
			return nil, errorf(CodeUnimplemented, "reRPC doesn't support watching health state")
		},
		opts...,
	)
	mux.HandleFunc(fmt.Sprintf("/%s/Watch", serviceFQN), func(w http.ResponseWriter, r *http.Request) {
		watch.Serve(w, r, &healthpb.HealthCheckRequest{})
	})

	return fmt.Sprintf("/%s/", serviceFQN), mux
}
