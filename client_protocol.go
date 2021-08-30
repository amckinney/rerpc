package rerpc

import (
	"context"
	"fmt"
	"runtime"
)

// ClientProtocol is used by the client to create Streams for outgoing
// HTTP requests.
type ClientProtocol interface {
	NewStream(ctx context.Context, url string) Stream
}

// NewGRPCClientProtocol creates a new ClientProtocol backed by gRPC.
func NewGRPCClientProtocol(doer Doer, opts ...GRPCClientProtocolOption) ClientProtocol {
	options := &grpcClientProtocolOptions{}
	for _, opt := range opts {
		opt(options)
	}
	var maxResponseSize int64
	if options.maxResponseSize != 0 {
		maxResponseSize = options.maxResponseSize
	}
	return &grpcClientProtocol{
		doer:              doer,
		enableGzipRequest: options.enableGzipRequest,
		maxResponseSize:   maxResponseSize,
	}
}

// GRPCClientProtocolOption is a gRPC ClientProtocol option.
type GRPCClientProtocolOption func(*grpcClientProtocolOptions)

// GRPCGzip enables or disables gzip compression based on the given value.
func GRPCGzip(enable bool) GRPCClientProtocolOption {
	return func(opts *grpcClientProtocolOptions) {
		opts.enableGzipRequest = enable
	}
}

// GRPCMaxResponseSize sets this clientProtocolection's max response bytes.
func GRPCMaxResponseSize(maxResponseSize int64) GRPCClientProtocolOption {
	return func(opts *grpcClientProtocolOptions) {
		opts.maxResponseSize = maxResponseSize
	}
}

type grpcClientProtocol struct {
	doer Doer

	enableGzipRequest bool
	maxResponseSize   int64
}

func (c *grpcClientProtocol) NewStream(ctx context.Context, url string) Stream {
	requestCompression := CompressionIdentity
	if c.enableGzipRequest {
		requestCompression = CompressionGzip
	}

	metadata, ok := CallMetadata(ctx)
	if !ok {
		// TODO(alex): Unreachable, but we should probably expose an internal error here.
	}

	// TODO(alex): If we attach headers at this layer, the previous interceptors
	// in the chain will NOT see them. Is this acceptable for ClientProtocol options (i.e.
	// header values)? This _might_ be considered a feature, but it could also
	// lead to confusing behavior in user-provided Interceptors.
	//
	// We don't want to attach at the client-level because these details are
	// gRPC-specific. However, given that the ClientProtocol is a property of the Client,
	// we _could_ expose a method on the ClientProtocol that would let these headers be set
	// in the NewCall method.
	//
	// This would look something along the lines of the following:
	//
	//  grpcClientProtocol, ok := c.clientProtocol.(*grpcClientProtocol)
	//  if ok {
	//    requestHeaders := grpcClientProtocol.requestHeaders(requestCompression)
	//  }
	//  ...
	//
	reqHeader := metadata.Request()

	reqHeader.Set("User-Agent", grpcGoUserAgent)
	reqHeader.Set("Content-Type", TypeDefaultGRPC)
	reqHeader.Set("Grpc-Accept-Encoding", acceptEncodingValue) // always advertise identity & gzip
	reqHeader.Set("Te", teTrailers)
	reqHeader.Set("Grpc-Encoding", requestCompression)

	return newClientStream(
		ctx,
		c.doer,
		url,
		c.maxResponseSize,
		c.enableGzipRequest,
	)
}

type grpcClientProtocolOptions struct {
	enableGzipRequest bool
	maxResponseSize   int64
}

var (
	grpcGoUserAgent = fmt.Sprintf("grpc-go-rerpc/%s (%s)", Version, runtime.Version())
	teTrailers      = "trailers"
)
