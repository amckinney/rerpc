package rerpc

import (
	"context"
	"net/http"
	"strings"
)

// Client represents a generic client.
type Client interface {
	// NewCall acts similarly to the the top-level NewCall constructor, but it doesn't
	// take any dependency on HTTP.
	NewCall(ctx context.Context, streamType StreamType, path string) (context.Context, StreamFunc)
}

// ClientOption configures some aspect of the client.
type ClientOption func(*clientOptions)

// WithClientProtocol registers the given Protocols. The match is exact, with the
// special case that the content type "*" is the fallback Protocol used
// when nothing matches.
func WithClientProtocol(clientProtocol ClientProtocol) ClientOption {
	return func(opts *clientOptions) {
		opts.clientProtocol = clientProtocol
	}
}

// NewHTTPClient returns a new Client backed by HTTP.
func NewHTTPClient(baseURL string, doer Doer, os ...ClientOption) Client {
	opts := clientOptions{clientProtocol: &grpcClientProtocol{doer: doer} /* gRPC with no options is the default ClientProtocol */}
	for _, o := range os {
		o(&opts)
	}
	return &client{
		baseURL:        strings.TrimRight(baseURL, "/"),
		clientProtocol: opts.clientProtocol,
	}
}

type client struct {
	baseURL        string
	clientProtocol ClientProtocol
}

func (c *client) NewCall(ctx context.Context, streamType StreamType, path string) (context.Context, StreamFunc) {
	spec := Specification{
		Type: streamType,
		Path: path,
	}

	// TODO(alex): This call makes HTTP-specific semantics bleed into generic API.
	// Fine for now, but we should change this. This would probably look something
	// like YARPC's Header abstraction.
	//
	// If we need to attach the headers here so that they're seen throughout the
	// interceptor chain, see the comment attached to the ClientProtocol implementation.
	ctx = NewCallContext(ctx, spec, make(http.Header), make(http.Header))
	sf := StreamFunc(func(ctx context.Context) Stream {
		return c.clientProtocol.NewStream(ctx, c.baseURL+"/"+path)
	})
	return ctx, sf
}

type clientOptions struct {
	clientProtocol ClientProtocol
}
