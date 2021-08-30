package rerpclocal

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rerpc/rerpc"
)

// ServeOption configures the server.
type ServeOption func(*serveOptions)

// Serve serves the plugins configured on the given handler.
func Serve(handler *Handler, opts ...ServeOption) {
	// We still need to translate stdin and stdout into the
	// shapes the handler expects.
	if len(os.Args) != 2 {
		// TODO(alex): For now we match on exactly two arguments.
		_, _ = os.Stderr.WriteString(fmt.Sprintf("expected exactly two arguments, but found %d\n", len(os.Args)))
		return
	}

	// TODO(alex): Make the timeout configurable on the server-side.
	// This is the only thing that could stop the process if it doesn't
	// complete on its own.
	//
	// How do we receive CallOptions from the client? We only have stdin,
	// so this would need to be potentially exposed as flags when the
	// host process invokes them? The same goes for other metadata like
	// Content-Type. For now it's hard-coded.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO(alex): The following structures should have real constructors,
	// and the StreamType and ContentType should be configurable (probably
	// compression, too).
	request := &Request{
		Type:        rerpc.StreamTypeUnary,
		Path:        os.Args[1],
		ContentType: rerpc.TypeProtoTwirp,
		Body:        os.Stdin,
	}

	response := &Response{
		Body: os.Stdout,
	}

	errChan := make(chan error, 1)
	go func() {
		// Execute the handler in a goroutine so that we can handle the timeout.
		errChan <- handler.Handle(ctx, request, response)
	}()

	select {
	case <-ctx.Done():
		if err := os.Stdin.Close(); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("shutting down: failed to close stdin: %v\n", err))
		}
		if err := os.Stdout.Close(); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("shutting down: failed to close stdout: %v\n", err))
		}
		os.Exit(1)
	case err := <-errChan:
		if err == nil {
			// The handler successfully completed.
			os.Exit(0)
		}
		_, _ = os.Stderr.WriteString(fmt.Sprintf("%v\n", err.Error()))
		os.Exit(1)
	}
}

type serveOptions struct{}
