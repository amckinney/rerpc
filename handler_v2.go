package rerpc

import (
	"context"
	"io"
	"net/http"
)

// HTTPHandlerOption configures some aspect of the handler.
type HTTPHandlerOption func(*httpHandlerOptions)

type httpHandlerOptions struct {
	protocols map[string]HTTPServerProtocol
}

// WithHTTPServerProtocols registers the given HTTPServerProtocols. The match is exact,
// with the special case that the content type "*" is the fallback HTTPServerProtocol used
// when nothing matches.
func WithHTTPServerProtocols(protocols ...HTTPServerProtocol) HTTPHandlerOption {
	return func(opts *httpHandlerOptions) {
		for _, protocol := range protocols {
			opts.protocols[protocol.ContentType()] = protocol
		}
	}
}

// HandlerV2 represents a generic handler.
type HandlerV2 interface {
	Handle(ctx context.Context, streamFunc StreamFunc, path string) error

	StreamType(path string) (StreamType, error)
}

// NewHTTPHandler returns a new http.Handler.
func NewHTTPHandler(handler HandlerV2, os ...HTTPHandlerOption) http.Handler {
	opts := httpHandlerOptions{protocols: defaultHTTPServerProtocols()}
	for _, o := range os {
		o(&opts)
	}
	return httpHandler{
		handler: handler,
		opts:    opts,
	}
}

// httpHandler implements a http.Handler by dispatching to the provided Route.
type httpHandler struct {
	handler HandlerV2
	opts    httpHandlerOptions
}

func (h httpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	protocol, ok := h.opts.protocols[r.Header.Get("Content-Type")]
	if !ok {
		protocol = h.opts.protocols["*"]
	}

	streamType, err := h.handler.StreamType(r.URL.Path)
	if err != nil {
		// TODO(alex): This 'is' checker could probably be un-exported to start
		// since it's largely used as an implementation detail / sentinel value.
		if IsUnknownPathError(err) {
			// In this case, we havent's sent any information over the stream,
			// so we need to manually set the response header here.
			rw.WriteHeader(http.StatusNotFound)
			io.WriteString(rw, err.Error())
			return
		}
	}

	// TODO(alex): We should trim down the rerpc.Specification so it doesn't
	// contain any Protocol-specific configuration values; the specification
	// should only contain Func metadata.
	//
	// If there's ever a need for Interceptors to access other Protocol-specific
	// configuration (e.g. MaxRequestSize), it can be added on an as-needed basis
	// (e.g. type-asserting an interface for a concrete type that we expose, or
	// something along those lines).
	//
	// Otherwise, Protocol-specific details bleed into the generic type structure,
	// which is not what we want.
	spec := Specification{
		Type:        streamType,
		Path:        r.URL.Path,
		ContentType: protocol.ContentType(),
	}

	// Instrument the stream's context.Context with Metadata.
	// All interceptors will receive the wrapped context.
	ctx := NewHandlerContext(r.Context(), spec, r.Header, rw.Header())
	r = r.WithContext(ctx)

	streamFunc, err := protocol.NewStream(streamType, rw, r)
	if err != nil {
		// The error will have already been sent over the stream,
		// so there's nothing to do here.
		//
		// TODO(alex): We might want to refactor this so that
		// if a certain set of structured errors are returned, we
		// can write the error to the ResponseWriter here (instead
		// of relying on specific implementations of NewStream to do so).
		return
	}
	if err := h.handler.Handle(ctx, streamFunc, spec.Path); err != nil {
		if IsUnknownPathError(err) {
			// In this case, we havent's sent any information over the stream,
			// so we need to manually set the response header here.
			rw.WriteHeader(http.StatusNotFound)
			io.WriteString(rw, err.Error())
			return
		}
	}
	return
}
