package rerpclocal

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/rerpc/rerpc"
)

// ClientOption configures a Client.
type ClientOption func(*clientOptions)

// Client represents a local client.
type Client struct {
	service string // e.g. the local binary (protoc-gen-go).
	codec   rerpc.Codec
}

// NewClient returns a new Client.
func NewClient(service string, opts ...ClientOption) *Client {
	options := clientOptions{}
	for _, o := range opts {
		o(&options)
	}
	return &Client{
		service: service,
		codec:   protobufCodec{}, // TODO(alex): The codec is written in-place for now, but it need to be exposed as an options.
	}
}

// NewCall issues a new RPC for the given path.
func (c *Client) NewCall(ctx context.Context, streamType rerpc.StreamType, path string) (context.Context, rerpc.StreamFunc) {
	// Set up the specification required for interceptors.
	spec := rerpc.Specification{
		Type: streamType,
		Path: path,
	}

	// TODO(alex): The NewCallContext constructor won't do as-is.
	// These headers need to be made generic so that HTTP semantics
	// aren't included here.
	ctx = rerpc.NewCallContext(ctx, spec, make(http.Header), make(http.Header))

	// Start a process for the given service.
	// This is equivalent to running the following:
	//
	//       ServiceName      |    RPC
	//  protoc-gen-go-service | /Generate
	//

	// TODO(alex): Stderr is currently unused, but we probably need this to handle errors.

	streamFunc := rerpc.StreamFunc(func(ctx context.Context) rerpc.Stream {
		command := exec.CommandContext(ctx, c.service, path /* Should the endpoint we're calling be the first (and only) argument? */)
		stdin, err := command.StdinPipe()
		if err != nil {
			// TODO(alex): How should we handle the error in these cases?
			// Should this be moved to the NewClient constructor, or should we expose an error in NewCall?
			// Maybe all of this should be encapsulated by the returned StreamFunc? (probably)
		}
		stdout, err := command.StdoutPipe()
		if err != nil {
			// TODO(alex): Handle error ...
		}
		go command.Run() // TODO(alex): This needs to be handled a little more gracefully, similar to what we do for the HTTP client stream.
		return newClientStream(
			ctx,
			c.codec,
			stdout,
			stdin,
		)
	})

	return ctx, streamFunc
}

func newClientStream(
	ctx context.Context,
	codec rerpc.Codec,
	reader io.ReadCloser,
	writer io.WriteCloser,
) *clientStream {
	return &clientStream{
		ctx:    ctx,
		codec:  codec,
		reader: reader,
		writer: writer,
	}
}

type clientStream struct {
	ctx    context.Context
	codec  rerpc.Codec
	reader io.ReadCloser
	writer io.WriteCloser
}

func (c *clientStream) Context() context.Context {
	return c.ctx
}

func (c *clientStream) Receive(msg interface{}) error {
	// TODO(alex): We need to add pooling here.
	buf := bytes.NewBuffer(nil)
	n, err := buf.ReadFrom(c.reader)
	if err != nil {
		return rerpc.Wrap(rerpc.CodeUnknown, err)
	}
	if n == 0 {
		return nil
	}
	if err := c.codec.Unmarshal(buf.Bytes(), msg); err != nil {
		return err
	}
	return nil
}

func (c *clientStream) CloseReceive() error {
	if err := c.reader.Close(); err != nil {
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

func (c *clientStream) Send(msg interface{}) error {
	bytes, err := c.codec.Marshal(msg)
	if err != nil {
		return rerpc.Wrap(rerpc.CodeInternal, err)
	}
	if _, err := c.writer.Write(bytes); err != nil {
		if rerr, ok := rerpc.AsError(err); ok {
			return rerr
		}
		return rerpc.Wrap(rerpc.CodeUnknown, err)
	}
	return nil
}

func (c *clientStream) CloseSend(_ error) error {
	if err := c.writer.Close(); err != nil {
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

type clientOptions struct{}
