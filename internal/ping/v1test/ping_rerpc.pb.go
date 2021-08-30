// Code generated by protoc-gen-go-rerpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-rerpc v0.0.1
// - protoc             v3.17.3
// source: internal/ping/v1test/ping.proto

package pingpb

import (
	context "context"
	errors "errors"
	strings "strings"

	rerpc "github.com/rerpc/rerpc"
)

// This is a compile-time assertion to ensure that this generated file and the
// rerpc package are compatible. If you get a compiler error that this constant
// isn't defined, this code was generated with a version of rerpc newer than the
// one compiled into your binary. You can fix the problem by either regenerating
// this code with an older version of rerpc or updating the rerpc version
// compiled into your binary.
const _ = rerpc.SupportsCodeGenV0 // requires reRPC v0.0.1 or later

// PingServiceClientReRPC is a client for the internal.ping.v1test.PingService
// service.
type PingServiceClientReRPC interface {
	Ping(ctx context.Context, req *PingRequest, opts ...rerpc.CallOption) (*PingResponse, error)
	Fail(ctx context.Context, req *FailRequest, opts ...rerpc.CallOption) (*FailResponse, error)
	Sum(ctx context.Context, opts ...rerpc.CallOption) *PingServiceClientReRPC_Sum
	CountUp(ctx context.Context, req *CountUpRequest, opts ...rerpc.CallOption) (*PingServiceClientReRPC_CountUp, error)
	CumSum(ctx context.Context, opts ...rerpc.CallOption) *PingServiceClientReRPC_CumSum
}

type pingServiceClientReRPC struct {
	doer    rerpc.Doer
	baseURL string
	options []rerpc.CallOption
}

// NewPingServiceClientReRPC constructs a client for the
// internal.ping.v1test.PingService service. Call options passed here apply to
// all calls made with this client.
//
// The URL supplied here should be the base URL for the gRPC server (e.g.,
// https://api.acme.com or https://acme.com/grpc).
func NewPingServiceClientReRPC(baseURL string, doer rerpc.Doer, opts ...rerpc.CallOption) PingServiceClientReRPC {
	return &pingServiceClientReRPC{
		baseURL: strings.TrimRight(baseURL, "/"),
		doer:    doer,
		options: opts,
	}
}

func (c *pingServiceClientReRPC) mergeOptions(opts []rerpc.CallOption) []rerpc.CallOption {
	merged := make([]rerpc.CallOption, 0, len(c.options)+len(opts))
	for _, o := range c.options {
		merged = append(merged, o)
	}
	for _, o := range opts {
		merged = append(merged, o)
	}
	return merged
}

// Ping calls internal.ping.v1test.PingService.Ping. Call options passed here
// apply only to this call.
func (c *pingServiceClientReRPC) Ping(ctx context.Context, req *PingRequest, opts ...rerpc.CallOption) (*PingResponse, error) {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged)
	ctx, call := rerpc.NewCall(
		ctx,
		c.doer,
		rerpc.StreamTypeUnary,
		c.baseURL,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"Ping",                 // protobuf method
		merged...,
	)
	wrapped := rerpc.Func(func(ctx context.Context, msg interface{}) (interface{}, error) {
		stream := call(ctx)
		if err := stream.Send(req); err != nil {
			_ = stream.CloseSend(err)
			_ = stream.CloseReceive()
			return nil, err
		}
		if err := stream.CloseSend(nil); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		var res PingResponse
		if err := stream.Receive(&res); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		return &res, stream.CloseReceive()
	})
	if ic != nil {
		wrapped = ic.Wrap(wrapped)
	}
	res, err := wrapped(ctx, req)
	if err != nil {
		return nil, err
	}
	typed, ok := res.(*PingResponse)
	if !ok {
		return nil, rerpc.Errorf(rerpc.CodeInternal, "expected response to be internal.ping.v1test.PingResponse, got %T", res)
	}
	return typed, nil
}

// Fail calls internal.ping.v1test.PingService.Fail. Call options passed here
// apply only to this call.
func (c *pingServiceClientReRPC) Fail(ctx context.Context, req *FailRequest, opts ...rerpc.CallOption) (*FailResponse, error) {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged)
	ctx, call := rerpc.NewCall(
		ctx,
		c.doer,
		rerpc.StreamTypeUnary,
		c.baseURL,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"Fail",                 // protobuf method
		merged...,
	)
	wrapped := rerpc.Func(func(ctx context.Context, msg interface{}) (interface{}, error) {
		stream := call(ctx)
		if err := stream.Send(req); err != nil {
			_ = stream.CloseSend(err)
			_ = stream.CloseReceive()
			return nil, err
		}
		if err := stream.CloseSend(nil); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		var res FailResponse
		if err := stream.Receive(&res); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		return &res, stream.CloseReceive()
	})
	if ic != nil {
		wrapped = ic.Wrap(wrapped)
	}
	res, err := wrapped(ctx, req)
	if err != nil {
		return nil, err
	}
	typed, ok := res.(*FailResponse)
	if !ok {
		return nil, rerpc.Errorf(rerpc.CodeInternal, "expected response to be internal.ping.v1test.FailResponse, got %T", res)
	}
	return typed, nil
}

// Sum calls internal.ping.v1test.PingService.Sum. Call options passed here
// apply only to this call.
func (c *pingServiceClientReRPC) Sum(ctx context.Context, opts ...rerpc.CallOption) *PingServiceClientReRPC_Sum {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged)
	ctx, call := rerpc.NewCall(
		ctx,
		c.doer,
		rerpc.StreamTypeClient,
		c.baseURL,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"Sum",                  // protobuf method
		merged...,
	)
	if ic != nil {
		call = ic.WrapStream(call)
	}
	stream := call(ctx)
	return NewPingServiceClientReRPC_Sum(stream)
}

// CountUp calls internal.ping.v1test.PingService.CountUp. Call options passed
// here apply only to this call.
func (c *pingServiceClientReRPC) CountUp(ctx context.Context, req *CountUpRequest, opts ...rerpc.CallOption) (*PingServiceClientReRPC_CountUp, error) {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged)
	ctx, call := rerpc.NewCall(
		ctx,
		c.doer,
		rerpc.StreamTypeServer,
		c.baseURL,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"CountUp",              // protobuf method
		merged...,
	)
	if ic != nil {
		call = ic.WrapStream(call)
	}
	stream := call(ctx)
	if err := stream.Send(req); err != nil {
		_ = stream.CloseSend(err)
		_ = stream.CloseReceive()
		return nil, err
	}
	if err := stream.CloseSend(nil); err != nil {
		_ = stream.CloseReceive()
		return nil, err
	}
	return NewPingServiceClientReRPC_CountUp(stream), nil
}

// CumSum calls internal.ping.v1test.PingService.CumSum. Call options passed
// here apply only to this call.
func (c *pingServiceClientReRPC) CumSum(ctx context.Context, opts ...rerpc.CallOption) *PingServiceClientReRPC_CumSum {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged)
	ctx, call := rerpc.NewCall(
		ctx,
		c.doer,
		rerpc.StreamTypeBidirectional,
		c.baseURL,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"CumSum",               // protobuf method
		merged...,
	)
	if ic != nil {
		call = ic.WrapStream(call)
	}
	stream := call(ctx)
	return NewPingServiceClientReRPC_CumSum(stream)
}

type pingServiceClientReRPCV2 struct {
	client  rerpc.Client
	options []rerpc.CallOption
}

// NewPingServiceClientReRPCV2 constructs a client for the
// internal.ping.v1test.PingService service. Call options passed here apply to
// all calls made with this client.
//
// The URL supplied here should be the base URL for the gRPC server (e.g.,
// https://api.acme.com or https://acme.com/grpc).
func NewPingServiceClientReRPCV2(client rerpc.Client, opts ...rerpc.CallOption) PingServiceClientReRPC {
	return &pingServiceClientReRPCV2{
		client:  client,
		options: opts,
	}
}

func (c *pingServiceClientReRPCV2) mergeOptions(opts []rerpc.CallOption) []rerpc.CallOption {
	merged := make([]rerpc.CallOption, 0, len(c.options)+len(opts))
	for _, o := range c.options {
		merged = append(merged, o)
	}
	for _, o := range opts {
		merged = append(merged, o)
	}
	return merged
}

// Ping calls internal.ping.v1test.PingService.Ping. Call options passed here
// apply only to this call.
func (c *pingServiceClientReRPCV2) Ping(ctx context.Context, req *PingRequest, opts ...rerpc.CallOption) (*PingResponse, error) {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged)

	ctx, call := c.client.NewCall(
		ctx,
		rerpc.StreamTypeUnary,
		"/internal.ping.v1test.PingService/Ping",
	)

	wrapped := rerpc.Func(func(ctx context.Context, msg interface{}) (interface{}, error) {
		stream := call(ctx)
		if err := stream.Send(req); err != nil {
			_ = stream.CloseSend(err)
			_ = stream.CloseReceive()
			return nil, err
		}
		if err := stream.CloseSend(nil); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		var res PingResponse
		if err := stream.Receive(&res); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		return &res, stream.CloseReceive()
	})
	if ic != nil {
		wrapped = ic.Wrap(wrapped)
	}
	res, err := wrapped(ctx, req)
	if err != nil {
		return nil, err
	}
	typed, ok := res.(*PingResponse)
	if !ok {
		return nil, rerpc.Errorf(rerpc.CodeInternal, "expected response to be internal.ping.v1test.PingResponse, got %T", res)
	}
	return typed, nil
}

// Fail calls internal.ping.v1test.PingService.Fail. Call options passed here
// apply only to this call.
func (c *pingServiceClientReRPCV2) Fail(ctx context.Context, req *FailRequest, opts ...rerpc.CallOption) (*FailResponse, error) {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged)

	ctx, call := c.client.NewCall(
		ctx,
		rerpc.StreamTypeUnary,
		"/internal.ping.v1test.PingService/Fail",
	)

	wrapped := rerpc.Func(func(ctx context.Context, msg interface{}) (interface{}, error) {
		stream := call(ctx)
		if err := stream.Send(req); err != nil {
			_ = stream.CloseSend(err)
			_ = stream.CloseReceive()
			return nil, err
		}
		if err := stream.CloseSend(nil); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		var res FailResponse
		if err := stream.Receive(&res); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		return &res, stream.CloseReceive()
	})
	if ic != nil {
		wrapped = ic.Wrap(wrapped)
	}
	res, err := wrapped(ctx, req)
	if err != nil {
		return nil, err
	}
	typed, ok := res.(*FailResponse)
	if !ok {
		return nil, rerpc.Errorf(rerpc.CodeInternal, "expected response to be internal.ping.v1test.FailResponse, got %T", res)
	}
	return typed, nil
}

// Sum calls internal.ping.v1test.PingService.Sum. Call options passed here
// apply only to this call.
func (c *pingServiceClientReRPCV2) Sum(ctx context.Context, opts ...rerpc.CallOption) *PingServiceClientReRPC_Sum {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged)
	ctx, call := c.client.NewCall(
		ctx,
		rerpc.StreamTypeClient,
		"/internal.ping.v1test.PingService/Sum",
	)
	if ic != nil {
		call = ic.WrapStream(call)
	}
	stream := call(ctx)
	return NewPingServiceClientReRPC_Sum(stream)
}

// CountUp calls internal.ping.v1test.PingService.CountUp. Call options passed
// here apply only to this call.
func (c *pingServiceClientReRPCV2) CountUp(ctx context.Context, req *CountUpRequest, opts ...rerpc.CallOption) (*PingServiceClientReRPC_CountUp, error) {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged)
	ctx, call := c.client.NewCall(
		ctx,
		rerpc.StreamTypeServer,
		"/internal.ping.v1test.PingService/CountUp",
	)
	if ic != nil {
		call = ic.WrapStream(call)
	}
	stream := call(ctx)
	if err := stream.Send(req); err != nil {
		_ = stream.CloseSend(err)
		_ = stream.CloseReceive()
		return nil, err
	}
	if err := stream.CloseSend(nil); err != nil {
		_ = stream.CloseReceive()
		return nil, err
	}
	return NewPingServiceClientReRPC_CountUp(stream), nil
}

// CumSum calls internal.ping.v1test.PingService.CumSum. Call options passed
// here apply only to this call.
func (c *pingServiceClientReRPCV2) CumSum(ctx context.Context, opts ...rerpc.CallOption) *PingServiceClientReRPC_CumSum {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged)
	ctx, call := c.client.NewCall(
		ctx,
		rerpc.StreamTypeBidirectional,
		"/internal.ping.v1test.PingService/CumSum",
	)
	if ic != nil {
		call = ic.WrapStream(call)
	}
	stream := call(ctx)
	return NewPingServiceClientReRPC_CumSum(stream)
}

// PingServiceReRPC is a server for the internal.ping.v1test.PingService
// service. To make sure that adding methods to this protobuf service doesn't
// break all implementations of this interface, all implementations must embed
// UnimplementedPingServiceReRPC.
//
// By default, recent versions of grpc-go have a similar forward compatibility
// requirement. See https://github.com/grpc/grpc-go/issues/3794 for a longer
// discussion.
type PingServiceReRPC interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Fail(context.Context, *FailRequest) (*FailResponse, error)
	Sum(context.Context, *PingServiceReRPC_Sum) error
	CountUp(context.Context, *CountUpRequest, *PingServiceReRPC_CountUp) error
	CumSum(context.Context, *PingServiceReRPC_CumSum) error
	mustEmbedUnimplementedPingServiceReRPC()
}

// NewPingServiceHandlerReRPC wraps each method on the service implementation in
// a *rerpc.Handler. The returned slice can be passed to rerpc.NewServeMux.
func NewPingServiceHandlerReRPC(svc PingServiceReRPC, opts ...rerpc.HandlerOption) []*rerpc.Handler {
	handlers := make([]*rerpc.Handler, 0, 5)
	ic := rerpc.ConfiguredHandlerInterceptor(opts)

	pingFunc := rerpc.Func(func(ctx context.Context, req interface{}) (interface{}, error) {
		typed, ok := req.(*PingRequest)
		if !ok {
			return nil, rerpc.Errorf(
				rerpc.CodeInternal,
				"can't call internal.ping.v1test.PingService.Ping with a %T",
				req,
			)
		}
		return svc.Ping(ctx, typed)
	})
	if ic != nil {
		pingFunc = ic.Wrap(pingFunc)
	}
	ping := rerpc.NewHandler(
		rerpc.StreamTypeUnary,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"Ping",                 // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			stream := sf(ctx)
			defer stream.CloseReceive()
			if err := ctx.Err(); err != nil {
				if errors.Is(err, context.Canceled) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeCanceled, err))
					return
				}
				if errors.Is(err, context.DeadlineExceeded) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeDeadlineExceeded, err))
					return
				}
				_ = stream.CloseSend(err) // unreachable per context docs
			}
			var req PingRequest
			if err := stream.Receive(&req); err != nil {
				_ = stream.CloseSend(err)
				return
			}
			res, err := pingFunc(ctx, &req)
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
				_ = stream.CloseSend(err)
				return
			}
			_ = stream.CloseSend(stream.Send(res))
		},
		opts...,
	)
	handlers = append(handlers, ping)

	failFunc := rerpc.Func(func(ctx context.Context, req interface{}) (interface{}, error) {
		typed, ok := req.(*FailRequest)
		if !ok {
			return nil, rerpc.Errorf(
				rerpc.CodeInternal,
				"can't call internal.ping.v1test.PingService.Fail with a %T",
				req,
			)
		}
		return svc.Fail(ctx, typed)
	})
	if ic != nil {
		failFunc = ic.Wrap(failFunc)
	}
	fail := rerpc.NewHandler(
		rerpc.StreamTypeUnary,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"Fail",                 // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			stream := sf(ctx)
			defer stream.CloseReceive()
			if err := ctx.Err(); err != nil {
				if errors.Is(err, context.Canceled) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeCanceled, err))
					return
				}
				if errors.Is(err, context.DeadlineExceeded) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeDeadlineExceeded, err))
					return
				}
				_ = stream.CloseSend(err) // unreachable per context docs
			}
			var req FailRequest
			if err := stream.Receive(&req); err != nil {
				_ = stream.CloseSend(err)
				return
			}
			res, err := failFunc(ctx, &req)
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
				_ = stream.CloseSend(err)
				return
			}
			_ = stream.CloseSend(stream.Send(res))
		},
		opts...,
	)
	handlers = append(handlers, fail)

	sum := rerpc.NewHandler(
		rerpc.StreamTypeClient,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"Sum",                  // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			if ic != nil {
				sf = ic.WrapStream(sf)
			}
			stream := sf(ctx)
			typed := NewPingServiceReRPC_Sum(stream)
			err := svc.Sum(stream.Context(), typed)
			_ = stream.CloseReceive()
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
			}
			_ = stream.CloseSend(err)
		},
		opts...,
	)
	handlers = append(handlers, sum)

	countUp := rerpc.NewHandler(
		rerpc.StreamTypeServer,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"CountUp",              // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			if ic != nil {
				sf = ic.WrapStream(sf)
			}
			stream := sf(ctx)
			typed := NewPingServiceReRPC_CountUp(stream)
			var req CountUpRequest
			if err := stream.Receive(&req); err != nil {
				_ = stream.CloseReceive()
				_ = stream.CloseSend(err)
				return
			}
			if err := stream.CloseReceive(); err != nil {
				_ = stream.CloseSend(err)
				return
			}
			err := svc.CountUp(stream.Context(), &req, typed)
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
			}
			_ = stream.CloseSend(err)
		},
		opts...,
	)
	handlers = append(handlers, countUp)

	cumSum := rerpc.NewHandler(
		rerpc.StreamTypeBidirectional,
		"internal.ping.v1test", // protobuf package
		"PingService",          // protobuf service
		"CumSum",               // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			if ic != nil {
				sf = ic.WrapStream(sf)
			}
			stream := sf(ctx)
			typed := NewPingServiceReRPC_CumSum(stream)
			err := svc.CumSum(stream.Context(), typed)
			_ = stream.CloseReceive()
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
			}
			_ = stream.CloseSend(err)
		},
		opts...,
	)
	handlers = append(handlers, cumSum)

	return handlers
}

// NewPingServiceHandlerReRPCV2 wraps each method on the service implementation in
// an *rerpc.Mux.
func NewPingServiceHandlerReRPCV2(svc PingServiceReRPC, opts ...rerpc.HandlerOption) *rerpc.Mux {
	routes := make([]rerpc.Route, 0, 5)
	ic := rerpc.ConfiguredHandlerInterceptor(opts)

	pingFunc := rerpc.Func(func(ctx context.Context, req interface{}) (interface{}, error) {
		typed, ok := req.(*PingRequest)
		if !ok {
			return nil, rerpc.Errorf(
				rerpc.CodeInternal,
				"can't call internal.ping.v1test.PingService.Ping with a %T",
				req,
			)
		}
		return svc.Ping(ctx, typed)
	})
	if ic != nil {
		pingFunc = ic.Wrap(pingFunc)
	}
	ping := rerpc.Route{
		Type: rerpc.StreamTypeUnary,
		Path: "/internal.ping.v1test.PingService/Ping",
		Implementation: func(ctx context.Context, streamFunc rerpc.StreamFunc) {
			stream := streamFunc(ctx)
			defer stream.CloseReceive()
			if err := ctx.Err(); err != nil {
				if errors.Is(err, context.Canceled) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeCanceled, err))
					return
				}
				if errors.Is(err, context.DeadlineExceeded) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeDeadlineExceeded, err))
					return
				}
				_ = stream.CloseSend(err) // unreachable per context docs
			}
			var req PingRequest
			if err := stream.Receive(&req); err != nil {
				_ = stream.CloseSend(err)
				return
			}
			res, err := pingFunc(ctx, &req)
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
				_ = stream.CloseSend(err)
				return
			}
			_ = stream.CloseSend(stream.Send(res))
			return
		},
	}
	routes = append(routes, ping)

	failFunc := rerpc.Func(func(ctx context.Context, req interface{}) (interface{}, error) {
		typed, ok := req.(*FailRequest)
		if !ok {
			return nil, rerpc.Errorf(
				rerpc.CodeInternal,
				"can't call internal.ping.v1test.PingService.Fail with a %T",
				req,
			)
		}
		return svc.Fail(ctx, typed)
	})
	if ic != nil {
		failFunc = ic.Wrap(failFunc)
	}
	fail := rerpc.Route{
		Type: rerpc.StreamTypeUnary,
		Path: "/internal.ping.v1test.PingService/Fail",
		Implementation: func(ctx context.Context, streamFunc rerpc.StreamFunc) {
			stream := streamFunc(ctx)
			defer stream.CloseReceive()
			if err := ctx.Err(); err != nil {
				if errors.Is(err, context.Canceled) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeCanceled, err))
					return
				}
				if errors.Is(err, context.DeadlineExceeded) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeDeadlineExceeded, err))
					return
				}
				_ = stream.CloseSend(err) // unreachable per context docs
			}
			var req FailRequest
			if err := stream.Receive(&req); err != nil {
				_ = stream.CloseSend(err)
				return
			}
			res, err := failFunc(ctx, &req)
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
				_ = stream.CloseSend(err)
				return
			}
			_ = stream.CloseSend(stream.Send(res))
			return
		},
	}
	routes = append(routes, fail)

	sum := rerpc.Route{
		Type: rerpc.StreamTypeClient,
		Path: "/internal.ping.v1test.PingService/Sum",
		Implementation: func(ctx context.Context, streamFunc rerpc.StreamFunc) {
			if ic != nil {
				streamFunc = ic.WrapStream(streamFunc)
			}
			stream := streamFunc(ctx)

			typed := NewPingServiceReRPC_Sum(stream)
			err := svc.Sum(stream.Context(), typed)
			_ = stream.CloseReceive()
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
			}
			_ = stream.CloseSend(err)
			return
		},
	}
	routes = append(routes, sum)

	countUp := rerpc.Route{
		Type: rerpc.StreamTypeServer,
		Path: "/internal.ping.v1test.PingService/CountUp",
		Implementation: func(ctx context.Context, streamFunc rerpc.StreamFunc) {
			if ic != nil {
				streamFunc = ic.WrapStream(streamFunc)
			}
			stream := streamFunc(ctx)

			typed := NewPingServiceReRPC_CountUp(stream)
			var req CountUpRequest
			if err := stream.Receive(&req); err != nil {
				_ = stream.CloseReceive()
				_ = stream.CloseSend(err)
				return
			}
			if err := stream.CloseReceive(); err != nil {
				_ = stream.CloseSend(err)
				return
			}
			err := svc.CountUp(stream.Context(), &req, typed)
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
			}
			_ = stream.CloseSend(err)
			return
		},
	}
	routes = append(routes, countUp)

	cumSum := rerpc.Route{
		Type: rerpc.StreamTypeBidirectional,
		Path: "/internal.ping.v1test.PingService/CumSum",
		Implementation: func(ctx context.Context, streamFunc rerpc.StreamFunc) {
			if ic != nil {
				streamFunc = ic.WrapStream(streamFunc)
			}
			stream := streamFunc(ctx)

			typed := NewPingServiceReRPC_CumSum(stream)
			err := svc.CumSum(stream.Context(), typed)
			_ = stream.CloseReceive()
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
			}
			_ = stream.CloseSend(err)
			return
		},
	}
	routes = append(routes, cumSum)

	return rerpc.NewMux(routes...)
}

var _ PingServiceReRPC = (*UnimplementedPingServiceReRPC)(nil) // verify interface implementation

// UnimplementedPingServiceReRPC returns CodeUnimplemented from all methods. To
// maintain forward compatibility, all implementations of PingServiceReRPC must
// embed UnimplementedPingServiceReRPC.
type UnimplementedPingServiceReRPC struct{}

func (UnimplementedPingServiceReRPC) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, rerpc.Errorf(rerpc.CodeUnimplemented, "internal.ping.v1test.PingService.Ping isn't implemented")
}

func (UnimplementedPingServiceReRPC) Fail(context.Context, *FailRequest) (*FailResponse, error) {
	return nil, rerpc.Errorf(rerpc.CodeUnimplemented, "internal.ping.v1test.PingService.Fail isn't implemented")
}

func (UnimplementedPingServiceReRPC) Sum(context.Context, *PingServiceReRPC_Sum) error {
	return rerpc.Errorf(rerpc.CodeUnimplemented, "internal.ping.v1test.PingService.Sum isn't implemented")
}

func (UnimplementedPingServiceReRPC) CountUp(context.Context, *CountUpRequest, *PingServiceReRPC_CountUp) error {
	return rerpc.Errorf(rerpc.CodeUnimplemented, "internal.ping.v1test.PingService.CountUp isn't implemented")
}

func (UnimplementedPingServiceReRPC) CumSum(context.Context, *PingServiceReRPC_CumSum) error {
	return rerpc.Errorf(rerpc.CodeUnimplemented, "internal.ping.v1test.PingService.CumSum isn't implemented")
}

func (UnimplementedPingServiceReRPC) mustEmbedUnimplementedPingServiceReRPC() {}

// PingServiceClientReRPC_Sum is the client-side stream for the
// internal.ping.v1test.PingService.Sum procedure.
type PingServiceClientReRPC_Sum struct {
	stream rerpc.Stream
}

func NewPingServiceClientReRPC_Sum(stream rerpc.Stream) *PingServiceClientReRPC_Sum {
	return &PingServiceClientReRPC_Sum{stream}
}

func (s *PingServiceClientReRPC_Sum) Send(msg *SumRequest) error {
	return s.stream.Send(msg)
}

func (s *PingServiceClientReRPC_Sum) CloseAndReceive() (*SumResponse, error) {
	if err := s.stream.CloseSend(nil); err != nil {
		return nil, err
	}
	var res SumResponse
	if err := s.stream.Receive(&res); err != nil {
		_ = s.stream.CloseReceive()
		return nil, err
	}
	if err := s.stream.CloseReceive(); err != nil {
		return nil, err
	}
	return &res, nil
}

// PingServiceClientReRPC_CountUp is the client-side stream for the
// internal.ping.v1test.PingService.CountUp procedure.
type PingServiceClientReRPC_CountUp struct {
	stream rerpc.Stream
}

func NewPingServiceClientReRPC_CountUp(stream rerpc.Stream) *PingServiceClientReRPC_CountUp {
	return &PingServiceClientReRPC_CountUp{stream}
}

func (s *PingServiceClientReRPC_CountUp) Receive() (*CountUpResponse, error) {
	var req CountUpResponse
	if err := s.stream.Receive(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *PingServiceClientReRPC_CountUp) Close() error {
	return s.stream.CloseReceive()
}

// PingServiceClientReRPC_CumSum is the client-side stream for the
// internal.ping.v1test.PingService.CumSum procedure.
type PingServiceClientReRPC_CumSum struct {
	stream rerpc.Stream
}

func NewPingServiceClientReRPC_CumSum(stream rerpc.Stream) *PingServiceClientReRPC_CumSum {
	return &PingServiceClientReRPC_CumSum{stream}
}

func (s *PingServiceClientReRPC_CumSum) Send(msg *CumSumRequest) error {
	return s.stream.Send(msg)
}

func (s *PingServiceClientReRPC_CumSum) CloseSend() error {
	return s.stream.CloseSend(nil)
}

func (s *PingServiceClientReRPC_CumSum) Receive() (*CumSumResponse, error) {
	var req CumSumResponse
	if err := s.stream.Receive(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *PingServiceClientReRPC_CumSum) CloseReceive() error {
	return s.stream.CloseReceive()
}

// PingServiceReRPC_Sum is the server-side stream for the
// internal.ping.v1test.PingService.Sum procedure.
type PingServiceReRPC_Sum struct {
	stream rerpc.Stream
}

func NewPingServiceReRPC_Sum(stream rerpc.Stream) *PingServiceReRPC_Sum {
	return &PingServiceReRPC_Sum{stream}
}

func (s *PingServiceReRPC_Sum) Receive() (*SumRequest, error) {
	var req SumRequest
	if err := s.stream.Receive(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *PingServiceReRPC_Sum) SendAndClose(msg *SumResponse) error {
	if err := s.stream.CloseReceive(); err != nil {
		return err
	}
	return s.stream.Send(msg)
}

// PingServiceReRPC_CountUp is the server-side stream for the
// internal.ping.v1test.PingService.CountUp procedure.
type PingServiceReRPC_CountUp struct {
	stream rerpc.Stream
}

func NewPingServiceReRPC_CountUp(stream rerpc.Stream) *PingServiceReRPC_CountUp {
	return &PingServiceReRPC_CountUp{stream}
}

func (s *PingServiceReRPC_CountUp) Send(msg *CountUpResponse) error {
	return s.stream.Send(msg)
}

// PingServiceReRPC_CumSum is the server-side stream for the
// internal.ping.v1test.PingService.CumSum procedure.
type PingServiceReRPC_CumSum struct {
	stream rerpc.Stream
}

func NewPingServiceReRPC_CumSum(stream rerpc.Stream) *PingServiceReRPC_CumSum {
	return &PingServiceReRPC_CumSum{stream}
}

func (s *PingServiceReRPC_CumSum) Receive() (*CumSumRequest, error) {
	var req CumSumRequest
	if err := s.stream.Receive(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *PingServiceReRPC_CumSum) Send(msg *CumSumResponse) error {
	return s.stream.Send(msg)
}
