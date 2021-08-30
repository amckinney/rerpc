package rerpc_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rerpc/rerpc"
	"github.com/rerpc/rerpc/internal/assert"
	pingpb "github.com/rerpc/rerpc/internal/ping/v1test"
	"google.golang.org/protobuf/proto"
)

func TestHandlerReadMaxBytesV2(t *testing.T) {
	// We only need to manually instantiate the Protocols when we have options
	// to configure. Each Protocol has its own set of options, all of which are
	// configured on the rerpc.Stream implementation. You can otherwise just use:
	//
	//  ping := rerpc.NewHTTPHandler(
	//    pingpb.NewPingServiceHandlerReRPCV2(
	//      &ExamplePingServer{},
	//    )
	//  )
	//
	// The rerpc.Specificaiton seen by Interceptors only contains metadata specific
	// to the rerpc.Route, including the Content-Type, Path, and StreamType.
	const readMaxBytes = 32
	grpcProtocols := rerpc.NewGPRCServerProtocols(rerpc.GRPCMaxRequestSize(readMaxBytes))
	ping := rerpc.NewHTTPHandler(
		pingpb.NewPingServiceHandlerReRPCV2(
			&ExamplePingServer{},
			// Interceptors can be passed in here.
		),
		rerpc.WithHTTPServerProtocols(grpcProtocols...),
	)

	t.Run("grpc", func(t *testing.T) {
		server := httptest.NewServer(ping)
		defer server.Close()

		// We only need to manually instantiate the Protocols when we have options
		// to configure. Each Protocol has its own set of options, all of which are
		// configured on the rerpc.Stream implementation. You can otherwise just use:
		//
		//  client := pingpb.NewPingServiceClientReRPCV2(
		//      server.URL,
		//      rerpc.NewHTTPClient(server.Client()),
		//    )
		//  )
		//
		grpcProtocol := rerpc.NewGRPCClientProtocol(server.Client(), rerpc.GRPCGzip(true))
		client := pingpb.NewPingServiceClientReRPCV2(rerpc.NewHTTPClient(server.URL, server.Client(), rerpc.WithClientProtocol(grpcProtocol)))

		padding := "padding                      "
		req := &pingpb.PingRequest{Number: 42, Msg: padding}
		// Ensure that the probe is actually too big.
		probeBytes, err := proto.Marshal(req)
		assert.Nil(t, err, "marshal request")
		assert.Equal(t, len(probeBytes), readMaxBytes+1, "probe size")

		_, err = client.Ping(context.Background(), req)

		assert.NotNil(t, err, "ping error")
		assert.Equal(t, rerpc.CodeOf(err), rerpc.CodeInvalidArgument, "error code")
		assert.True(
			t,
			strings.Contains(err.Error(), "larger than configured max"),
			`error msg contains "larger than configured max"`,
		)
	})
}
