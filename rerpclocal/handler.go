package rerpclocal

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/rerpc/rerpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// HandlerOption configures a Handler.
type HandlerOption func(*handlerOptions)

// Handler represents a local handler.
type Handler struct {
	mux    *rerpc.Mux
	codecs map[string]rerpc.Codec // TODO(alex): We probably want something a little more sophisticated than a map here.
}

// NewHandler returns a new Handler.
func NewHandler(mux *rerpc.Mux, opts ...HandlerOption) *Handler {
	options := handlerOptions{}
	for _, o := range opts {
		o(&options)
	}
	return &Handler{
		mux: mux,
		// TODO(alex): The codecs are written in-place for now, but they need to be exposed as options.
		codecs: map[string]rerpc.Codec{
			rerpc.TypeJSON: jsonProtobufCodec{
				marshaler:   protojson.MarshalOptions{UseProtoNames: true},
				unmarshaler: protojson.UnmarshalOptions{DiscardUnknown: true},
			},
			rerpc.TypeProtoTwirp: protobufCodec{},
		},
	}
}

// Handle executes the local handler for the given request.
func (h *Handler) Handle(ctx context.Context, req *Request, res *Response) error {
	streamType, err := h.mux.StreamType(req.Path)
	if err != nil {
		return err
	}

	spec := rerpc.Specification{
		Type:        streamType,
		Path:        req.Path,
		ContentType: req.ContentType,
	}

	codec, ok := h.codecs[req.ContentType]
	if !ok {
		return rerpc.Errorf(rerpc.CodeInvalidArgument, "unsupported content type %q", req.ContentType)
	}

	// TODO(alex): The NewHandlerContext constructor won't do as-is.
	// These headers need to be made generic so that HTTP semantics
	// aren't included here.
	ctx = rerpc.NewHandlerContext(ctx, spec, make(http.Header), make(http.Header))

	streamFunc := rerpc.StreamFunc(func(ctx context.Context) rerpc.Stream {
		return newServerStream(
			ctx,
			codec,
			req.Body,
			res.Body,
		)
	})
	return h.mux.Handle(ctx, streamFunc, spec.Path)
}

func newServerStream(
	ctx context.Context,
	codec rerpc.Codec,
	reader io.ReadCloser,
	writer io.WriteCloser,
) *serverStream {
	return &serverStream{
		ctx:    ctx,
		codec:  codec,
		reader: reader,
		writer: writer,
	}
}

type serverStream struct {
	ctx    context.Context
	codec  rerpc.Codec
	reader io.ReadCloser
	writer io.WriteCloser
}

func (s *serverStream) Context() context.Context {
	return s.ctx
}

func (s *serverStream) Receive(msg interface{}) error {
	// TODO(alex): We need to add pooling here.
	buf := bytes.NewBuffer(nil)
	n, err := buf.ReadFrom(s.reader)
	if err != nil {
		return rerpc.Wrap(rerpc.CodeUnknown, err)
	}
	if n == 0 {
		return nil
	}
	if err := s.codec.Unmarshal(buf.Bytes(), msg); err != nil {
		return err
	}
	return nil
}

func (s *serverStream) CloseReceive() error {
	if err := s.reader.Close(); err != nil {
		if errors.Is(err, os.ErrClosed) {
			// TODO(alex): Do we need to do anything else special here,
			// such as checking Stderr?
			return nil
		}
		if rerr, ok := rerpc.AsError(err); ok {
			return rerr
		}
		return rerpc.Wrap(rerpc.CodeUnknown, err)
	}
	return nil
}

func (s *serverStream) Send(msg interface{}) error {
	bytes, err := s.codec.Marshal(msg)
	if err != nil {
		return rerpc.Wrap(rerpc.CodeInternal, err)
	}
	if _, err := s.writer.Write(bytes); err != nil {
		if rerr, ok := rerpc.AsError(err); ok {
			return rerr
		}
		return rerpc.Wrap(rerpc.CodeUnknown, err)
	}
	return nil
}

func (s *serverStream) CloseSend(err error) error {
	if err == nil {
		return nil
	}
	// TODO(alex): We probably want to send errors
	// over Stderr here.
	return err
}

type handlerOptions struct{}
