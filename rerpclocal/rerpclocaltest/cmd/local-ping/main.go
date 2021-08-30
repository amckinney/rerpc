package main

import (
	"context"
	"errors"
	"io"

	"github.com/rerpc/rerpc"
	pingpb "github.com/rerpc/rerpc/internal/ping/v1test"
	"github.com/rerpc/rerpc/rerpclocal"
)

func main() {
	rerpclocal.Serve(
		rerpclocal.NewHandler(
			pingpb.NewPingServiceHandlerReRPCV2(
				&ExamplePingServer{},
			),
		),
	)
}

// ExamplePingServer implements some trivial business logic. The protobuf
// definition for this API is in internal/ping/v1test/ping.proto.
type ExamplePingServer struct {
	pingpb.UnimplementedPingServiceReRPC
}

// Ping implements pingpb.PingServiceReRPC.
func (*ExamplePingServer) Ping(ctx context.Context, req *pingpb.PingRequest) (*pingpb.PingResponse, error) {
	return &pingpb.PingResponse{Number: req.Number, Msg: req.Msg}, nil
}

// Fail implements pingpb.PingServiceReRPC.
func (*ExamplePingServer) Fail(ctx context.Context, req *pingpb.FailRequest) (*pingpb.FailResponse, error) {
	return nil, rerpc.Errorf(rerpc.CodeResourceExhausted, "err")
}

// Sum implements pingpb.PingServiceReRPC.
func (*ExamplePingServer) Sum(ctx context.Context, stream *pingpb.PingServiceReRPC_Sum) error {
	var sum int64
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		msg, err := stream.Receive()
		// TODO(alex): Write the messages received to a file so we can introspect it.
		if errors.Is(err, io.EOF) {
			return stream.SendAndClose(&pingpb.SumResponse{
				Sum: sum,
			})
		} else if err != nil {
			return err
		}
		sum += msg.Number
	}
}

// CountUp implements pingpb.PingServiceReRPC.
func (*ExamplePingServer) CountUp(ctx context.Context, req *pingpb.CountUpRequest, stream *pingpb.PingServiceReRPC_CountUp) error {
	if req.Number <= 0 {
		return rerpc.Errorf(rerpc.CodeInvalidArgument, "number must be positive: got %v", req.Number)
	}
	for i := int64(1); i <= req.Number; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := stream.Send(&pingpb.CountUpResponse{Number: i}); err != nil {
			return err
		}
	}
	return nil
}

// CumSum implements pingpb.PingServiceReRPC.
func (*ExamplePingServer) CumSum(ctx context.Context, stream *pingpb.PingServiceReRPC_CumSum) error {
	var sum int64
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		msg, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}
		sum += msg.Number
		if err := stream.Send(&pingpb.CumSumResponse{Sum: sum}); err != nil {
			return err
		}
	}
}
